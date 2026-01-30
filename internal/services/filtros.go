package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"radarfutebol-sse/internal/models"
)

// Cache de alertas de gol por usuario (evita som repetido)
// Chave: alerta-gol-usuario-{userId}, Valor: map[idEvento]tempo
func getAlertasGolUsuario(userID int) map[string]string {
	if rdbPrefs == nil || userID == 0 {
		return make(map[string]string)
	}

	key := fmt.Sprintf("alerta-gol-usuario-%d", userID)
	data, err := rdbPrefs.Get(ctx, key).Result()
	if err != nil || data == "" {
		return make(map[string]string)
	}

	var alertas map[string]string
	if err := json.Unmarshal([]byte(data), &alertas); err != nil {
		return make(map[string]string)
	}
	return alertas
}

func setAlertasGolUsuario(userID int, alertas map[string]string) {
	if rdbPrefs == nil || userID == 0 {
		return
	}

	key := fmt.Sprintf("alerta-gol-usuario-%d", userID)
	data, err := json.Marshal(alertas)
	if err != nil {
		return
	}
	rdbPrefs.Set(ctx, key, string(data), 5*time.Minute)
}

// PreferenciasUsuario preferencias de campeonatos e jogos favoritos
type PreferenciasUsuario struct {
	CampeonatosFavoritos map[string]bool
	JogosFavoritos       map[string]bool
}

// FiltrarEventosPainel filtra eventos para o painel igual ao Laravel
func FiltrarEventosPainel(eventos []*models.Evento, filtro *models.Filtro, prefs *PreferenciasUsuario) (*models.PainelResponse, error) {
	// Inicializa como slice vazio (nunca nil) para JSON serializar como [] ao inves de null
	jogosFiltrados := make([]*models.Evento, 0)
	countJogosLive := 0
	countJogosTotal := 0
	countGols := 0

	// Cache de alertas de gol do usuario (evita som repetido)
	var alertasGol map[string]string
	if filtro.SomLigado && filtro.IdUsuario > 0 {
		alertasGol = getAlertasGolUsuario(filtro.IdUsuario)
	}
	alertasGolModificado := false

	for _, evento := range eventos {
		if !aplicaFiltros(evento, filtro, prefs) {
			continue
		}

		if evento.Status == "inprogress" {
			countJogosLive++
		}
		countJogosTotal++

		// Conta gols (alertarSomGol ativo e som ligado)
		// Só conta se ainda não notificou este gol (baseado no tempo do jogo)
		if filtro.SomLigado && filtro.IdUsuario > 0 && evento.AlertarSomGol.Bool() {
			idEventoStr := strconv.Itoa(evento.IdEvento)
			tempoAtual := evento.TempoAtual
			tempoAnterior, jaNotificou := alertasGol[idEventoStr]

			// Só conta como novo gol se não notificou ou se o tempo mudou
			if !jaNotificou || tempoAnterior != tempoAtual {
				countGols++
				alertasGol[idEventoStr] = tempoAtual
				alertasGolModificado = true
			}
		}

		// Marca favoritos
		if prefs != nil {
			idEventoStr := strconv.Itoa(evento.IdEvento)
			evento.Favorito = models.FlexBool(prefs.JogosFavoritos[idEventoStr])
			evento.CampeonatoFavorito = models.FlexBool(prefs.CampeonatosFavoritos[evento.IdCampeonatoUnico])
		}

		jogosFiltrados = append(jogosFiltrados, evento)
	}

	// Salva cache de alertas se foi modificado
	if alertasGolModificado {
		setAlertasGolUsuario(filtro.IdUsuario, alertasGol)
	}

	// Ordenar
	ordenarEventosPainel(jogosFiltrados, filtro.OrdemInicio)

	// Limitar quantidade
	if filtro.CountJogosMostrar > 0 && len(jogosFiltrados) > filtro.CountJogosMostrar {
		jogosFiltrados = jogosFiltrados[:filtro.CountJogosMostrar]
	}

	// Garante que o slice nunca seja nil para JSON serializar como []
	if jogosFiltrados == nil {
		jogosFiltrados = []*models.Evento{}
	}

	return &models.PainelResponse{
		Eventos: jogosFiltrados,
		Counts: models.Counts{
			Live:  countJogosLive,
			Total: countJogosTotal,
			Gols:  countGols,
		},
	}, nil
}

