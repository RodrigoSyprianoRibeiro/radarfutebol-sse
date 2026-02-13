package models

// Evento representa um jogo/evento de futebol
// Nota: Campos de estatísticas usam FlexValue pois PHP envia como string ou number
type Evento struct {
	// Identificacao
	IdEvento      int    `json:"idEvento"`
	IdWilliamhill string `json:"idWilliamhill"`
	IdBetfair     string `json:"idBetfair"`
	SlugEvento    string `json:"slugEvento"`

	// Time Casa
	IdTimeCasa             int    `json:"idTimeCasa"`
	TimeCasa               string `json:"timeCasa"`
	SlugTimeCasa           string `json:"slugTimeCasa"`
	GolTimeCasaFt          *int   `json:"golTimeCasaFt"`
	GolTimeCasaHt          *int   `json:"golTimeCasaHt"`
	CartaoVermelhoTimeCasa *int   `json:"cartaoVermelhoTimeCasa"`
	OddTimeCasa            string `json:"oddTimeCasa"`
	ClassOddTimeCasa       string `json:"classOddTimeCasa"`

	// Time Fora
	IdTimeFora             int    `json:"idTimeFora"`
	TimeFora               string `json:"timeFora"`
	SlugTimeFora           string `json:"slugTimeFora"`
	GolTimeForaFt          *int   `json:"golTimeForaFt"`
	GolTimeForaHt          *int   `json:"golTimeForaHt"`
	CartaoVermelhoTimeFora *int   `json:"cartaoVermelhoTimeFora"`
	OddTimeFora            string `json:"oddTimeFora"`
	ClassOddTimeFora       string `json:"classOddTimeFora"`

	// Status do jogo
	Status             string `json:"status"`
	TempoAtual         string `json:"tempoAtual"`
	Inicio             string `json:"inicio"`
	Oraculo            int    `json:"oraculo"`
	OraculoFree        int    `json:"oraculoFree"`
	OverEvento         int    `json:"overEvento"`
	LayCsEvento        int    `json:"layCsEvento"`
	ProblemaRadar      int    `json:"problemaRadar"`
	TemEscalacao       int    `json:"temEscalacao"`
	DescontoHt         *int   `json:"descontoHt"`
	DescontoFt         *int   `json:"descontoFt"`
	WilliamhillIvertido int   `json:"williamhillIvertido"`

	// Campeonato
	IdCampeonato           int    `json:"idCampeonato"`
	IdCampeonatoUnico      string `json:"idCampeonatoUnico"`
	NomeCampeonato         string `json:"nomeCampeonato"`
	NomeCampeonatoReduzido string `json:"nomeCampeonatoReduzido"`
	SlugCampeonato         string `json:"slugCampeonato"`
	NomeCategoria          string `json:"nomeCategoria"`
	SlugCategoria          string `json:"slugCategoria"`
	Flag                   string `json:"flag"`
	Prioridade             int    `json:"prioridade"`
	TemClassificacao       int    `json:"temClassificacao"`
	IdTemporada            string `json:"idTemporada"`
	AnoTemporada           string `json:"anoTemporada"`

	// Odds
	OddEmpate    string `json:"oddEmpate"`
	OddUnder15FT string `json:"oddUnder15FT"`
	OddOver15FT  string `json:"oddOver15FT"`
	OddUnder25FT string `json:"oddUnder25FT"`
	OddOver25FT  string `json:"oddOver25FT"`
	OddBttsSim   string `json:"oddBttsSim"`
	OddBttsNao   string `json:"oddBttsNao"`

	// Classes CSS Odds
	ClassOddEmpate    string `json:"classOddEmpate"`
	ClassOddUnder15FT string `json:"classOddUnder15FT"`
	ClassOddOver15FT  string `json:"classOddOver15FT"`
	ClassOddUnder25FT string `json:"classOddUnder25FT"`
	ClassOddOver25FT  string `json:"classOddOver25FT"`
	ClassOddBttsSim   string `json:"classOddBttsSim"`
	ClassOddBttsNao   string `json:"classOddBttsNao"`

	// Links
	LinkWilliamhill    string `json:"linkWilliamhill"`
	LinkBetfair        string `json:"linkBetfair"`
	LinkOddjusta       string `json:"linkOddjusta"`
	LinkBolsadeaposta  string `json:"linkBolsadeaposta"`
	LinkFulltbet       string `json:"linkFulltbet"`
	LinkOrbit          string `json:"linkOrbit"`

	// Estatisticas Time Casa (do Redis/oraculo)
	PosseBolaTimeCasa           FlexValue `json:"posseBolaTimeCasa"`
	ChutesGolTimeCasa           FlexValue `json:"chutesGolTimeCasa"`
	ChutesForaTimeCasa          FlexValue `json:"chutesForaTimeCasa"`
	ChutesTraveTimeCasa         FlexValue `json:"chutesTraveTimeCasa"`
	ChutesBloqueadoTimeCasa     FlexValue `json:"chutesBloqueadoTimeCasa"`
	EscanteiosTimeCasa          FlexValue `json:"escanteiosTimeCasa"`
	AtaquesPerigososTimeCasa    FlexValue `json:"ataquesPerigososTimeCasa"`
	PenalidadesTimeCasa         FlexValue `json:"penalidadesTimeCasa"`
	ProbabilidadesTimeCasa      FlexValue `json:"probabilidadesTimeCasa"`
	Pontos10MinTimeCasa         FlexValue `json:"pontos10MinTimeCasa"`

	// Estatisticas 1 Tempo Casa
	ChutesGolTimeCasa1Tempo        FlexValue `json:"chutesGolTimeCasa1Tempo"`
	ChutesForaTimeCasa1Tempo       FlexValue `json:"chutesForaTimeCasa1Tempo"`
	ChutesTraveTimeCasa1Tempo      FlexValue `json:"chutesTraveTimeCasa1Tempo"`
	ChutesBloqueadoTimeCasa1Tempo  FlexValue `json:"chutesBloqueadoTimeCasa1Tempo"`
	EscanteiosTimeCasa1Tempo       FlexValue `json:"escanteiosTimeCasa1Tempo"`
	AtaquesPerigososTimeCasa1Tempo FlexValue `json:"ataquesPerigososTimeCasa1Tempo"`
	PenalidadesTimeCasa1Tempo      FlexValue `json:"penalidadesTimeCasa1Tempo"`

	// Estatisticas 2 Tempo Casa
	ChutesGolTimeCasa2Tempo        FlexValue `json:"chutesGolTimeCasa2Tempo"`
	ChutesForaTimeCasa2Tempo       FlexValue `json:"chutesForaTimeCasa2Tempo"`
	ChutesTraveTimeCasa2Tempo      FlexValue `json:"chutesTraveTimeCasa2Tempo"`
	ChutesBloqueadoTimeCasa2Tempo  FlexValue `json:"chutesBloqueadoTimeCasa2Tempo"`
	EscanteiosTimeCasa2Tempo       FlexValue `json:"escanteiosTimeCasa2Tempo"`
	AtaquesPerigososTimeCasa2Tempo FlexValue `json:"ataquesPerigososTimeCasa2Tempo"`
	PenalidadesTimeCasa2Tempo      FlexValue `json:"penalidadesTimeCasa2Tempo"`

	// Estatisticas 5 Min Casa
	ChutesGolTimeCasa5Min        FlexValue `json:"chutesGolTimeCasa5Min"`
	ChutesForaTimeCasa5Min       FlexValue `json:"chutesForaTimeCasa5Min"`
	ChutesTraveTimeCasa5Min      FlexValue `json:"chutesTraveTimeCasa5Min"`
	ChutesBloqueadoTimeCasa5Min  FlexValue `json:"chutesBloqueadoTimeCasa5Min"`
	EscanteiosTimeCasa5Min       FlexValue `json:"escanteiosTimeCasa5Min"`
	AtaquesPerigososTimeCasa5Min FlexValue `json:"ataquesPerigososTimeCasa5Min"`
	PenalidadesTimeCasa5Min      FlexValue `json:"penalidadesTimeCasa5Min"`

	// Estatisticas 10 Min Casa
	ChutesGolTimeCasa10Min        FlexValue `json:"chutesGolTimeCasa10Min"`
	ChutesForaTimeCasa10Min       FlexValue `json:"chutesForaTimeCasa10Min"`
	ChutesTraveTimeCasa10Min      FlexValue `json:"chutesTraveTimeCasa10Min"`
	ChutesBloqueadoTimeCasa10Min  FlexValue `json:"chutesBloqueadoTimeCasa10Min"`
	EscanteiosTimeCasa10Min       FlexValue `json:"escanteiosTimeCasa10Min"`
	AtaquesPerigososTimeCasa10Min FlexValue `json:"ataquesPerigososTimeCasa10Min"`
	PenalidadesTimeCasa10Min      FlexValue `json:"penalidadesTimeCasa10Min"`

	// Classes CSS Time Casa
	ClassPosseBolaTimeCasa           string `json:"classPosseBolaTimeCasa"`
	ClassChutesGolTimeCasa           string `json:"classChutesGolTimeCasa"`
	ClassChutesForaTimeCasa          string `json:"classChutesForaTimeCasa"`
	ClassChutesTraveTimeCasa         string `json:"classChutesTraveTimeCasa"`
	ClassChutesBloqueadoTimeCasa     string `json:"classChutesBloqueadoTimeCasa"`
	ClassEscanteiosTimeCasa          string `json:"classEscanteiosTimeCasa"`
	ClassAtaquesPerigososTimeCasa    string `json:"classAtaquesPerigososTimeCasa"`
	ClassPenalidadesTimeCasa         string `json:"classPenalidadesTimeCasa"`
	ClassProbabilidadesTimeCasa      string `json:"classProbabilidadesTimeCasa"`
	ClassPontos10MinTimeCasa         string `json:"classPontos10MinTimeCasa"`

	// Classes CSS 1 Tempo Casa
	ClassChutesGolTimeCasa1Tempo        string `json:"classChutesGolTimeCasa1Tempo"`
	ClassChutesForaTimeCasa1Tempo       string `json:"classChutesForaTimeCasa1Tempo"`
	ClassChutesTraveTimeCasa1Tempo      string `json:"classChutesTraveTimeCasa1Tempo"`
	ClassChutesBloqueadoTimeCasa1Tempo  string `json:"classChutesBloqueadoTimeCasa1Tempo"`
	ClassEscanteiosTimeCasa1Tempo       string `json:"classEscanteiosTimeCasa1Tempo"`
	ClassAtaquesPerigososTimeCasa1Tempo string `json:"classAtaquesPerigososTimeCasa1Tempo"`
	ClassPenalidadesTimeCasa1Tempo      string `json:"classPenalidadesTimeCasa1Tempo"`

	// Classes CSS 2 Tempo Casa
	ClassChutesGolTimeCasa2Tempo        string `json:"classChutesGolTimeCasa2Tempo"`
	ClassChutesForaTimeCasa2Tempo       string `json:"classChutesForaTimeCasa2Tempo"`
	ClassChutesTraveTimeCasa2Tempo      string `json:"classChutesTraveTimeCasa2Tempo"`
	ClassChutesBloqueadoTimeCasa2Tempo  string `json:"classChutesBloqueadoTimeCasa2Tempo"`
	ClassEscanteiosTimeCasa2Tempo       string `json:"classEscanteiosTimeCasa2Tempo"`
	ClassAtaquesPerigososTimeCasa2Tempo string `json:"classAtaquesPerigososTimeCasa2Tempo"`
	ClassPenalidadesTimeCasa2Tempo      string `json:"classPenalidadesTimeCasa2Tempo"`

	// Classes CSS 5 Min Casa
	ClassChutesGolTimeCasa5Min        string `json:"classChutesGolTimeCasa5Min"`
	ClassChutesForaTimeCasa5Min       string `json:"classChutesForaTimeCasa5Min"`
	ClassChutesTraveTimeCasa5Min      string `json:"classChutesTraveTimeCasa5Min"`
	ClassChutesBloqueadoTimeCasa5Min  string `json:"classChutesBloqueadoTimeCasa5Min"`
	ClassEscanteiosTimeCasa5Min       string `json:"classEscanteiosTimeCasa5Min"`
	ClassAtaquesPerigososTimeCasa5Min string `json:"classAtaquesPerigososTimeCasa5Min"`
	ClassPenalidadesTimeCasa5Min      string `json:"classPenalidadesTimeCasa5Min"`

	// Classes CSS 10 Min Casa
	ClassChutesGolTimeCasa10Min        string `json:"classChutesGolTimeCasa10Min"`
	ClassChutesForaTimeCasa10Min       string `json:"classChutesForaTimeCasa10Min"`
	ClassChutesTraveTimeCasa10Min      string `json:"classChutesTraveTimeCasa10Min"`
	ClassChutesBloqueadoTimeCasa10Min  string `json:"classChutesBloqueadoTimeCasa10Min"`
	ClassEscanteiosTimeCasa10Min       string `json:"classEscanteiosTimeCasa10Min"`
	ClassAtaquesPerigososTimeCasa10Min string `json:"classAtaquesPerigososTimeCasa10Min"`
	ClassPenalidadesTimeCasa10Min      string `json:"classPenalidadesTimeCasa10Min"`

	// Estatisticas Time Fora (do Redis/oraculo)
	PosseBolaTimeFora           FlexValue `json:"posseBolaTimeFora"`
	ChutesGolTimeFora           FlexValue `json:"chutesGolTimeFora"`
	ChutesForaTimeFora          FlexValue `json:"chutesForaTimeFora"`
	ChutesTraveTimeFora         FlexValue `json:"chutesTraveTimeFora"`
	ChutesBloqueadoTimeFora     FlexValue `json:"chutesBloqueadoTimeFora"`
	EscanteiosTimeFora          FlexValue `json:"escanteiosTimeFora"`
	AtaquesPerigososTimeFora    FlexValue `json:"ataquesPerigososTimeFora"`
	PenalidadesTimeFora         FlexValue `json:"penalidadesTimeFora"`
	ProbabilidadesTimeFora      FlexValue `json:"probabilidadesTimeFora"`
	Pontos10MinTimeFora         FlexValue `json:"pontos10MinTimeFora"`

	// Estatisticas 1 Tempo Fora
	ChutesGolTimeFora1Tempo        FlexValue `json:"chutesGolTimeFora1Tempo"`
	ChutesForaTimeFora1Tempo       FlexValue `json:"chutesForaTimeFora1Tempo"`
	ChutesTraveTimeFora1Tempo      FlexValue `json:"chutesTraveTimeFora1Tempo"`
	ChutesBloqueadoTimeFora1Tempo  FlexValue `json:"chutesBloqueadoTimeFora1Tempo"`
	EscanteiosTimeFora1Tempo       FlexValue `json:"escanteiosTimeFora1Tempo"`
	AtaquesPerigososTimeFora1Tempo FlexValue `json:"ataquesPerigososTimeFora1Tempo"`
	PenalidadesTimeFora1Tempo      FlexValue `json:"penalidadesTimeFora1Tempo"`

	// Estatisticas 2 Tempo Fora
	ChutesGolTimeFora2Tempo        FlexValue `json:"chutesGolTimeFora2Tempo"`
	ChutesForaTimeFora2Tempo       FlexValue `json:"chutesForaTimeFora2Tempo"`
	ChutesTraveTimeFora2Tempo      FlexValue `json:"chutesTraveTimeFora2Tempo"`
	ChutesBloqueadoTimeFora2Tempo  FlexValue `json:"chutesBloqueadoTimeFora2Tempo"`
	EscanteiosTimeFora2Tempo       FlexValue `json:"escanteiosTimeFora2Tempo"`
	AtaquesPerigososTimeFora2Tempo FlexValue `json:"ataquesPerigososTimeFora2Tempo"`
	PenalidadesTimeFora2Tempo      FlexValue `json:"penalidadesTimeFora2Tempo"`

	// Estatisticas 5 Min Fora
	ChutesGolTimeFora5Min        FlexValue `json:"chutesGolTimeFora5Min"`
	ChutesForaTimeFora5Min       FlexValue `json:"chutesForaTimeFora5Min"`
	ChutesTraveTimeFora5Min      FlexValue `json:"chutesTraveTimeFora5Min"`
	ChutesBloqueadoTimeFora5Min  FlexValue `json:"chutesBloqueadoTimeFora5Min"`
	EscanteiosTimeFora5Min       FlexValue `json:"escanteiosTimeFora5Min"`
	AtaquesPerigososTimeFora5Min FlexValue `json:"ataquesPerigososTimeFora5Min"`
	PenalidadesTimeFora5Min      FlexValue `json:"penalidadesTimeFora5Min"`

	// Estatisticas 10 Min Fora
	ChutesGolTimeFora10Min        FlexValue `json:"chutesGolTimeFora10Min"`
	ChutesForaTimeFora10Min       FlexValue `json:"chutesForaTimeFora10Min"`
	ChutesTraveTimeFora10Min      FlexValue `json:"chutesTraveTimeFora10Min"`
	ChutesBloqueadoTimeFora10Min  FlexValue `json:"chutesBloqueadoTimeFora10Min"`
	EscanteiosTimeFora10Min       FlexValue `json:"escanteiosTimeFora10Min"`
	AtaquesPerigososTimeFora10Min FlexValue `json:"ataquesPerigososTimeFora10Min"`
	PenalidadesTimeFora10Min      FlexValue `json:"penalidadesTimeFora10Min"`

	// Classes CSS Time Fora
	ClassPosseBolaTimeFora           string `json:"classPosseBolaTimeFora"`
	ClassChutesGolTimeFora           string `json:"classChutesGolTimeFora"`
	ClassChutesForaTimeFora          string `json:"classChutesForaTimeFora"`
	ClassChutesTraveTimeFora         string `json:"classChutesTraveTimeFora"`
	ClassChutesBloqueadoTimeFora     string `json:"classChutesBloqueadoTimeFora"`
	ClassEscanteiosTimeFora          string `json:"classEscanteiosTimeFora"`
	ClassAtaquesPerigososTimeFora    string `json:"classAtaquesPerigososTimeFora"`
	ClassPenalidadesTimeFora         string `json:"classPenalidadesTimeFora"`
	ClassProbabilidadesTimeFora      string `json:"classProbabilidadesTimeFora"`
	ClassPontos10MinTimeFora         string `json:"classPontos10MinTimeFora"`

	// Classes CSS 1 Tempo Fora
	ClassChutesGolTimeFora1Tempo        string `json:"classChutesGolTimeFora1Tempo"`
	ClassChutesForaTimeFora1Tempo       string `json:"classChutesForaTimeFora1Tempo"`
	ClassChutesTraveTimeFora1Tempo      string `json:"classChutesTraveTimeFora1Tempo"`
	ClassChutesBloqueadoTimeFora1Tempo  string `json:"classChutesBloqueadoTimeFora1Tempo"`
	ClassEscanteiosTimeFora1Tempo       string `json:"classEscanteiosTimeFora1Tempo"`
	ClassAtaquesPerigososTimeFora1Tempo string `json:"classAtaquesPerigososTimeFora1Tempo"`
	ClassPenalidadesTimeFora1Tempo      string `json:"classPenalidadesTimeFora1Tempo"`

	// Classes CSS 2 Tempo Fora
	ClassChutesGolTimeFora2Tempo        string `json:"classChutesGolTimeFora2Tempo"`
	ClassChutesForaTimeFora2Tempo       string `json:"classChutesForaTimeFora2Tempo"`
	ClassChutesTraveTimeFora2Tempo      string `json:"classChutesTraveTimeFora2Tempo"`
	ClassChutesBloqueadoTimeFora2Tempo  string `json:"classChutesBloqueadoTimeFora2Tempo"`
	ClassEscanteiosTimeFora2Tempo       string `json:"classEscanteiosTimeFora2Tempo"`
	ClassAtaquesPerigososTimeFora2Tempo string `json:"classAtaquesPerigososTimeFora2Tempo"`
	ClassPenalidadesTimeFora2Tempo      string `json:"classPenalidadesTimeFora2Tempo"`

	// Classes CSS 5 Min Fora
	ClassChutesGolTimeFora5Min        string `json:"classChutesGolTimeFora5Min"`
	ClassChutesForaTimeFora5Min       string `json:"classChutesForaTimeFora5Min"`
	ClassChutesTraveTimeFora5Min      string `json:"classChutesTraveTimeFora5Min"`
	ClassChutesBloqueadoTimeFora5Min  string `json:"classChutesBloqueadoTimeFora5Min"`
	ClassEscanteiosTimeFora5Min       string `json:"classEscanteiosTimeFora5Min"`
	ClassAtaquesPerigososTimeFora5Min string `json:"classAtaquesPerigososTimeFora5Min"`
	ClassPenalidadesTimeFora5Min      string `json:"classPenalidadesTimeFora5Min"`

	// Classes CSS 10 Min Fora
	ClassChutesGolTimeFora10Min        string `json:"classChutesGolTimeFora10Min"`
	ClassChutesForaTimeFora10Min       string `json:"classChutesForaTimeFora10Min"`
	ClassChutesTraveTimeFora10Min      string `json:"classChutesTraveTimeFora10Min"`
	ClassChutesBloqueadoTimeFora10Min  string `json:"classChutesBloqueadoTimeFora10Min"`
	ClassEscanteiosTimeFora10Min       string `json:"classEscanteiosTimeFora10Min"`
	ClassAtaquesPerigososTimeFora10Min string `json:"classAtaquesPerigososTimeFora10Min"`
	ClassPenalidadesTimeFora10Min      string `json:"classPenalidadesTimeFora10Min"`

	// Score de Lances (SL) - usado para alertas de pressao
	ScoreLances10MinTimeCasa      FlexValue `json:"scoreLances10MinTimeCasa"`
	ScoreLances10MinTimeFora      FlexValue `json:"scoreLances10MinTimeFora"`
	ClassScoreLances10MinTimeCasa string    `json:"classScoreLances10MinTimeCasa"`
	ClassScoreLances10MinTimeFora string    `json:"classScoreLances10MinTimeFora"`
	ScoreLances5MinTimeCasa       FlexValue `json:"scoreLances5MinTimeCasa"`
	ScoreLances5MinTimeFora       FlexValue `json:"scoreLances5MinTimeFora"`
	ClassScoreLances5MinTimeCasa  string    `json:"classScoreLances5MinTimeCasa"`
	ClassScoreLances5MinTimeFora  string    `json:"classScoreLances5MinTimeFora"`

	// Alertas (FlexBool pois PHP envia 0/1)
	AlertarGolTimeCasa           FlexBool `json:"alertarGolTimeCasa"`
	AlertarPenalTimeCasa         FlexBool `json:"alertarPenalTimeCasa"`
	AlertarGolTimeFora           FlexBool `json:"alertarGolTimeFora"`
	AlertarPenalTimeFora         FlexBool `json:"alertarPenalTimeFora"`
	AlertarSomGol                FlexBool `json:"alertarSomGol"`
	Cuidado                      FlexBool `json:"cuidado"`
	AlertaMomentoGolAtivo        FlexBool  `json:"alertaMomentoGolAtivo"`
	AlertaMomentoGolValor        FlexValue `json:"alertaMomentoGolValor"`
	AlertaPressaoIndividualAtivo FlexBool  `json:"alertaPressaoIndividualAtivo"`
	AlertaPressaoIndividualTime  string    `json:"alertaPressaoIndividualTime"`
	AlertaPressaoIndividualNome  string    `json:"alertaPressaoIndividualNome"`
	AlertaPressaoIndividualValor FlexValue `json:"alertaPressaoIndividualValor"`
	PressaoTimeCasa              FlexValue `json:"pressaoTimeCasa"`
	PressaoTimeFora              FlexValue `json:"pressaoTimeFora"`
	ClassPressaoTimeCasa         string    `json:"classPressaoTimeCasa"`
	ClassPressaoTimeFora         string    `json:"classPressaoTimeFora"`
	SomaPressao                  FlexValue `json:"somaPressao"`

	// Icones
	IconeComentarioTimeCasa string `json:"iconeComentarioTimeCasa"`
	IconeComentarioTimeFora string `json:"iconeComentarioTimeFora"`

	// Acrescimos
	Acrescimo1Tempo              FlexValue `json:"acrescimo1Tempo"`
	Acrescimo2Tempo              FlexValue `json:"acrescimo2Tempo"`
	ClassAcrescimo1Tempo         string    `json:"classAcrescimo1Tempo"`
	ClassAcrescimo2Tempo         string    `json:"classAcrescimo2Tempo"`
	PrevisaoAcrescimo1Tempo      FlexValue `json:"previsaoAcrescimo1Tempo"`
	PrevisaoAcrescimo2Tempo      FlexValue `json:"previsaoAcrescimo2Tempo"`
	ClassPrevisaoAcrescimo1Tempo string    `json:"classPrevisaoAcrescimo1Tempo"`
	ClassPrevisaoAcrescimo2Tempo string    `json:"classPrevisaoAcrescimo2Tempo"`

	// Extras
	AnaliseIA          string   `json:"analiseIA"`
	TemAnaliseIA       bool     `json:"temAnaliseIA"`
	TeamStreaks        []any    `json:"teamStreaks"`
	Favorito           FlexBool `json:"favorito"`
	CampeonatoFavorito FlexBool `json:"campeonatoFavorito"`
}

