package services

import (
	"testing"

	"radarfutebol-sse/internal/models"
)

// Helper para criar evento de teste
func criarEvento(id int, timeCasa, timeFora string, status string) *models.Evento {
	golCasa := 0
	golFora := 0
	return &models.Evento{
		IdEvento:          id,
		TimeCasa:          timeCasa,
		TimeFora:          timeFora,
		Status:            status,
		IdCampeonatoUnico: "br-serie-a",
		NomeCampeonato:    "Brasileirao Serie A",
		Prioridade:        1,
		Inicio:            "2026-01-31 15:00",
		GolTimeCasaFt:     &golCasa,
		GolTimeForaFt:     &golFora,
		OddTimeCasa:       "1.50",
		OddTimeFora:       "3.00",
		Oraculo:           1,
		IdBetfair:         "123456",
	}
}

// Helper para criar evento com gols
func criarEventoComGols(id int, golCasa, golFora int) *models.Evento {
	evento := criarEvento(id, "Time Casa", "Time Fora", "inprogress")
	evento.GolTimeCasaFt = &golCasa
	evento.GolTimeForaFt = &golFora
	return evento
}

// =============================================================================
// TESTES DE DUPLICACAO - Verifica se eventos nao sao duplicados
// =============================================================================

func TestFiltrarEventosPainel_SemDuplicatas(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "inprogress"),
		criarEvento(3, "Santos", "Gremio", "inprogress"),
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, err := FiltrarEventosPainel(eventos, filtro, nil)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verifica duplicatas
	idsVistos := make(map[int]bool)
	for _, evento := range resultado.Eventos {
		if idsVistos[evento.IdEvento] {
			t.Errorf("Evento duplicado encontrado: ID=%d", evento.IdEvento)
		}
		idsVistos[evento.IdEvento] = true
	}

	if len(resultado.Eventos) != 3 {
		t.Errorf("Esperado 3 eventos, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltrarEventosHome_SemDuplicatas(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "inprogress"),
		criarEvento(3, "Santos", "Gremio", "inprogress"),
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, err := FiltrarEventosHome(eventos, filtro, nil)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verifica duplicatas nos campeonatos
	idsVistos := make(map[int]bool)
	for _, camp := range resultado.Campeonatos {
		for _, evento := range camp.Eventos {
			if idsVistos[evento.IdEvento] {
				t.Errorf("Evento duplicado encontrado: ID=%d", evento.IdEvento)
			}
			idsVistos[evento.IdEvento] = true
		}
	}
}

func TestFiltrarEventosPainel_EventosDuplicadosNaEntrada(t *testing.T) {
	// Simula entrada com duplicatas (bug no cache)
	evento := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	eventos := []*models.Evento{evento, evento, evento}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, err := FiltrarEventosPainel(eventos, filtro, nil)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Mesmo com duplicatas na entrada, cada ponteiro e copiado
	// Mas os IDs serao iguais - isso indica problema na origem
	if len(resultado.Eventos) != 3 {
		t.Logf("AVISO: Entrada com duplicatas resulta em %d eventos", len(resultado.Eventos))
	}
}

// =============================================================================
// TESTES DE ISOLAMENTO - Verifica que modificacoes nao afetam cache
// =============================================================================

func TestFiltrarEventosPainel_IsolamentoFavoritos(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
	}

	// Usuario 1 com favorito
	prefs1 := &PreferenciasUsuario{
		JogosFavoritos: map[string]bool{"1": true},
	}

	// Usuario 2 sem favorito
	prefs2 := &PreferenciasUsuario{
		JogosFavoritos: map[string]bool{},
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	// Filtra para usuario 1
	resultado1, _ := FiltrarEventosPainel(eventos, filtro, prefs1)

	// Filtra para usuario 2
	resultado2, _ := FiltrarEventosPainel(eventos, filtro, prefs2)

	// Verifica que favorito de usuario 1 nao vazou para usuario 2
	if resultado1.Eventos[0].Favorito.Bool() != true {
		t.Error("Usuario 1 deveria ter evento como favorito")
	}

	if resultado2.Eventos[0].Favorito.Bool() != false {
		t.Error("Usuario 2 NAO deveria ter evento como favorito (vazamento de estado)")
	}

	// Verifica que o evento original nao foi modificado
	if eventos[0].Favorito.Bool() != false {
		t.Error("Evento original foi modificado (cache corrompido)")
	}
}

func TestFiltrarEventosPainel_IsolamentoCampeonatoFavorito(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
	}

	// Usuario 1 com campeonato favorito
	prefs1 := &PreferenciasUsuario{
		CampeonatosFavoritos: map[string]bool{"br-serie-a": true},
	}

	// Usuario 2 sem campeonato favorito
	prefs2 := &PreferenciasUsuario{
		CampeonatosFavoritos: map[string]bool{},
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado1, _ := FiltrarEventosPainel(eventos, filtro, prefs1)
	resultado2, _ := FiltrarEventosPainel(eventos, filtro, prefs2)

	if resultado1.Eventos[0].CampeonatoFavorito.Bool() != true {
		t.Error("Usuario 1 deveria ter campeonato como favorito")
	}

	if resultado2.Eventos[0].CampeonatoFavorito.Bool() != false {
		t.Error("Usuario 2 NAO deveria ter campeonato como favorito (vazamento)")
	}
}

func TestFiltrarEventosPainel_MultiUsuariosConcorrentes(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "inprogress"),
	}

	// Simula 10 usuarios com preferencias diferentes
	for i := 0; i < 10; i++ {
		prefs := &PreferenciasUsuario{
			JogosFavoritos: map[string]bool{},
		}
		if i%2 == 0 {
			prefs.JogosFavoritos["1"] = true
		}

		filtro := &models.Filtro{
			IdUsuario:         i,
			CountJogosMostrar: 100,
		}

		resultado, _ := FiltrarEventosPainel(eventos, filtro, prefs)

		// Verifica resultado correto para cada usuario
		for _, evento := range resultado.Eventos {
			if evento.IdEvento == 1 {
				esperado := i%2 == 0
				if evento.Favorito.Bool() != esperado {
					t.Errorf("Usuario %d: favorito incorreto. Esperado=%v, Recebeu=%v",
						i, esperado, evento.Favorito.Bool())
				}
			}
		}
	}

	// Verifica que eventos originais nao foram modificados
	for _, evento := range eventos {
		if evento.Favorito.Bool() {
			t.Error("Evento original foi modificado durante processamento multi-usuario")
		}
	}
}

