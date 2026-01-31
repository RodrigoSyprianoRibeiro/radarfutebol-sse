package services

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// TokenCache armazena tokens validados em memoria para evitar queries repetidas
type TokenCache struct {
	mu      sync.RWMutex
	tokens  map[string]*CachedToken
	maxAge  time.Duration
}

// CachedToken representa um token em cache
type CachedToken struct {
	UserID    int
	Valid     bool
	CachedAt  time.Time
}

var tokenCache = &TokenCache{
	tokens: make(map[string]*CachedToken),
	maxAge: 5 * time.Minute, // Cache por 5 minutos
}

// ValidarToken verifica se o token e valido para o usuario
// Retorna (valido, erro)
// NOTA: Esta funcao so deve ser chamada quando token NAO e vazio
// O handler deve tratar token vazio antes de chamar esta funcao
func ValidarToken(userID int, token string) (bool, error) {
	// Verifica cache primeiro
	cacheKey := fmt.Sprintf("%d:%s", userID, token)
	if cached := tokenCache.get(cacheKey); cached != nil {
		return cached.Valid, nil
	}

	// Busca no banco
	valid, err := validarTokenNoBanco(userID, token)
	if err != nil {
		log.Printf("Erro ao validar token (user=%d): %v", userID, err)
		// Retorna erro - o handler decide se trata como anonimo ou nao
		return false, err
	}

	// Cacheia resultado
	tokenCache.set(cacheKey, &CachedToken{
		UserID:   userID,
		Valid:    valid,
		CachedAt: time.Now(),
	})

	if !valid {
		log.Printf("Token invalido para usuario %d", userID)
	}

	return valid, nil
}

// validarTokenNoBanco verifica o token diretamente no MySQL
func validarTokenNoBanco(userID int, token string) (bool, error) {
	if db == nil {
		// Se MySQL nao esta disponivel, permite acesso
		return true, nil
	}

	var count int
	query := `SELECT COUNT(*) FROM users WHERE id = ? AND token_access = ?`

	err := db.QueryRow(query, userID, token).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("erro ao validar token: %w", err)
	}

	return count > 0, nil
}

// get busca um token no cache
func (tc *TokenCache) get(key string) *CachedToken {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	cached, exists := tc.tokens[key]
	if !exists {
		return nil
	}

	// Verifica se expirou
	if time.Since(cached.CachedAt) > tc.maxAge {
		return nil
	}

	return cached
}

// set adiciona um token ao cache
func (tc *TokenCache) set(key string, token *CachedToken) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.tokens[key] = token

	// Limpa cache antigo periodicamente (a cada 1000 entradas)
	if len(tc.tokens) > 1000 {
		tc.cleanup()
	}
}

// cleanup remove tokens expirados do cache
func (tc *TokenCache) cleanup() {
	now := time.Now()
	for key, cached := range tc.tokens {
		if now.Sub(cached.CachedAt) > tc.maxAge {
			delete(tc.tokens, key)
		}
	}
}

// InvalidarTokenUsuario remove o cache de um usuario especifico
func InvalidarTokenUsuario(userID int) {
	tokenCache.mu.Lock()
	defer tokenCache.mu.Unlock()

	prefix := fmt.Sprintf("%d:", userID)
	for key := range tokenCache.tokens {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(tokenCache.tokens, key)
		}
	}
}
