package services

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"radarfutebol-sse/internal/models"
)

// Broadcaster gerencia cache em memoria e broadcast para todos os clientes
// Evita que cada conexao consulte o Redis individualmente
type Broadcaster struct {
	mu sync.RWMutex

	// Cache de eventos em memoria (atualizado a cada 2s por uma unica goroutine)
	eventosCache    []*models.Evento
	eventosCacheAt  time.Time
	eventosCacheTTL time.Duration

	// Cache de oraculo por jogo
	oraculoCache   map[string]*OraculoCache
	oraculoCacheMu sync.RWMutex

	// Controle
	stopChan chan struct{}
	running  bool
}

// OraculoCache cache individual de cada jogo do oraculo
type OraculoCache struct {
	Data      map[string]interface{}
	UpdatedAt time.Time
}

var broadcaster *Broadcaster
var broadcasterOnce sync.Once

// GetBroadcaster retorna singleton do broadcaster
func GetBroadcaster() *Broadcaster {
	broadcasterOnce.Do(func() {
		broadcaster = &Broadcaster{
			eventosCacheTTL: 2 * time.Second,
			oraculoCache:    make(map[string]*OraculoCache),
			stopChan:        make(chan struct{}),
		}
	})
	return broadcaster
}

// Start inicia o broadcaster em background
func (b *Broadcaster) Start() {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return
	}
	b.running = true
	b.mu.Unlock()

	// Goroutine que atualiza cache de eventos a cada 2 segundos
	go b.eventosUpdater()

	// Goroutine que limpa cache de oraculo antigo
	go b.oraculoCleaner()

	log.Println("Broadcaster iniciado")
}

// Stop para o broadcaster
func (b *Broadcaster) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.running {
		close(b.stopChan)
		b.running = false
	}
}

// eventosUpdater atualiza cache de eventos periodicamente
func (b *Broadcaster) eventosUpdater() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Primeira carga imediata
	b.refreshEventosCache()

	for {
		select {
		case <-b.stopChan:
			return
		case <-ticker.C:
			b.refreshEventosCache()
		}
	}
}

// refreshEventosCache atualiza o cache de eventos do Redis
func (b *Broadcaster) refreshEventosCache() {
	eventos, err := getEventosFromRedis()
	if err != nil {
		log.Printf("Broadcaster: erro ao buscar eventos: %v", err)
		return
	}

	b.mu.Lock()
	b.eventosCache = eventos
	b.eventosCacheAt = time.Now()
	b.mu.Unlock()
}

// GetEventosCache retorna eventos do cache em memoria
func (b *Broadcaster) GetEventosCache() []*models.Evento {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Retorna copia para evitar race conditions
	if b.eventosCache == nil {
		return nil
	}

	return b.eventosCache
}

// GetEventosPainelFiltradoCached aplica filtros sobre cache em memoria
func (b *Broadcaster) GetEventosPainelFiltradoCached(filtro *models.Filtro) ([]byte, error) {
	eventos := b.GetEventosCache()

	if len(eventos) == 0 {
		return json.Marshal(&models.PainelResponse{
			Eventos: []*models.Evento{},
			Counts:  models.Counts{Live: 0, Total: 0, Gols: 0},
		})
	}

	// Busca preferencias do usuario (ainda do Redis, mas e pequeno)
	var prefs *PreferenciasUsuario
	var err error
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

// GetEventosHomeFiltradoCached aplica filtros sobre cache em memoria
func (b *Broadcaster) GetEventosHomeFiltradoCached(filtro *models.Filtro) ([]byte, error) {
	eventos := b.GetEventosCache()

	if len(eventos) == 0 {
		return json.Marshal(&models.HomeResponse{
			Campeonatos: []*models.Campeonato{},
			Counts:      models.Counts{Live: 0, Total: 0, Gols: 0},
		})
	}

	// Busca preferencias do usuario
	var prefs *PreferenciasUsuario
	var err error
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

// GetOraculoCached busca oraculo do cache ou Redis
func (b *Broadcaster) GetOraculoCached(idWilliamhill string) (map[string]interface{}, error) {
	b.oraculoCacheMu.RLock()
	cached, exists := b.oraculoCache[idWilliamhill]
	b.oraculoCacheMu.RUnlock()

	// Se cache existe e tem menos de 2 segundos, usa cache
	if exists && time.Since(cached.UpdatedAt) < 2*time.Second {
		return cached.Data, nil
	}

	// Busca do Redis
	data, err := GetOraculoCache(idWilliamhill)
	if err != nil {
		return nil, err
	}

	// Atualiza cache
	if data != nil {
		b.oraculoCacheMu.Lock()
		b.oraculoCache[idWilliamhill] = &OraculoCache{
			Data:      data,
			UpdatedAt: time.Now(),
		}
		b.oraculoCacheMu.Unlock()
	}

	return data, nil
}

// oraculoCleaner limpa cache de oraculo antigo periodicamente
func (b *Broadcaster) oraculoCleaner() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-b.stopChan:
			return
		case <-ticker.C:
			b.cleanOldOraculoCache()
		}
	}
}

// cleanOldOraculoCache remove entradas antigas do cache de oraculo
func (b *Broadcaster) cleanOldOraculoCache() {
	b.oraculoCacheMu.Lock()
	defer b.oraculoCacheMu.Unlock()

	now := time.Now()
	for id, cached := range b.oraculoCache {
		// Remove se tem mais de 10 minutos sem uso
		if now.Sub(cached.UpdatedAt) > 10*time.Minute {
			delete(b.oraculoCache, id)
		}
	}
}