// =============================================================================
// TESTES DE FILTROS INDIVIDUAIS
// =============================================================================

func TestFiltro_MostrarApenasJogosLive(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "scheduled"),
		criarEvento(3, "Santos", "Gremio", "finished"),
	}

	filtro := &models.Filtro{
		MostrarApenasJogosLive: true,
		CountJogosMostrar:      100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento live, recebeu %d", len(resultado.Eventos))
	}

	if resultado.Eventos[0].IdEvento != 1 {
		t.Error("Evento retornado deveria ser o de status inprogress")
	}
}

func TestFiltro_MostrarApenasJogosFavoritos(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "inprogress"),
		criarEvento(3, "Santos", "Gremio", "inprogress"),
	}

	prefs := &PreferenciasUsuario{
		JogosFavoritos: map[string]bool{"2": true},
	}

	filtro := &models.Filtro{
		MostrarApenasJogosFavoritos: true,
		CountJogosMostrar:           100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, prefs)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento favorito, recebeu %d", len(resultado.Eventos))
	}

	if resultado.Eventos[0].IdEvento != 2 {
		t.Error("Evento retornado deveria ser o favorito (ID=2)")
	}
}

func TestFiltro_MostrarApenasJogosFavoritos_CampeonatoFavorito(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
	}

	prefs := &PreferenciasUsuario{
		JogosFavoritos:       map[string]bool{},
		CampeonatosFavoritos: map[string]bool{"br-serie-a": true},
	}

	filtro := &models.Filtro{
		MostrarApenasJogosFavoritos: true,
		CountJogosMostrar:           100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, prefs)

	if len(resultado.Eventos) != 1 {
		t.Error("Evento de campeonato favorito deveria aparecer no filtro de favoritos")
	}
}

