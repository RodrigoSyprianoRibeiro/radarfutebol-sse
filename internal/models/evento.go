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

	// Classes CSS 10 Min Fora
	ClassChutesGolTimeFora10Min        string `json:"classChutesGolTimeFora10Min"`
	ClassChutesForaTimeFora10Min       string `json:"classChutesForaTimeFora10Min"`
	ClassChutesTraveTimeFora10Min      string `json:"classChutesTraveTimeFora10Min"`
	ClassChutesBloqueadoTimeFora10Min  string `json:"classChutesBloqueadoTimeFora10Min"`
	ClassEscanteiosTimeFora10Min       string `json:"classEscanteiosTimeFora10Min"`
	ClassAtaquesPerigososTimeFora10Min string `json:"classAtaquesPerigososTimeFora10Min"`
	ClassPenalidadesTimeFora10Min      string `json:"classPenalidadesTimeFora10Min"`

	// Alertas (FlexBool pois PHP envia 0/1)
	AlertarGolTimeCasa           FlexBool `json:"alertarGolTimeCasa"`
	AlertarPenalTimeCasa         FlexBool `json:"alertarPenalTimeCasa"`
	AlertarGolTimeFora           FlexBool `json:"alertarGolTimeFora"`
	AlertarPenalTimeFora         FlexBool `json:"alertarPenalTimeFora"`
	AlertarSomGol                FlexBool `json:"alertarSomGol"`
	Cuidado                      FlexBool `json:"cuidado"`
	AlertaMomentoGolAtivo        FlexBool `json:"alertaMomentoGolAtivo"`
	AlertaPressaoIndividualAtivo FlexBool `json:"alertaPressaoIndividualAtivo"`

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
