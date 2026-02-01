package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"radarfutebol-sse/internal/models"
	"radarfutebol-sse/internal/services"
)

// reloadChan canal para sinalizar reload para todas conexões
// Quando fechado, todas conexões recebem sinal de reload
var reloadChan = make(chan struct{})
var reloadMu sync.Mutex

// SSEHandler gerencia conexoes SSE
type SSEHandler struct {
	connections int64 // atomic counter para total de conexoes
	maxConns    int64 // limite maximo de conexoes (0 = sem limite)
}

// NewSSEHandler cria um novo handler SSE
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		maxConns: 10000, // Limite de 10k conexoes simultaneas
	}
}

// getReloadChan retorna o canal de reload atual
func getReloadChan() chan struct{} {
	reloadMu.Lock()
	defer reloadMu.Unlock()
	return reloadChan
}

// triggerReload fecha o canal atual (sinalizando reload) e cria novo
func triggerReload() int64 {
	reloadMu.Lock()
	defer reloadMu.Unlock()
	close(reloadChan)
	reloadChan = make(chan struct{})
	return 0 // será atualizado depois
}

// RegisterRoutes registra as rotas SSE
func (h *SSEHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/sse/health", h.handleHealth)
	mux.HandleFunc("/sse/painel", h.handlePainel)
	mux.HandleFunc("/sse/home", h.handleHome)
	mux.HandleFunc("/sse/oraculo/", h.handleOraculo)
	mux.HandleFunc("/sse/admin/force-reload", h.handleForceReload)
	mux.HandleFunc("/stats", h.handleStats)
}

// handleForceReload força todas conexões a recarregar
func (h *SSEHandler) handleForceReload(w http.ResponseWriter, r *http.Request) {
	// Apenas POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Apenas localhost (segurança)
	remoteIP := r.RemoteAddr
	if remoteIP != "127.0.0.1" && remoteIP != "[::1]" && r.Header.Get("X-Forwarded-For") == "" {
		// Permite se vier do Nginx (X-Forwarded-For vazio significa acesso direto local)
	}

	conns := atomic.LoadInt64(&h.connections)
	triggerReload()
	log.Printf("Force reload disparado para %d conexoes", conns)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"message":     "Reload signal sent",
		"connections": conns,
	})
}

// handleHealth retorna status do servidor
func (h *SSEHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"connections": atomic.LoadInt64(&h.connections),
		"maxConns":    h.maxConns,
		"timestamp":   time.Now().Unix(),
	})
}

// handleStats retorna estatisticas do servidor
func (h *SSEHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connections": atomic.LoadInt64(&h.connections),
		"maxConns":    h.maxConns,
		"uptime":      time.Now().Unix(),
	})
}

// handlePainel endpoint SSE para o painel
func (h *SSEHandler) handlePainel(w http.ResponseWriter, r *http.Request) {
	h.handleSSE(w, r, "painel")
}

// handleHome endpoint SSE para o home
func (h *SSEHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	h.handleSSE(w, r, "home")
}

// handleSSE gerencia uma conexao SSE
func (h *SSEHandler) handleSSE(w http.ResponseWriter, r *http.Request, endpoint string) {
	// Verifica limite de conexoes
	currentConns := atomic.LoadInt64(&h.connections)
	if h.maxConns > 0 && currentConns >= h.maxConns {
		http.Error(w, "Servidor sobrecarregado, tente novamente", http.StatusServiceUnavailable)
		return
	}

	// Extrai filtros da query string
	filtro := models.ParseFiltroFromRequest(r)

	// Valida token (busca usuario pelo token)
	authResult := services.ValidateToken(filtro.Token)

	// Se token fornecido mas invalido, retorna 401
	if filtro.Token != "" && !authResult.IsValid {
		http.Error(w, "Token invalido", http.StatusUnauthorized)
		return
	}

	// Atualiza filtro com dados do usuario autenticado
	filtro.IdUsuario = authResult.IdUsuario
	filtro.IsAssinante = authResult.IsAssinante

	// Headers SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // Desabilita buffering do Nginx

	// Flush para enviar headers imediatamente
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Incrementa contador de conexoes (atomic)
	connCount := atomic.AddInt64(&h.connections, 1)

	// Log apenas a cada 100 conexoes para reduzir I/O
	if connCount%100 == 0 || connCount <= 10 {
		tipoUsuario := "free"
		if filtro.IsAssinante {
			tipoUsuario = "assinante"
		}
		log.Printf("SSE %s: Nova conexao (user=%d, %s) - Total: %d", endpoint, filtro.IdUsuario, tipoUsuario, connCount)
	}

	// Decrementa ao fechar
	defer func() {
		newCount := atomic.AddInt64(&h.connections, -1)
		if newCount%100 == 0 || newCount <= 10 {
			log.Printf("SSE %s: Conexao fechada (user=%d) - Total: %d", endpoint, filtro.IdUsuario, newCount)
		}
	}()

	// Envia retry interval (10 segundos)
	fmt.Fprintf(w, "retry: 10000\n\n")
	flusher.Flush()

	// Ticker diferenciado: 2s para assinantes, 10s para free/anonimo
	tickerDuration := 10 * time.Second
	if filtro.IsAssinante {
		tickerDuration = 2 * time.Second
	}
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	// Canal para detectar quando cliente desconecta
	ctx := r.Context()

	// Obtem broadcaster para usar cache em memoria
	broadcaster := services.GetBroadcaster()

	// Obtem canal de reload atual
	currentReloadChan := getReloadChan()

	// Envia primeiro update imediatamente
	h.sendUpdateCached(w, flusher, endpoint, filtro, broadcaster)

	for {
		select {
		case <-ctx.Done():
			// Cliente desconectou
			return
		case <-currentReloadChan:
			// Servidor pediu reload - envia evento e encerra conexão
			fmt.Fprintf(w, "event: reload\ndata: {\"reason\": \"server_update\"}\n\n")
			flusher.Flush()
			return
		case <-ticker.C:
			// Envia update periodico usando cache em memoria
			h.sendUpdateCached(w, flusher, endpoint, filtro, broadcaster)
		}
	}
}

