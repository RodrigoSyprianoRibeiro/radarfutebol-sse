package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"radarfutebol-sse/internal/config"
	"radarfutebol-sse/internal/handlers"
	"radarfutebol-sse/internal/services"
)

func main() {
	// Carrega .env (ignora erro se nao existir)
	godotenv.Load()

	// Configura GOMAXPROCS para usar todos os CPUs disponiveis
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	log.Printf("Iniciando servidor SSE Go (CPUs: %d, GOMAXPROCS: %d)...", numCPU, runtime.GOMAXPROCS(0))

	// Carrega configuracoes
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuracoes: %v", err)
	}

	// Inicializa MySQL (opcional - dados vêm do Redis)
	if err := services.InitMySQL(cfg.MySQL); err != nil {
		log.Printf("Aviso: MySQL não disponível: %v", err)
		log.Println("Continuando sem MySQL (dados vêm do Redis)")
	} else {
		log.Println("MySQL inicializado")
		// Inicializa cache de autenticacao (requer MySQL)
		services.InitAuthCache()
	}

	// Inicializa Redis (database 0 - cache principal)
	if err := services.InitRedis(cfg.Redis); err != nil {
		log.Fatalf("Erro ao inicializar Redis: %v", err)
	}

	// Inicializa Redis preferencias (database 2 - favoritos do usuario)
	if err := services.InitRedisPreferencias(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password); err != nil {
		log.Printf("Aviso: Erro ao inicializar Redis preferencias: %v", err)
		// Nao fatal - continua sem preferencias
	}

	// Inicia o Broadcaster (cache em memoria + atualizacao periodica)
	broadcaster := services.GetBroadcaster()
	broadcaster.Start()
	defer broadcaster.Stop()
	log.Println("Broadcaster iniciado (cache em memoria)")

	// Cria handler SSE
	sseHandler := handlers.NewSSEHandler()

	// Configura rotas
	mux := http.NewServeMux()
	sseHandler.RegisterRoutes(mux)

	// Middleware de CORS e logging
	handler := corsMiddleware(mux)

	// Configura servidor HTTP otimizado para SSE de alta concorrencia
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler,

		// Timeouts otimizados para SSE
		ReadTimeout:       30 * time.Second,  // Tempo para ler request completo
		ReadHeaderTimeout: 10 * time.Second,  // Tempo para ler headers
		WriteTimeout:      0,                 // Sem timeout para SSE (conexoes longas)
		IdleTimeout:       120 * time.Second, // Conexoes idle

		// Limites de tamanho
		MaxHeaderBytes: 1 << 16, // 64KB max headers
	}

	// Canal para shutdown graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Inicia servidor em goroutine
	go func() {
		log.Printf("Servidor SSE rodando na porta %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguarda sinal de shutdown
	sig := <-stop
	log.Printf("Recebido sinal de shutdown (%v)...", sig)

	// Graceful shutdown com timeout de 30 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Para de aceitar novas conexoes e espera as existentes terminarem
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Erro no shutdown graceful: %v", err)
	}

	// Para o broadcaster
	broadcaster.Stop()

	log.Println("Servidor SSE encerrado graciosamente")
}

// corsMiddleware adiciona headers CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
