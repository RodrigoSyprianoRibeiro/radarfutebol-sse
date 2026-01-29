package services

import (
	"encoding/json"
	"fmt"
	"log"

	"radarfutebol-sse/internal/models"
)

// Prefix usado pelo Laravel Redis (database 1)
const redisPrefix = "radarfutebolcom_database_"

// Chave JSON pura criada pelo Laravel para o Go
const eventosJsonKey = "eventos-painel-json"

// GetEventosPainelFiltrado busca eventos do Redis e aplica filtros
func GetEventosPainelFiltrado(filtro *models.Filtro) ([]byte, error) {
	// Busca eventos raw do Redis
	eventos, err := getEventosFromRedis()
	if err != nil {
		return nil, err
	}

	if len(eventos) == 0 {
		return json.Marshal(&models.PainelResponse{
			Eventos: []*models.Evento{},
			Counts:  models.Counts{Live: 0, Total: 0, Gols: 0},
		})
	}

	// Busca preferencias do usuario
	var prefs *PreferenciasUsuario
	if filtro.IdUsuario > 0 {
		prefs, err = GetPreferenciasUsuarioCompletas(filtro.IdUsuario)
		if err != nil {
			log.Printf("Erro ao buscar preferencias do usuario %d: %v", filtro.IdUsuario, err)
		}
	}

	// Aplica filtros
	response, err := FiltrarEventosPainel(eventos, filtro, prefs)
	if err != nil {
		return nil, err
	}

	return json.Marshal(response)
}

// GetEventosHomeFiltrado busca eventos do Redis e aplica filtros para Home
func GetEventosHomeFiltrado(filtro *models.Filtro) ([]byte, error) {
	// Busca eventos raw do Redis
	eventos, err := getEventosFromRedis()
	if err != nil {
		return nil, err
	}

	if len(eventos) == 0 {
		return json.Marshal(&models.HomeResponse{
			Campeonatos: []*models.Campeonato{},
			Counts:      models.Counts{Live: 0, Total: 0, Gols: 0},
		})
	}

	// Busca preferencias do usuario
	var prefs *PreferenciasUsuario
	if filtro.IdUsuario > 0 {
		prefs, err = GetPreferenciasUsuarioCompletas(filtro.IdUsuario)
		if err != nil {
			log.Printf("Erro ao buscar preferencias do usuario %d: %v", filtro.IdUsuario, err)
		}
	}

	// Aplica filtros
	response, err := FiltrarEventosHome(eventos, filtro, prefs)
	if err != nil {
		return nil, err
	}

	return json.Marshal(response)
}

// getEventosFromRedis busca eventos do cache Redis
func getEventosFromRedis() ([]*models.Evento, error) {
	// Busca da chave JSON pura criada pelo Laravel para o Go
	// Chave: eventos-painel-json (sem prefixo, JSON puro)
	data, err := GetString(eventosJsonKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos do Redis: %w", err)
	}

	if data == "" {
		log.Println("SSE: Chave eventos-painel-json nao encontrada no Redis")
		return nil, nil
	}

	var eventos []*models.Evento
	if err := json.Unmarshal([]byte(data), &eventos); err != nil {
		return nil, fmt.Errorf("erro ao decodificar eventos JSON: %w", err)
	}

	log.Printf("SSE: Carregados %d eventos do Redis", len(eventos))
	return eventos, nil
}

// GetEventosPainelRaw retorna JSON raw do Redis (sem filtros - fallback)
func GetEventosPainelRaw() ([]byte, error) {
	data, err := GetString(redisPrefix + "sse:painel")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar sse:painel do Redis: %w", err)
	}

	if data == "" {
		log.Println("SSE: Chave sse:painel nao encontrada no Redis")
		return []byte(`{"eventos":[],"counts":{"live":0,"total":0,"gols":0}}`), nil
	}

	return []byte(data), nil
}

// GetEventosHomeRaw retorna JSON raw do Redis (sem filtros - fallback)
func GetEventosHomeRaw() ([]byte, error) {
	data, err := GetString(redisPrefix + "sse:home")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar sse:home do Redis: %w", err)
	}

	if data == "" {
		log.Println("SSE: Chave sse:home nao encontrada no Redis")
		return []byte(`{"campeonatos":[],"counts":{"live":0,"total":0,"gols":0}}`), nil
	}

	return []byte(data), nil
}

// ValidateJSON verifica se o JSON e valido
func ValidateJSON(data []byte) bool {
	return json.Valid(data)
}