// FiltrarEventosHome filtra eventos para a home igual ao Laravel
func FiltrarEventosHome(eventos []*models.Evento, filtro *models.Filtro, prefs *PreferenciasUsuario) (*models.HomeResponse, error) {
	// Inicializa como slice vazio (nunca nil) para JSON serializar como [] ao inves de null
	jogosFiltrados := make([]*models.Evento, 0)
	countJogosLive := 0
	countJogosTotal := 0
	countGols := 0

	// Cache de alertas de gol do usuario (evita som repetido)
	var alertasGol map[string]string
	if filtro.SomLigado && filtro.IdUsuario > 0 {
		alertasGol = getAlertasGolUsuario(filtro.IdUsuario)
	}
	alertasGolModificado := false

	for _, evento := range eventos {
		if !aplicaFiltros(evento, filtro, prefs) {
			continue
		}

		if evento.Status == "inprogress" {
			countJogosLive++
		}
		countJogosTotal++

		// Conta gols (alertarSomGol ativo e som ligado)
		// Só conta se ainda não notificou este gol (baseado no tempo do jogo)
		if filtro.SomLigado && filtro.IdUsuario > 0 && evento.AlertarSomGol.Bool() {
			idEventoStr := strconv.Itoa(evento.IdEvento)
			tempoAtual := evento.TempoAtual
			tempoAnterior, jaNotificou := alertasGol[idEventoStr]

			// Só conta como novo gol se não notificou ou se o tempo mudou
			if !jaNotificou || tempoAnterior != tempoAtual {
				countGols++
				alertasGol[idEventoStr] = tempoAtual
				alertasGolModificado = true
			}
		}

		// Marca favoritos
		if prefs != nil {
			idEventoStr := strconv.Itoa(evento.IdEvento)
			evento.Favorito = models.FlexBool(prefs.JogosFavoritos[idEventoStr])
			evento.CampeonatoFavorito = models.FlexBool(prefs.CampeonatosFavoritos[evento.IdCampeonatoUnico])
		}

		jogosFiltrados = append(jogosFiltrados, evento)
	}

	// Salva cache de alertas se foi modificado
	if alertasGolModificado {
		setAlertasGolUsuario(filtro.IdUsuario, alertasGol)
	}

	// Ordenar
	ordenarEventosHome(jogosFiltrados, filtro.OrdemInicio)

	// Agrupar por campeonato e limitar
	campeonatos := agruparPorCampeonato(jogosFiltrados, filtro.CountJogosMostrar, prefs)

	// Garante que o slice nunca seja nil para JSON serializar como []
	if campeonatos == nil {
		campeonatos = []*models.Campeonato{}
	}

	return &models.HomeResponse{
		Campeonatos: campeonatos,
		Counts: models.Counts{
			Live:  countJogosLive,
			Total: countJogosTotal,
			Gols:  countGols,
		},
	}, nil
}

