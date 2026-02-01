package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// authCacheTTL tempo de vida do cache no Redis (5 minutos)
const authCacheTTL = 5 * time.Minute

// authCacheEntry entrada do cache de autenticacao
type authCacheEntry struct {
	IdUsuario   int  `json:"id_usuario"`
	IsAssinante bool `json:"is_assinante"`
	TeamId      int  `json:"team_id"`
}

// InitAuthCache inicializa o cache de autenticacao
func InitAuthCache() {
	log.Println("Cache de autenticacao inicializado (usando Redis)")
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
// Usa Redis como cache para evitar consultas frequentes ao MySQL
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

	// Verifica cache no Redis primeiro
	cacheKey := getCacheKey(token)
	if cached := getAuthFromRedis(cacheKey); cached != nil {
		return AuthResult{
			IdUsuario:   cached.IdUsuario,
			IsValid:     true,
			IsAssinante: cached.IsAssinante,
			TeamId:      cached.TeamId,
		}
	}

	// Consulta MySQL
	result := queryToken(token)

	// Se token valido, salva no cache Redis
	if result.IsValid && result.IdUsuario > 0 {
		saveAuthToRedis(cacheKey, &authCacheEntry{
			IdUsuario:   result.IdUsuario,
			IsAssinante: result.IsAssinante,
			TeamId:      result.TeamId,
		})
	}

	return result
}

// getCacheKey gera a chave do cache para o token
func getCacheKey(token string) string {
	// Usa apenas os primeiros 20 caracteres do token para a chave
	tokenPrefix := token
	if len(token) > 20 {
		tokenPrefix = token[:20]
	}
	return fmt.Sprintf("sse-auth:%s", tokenPrefix)
}

// getAuthFromRedis busca dados de autenticacao do Redis
func getAuthFromRedis(key string) *authCacheEntry {
	if rdb == nil {
		return nil
	}

	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // Chave nao existe
	}
	if err != nil {
		log.Printf("Auth: Erro ao buscar cache Redis: %v", err)
		return nil
	}

	var entry authCacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		log.Printf("Auth: Erro ao decodificar cache: %v", err)
		return nil
	}

	return &entry
}

// saveAuthToRedis salva dados de autenticacao no Redis
func saveAuthToRedis(key string, entry *authCacheEntry) {
	if rdb == nil {
		return
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Auth: Erro ao serializar cache: %v", err)
		return
	}

	if err := rdb.Set(ctx, key, data, authCacheTTL).Err(); err != nil {
		log.Printf("Auth: Erro ao salvar cache Redis: %v", err)
	}
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
