package models

import (
	"net/http"
	"strconv"
	"strings"
)

// Filtro parametros de filtro da requisicao SSE
type Filtro struct {
	IdUsuario                   int
	Token                       string // Token de acesso do usuario
	IsAssinante                 bool   // Se o usuario é assinante (team_id 1-4)
	SomLigado                   bool
	OrdemInicio                 bool
	CampoBusca                  string
	MostrarApenasJogosLive      bool
	MostrarApenasJogosFavoritos bool
	CountJogosMostrar           int
	MostrarFiltroAcrescimo      bool
	FiltroAcrescimoHt           int
	FiltroAcrescimoFt           int
	FiltroAcrescimoHtOperador   string
	FiltroAcrescimoFtOperador   string
	FiltroAcrescimoCondicao     string
	MostrarApenasJogosOraculo   bool
	MostrarApenasJogosBetfair   bool
	MostrarApenasJogosOver      bool
	MostrarApenasJogosLayCs     bool
	FavoritoVencendo            bool
	FavoritoPerdendo            bool
	CasaVencendo                bool
	VisitanteVencendo           bool
	Empatado                    bool
	FiltroMomentoGol            bool
	FiltroPressao               bool
	FiltroAlertas               bool
	FiltroDiferencaXg           bool
}

// ParseFiltroFromRequest extrai filtros da query string igual ao Laravel
func ParseFiltroFromRequest(r *http.Request) *Filtro {
	q := r.URL.Query()

	return &Filtro{
		IdUsuario:                   getIntParam(q.Get("idUsuario"), 0),
		Token:                       strings.TrimSpace(q.Get("token")),
		IsAssinante:                 false, // Sera definido pelo handler apos validacao
		SomLigado:                   getBoolParam(q.Get("somLigado")),
		OrdemInicio:                 getBoolParam(q.Get("ordemInicio")),
		CampoBusca:                  strings.TrimSpace(q.Get("campoBusca")),
		MostrarApenasJogosLive:      getBoolParam(q.Get("mostrarApenasJogosLive")),
		MostrarApenasJogosFavoritos: getBoolParam(q.Get("mostrarApenasJogosFavoritos")),
		CountJogosMostrar:           getIntParam(q.Get("countJogosMostrar"), 25),
		MostrarFiltroAcrescimo:      getBoolParam(q.Get("mostrarFiltroAcrescimo")),
		FiltroAcrescimoHt:           getIntParam(q.Get("filtroAcrescimoHt"), 1),
		FiltroAcrescimoFt:           getIntParam(q.Get("filtroAcrescimoFt"), 1),
		FiltroAcrescimoHtOperador:   getOperadorParam(q.Get("filtroAcrescimoHtOperador")),
		FiltroAcrescimoFtOperador:   getOperadorParam(q.Get("filtroAcrescimoFtOperador")),
		FiltroAcrescimoCondicao:     getCondicaoParam(q.Get("filtroAcrescimoCondicao")),
		MostrarApenasJogosOraculo:   getBoolParam(q.Get("mostrarApenasJogosOraculo")),
		MostrarApenasJogosBetfair:   getBoolParam(q.Get("mostrarApenasJogosBetfair")),
		MostrarApenasJogosOver:      getBoolParam(q.Get("mostrarApenasJogosOver")),
		MostrarApenasJogosLayCs:     getBoolParam(q.Get("mostrarApenasJogosLayCs")),
		FavoritoVencendo:            getBoolParam(q.Get("favoritoVencendo")),
		FavoritoPerdendo:            getBoolParam(q.Get("favoritoPerdendo")),
		CasaVencendo:                getBoolParam(q.Get("casaVencendo")),
		VisitanteVencendo:           getBoolParam(q.Get("visitanteVencendo")),
		Empatado:                    getBoolParam(q.Get("empatado")),
		FiltroMomentoGol:            getBoolParam(q.Get("filtroMomentoGol")),
		FiltroPressao:               getBoolParam(q.Get("filtroPressao")),
		FiltroAlertas:               getBoolParam(q.Get("filtroAlertas")),
		FiltroDiferencaXg:           getBoolParam(q.Get("filtroDiferencaXg")),
	}
}

// getBoolParam converte string para bool (igual filter_var do PHP)
func getBoolParam(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	return val == "true" || val == "1" || val == "yes" || val == "on"
}

// getCondicaoParam valida e retorna condicao 'e' ou 'ou'
func getCondicaoParam(val string) string {
	val = strings.ToLower(strings.TrimSpace(val))
	if val == "e" {
		return "e"
	}
	return "ou"
}

// getOperadorParam valida e retorna operador >= ou <=
func getOperadorParam(val string) string {
	val = strings.TrimSpace(val)
	if val == "<=" {
		return "<="
	}
	return ">="
}

// getIntParam converte string para int com valor default
func getIntParam(val string, defaultVal int) int {
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