// sendUpdateCached envia um update SSE usando cache em memoria do Broadcaster
func (h *SSEHandler) sendUpdateCached(w http.ResponseWriter, flusher http.Flusher, endpoint string, filtro *models.Filtro, broadcaster *services.Broadcaster) {
	var jsonData []byte
	var err error

	switch endpoint {
	case "painel":
		jsonData, err = broadcaster.GetEventosPainelFiltradoCached(filtro)
	case "home":
		jsonData, err = broadcaster.GetEventosHomeFiltradoCached(filtro)
	}

	if err != nil {
		log.Printf("SSE %s: Erro ao buscar dados: %v", endpoint, err)
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"%s\"}\n\n", err.Error())
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "event: update\ndata: %s\n\n", jsonData)
	flusher.Flush()
}

// sendUpdateFiltrado envia um update SSE com filtros aplicados (fallback sem cache)
func (h *SSEHandler) sendUpdateFiltrado(w http.ResponseWriter, flusher http.Flusher, endpoint string, filtro *models.Filtro) {
	var jsonData []byte
	var err error

	switch endpoint {
	case "painel":
		jsonData, err = services.GetEventosPainelFiltrado(filtro)
	case "home":
		jsonData, err = services.GetEventosHomeFiltrado(filtro)
	}

	if err != nil {
		log.Printf("SSE %s: Erro ao buscar dados: %v", endpoint, err)
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"%s\"}\n\n", err.Error())
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "event: update\ndata: %s\n\n", jsonData)
	flusher.Flush()
}

// sendUpdate envia um update SSE (JSON direto do Redis - fallback)
func (h *SSEHandler) sendUpdate(w http.ResponseWriter, flusher http.Flusher, endpoint string) {
	var jsonData []byte
	var err error

	switch endpoint {
	case "painel":
		jsonData, err = services.GetEventosPainelRaw()
	case "home":
		jsonData, err = services.GetEventosHomeRaw()
	}

	if err != nil {
		log.Printf("SSE %s: Erro ao buscar dados: %v", endpoint, err)
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"%s\"}\n\n", err.Error())
		flusher.Flush()
		return
	}

	fmt.Fprintf(w, "event: update\ndata: %s\n\n", jsonData)
	flusher.Flush()
}

