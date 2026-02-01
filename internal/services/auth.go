package services

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

// authCacheEntry entrada do cache de autenticacao
type authCacheEntry struct {
	idUsuario   int
	isValid     bool
	isAssinante bool
	teamId      int
	expiresAt   time.Time
}

// authCache cache de validacao de tokens
var authCache sync.Map

// authCacheTTL tempo de vida do cache (5 minutos)
const authCacheTTL = 5 * time.Minute

// InitAuthCache inicializa o cache de autenticacao
// Inicia uma goroutine para limpeza periodica de entradas expiradas
func InitAuthCache() {
	go cleanupAuthCache()
	log.Println("Cache de autenticacao inicializado")
}

// cleanupAuthCache limpa entradas expiradas do cache a cada 5 minutos
func cleanupAuthCache() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		authCache.Range(func(key, value interface{}) bool {
			entry := value.(authCacheEntry)
			if now.After(entry.expiresAt) {
				authCache.Delete(key)
			}
			return true
		})
	}
}

// AuthResult resultado da validacao de token
type AuthResult struct {
	IdUsuario   int
	IsValid     bool
	IsAssinante bool
	TeamId      int
}

// ValidateToken valida o token e retorna dados do usuario
// Busca apenas pelo token, sem precisar do idUsuario
func ValidateToken(token string) AuthResult {
	// Sem token - usuario anonimo
	if token == "" {
		return AuthResult{
			IdUsuario:   0,
			IsValid:     true, // Anonimo sempre valido
			IsAssinante: false,
			TeamId:      0,
		}
	}

	// Verifica cache primeiro (chave é o proprio token)
	cacheKey := getCacheKey(token)
	if cached, ok := authCache.Load(cacheKey); ok {
		entry := cached.(authCacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return AuthResult{
				IdUsuario:   entry.idUsuario,
				IsValid:     entry.isValid,
				IsAssinante: entry.isAssinante,
				TeamId:      entry.teamId,
			}
		}
		// Cache expirado, remove
		authCache.Delete(cacheKey)
	}

	// Consulta MySQL
	result := queryToken(token)

	// Salva no cache
	authCache.Store(cacheKey, authCacheEntry{
		idUsuario:   result.IdUsuario,
		isValid:     result.IsValid,
		isAssinante: result.IsAssinante,
		teamId:      result.TeamId,
		expiresAt:   time.Now().Add(authCacheTTL),
	})

	return result
}

// getCacheKey gera a chave do cache para o token
func getCacheKey(token string) string {
	// Usa apenas os primeiros 20 caracteres do token para a chave
	// (suficiente para identificacao sem usar muita memoria)
	if len(token) > 20 {
		return token[:20]
	}
	return token
}

// queryToken consulta o banco para validar token (busca apenas pelo token)
func queryToken(token string) AuthResult {
	if db == nil {
		log.Printf("Auth: MySQL nao disponivel, tratando como anonimo")
		return AuthResult{
			IdUsuario:   0,
			IsValid:     true, // Se nao tem MySQL, permite mas como anonimo
			IsAssinante: false,
			TeamId:      0,
		}
	}

	var idUsuario, teamId int
	err := db.QueryRow(
		"SELECT id, current_team_id FROM users WHERE token_access = ?",
		token,
	).Scan(&idUsuario, &teamId)

	if err != nil {
		if err == sql.ErrNoRows {
			// Token invalido - nao existe
			log.Printf("Auth: Token nao encontrado")
			return AuthResult{
				IdUsuario:   0,
				IsValid:     false,
				IsAssinante: false,
				TeamId:      0,
			}
		}
		// Erro de conexao - trata como anonimo para nao bloquear
		log.Printf("Auth: Erro ao consultar token: %v", err)
		return AuthResult{
			IdUsuario:   0,
			IsValid:     true,
			IsAssinante: false,
			TeamId:      0,
		}
	}

	// Token valido - verifica se é assinante (team_id 1-4)
	isAssinante := teamId >= 1 && teamId <= 4

	return AuthResult{
		IdUsuario:   idUsuario,
		IsValid:     true,
		IsAssinante: isAssinante,
		TeamId:      teamId,
	}
}

// IsAssinanteTeamId verifica se um team_id corresponde a assinante
// Admin Root (1), Admin (2), VIP (3), Assinante (4) = assinante
// Free (5), outros = nao assinante
func IsAssinanteTeamId(teamId int) bool {
	return teamId >= 1 && teamId <= 4
}