// Campeonato representa um campeonato com seus eventos
type Campeonato struct {
	Id                     int                `json:"id"`
	IdUnico                string             `json:"idUnico"`
	IdTemporada            string             `json:"idTemporada"`
	Flag                   string             `json:"flag"`
	NomeCategoria          string             `json:"nomeCategoria"`
	SlugCategoria          string             `json:"slugCategoria"`
	NomeCampeonato         string             `json:"nomeCampeonato"`
	NomeCampeonatoReduzido string             `json:"nomeCampeonatoReduzido"`
	SlugCampeonato         string             `json:"slugCampeonato"`
	TemClassificacao       int                `json:"temClassificacao"`
	Prioridade             int                `json:"prioridade"`
	Favorito               FlexBool           `json:"favorito"`
	Sequencia              int                `json:"sequencia"`
	Eventos                map[string]*Evento `json:"eventos"`
}

// Counts contadores para o response
type Counts struct {
	Live  int `json:"live"`
	Total int `json:"total"`
	Gols  int `json:"gols"`
}

// PainelResponse response do endpoint /sse/painel
type PainelResponse struct {
	Eventos []*Evento `json:"eventos"`
	Counts  Counts    `json:"counts"`
}

// HomeResponse response do endpoint /sse/home
type HomeResponse struct {
	Campeonatos []*Campeonato `json:"campeonatos"`
	Counts      Counts        `json:"counts"`
}

