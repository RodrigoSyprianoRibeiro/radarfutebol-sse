package services

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

// authCacheEntry entrada do cache de autenticacao
type authCacheEntry struct {
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
	IsValid     bool
	IsAssinante bool
	TeamId      int
}

// ValidateUserToken valida o token de um usuario
// Retorna se o token é valido e se o usuario é assinante (team_id 1-4)
func ValidateUserToken(idUsuario int, token string) AuthResult {
	// Usuario anonimo (idUsuario = 0) - nao precisa validar token
	if idUsuario == 0 {
		return AuthResult{
			IsValid:     true, // Anonimo sempre valido
			IsAssinante: false,
			TeamId:      0,
		}
	}

	// idUsuario > 0 mas sem token - trata como anonimo
	if token == "" {
		return AuthResult{
			IsValid:     true,
			IsAssinante: false,
			TeamId:      0,
		}
	}

	// Verifica cache primeiro
	cacheKey := getCacheKey(idUsuario, token)
	if cached, ok := authCache.Load(cacheKey); ok {
		entry := cached.(authCacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return AuthResult{
				IsValid:     entry.isValid,
				IsAssinante: entry.isAssinante,
				TeamId:      entry.teamId,
			}
		}
		// Cache expirado, remove
		authCache.Delete(cacheKey)
	}

	// Consulta MySQL
	result := queryUserToken(idUsuario, token)

	// Salva no cache
	authCache.Store(cacheKey, authCacheEntry{
		isValid:     result.IsValid,
		isAssinante: result.IsAssinante,
		teamId:      result.TeamId,
		expiresAt:   time.Now().Add(authCacheTTL),
	})

	return result
}

// getCacheKey gera a chave do cache para o par usuario/token
func getCacheKey(idUsuario int, token string) string {
	// Usa apenas os primeiros 16 caracteres do token para a chave
	// (suficiente para identificacao sem usar muita memoria)
	tokenPrefix := token
	if len(token) > 16 {
		tokenPrefix = token[:16]
	}
	return string(rune(idUsuario)) + ":" + tokenPrefix
}

// queryUserToken consulta o banco para validar token
func queryUserToken(idUsuario int, token string) AuthResult {
	if db == nil {
		log.Printf("Auth: MySQL nao disponivel, tratando usuario %d como anonimo", idUsuario)
		return AuthResult{
			IsValid:     true, // Se nao tem MySQL, permite mas como anonimo
			IsAssinante: false,
			TeamId:      0,
		}
	}

	var teamId int
	err := db.QueryRow(
		"SELECT current_team_id FROM users WHERE id = ? AND token_access = ?",
		idUsuario, token,
	).Scan(&teamId)

	if err != nil {
		if err == sql.ErrNoRows {
			// Token invalido - usuario existe mas token nao bate
			log.Printf("Auth: Token invalido para usuario %d", idUsuario)
			return AuthResult{
				IsValid:     false,
				IsAssinante: false,
				TeamId:      0,
			}
		}
		// Erro de conexao - trata como anonimo para nao bloquear
		log.Printf("Auth: Erro ao consultar usuario %d: %v", idUsuario, err)
		return AuthResult{
			IsValid:     true,
			IsAssinante: false,
			TeamId:      0,
		}
	}

	// Token valido - verifica se é assinante (team_id 1-4)
	isAssinante := teamId >= 1 && teamId <= 4

	return AuthResult{
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