func TestFiltro_MostrarApenasJogosOraculo(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.Oraculo = 1
	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.Oraculo = 0

	eventos := []*models.Evento{evento1, evento2}

	filtro := &models.Filtro{
		MostrarApenasJogosOraculo: true,
		CountJogosMostrar:         100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento com oraculo, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltro_MostrarApenasJogosBetfair(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.IdBetfair = "123"
	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.IdBetfair = ""

	eventos := []*models.Evento{evento1, evento2}

	filtro := &models.Filtro{
		MostrarApenasJogosBetfair: true,
		CountJogosMostrar:         100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento betfair, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltro_CasaVencendo(t *testing.T) {
	eventos := []*models.Evento{
		criarEventoComGols(1, 2, 1), // Casa vencendo
		criarEventoComGols(2, 1, 2), // Fora vencendo
		criarEventoComGols(3, 1, 1), // Empate
	}

	filtro := &models.Filtro{
		CasaVencendo:      true,
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento com casa vencendo, recebeu %d", len(resultado.Eventos))
	}

	if resultado.Eventos[0].IdEvento != 1 {
		t.Error("Evento retornado deveria ser ID=1 (casa vencendo)")
	}
}

func TestFiltro_VisitanteVencendo(t *testing.T) {
	eventos := []*models.Evento{
		criarEventoComGols(1, 2, 1), // Casa vencendo
		criarEventoComGols(2, 1, 2), // Fora vencendo
		criarEventoComGols(3, 1, 1), // Empate
	}

	filtro := &models.Filtro{
		VisitanteVencendo: true,
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento com visitante vencendo, recebeu %d", len(resultado.Eventos))
	}

	if resultado.Eventos[0].IdEvento != 2 {
		t.Error("Evento retornado deveria ser ID=2 (visitante vencendo)")
	}
}

func TestFiltro_Empatado(t *testing.T) {
	eventos := []*models.Evento{
		criarEventoComGols(1, 2, 1), // Casa vencendo
		criarEventoComGols(2, 1, 2), // Fora vencendo
		criarEventoComGols(3, 1, 1), // Empate
		criarEventoComGols(4, 0, 0), // Empate 0x0
	}

	filtro := &models.Filtro{
		Empatado:          true,
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 2 {
		t.Errorf("Esperado 2 eventos empatados, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltro_FavoritoVencendo(t *testing.T) {
	// Evento 1: Casa favorita (odd menor) e vencendo
	evento1 := criarEventoComGols(1, 2, 1)
	evento1.OddTimeCasa = "1.50" // Favorito
	evento1.OddTimeFora = "3.00"

	// Evento 2: Fora favorito e vencendo
	evento2 := criarEventoComGols(2, 0, 1)
	evento2.OddTimeCasa = "3.00"
	evento2.OddTimeFora = "1.50" // Favorito

	// Evento 3: Casa favorita mas perdendo
	evento3 := criarEventoComGols(3, 0, 1)
	evento3.OddTimeCasa = "1.50" // Favorito
	evento3.OddTimeFora = "3.00"

	eventos := []*models.Evento{evento1, evento2, evento3}

	filtro := &models.Filtro{
		FavoritoVencendo:  true,
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 2 {
		t.Errorf("Esperado 2 eventos com favorito vencendo, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltro_FavoritoPerdendo(t *testing.T) {
	// Evento 1: Casa favorita mas perdendo
	evento1 := criarEventoComGols(1, 0, 1)
	evento1.OddTimeCasa = "1.50" // Favorito
	evento1.OddTimeFora = "3.00"

	// Evento 2: Casa favorita e vencendo
	evento2 := criarEventoComGols(2, 2, 0)
	evento2.OddTimeCasa = "1.50" // Favorito
	evento2.OddTimeFora = "3.00"

	eventos := []*models.Evento{evento1, evento2}

	filtro := &models.Filtro{
		FavoritoPerdendo:  true,
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento com favorito perdendo, recebeu %d", len(resultado.Eventos))
	}

	if resultado.Eventos[0].IdEvento != 1 {
		t.Error("Evento retornado deveria ser ID=1 (favorito perdendo)")
	}
}

func TestFiltro_CampoBusca(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "inprogress"),
		criarEvento(3, "Santos", "Gremio", "inprogress"),
	}

	testCases := []struct {
		busca    string
		esperado int
	}{
		{"flamengo", 1},
		{"Palmeiras", 1},
		{"CORINTHIANS", 1},
		{"santos", 1},
		{"xyz", 0},
		{"fla", 1}, // Substring
	}

	for _, tc := range testCases {
		filtro := &models.Filtro{
			CampoBusca:        tc.busca,
			CountJogosMostrar: 100,
		}

		resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

		if len(resultado.Eventos) != tc.esperado {
			t.Errorf("Busca '%s': esperado %d eventos, recebeu %d",
				tc.busca, tc.esperado, len(resultado.Eventos))
		}
	}
}

func TestFiltro_MostrarApenasJogosOver(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.OverEvento = 1
	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.OverEvento = 0

	eventos := []*models.Evento{evento1, evento2}

	filtro := &models.Filtro{
		MostrarApenasJogosOver: true,
		CountJogosMostrar:      100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento over, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltro_MostrarApenasJogosLayCs(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.LayCsEvento = 1
	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.LayCsEvento = 0

	eventos := []*models.Evento{evento1, evento2}

	filtro := &models.Filtro{
		MostrarApenasJogosLayCs: true,
		CountJogosMostrar:       100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento lay cs, recebeu %d", len(resultado.Eventos))
	}
}

func TestFiltro_Alertas(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.AlertaMomentoGolAtivo = models.FlexBool(true)

	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.AlertaPressaoIndividualAtivo = models.FlexBool(true)

	evento3 := criarEvento(3, "Santos", "Gremio", "inprogress")
	// Sem alertas

	eventos := []*models.Evento{evento1, evento2, evento3}

	filtro := &models.Filtro{
		FiltroAlertas:     true,
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 2 {
		t.Errorf("Esperado 2 eventos com alertas, recebeu %d", len(resultado.Eventos))
	}
}

// =============================================================================
// TESTES DE COMBINACAO DE FILTROS
// =============================================================================

func TestFiltro_Combinacao_LiveEFavorito(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "scheduled"),
		criarEvento(3, "Santos", "Gremio", "inprogress"),
	}

	prefs := &PreferenciasUsuario{
		JogosFavoritos: map[string]bool{"1": true, "2": true},
	}

	filtro := &models.Filtro{
		MostrarApenasJogosLive:      true,
		MostrarApenasJogosFavoritos: true,
		CountJogosMostrar:           100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, prefs)

	// Deve retornar apenas ID=1 (live E favorito)
	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento (live E favorito), recebeu %d", len(resultado.Eventos))
	}

	if len(resultado.Eventos) > 0 && resultado.Eventos[0].IdEvento != 1 {
		t.Error("Evento deveria ser ID=1")
	}
}

func TestFiltro_Combinacao_LiveOraculoBetfair(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.Oraculo = 1
	evento1.IdBetfair = "123"

	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.Oraculo = 1
	evento2.IdBetfair = "" // Sem betfair

	evento3 := criarEvento(3, "Santos", "Gremio", "scheduled")
	evento3.Oraculo = 1
	evento3.IdBetfair = "456"

	eventos := []*models.Evento{evento1, evento2, evento3}

	filtro := &models.Filtro{
		MostrarApenasJogosLive:    true,
		MostrarApenasJogosOraculo: true,
		MostrarApenasJogosBetfair: true,
		CountJogosMostrar:         100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento (live E oraculo E betfair), recebeu %d", len(resultado.Eventos))
	}
}

// =============================================================================
// TESTES DE CONTADORES
// =============================================================================

func TestContadores_JogosLive(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
		criarEvento(2, "Corinthians", "Sao Paulo", "inprogress"),
		criarEvento(3, "Santos", "Gremio", "scheduled"),
		criarEvento(4, "Atletico", "Cruzeiro", "finished"),
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if resultado.Counts.Live != 2 {
		t.Errorf("Esperado 2 jogos live, recebeu %d", resultado.Counts.Live)
	}

	if resultado.Counts.Total != 4 {
		t.Errorf("Esperado 4 jogos total, recebeu %d", resultado.Counts.Total)
	}
}

// =============================================================================
// TESTES DE LIMITE
// =============================================================================

func TestLimite_CountJogosMostrar(t *testing.T) {
	eventos := make([]*models.Evento, 100)
	for i := 0; i < 100; i++ {
		eventos[i] = criarEvento(i+1, "Time A", "Time B", "inprogress")
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 25,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, nil)

	if len(resultado.Eventos) != 25 {
		t.Errorf("Esperado limite de 25 eventos, recebeu %d", len(resultado.Eventos))
	}
}

// =============================================================================
// TESTES DE AGRUPAMENTO (Home)
// =============================================================================

func TestAgruparPorCampeonato_SemDuplicatas(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.IdCampeonatoUnico = "br-serie-a"

	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.IdCampeonatoUnico = "br-serie-a"

	evento3 := criarEvento(3, "Real Madrid", "Barcelona", "inprogress")
	evento3.IdCampeonatoUnico = "es-la-liga"

	eventos := []*models.Evento{evento1, evento2, evento3}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosHome(eventos, filtro, nil)

	if len(resultado.Campeonatos) != 2 {
		t.Errorf("Esperado 2 campeonatos, recebeu %d", len(resultado.Campeonatos))
	}

	// Verifica que nao ha duplicatas em nenhum campeonato
	for _, camp := range resultado.Campeonatos {
		idsVistos := make(map[int]bool)
		for _, evento := range camp.Eventos {
			if idsVistos[evento.IdEvento] {
				t.Errorf("Evento duplicado no campeonato %s: ID=%d", camp.IdUnico, evento.IdEvento)
			}
			idsVistos[evento.IdEvento] = true
		}
	}
}

func TestAgruparPorCampeonato_LimiteRespeitado(t *testing.T) {
	eventos := make([]*models.Evento, 50)
	for i := 0; i < 50; i++ {
		eventos[i] = criarEvento(i+1, "Time A", "Time B", "inprogress")
		eventos[i].IdCampeonatoUnico = "br-serie-a"
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 10,
	}

	resultado, _ := FiltrarEventosHome(eventos, filtro, nil)

	totalEventos := 0
	for _, camp := range resultado.Campeonatos {
		totalEventos += len(camp.Eventos)
	}

	if totalEventos != 10 {
		t.Errorf("Esperado limite de 10 eventos, recebeu %d", totalEventos)
	}
}

// =============================================================================
// TESTES DE ORDENACAO
// =============================================================================

func TestOrdenacao_FavoritosPrimeiro(t *testing.T) {
	evento1 := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento1.Prioridade = 1

	evento2 := criarEvento(2, "Corinthians", "Sao Paulo", "inprogress")
	evento2.Prioridade = 2

	eventos := []*models.Evento{evento1, evento2}

	prefs := &PreferenciasUsuario{
		JogosFavoritos: map[string]bool{"2": true}, // ID 2 e favorito
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, _ := FiltrarEventosPainel(eventos, filtro, prefs)

	// Favorito deve vir primeiro, mesmo com prioridade maior
	if resultado.Eventos[0].IdEvento != 2 {
		t.Error("Evento favorito deveria estar primeiro na lista")
	}
}

// =============================================================================
// TESTES DE EDGE CASES
// =============================================================================

func TestEdgeCase_EventosVazios(t *testing.T) {
	eventos := []*models.Evento{}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, err := FiltrarEventosPainel(eventos, filtro, nil)

	if err != nil {
		t.Fatalf("Erro inesperado com lista vazia: %v", err)
	}

	if resultado.Eventos == nil {
		t.Error("Eventos nao deveria ser nil, deveria ser slice vazio")
	}

	if len(resultado.Eventos) != 0 {
		t.Errorf("Esperado 0 eventos, recebeu %d", len(resultado.Eventos))
	}
}

func TestEdgeCase_PreferenciasNil(t *testing.T) {
	eventos := []*models.Evento{
		criarEvento(1, "Flamengo", "Palmeiras", "inprogress"),
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 100,
	}

	resultado, err := FiltrarEventosPainel(eventos, filtro, nil)

	if err != nil {
		t.Fatalf("Erro inesperado com preferencias nil: %v", err)
	}

	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento, recebeu %d", len(resultado.Eventos))
	}

	// Favorito deve ser false quando prefs e nil
	if resultado.Eventos[0].Favorito.Bool() {
		t.Error("Favorito deveria ser false quando preferencias sao nil")
	}
}

func TestEdgeCase_GolsNil(t *testing.T) {
	evento := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento.GolTimeCasaFt = nil
	evento.GolTimeForaFt = nil

	eventos := []*models.Evento{evento}

	filtro := &models.Filtro{
		Empatado:          true, // Filtro que usa gols
		CountJogosMostrar: 100,
	}

	resultado, err := FiltrarEventosPainel(eventos, filtro, nil)

	if err != nil {
		t.Fatalf("Erro inesperado com gols nil: %v", err)
	}

	// 0 == 0, entao deve passar no filtro de empate
	if len(resultado.Eventos) != 1 {
		t.Errorf("Esperado 1 evento (nil gols = empate 0x0), recebeu %d", len(resultado.Eventos))
	}
}

func TestEdgeCase_OddsVazias(t *testing.T) {
	evento := criarEvento(1, "Flamengo", "Palmeiras", "inprogress")
	evento.OddTimeCasa = ""
	evento.OddTimeFora = ""
	gol := 1
	evento.GolTimeCasaFt = &gol
	evento.GolTimeForaFt = nil

	eventos := []*models.Evento{evento}

	filtro := &models.Filtro{
		FavoritoVencendo:  true,
		CountJogosMostrar: 100,
	}

	// Nao deve dar panic com odds vazias
	resultado, err := FiltrarEventosPainel(eventos, filtro, nil)

	if err != nil {
		t.Fatalf("Erro inesperado com odds vazias: %v", err)
	}

	// Com odds 0, nenhum time e favorito, entao nenhum passa
	if len(resultado.Eventos) != 0 {
		t.Logf("Eventos com odds vazias: %d", len(resultado.Eventos))
	}
}

// =============================================================================
// BENCHMARK - Performance
// =============================================================================

func BenchmarkFiltrarEventosPainel_100Eventos(b *testing.B) {
	eventos := make([]*models.Evento, 100)
	for i := 0; i < 100; i++ {
		eventos[i] = criarEvento(i+1, "Time A", "Time B", "inprogress")
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 50,
	}

	prefs := &PreferenciasUsuario{
		JogosFavoritos: map[string]bool{"1": true, "5": true, "10": true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FiltrarEventosPainel(eventos, filtro, prefs)
	}
}

func BenchmarkFiltrarEventosHome_100Eventos(b *testing.B) {
	eventos := make([]*models.Evento, 100)
	for i := 0; i < 100; i++ {
		eventos[i] = criarEvento(i+1, "Time A", "Time B", "inprogress")
		eventos[i].IdCampeonatoUnico = "camp-" + string(rune('a'+i%10))
	}

	filtro := &models.Filtro{
		CountJogosMostrar: 50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FiltrarEventosHome(eventos, filtro, nil)
	}
}