// aplicaFiltros verifica se evento passa nos filtros (igual ao PHP)
func aplicaFiltros(evento *models.Evento, filtro *models.Filtro, prefs *PreferenciasUsuario) bool {
	// Filtro acrescimo
	if filtro.MostrarFiltroAcrescimo {
		prev1, _ := strconv.Atoi(evento.PrevisaoAcrescimo1Tempo.String())
		prev2, _ := strconv.Atoi(evento.PrevisaoAcrescimo2Tempo.String())
		if prev1 < filtro.FiltroAcrescimoHt && prev2 < filtro.FiltroAcrescimoFt {
			return false
		}
	}

	// Filtro over
	if filtro.MostrarApenasJogosOver && evento.OverEvento == 0 {
		return false
	}

	// Filtro lay cs
	if filtro.MostrarApenasJogosLayCs && evento.LayCsEvento == 0 {
		return false
	}

	// Filtro jogos live
	if filtro.MostrarApenasJogosLive && evento.Status != "inprogress" {
		return false
	}

	// Filtro favoritos
	if filtro.MostrarApenasJogosFavoritos && prefs != nil {
		idEventoStr := strconv.Itoa(evento.IdEvento)
		jogoFav := prefs.JogosFavoritos[idEventoStr]
		campFav := prefs.CampeonatosFavoritos[evento.IdCampeonatoUnico]
		if !jogoFav && !campFav {
			return false
		}
	}

	// Filtro oraculo
	if filtro.MostrarApenasJogosOraculo && evento.Oraculo == 0 {
		return false
	}

	// Filtro betfair
	if filtro.MostrarApenasJogosBetfair && evento.IdBetfair == "" {
		return false
	}

	// Filtro favorito vencendo
	if filtro.FavoritoVencendo {
		oddCasa := parseOdd(evento.OddTimeCasa)
		oddFora := parseOdd(evento.OddTimeFora)
		golCasa := getGol(evento.GolTimeCasaFt)
		golFora := getGol(evento.GolTimeForaFt)

		favCasaVence := oddCasa < oddFora && golCasa > golFora
		favForaVence := oddFora < oddCasa && golFora > golCasa

		if !favCasaVence && !favForaVence {
			return false
		}
	}

	// Filtro favorito perdendo
	if filtro.FavoritoPerdendo {
		oddCasa := parseOdd(evento.OddTimeCasa)
		oddFora := parseOdd(evento.OddTimeFora)
		golCasa := getGol(evento.GolTimeCasaFt)
		golFora := getGol(evento.GolTimeForaFt)

		favCasaPerde := oddCasa < oddFora && golCasa < golFora
		favForaPerde := oddFora < oddCasa && golFora < golCasa

		if !favCasaPerde && !favForaPerde {
			return false
		}
	}

	// Filtro casa vencendo
	if filtro.CasaVencendo {
		golCasa := getGol(evento.GolTimeCasaFt)
		golFora := getGol(evento.GolTimeForaFt)
		if golCasa <= golFora {
			return false
		}
	}

	// Filtro visitante vencendo
	if filtro.VisitanteVencendo {
		golCasa := getGol(evento.GolTimeCasaFt)
		golFora := getGol(evento.GolTimeForaFt)
		if golFora <= golCasa {
			return false
		}
	}

	// Filtro empatado
	if filtro.Empatado {
		golCasa := getGol(evento.GolTimeCasaFt)
		golFora := getGol(evento.GolTimeForaFt)
		if golCasa != golFora {
			return false
		}
	}

	// Filtro momento gol
	if filtro.FiltroMomentoGol && !evento.AlertaMomentoGolAtivo.Bool() {
		return false
	}

	// Filtro pressao
	if filtro.FiltroPressao && !evento.AlertaPressaoIndividualAtivo.Bool() {
		return false
	}

	// Filtro busca
	if filtro.CampoBusca != "" {
		busca := strings.ToLower(filtro.CampoBusca)
		timeCasa := strings.ToLower(evento.TimeCasa)
		timeFora := strings.ToLower(evento.TimeFora)
		campeonato := strings.ToLower(evento.NomeCampeonato)

		// Busca por substring ou similaridade
		matchCasa := strings.Contains(timeCasa, busca) || similaridadePalavras(busca, timeCasa) >= 0.4
		matchFora := strings.Contains(timeFora, busca) || similaridadePalavras(busca, timeFora) >= 0.4
		matchCamp := strings.Contains(campeonato, busca) || similaridadePalavras(busca, campeonato) >= 0.4

		if !matchCasa && !matchFora && !matchCamp {
			return false
		}
	}

	return true
}

// ordenarEventosPainel ordena eventos do painel
func ordenarEventosPainel(eventos []*models.Evento, ordemInicio bool) {
	sort.SliceStable(eventos, func(i, j int) bool {
		// Favoritos primeiro
		if eventos[i].Favorito.Bool() != eventos[j].Favorito.Bool() {
			return eventos[i].Favorito.Bool()
		}
		if eventos[i].CampeonatoFavorito.Bool() != eventos[j].CampeonatoFavorito.Bool() {
			return eventos[i].CampeonatoFavorito.Bool()
		}

		if ordemInicio {
			// Ordena por inicio, depois prioridade
			if eventos[i].Inicio != eventos[j].Inicio {
				return eventos[i].Inicio < eventos[j].Inicio
			}
			return eventos[i].Prioridade < eventos[j].Prioridade
		}

		// Ordena por prioridade, depois inicio
		if eventos[i].Prioridade != eventos[j].Prioridade {
			return eventos[i].Prioridade < eventos[j].Prioridade
		}
		return eventos[i].Inicio < eventos[j].Inicio
	})
}

