package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"radarfutebol-sse/internal/models"
	"radarfutebol-sse/internal/services"
)

// SSEHandler gerencia conexoes SSE
type SSEHandler struct {
	connections int64 // atomic counter
}

// NewSSEHandler cria um novo handler SSE
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{}
}

// RegisterRoutes registra as rotas SSE
func (h *SSEHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/sse/health", h.handleHealth)
	mux.HandleFunc("/sse/painel", h.handlePainel)
	mux.HandleFunc("/sse/home", h.handleHome)
	mux.HandleFunc("/stats", h.handleStats)
}

// handleHealth retorna status do servidor
func (h *SSEHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"connections": atomic.LoadInt64(&h.connections),
		"timestamp":   time.Now().Unix(),
	})
}

// handleStats retorna estatisticas do servidor
func (h *SSEHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connections": atomic.LoadInt64(&h.connections),
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

	// Extrai filtros da query string
	filtro := models.ParseFiltroFromRequest(r)

	// Incrementa contador de conexoes (atomic)
	connCount := atomic.AddInt64(&h.connections, 1)
	log.Printf("SSE %s: Nova conexao (user=%d) - Total: %d", endpoint, filtro.IdUsuario, connCount)

	// Decrementa ao fechar
	defer func() {
		newCount := atomic.AddInt64(&h.connections, -1)
		log.Printf("SSE %s: Conexao fechada (user=%d) - Total: %d", endpoint, filtro.IdUsuario, newCount)
	}()

	// Envia retry interval (10 segundos)
	fmt.Fprintf(w, "retry: 10000\n\n")
	flusher.Flush()

	// Ticker para enviar updates a cada 2 segundos
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Canal para detectar quando cliente desconecta
	ctx := r.Context()

	// Envia primeiro update imediatamente
	h.sendUpdateFiltrado(w, flusher, endpoint, filtro)

	for {
		select {
		case <-ctx.Done():
			// Cliente desconectou
			return
		case <-ticker.C:
			// Envia update periodico com filtros
			h.sendUpdateFiltrado(w, flusher, endpoint, filtro)
		}
	}
}

// sendUpdateFiltrado envia um update SSE com filtros aplicados
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
