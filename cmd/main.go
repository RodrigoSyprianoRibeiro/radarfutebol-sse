package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	log.Println("Iniciando servidor SSE Go...")

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
	}

	// Inicializa Redis (database 0 - cache principal)
	if err := services.InitRedis(cfg.Redis); err != nil {
		log.Fatalf("Erro ao inicializar Redis: %v", err)
	}
	log.Println("Redis inicializado")

	// Inicializa Redis preferencias (database 2 - favoritos do usuario)
	if err := services.InitRedisPreferencias(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password); err != nil {
		log.Printf("Aviso: Erro ao inicializar Redis preferencias: %v", err)
		// Nao fatal - continua sem preferencias
	}

	// Cria handler SSE
	sseHandler := handlers.NewSSEHandler()

	// Configura rotas
	mux := http.NewServeMux()
	sseHandler.RegisterRoutes(mux)

	// Middleware de CORS
	handler := corsMiddleware(mux)

	// Configura servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // Sem timeout para SSE
		IdleTimeout:  120 * time.Second,
	}

	// Canal para shutdown graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Inicia servidor em goroutine
	go func() {
		log.Printf("Servidor SSE rodando na porta %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguarda sinal de shutdown
	<-stop
	log.Println("Recebido sinal de shutdown...")

	// Fecha conexoes graciosamente
	log.Println("Servidor SSE encerrado")
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