// ordenarEventosHome ordena eventos da home
func ordenarEventosHome(eventos []*models.Evento, ordemInicio bool) {
	sort.SliceStable(eventos, func(i, j int) bool {
		// Campeonato favorito primeiro
		if eventos[i].CampeonatoFavorito.Bool() != eventos[j].CampeonatoFavorito.Bool() {
			return eventos[i].CampeonatoFavorito.Bool()
		}

		if ordemInicio {
			if eventos[i].Inicio != eventos[j].Inicio {
				return eventos[i].Inicio < eventos[j].Inicio
			}
			return eventos[i].Prioridade < eventos[j].Prioridade
		}

		if eventos[i].Prioridade != eventos[j].Prioridade {
			return eventos[i].Prioridade < eventos[j].Prioridade
		}
		return eventos[i].Inicio < eventos[j].Inicio
	})
}

// agruparPorCampeonato agrupa eventos por campeonato
func agruparPorCampeonato(eventos []*models.Evento, maxJogos int, prefs *PreferenciasUsuario) []*models.Campeonato {
	campeonatosMap := make(map[string]*models.Campeonato)
	var campeonatosOrdem []string
	jogosExibidos := 0

	for _, evento := range eventos {
		if maxJogos > 0 && jogosExibidos >= maxJogos {
			break
		}

		idUnico := evento.IdCampeonatoUnico

		if _, exists := campeonatosMap[idUnico]; !exists {
			campFav := false
			if prefs != nil {
				campFav = prefs.CampeonatosFavoritos[idUnico]
			}

			campeonatosMap[idUnico] = &models.Campeonato{
				Id:                     evento.IdCampeonato,
				IdUnico:                idUnico,
				IdTemporada:            evento.IdTemporada,
				Flag:                   evento.Flag,
				NomeCategoria:          evento.NomeCategoria,
				SlugCategoria:          evento.SlugCategoria,
				NomeCampeonato:         evento.NomeCampeonato,
				NomeCampeonatoReduzido: evento.NomeCampeonatoReduzido,
				SlugCampeonato:         evento.SlugCampeonato,
				TemClassificacao:       evento.TemClassificacao,
				Prioridade:             evento.Prioridade,
				Favorito:               models.FlexBool(campFav),
				Sequencia:              len(campeonatosOrdem),
				Eventos:                make(map[string]*models.Evento),
			}
			campeonatosOrdem = append(campeonatosOrdem, idUnico)
		}

		idEventoStr := strconv.Itoa(evento.IdEvento)
		campeonatosMap[idUnico].Eventos[idEventoStr] = evento
		jogosExibidos++
	}

	// Converte map para slice na ordem correta
	result := make([]*models.Campeonato, len(campeonatosOrdem))
	for i, id := range campeonatosOrdem {
		result[i] = campeonatosMap[id]
	}

	return result
}

// parseOdd converte string de odd para float
func parseOdd(oddStr string) float64 {
	if oddStr == "" {
		return 0
	}
	odd, _ := strconv.ParseFloat(oddStr, 64)
	return odd
}

// getGol retorna valor do gol ou 0
func getGol(gol *int) int {
	if gol == nil {
		return 0
	}
	return *gol
}

// similaridadePalavras calcula similaridade entre duas strings (Levenshtein simplificado)
func similaridadePalavras(a, b string) float64 {
	a = normalizarTexto(a)
	b = normalizarTexto(b)

	if a == b {
		return 1.0
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	// Verifica se uma contem a outra
	if strings.Contains(b, a) || strings.Contains(a, b) {
		return 0.8
	}

	// Calcula distancia Levenshtein
	dist := levenshteinDistance(a, b)
	maxLen := max(len(a), len(b))
	return 1.0 - float64(dist)/float64(maxLen)
}

// normalizarTexto remove acentos e converte para minusculo
func normalizarTexto(s string) string {
	s = strings.ToLower(s)
	// Remove acentos basicos
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "õ", "o", "ô", "o", "ö", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u",
		"ç", "c", "ñ", "n",
	)
	s = replacer.Replace(s)

	// Remove caracteres nao alfanumericos
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}

// levenshteinDistance calcula distancia de Levenshtein entre duas strings
func levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Matriz de distancias
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(a)][len(b)]
}
