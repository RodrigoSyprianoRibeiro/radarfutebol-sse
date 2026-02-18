package services

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"radarfutebol-sse/internal/config"
)

var db *sql.DB

// InitMySQL inicializa a conexao com o MySQL
func InitMySQL(cfg config.MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao conectar MySQL: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return fmt.Errorf("erro ao ping MySQL: %w", err)
	}

	log.Println("MySQL conectado com sucesso")
	return nil
}

// GetDB retorna a conexao do banco
func GetDB() *sql.DB {
	return db
}

// GetEventosBase busca eventos base do banco de dados
func GetEventosBase(filtro map[string]interface{}) ([]map[string]interface{}, error) {
	query := `
		SELECT
			e.id,
			e.id_williamhill,
			e.id_betfair,
			e.status,
			e.escalacao as tem_escalacao,
			e.desconto_ht,
			e.desconto_ft,
			e.ativo,
			e.oraculo,
			e.oraculo_free,
			e.over as over_evento,
			e.laycs as lay_cs_evento,
			e.problema_radar,
			e.williamhill_ivertido,
			e.inicio,
			e.id_campeonato,
			e.id_time_casa,
			e.id_time_fora,
			tc.name as time_casa,
			tc.slug as slug_time_casa,
			tc.tem_imagem as tem_imagem_time_casa,
			tf.name as time_fora,
			tf.slug as slug_time_fora,
			tf.tem_imagem as tem_imagem_time_fora,
			c.name as nome_campeonato,
			c.name as nome_campeonato_reduzido,
			c.slug as slug_campeonato,
			c.prioridade,
			c.classificacao as tem_classificacao,
			cat.name as nome_categoria,
			cat.slug as slug_categoria,
			cat.flag
		FROM eventos e
		LEFT JOIN times tc ON e.id_time_casa = tc.id
		LEFT JOIN times tf ON e.id_time_fora = tf.id
		LEFT JOIN campeonatos c ON e.id_campeonato = c.id
		LEFT JOIN categorias cat ON c.id_categoria = cat.id
		WHERE 1=1
		ORDER BY
			CASE e.status
				WHEN 'first_half' THEN 1
				WHEN 'second_half' THEN 2
				WHEN 'half_time' THEN 3
				WHEN 'extra_time' THEN 4
				WHEN 'penalties' THEN 5
				ELSE 6
			END,
			c.prioridade ASC,
			e.inicio ASC
		LIMIT 50
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos: %w", err)
	}
	defer rows.Close()

	var eventos []map[string]interface{}

	for rows.Next() {
		var (
			id, idTimeCasa, idTimeFora, idCampeonato int
			idWilliamhill, idBetfair sql.NullString
			status string
			temEscalacao, ativo, oraculo, oraculoFree int
			overEvento, layCsEvento, problemaRadar, williamhillIvertido int
			descontoHt, descontoFt sql.NullInt64
			inicio string
			timeCasa, slugTimeCasa string
			temImagemTimeCasa int
			timeFora, slugTimeFora string
			temImagemTimeFora int
			nomeCampeonato, nomeCampeonatoReduzido, slugCampeonato string
			prioridade, temClassificacao int
			nomeCategoria, slugCategoria, flag sql.NullString
		)

		err := rows.Scan(
			&id, &idWilliamhill, &idBetfair, &status, &temEscalacao,
			&descontoHt, &descontoFt, &ativo, &oraculo, &oraculoFree,
			&overEvento, &layCsEvento, &problemaRadar, &williamhillIvertido,
			&inicio, &idCampeonato, &idTimeCasa, &idTimeFora,
			&timeCasa, &slugTimeCasa, &temImagemTimeCasa,
			&timeFora, &slugTimeFora, &temImagemTimeFora,
			&nomeCampeonato, &nomeCampeonatoReduzido, &slugCampeonato,
			&prioridade, &temClassificacao,
			&nomeCategoria, &slugCategoria, &flag,
		)
		if err != nil {
			log.Printf("Erro ao scan evento: %v", err)
			continue
		}

		evento := map[string]interface{}{
			"id":               id,
			"idWilliamhill":    nullStringToString(idWilliamhill),
			"idBetfair":        nullStringToString(idBetfair),
			"status":           status,
			"temEscalacao":     temEscalacao,
			"descontoHt":       nullIntToPtr(descontoHt),
			"descontoFt":       nullIntToPtr(descontoFt),
			"ativo":            ativo,
			"oraculo":          oraculo,
			"oraculoFree":      oraculoFree,
			"overEvento":       overEvento,
			"layCsEvento":      layCsEvento,
			"problemaRadar":    problemaRadar,
			"williamhillIvertido": williamhillIvertido,
			"inicio":           inicio,
			"idCampeonato":     idCampeonato,
			"idTimeCasa":       idTimeCasa,
			"idTimeFora":       idTimeFora,
			"timeCasa":         timeCasa,
			"slugTimeCasa":     slugTimeCasa,
			"timeFora":         timeFora,
			"slugTimeFora":     slugTimeFora,
			"nomeCampeonato":   nomeCampeonato,
			"nomeCampeonatoReduzido": nomeCampeonatoReduzido,
			"slugCampeonato":   slugCampeonato,
			"prioridade":       prioridade,
			"temClassificacao": temClassificacao,
			"nomeCategoria":    nullStringToString(nomeCategoria),
			"slugCategoria":    nullStringToString(slugCategoria),
			"flag":             nullStringToString(flag),
		}

		eventos = append(eventos, evento)
	}

	return eventos, nil
}

// GetCountsByStatus retorna contadores de eventos
func GetCountsByStatus() (live, total, gols int, err error) {
	query := `
		SELECT
			SUM(CASE WHEN status IN ('first_half', 'second_half', 'half_time', 'extra_time', 'penalties') THEN 1 ELSE 0 END) as ao_vivo,
			COUNT(*) as total
		FROM eventos
	`

	err = db.QueryRow(query).Scan(&live, &total)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("erro ao buscar counts: %w", err)
	}

	return live, total, 0, nil
}

// getEventoInfoFromDB busca status, escalacao, problemaRadar e acrescimos do evento pelo idWilliamhill
func getEventoInfoFromDB(idWilliamhill string) (*EventoInfo, error) {
	query := `SELECT status, escalacao, problema_radar, desconto_ht, desconto_ft FROM eventos WHERE id_williamhill = ? LIMIT 1`

	var (
		status        string
		temEscalacao  int
		problemaRadar int
		descontoHt    sql.NullInt64
		descontoFt    sql.NullInt64
	)

	err := db.QueryRow(query, idWilliamhill).Scan(&status, &temEscalacao, &problemaRadar, &descontoHt, &descontoFt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar evento: %w", err)
	}

	return &EventoInfo{
		Status:        status,
		TemEscalacao:  temEscalacao,
		ProblemaRadar: problemaRadar,
		DescontoHt:    nullIntToPtr(descontoHt),
		DescontoFt:    nullIntToPtr(descontoFt),
	}, nil
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullIntToPtr(ni sql.NullInt64) *int {
	if ni.Valid {
		v := int(ni.Int64)
		return &v
	}
	return nil
}