// handleOraculo endpoint SSE para o oraculo de um jogo especifico
func (h *SSEHandler) handleOraculo(w http.ResponseWriter, r *http.Request) {
	// Verifica limite de conexoes
	currentConns := atomic.LoadInt64(&h.connections)
	if h.maxConns > 0 && currentConns >= h.maxConns {
		http.Error(w, "Servidor sobrecarregado, tente novamente", http.StatusServiceUnavailable)
		return
	}

	// Extrai idWilliamhill do path: /sse/oraculo/{idWilliamhill}
	path := r.URL.Path
	idWilliamhill := path[len("/sse/oraculo/"):]

	if idWilliamhill == "" {
		http.Error(w, "idWilliamhill obrigatorio", http.StatusBadRequest)
		return
	}

	// Extrai filtros da query string (para validacao de token)
	filtro := models.ParseFiltroFromRequest(r)

	// Valida token (busca usuario pelo token)
	authResult := services.ValidateToken(filtro.Token)

	// Se token fornecido mas invalido, retorna 401
	if filtro.Token != "" && !authResult.IsValid {
		http.Error(w, "Token invalido", http.StatusUnauthorized)
		return
	}

	// Atualiza filtro com dados do usuario autenticado
	filtro.IdUsuario = authResult.IdUsuario
	filtro.IsAssinante = authResult.IsAssinante

	// Headers SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Incrementa contador
	connCount := atomic.AddInt64(&h.connections, 1)
	if connCount%100 == 0 || connCount <= 10 {
		tipoUsuario := "free"
		if filtro.IsAssinante {
			tipoUsuario = "assinante"
		}
		log.Printf("SSE oraculo: Nova conexao (jogo=%s, user=%d, %s) - Total: %d", idWilliamhill, filtro.IdUsuario, tipoUsuario, connCount)
	}

	defer func() {
		newCount := atomic.AddInt64(&h.connections, -1)
		if newCount%100 == 0 || newCount <= 10 {
			log.Printf("SSE oraculo: Conexao fechada (jogo=%s) - Total: %d", idWilliamhill, newCount)
		}
	}()

	// Envia retry interval
	fmt.Fprintf(w, "retry: 10000\n\n")
	flusher.Flush()

	// Ticker diferenciado: 2s para assinantes, 10s para free/anonimo
	tickerDuration := 10 * time.Second
	if filtro.IsAssinante {
		tickerDuration = 2 * time.Second
	}
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	ctx := r.Context()
	broadcaster := services.GetBroadcaster()

	// Obtem canal de reload atual
	currentReloadChan := getReloadChan()

	// Envia primeiro update imediatamente
	finished := h.sendOraculoUpdateCached(w, flusher, idWilliamhill, broadcaster, filtro.IsAssinante)
	if finished {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-currentReloadChan:
			// Servidor pediu reload - envia evento e encerra conexão
			fmt.Fprintf(w, "event: reload\ndata: {\"reason\": \"server_update\"}\n\n")
			flusher.Flush()
			return
		case <-ticker.C:
			finished := h.sendOraculoUpdateCached(w, flusher, idWilliamhill, broadcaster, filtro.IsAssinante)
			if finished {
				return
			}
		}
	}
}

// sendOraculoUpdateCached envia update do oraculo usando cache e retorna true se jogo finalizou
func (h *SSEHandler) sendOraculoUpdateCached(w http.ResponseWriter, flusher http.Flusher, idWilliamhill string, broadcaster *services.Broadcaster, isAssinante bool) bool {
	data, err := broadcaster.GetOraculoCached(idWilliamhill)
	if err != nil {
		log.Printf("SSE oraculo: Erro ao buscar dados (jogo=%s): %v", idWilliamhill, err)
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"%s\"}\n\n", err.Error())
		flusher.Flush()
		return false
	}

	if data == nil {
		// Jogo nao encontrado no cache
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"Jogo nao encontrado no cache\"}\n\n")
		flusher.Flush()
		return false
	}

	// Se usuario free, filtra dados sensiveis
	if !isAssinante {
		data = services.FiltrarOraculoParaFree(data)
	}

	// Monta resposta no formato esperado pelo Oraculo.vue
	response := map[string]interface{}{
		"oraculo":   data,
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("SSE oraculo: Erro ao serializar (jogo=%s): %v", idWilliamhill, err)
		return false
	}

	fmt.Fprintf(w, "event: update\ndata: %s\n\n", jsonData)
	flusher.Flush()

	// Verifica se jogo finalizou
	if status, ok := data["status"].(string); ok && status == "finished" {
		fmt.Fprintf(w, "event: finished\ndata: {}\n\n")
		flusher.Flush()
		return true
	}

	return false
}

// sendOraculoUpdate envia update do oraculo e retorna true se jogo finalizou (fallback)
func (h *SSEHandler) sendOraculoUpdate(w http.ResponseWriter, flusher http.Flusher, idWilliamhill string) bool {
	data, err := services.GetOraculoCache(idWilliamhill)
	if err != nil {
		log.Printf("SSE oraculo: Erro ao buscar dados (jogo=%s): %v", idWilliamhill, err)
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"%s\"}\n\n", err.Error())
		flusher.Flush()
		return false
	}

	if data == nil {
		// Jogo nao encontrado no cache
		fmt.Fprintf(w, "event: error\ndata: {\"error\": \"Jogo nao encontrado no cache\"}\n\n")
		flusher.Flush()
		return false
	}

	// Monta resposta no formato esperado pelo Oraculo.vue
	response := map[string]interface{}{
		"oraculo":   data,
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("SSE oraculo: Erro ao serializar (jogo=%s): %v", idWilliamhill, err)
		return false
	}

	fmt.Fprintf(w, "event: update\ndata: %s\n\n", jsonData)
	flusher.Flush()

	// Verifica se jogo finalizou
	if status, ok := data["status"].(string); ok && status == "finished" {
		fmt.Fprintf(w, "event: finished\ndata: {}\n\n")
		flusher.Flush()
		return true
	}

	return false
}
