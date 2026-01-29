package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"radarfutebol-sse/internal/config"
)

var rdb *redis.Client
var ctx = context.Background()

// InitRedis inicializa a conexao com o Redis
func InitRedis(cfg config.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0, // Database 0 para cache do oraculo
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("erro ao conectar Redis: %w", err)
	}

	log.Println("Redis conectado com sucesso")
	return nil
}

// GetRedis retorna o cliente Redis
func GetRedis() *redis.Client {
	return rdb
}

// GetOraculoCache busca dados do oraculo do cache Redis
func GetOraculoCache(idWilliamhill string) (map[string]interface{}, error) {
	if idWilliamhill == "" {
		return nil, nil
	}

	key := fmt.Sprintf("oraculo-cache:idJogo-%s", idWilliamhill)
	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Chave nao existe
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cache: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar cache: %w", err)
	}

	return result, nil
}

// rdbPrefs cliente Redis para preferencias (database 2)
var rdbPrefs *redis.Client

// InitRedisPreferencias inicializa conexao Redis para preferencias
func InitRedisPreferencias(host string, port int, password string) error {
	rdbPrefs = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       2, // Database 2 para preferencias
	})

	_, err := rdbPrefs.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("erro ao conectar Redis preferencias: %w", err)
	}

	log.Println("Redis preferencias (DB 2) conectado com sucesso")
	return nil
}

// GetPreferenciasUsuarioCompletas busca campeonatos e jogos favoritos
func GetPreferenciasUsuarioCompletas(userID int) (*PreferenciasUsuario, error) {
	if rdbPrefs == nil {
		return nil, fmt.Errorf("redis preferencias nao inicializado")
	}

	prefs := &PreferenciasUsuario{
		CampeonatosFavoritos: make(map[string]bool),
		JogosFavoritos:       make(map[string]bool),
	}

	// Busca campeonatos favoritos
	// Chave: preferencias:campeonatos-favoritos-{userId}
	keyCamp := fmt.Sprintf("preferencias:campeonatos-favoritos-%d", userID)
	dataCamp, err := rdbPrefs.Get(ctx, keyCamp).Result()
	if err == nil && dataCamp != "" {
		var campFavs map[string]bool
		if err := json.Unmarshal([]byte(dataCamp), &campFavs); err == nil {
			prefs.CampeonatosFavoritos = campFavs
		}
	}

	// Busca jogos favoritos
	// Chave: preferencias:jogos-favoritos-{userId}
	keyJogos := fmt.Sprintf("preferencias:jogos-favoritos-%d", userID)
	dataJogos, err := rdbPrefs.Get(ctx, keyJogos).Result()
	if err == nil && dataJogos != "" {
		var jogosFavs map[string]bool
		if err := json.Unmarshal([]byte(dataJogos), &jogosFavs); err == nil {
			prefs.JogosFavoritos = jogosFavs
		}
	}

	return prefs, nil
}

// GetPreferenciasUsuario busca preferencias do usuario do Redis (database 2)
func GetPreferenciasUsuario(userID int) (map[string]interface{}, error) {
	if rdbPrefs == nil {
		return nil, fmt.Errorf("redis preferencias nao inicializado")
	}

	key := fmt.Sprintf("preferencias-usuario:%d", userID)
	data, err := rdbPrefs.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar preferencias: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar preferencias: %w", err)
	}

	return result, nil
}

// GetString busca uma string do Redis
func GetString(key string) (string, error) {
	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return data, nil
}

// SetString salva uma string no Redis
func SetString(key string, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}