// Filtro esta definido em filtro.go

// FiltrarParaFree retorna uma copia do evento com apenas campos liberados para usuarios free/anonimos
// Remove: SL, alertas avancados, estatisticas, probabilidades, xG, pontos10Min, pressao, analise IA
func (e *Evento) FiltrarParaFree() *Evento {
	// Cria copia do evento
	copia := *e

	// Zera Score de Lances (SL)
	copia.ScoreLances10MinTimeCasa = ""
	copia.ScoreLances10MinTimeFora = ""
	copia.ClassScoreLances10MinTimeCasa = ""
	copia.ClassScoreLances10MinTimeFora = ""
	copia.ScoreLances5MinTimeCasa = ""
	copia.ScoreLances5MinTimeFora = ""
	copia.ClassScoreLances5MinTimeCasa = ""
	copia.ClassScoreLances5MinTimeFora = ""

	// Zera alertas avancados (momento gol e pressao individual)
	copia.AlertaMomentoGolAtivo = FlexBool(false)
	copia.AlertaMomentoGolValor = ""
	copia.AlertaPressaoIndividualAtivo = FlexBool(false)
	copia.AlertaPressaoIndividualTime = ""
	copia.AlertaPressaoIndividualNome = ""
	copia.AlertaPressaoIndividualValor = ""
	copia.PressaoTimeCasa = ""
	copia.PressaoTimeFora = ""
	copia.ClassPressaoTimeCasa = ""
	copia.ClassPressaoTimeFora = ""
	copia.SomaPressao = ""

	// Zera estatisticas Time Casa
	copia.PosseBolaTimeCasa = ""
	copia.ChutesGolTimeCasa = ""
	copia.ChutesForaTimeCasa = ""
	copia.ChutesTraveTimeCasa = ""
	copia.ChutesBloqueadoTimeCasa = ""
	copia.EscanteiosTimeCasa = ""
	copia.AtaquesPerigososTimeCasa = ""
	copia.PenalidadesTimeCasa = ""
	copia.ProbabilidadesTimeCasa = ""
	copia.Pontos10MinTimeCasa = ""

	// Zera estatisticas 1 Tempo Casa
	copia.ChutesGolTimeCasa1Tempo = ""
	copia.ChutesForaTimeCasa1Tempo = ""
	copia.ChutesTraveTimeCasa1Tempo = ""
	copia.ChutesBloqueadoTimeCasa1Tempo = ""
	copia.EscanteiosTimeCasa1Tempo = ""
	copia.AtaquesPerigososTimeCasa1Tempo = ""
	copia.PenalidadesTimeCasa1Tempo = ""

	// Zera estatisticas 2 Tempo Casa
	copia.ChutesGolTimeCasa2Tempo = ""
	copia.ChutesForaTimeCasa2Tempo = ""
	copia.ChutesTraveTimeCasa2Tempo = ""
	copia.ChutesBloqueadoTimeCasa2Tempo = ""
	copia.EscanteiosTimeCasa2Tempo = ""
	copia.AtaquesPerigososTimeCasa2Tempo = ""
	copia.PenalidadesTimeCasa2Tempo = ""

	// Zera estatisticas 5 Min Casa
	copia.ChutesGolTimeCasa5Min = ""
	copia.ChutesForaTimeCasa5Min = ""
	copia.ChutesTraveTimeCasa5Min = ""
	copia.ChutesBloqueadoTimeCasa5Min = ""
	copia.EscanteiosTimeCasa5Min = ""
	copia.AtaquesPerigososTimeCasa5Min = ""
	copia.PenalidadesTimeCasa5Min = ""

	// Zera estatisticas 10 Min Casa
	copia.ChutesGolTimeCasa10Min = ""
	copia.ChutesForaTimeCasa10Min = ""
	copia.ChutesTraveTimeCasa10Min = ""
	copia.ChutesBloqueadoTimeCasa10Min = ""
	copia.EscanteiosTimeCasa10Min = ""
	copia.AtaquesPerigososTimeCasa10Min = ""
	copia.PenalidadesTimeCasa10Min = ""

	// Zera classes CSS Time Casa
	copia.ClassPosseBolaTimeCasa = ""
	copia.ClassChutesGolTimeCasa = ""
	copia.ClassChutesForaTimeCasa = ""
	copia.ClassChutesTraveTimeCasa = ""
	copia.ClassChutesBloqueadoTimeCasa = ""
	copia.ClassEscanteiosTimeCasa = ""
	copia.ClassAtaquesPerigososTimeCasa = ""
	copia.ClassPenalidadesTimeCasa = ""
	copia.ClassProbabilidadesTimeCasa = ""
	copia.ClassPontos10MinTimeCasa = ""

	// Zera classes CSS 1 Tempo Casa
	copia.ClassChutesGolTimeCasa1Tempo = ""
	copia.ClassChutesForaTimeCasa1Tempo = ""
	copia.ClassChutesTraveTimeCasa1Tempo = ""
	copia.ClassChutesBloqueadoTimeCasa1Tempo = ""
	copia.ClassEscanteiosTimeCasa1Tempo = ""
	copia.ClassAtaquesPerigososTimeCasa1Tempo = ""
	copia.ClassPenalidadesTimeCasa1Tempo = ""

	// Zera classes CSS 2 Tempo Casa
	copia.ClassChutesGolTimeCasa2Tempo = ""
	copia.ClassChutesForaTimeCasa2Tempo = ""
	copia.ClassChutesTraveTimeCasa2Tempo = ""
	copia.ClassChutesBloqueadoTimeCasa2Tempo = ""
	copia.ClassEscanteiosTimeCasa2Tempo = ""
	copia.ClassAtaquesPerigososTimeCasa2Tempo = ""
	copia.ClassPenalidadesTimeCasa2Tempo = ""

	// Zera classes CSS 5 Min Casa
	copia.ClassChutesGolTimeCasa5Min = ""
	copia.ClassChutesForaTimeCasa5Min = ""
	copia.ClassChutesTraveTimeCasa5Min = ""
	copia.ClassChutesBloqueadoTimeCasa5Min = ""
	copia.ClassEscanteiosTimeCasa5Min = ""
	copia.ClassAtaquesPerigososTimeCasa5Min = ""
	copia.ClassPenalidadesTimeCasa5Min = ""

	// Zera classes CSS 10 Min Casa
	copia.ClassChutesGolTimeCasa10Min = ""
	copia.ClassChutesForaTimeCasa10Min = ""
	copia.ClassChutesTraveTimeCasa10Min = ""
	copia.ClassChutesBloqueadoTimeCasa10Min = ""
	copia.ClassEscanteiosTimeCasa10Min = ""
	copia.ClassAtaquesPerigososTimeCasa10Min = ""
	copia.ClassPenalidadesTimeCasa10Min = ""

	// Zera estatisticas Time Fora
	copia.PosseBolaTimeFora = ""
	copia.ChutesGolTimeFora = ""
	copia.ChutesForaTimeFora = ""
	copia.ChutesTraveTimeFora = ""
	copia.ChutesBloqueadoTimeFora = ""
	copia.EscanteiosTimeFora = ""
	copia.AtaquesPerigososTimeFora = ""
	copia.PenalidadesTimeFora = ""
	copia.ProbabilidadesTimeFora = ""
	copia.Pontos10MinTimeFora = ""

	// Zera estatisticas 1 Tempo Fora
	copia.ChutesGolTimeFora1Tempo = ""
	copia.ChutesForaTimeFora1Tempo = ""
	copia.ChutesTraveTimeFora1Tempo = ""
	copia.ChutesBloqueadoTimeFora1Tempo = ""
	copia.EscanteiosTimeFora1Tempo = ""
	copia.AtaquesPerigososTimeFora1Tempo = ""
	copia.PenalidadesTimeFora1Tempo = ""

	// Zera estatisticas 2 Tempo Fora
	copia.ChutesGolTimeFora2Tempo = ""
	copia.ChutesForaTimeFora2Tempo = ""
	copia.ChutesTraveTimeFora2Tempo = ""
	copia.ChutesBloqueadoTimeFora2Tempo = ""
	copia.EscanteiosTimeFora2Tempo = ""
	copia.AtaquesPerigososTimeFora2Tempo = ""
	copia.PenalidadesTimeFora2Tempo = ""

	// Zera estatisticas 5 Min Fora
	copia.ChutesGolTimeFora5Min = ""
	copia.ChutesForaTimeFora5Min = ""
	copia.ChutesTraveTimeFora5Min = ""
	copia.ChutesBloqueadoTimeFora5Min = ""
	copia.EscanteiosTimeFora5Min = ""
	copia.AtaquesPerigososTimeFora5Min = ""
	copia.PenalidadesTimeFora5Min = ""

	// Zera estatisticas 10 Min Fora
	copia.ChutesGolTimeFora10Min = ""
	copia.ChutesForaTimeFora10Min = ""
	copia.ChutesTraveTimeFora10Min = ""
	copia.ChutesBloqueadoTimeFora10Min = ""
	copia.EscanteiosTimeFora10Min = ""
	copia.AtaquesPerigososTimeFora10Min = ""
	copia.PenalidadesTimeFora10Min = ""

	// Zera classes CSS Time Fora
	copia.ClassPosseBolaTimeFora = ""
	copia.ClassChutesGolTimeFora = ""
	copia.ClassChutesForaTimeFora = ""
	copia.ClassChutesTraveTimeFora = ""
	copia.ClassChutesBloqueadoTimeFora = ""
	copia.ClassEscanteiosTimeFora = ""
	copia.ClassAtaquesPerigososTimeFora = ""
	copia.ClassPenalidadesTimeFora = ""
	copia.ClassProbabilidadesTimeFora = ""
	copia.ClassPontos10MinTimeFora = ""

	// Zera classes CSS 1 Tempo Fora
	copia.ClassChutesGolTimeFora1Tempo = ""
	copia.ClassChutesForaTimeFora1Tempo = ""
	copia.ClassChutesTraveTimeFora1Tempo = ""
	copia.ClassChutesBloqueadoTimeFora1Tempo = ""
	copia.ClassEscanteiosTimeFora1Tempo = ""
	copia.ClassAtaquesPerigososTimeFora1Tempo = ""
	copia.ClassPenalidadesTimeFora1Tempo = ""

	// Zera classes CSS 2 Tempo Fora
	copia.ClassChutesGolTimeFora2Tempo = ""
	copia.ClassChutesForaTimeFora2Tempo = ""
	copia.ClassChutesTraveTimeFora2Tempo = ""
	copia.ClassChutesBloqueadoTimeFora2Tempo = ""
	copia.ClassEscanteiosTimeFora2Tempo = ""
	copia.ClassAtaquesPerigososTimeFora2Tempo = ""
	copia.ClassPenalidadesTimeFora2Tempo = ""

	// Zera classes CSS 5 Min Fora
	copia.ClassChutesGolTimeFora5Min = ""
	copia.ClassChutesForaTimeFora5Min = ""
	copia.ClassChutesTraveTimeFora5Min = ""
	copia.ClassChutesBloqueadoTimeFora5Min = ""
	copia.ClassEscanteiosTimeFora5Min = ""
	copia.ClassAtaquesPerigososTimeFora5Min = ""
	copia.ClassPenalidadesTimeFora5Min = ""

	// Zera classes CSS 10 Min Fora
	copia.ClassChutesGolTimeFora10Min = ""
	copia.ClassChutesForaTimeFora10Min = ""
	copia.ClassChutesTraveTimeFora10Min = ""
	copia.ClassChutesBloqueadoTimeFora10Min = ""
	copia.ClassEscanteiosTimeFora10Min = ""
	copia.ClassAtaquesPerigososTimeFora10Min = ""
	copia.ClassPenalidadesTimeFora10Min = ""

	// Zera acrescimos e previsao de acrescimos
	copia.Acrescimo1Tempo = ""
	copia.Acrescimo2Tempo = ""
	copia.ClassAcrescimo1Tempo = ""
	copia.ClassAcrescimo2Tempo = ""
	copia.PrevisaoAcrescimo1Tempo = ""
	copia.PrevisaoAcrescimo2Tempo = ""
	copia.ClassPrevisaoAcrescimo1Tempo = ""
	copia.ClassPrevisaoAcrescimo2Tempo = ""

	// Zera analise IA
	copia.AnaliseIA = ""

	return &copia
}
