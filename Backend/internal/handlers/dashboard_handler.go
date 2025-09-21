package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"example.com/m/v2/internal/services"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	dataService        *services.DataService
	geolocationService *services.GeolocationService
	db                 *gorm.DB
}

// Estruturas de resposta para o dashboard
type DashboardStats struct {
	Totalizadores Totalizadores `json:"totalizadores"`
	Graficos      Graficos      `json:"graficos"`
	Regioes       Regioes       `json:"regioes"`
}

type Totalizadores struct {
	TotalHospitais        int     `json:"total_hospitais"`
	TotalMedicos          int     `json:"total_medicos"`
	TotalPacientes        int     `json:"total_pacientes"`
	TotalLeitos           int     `json:"total_leitos"`
	MediaLeitosHospital   float64 `json:"media_leitos_hospital"`
	PercentualConvenio    float64 `json:"percentual_convenio"`
}

type Graficos struct {
	Top10Doencas        []DoencaCount `json:"top_10_doencas"`
	DoencaPorEstado     []EstadoDoenca `json:"doenca_por_estado"`
	PercentualPopulacao float64        `json:"percentual_populacao_doente"`
}

type DoencaCount struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Quantidade int   `json:"quantidade"`
}

type EstadoDoenca struct {
	Estado     string `json:"estado"`
	Doenca     string `json:"doenca"`
	Quantidade int    `json:"quantidade"`
}

type Regioes struct {
	RegiaoMaisDoente string  `json:"regiao_mais_doente"`
	PercentualRegiao float64 `json:"percentual_regiao"`
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{
		dataService:        services.NewDataService(db),
		geolocationService: services.NewGeolocationService(),
		db:                 db,
	}
}

// GetDashboardStats - Endpoint principal para estatísticas do dashboard
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Buscar totalizadores
	totalizadores, err := h.getTotalizadores()
	if err != nil {
		h.sendErrorResponse(w, "Erro ao buscar totalizadores", err)
		return
	}

	// Buscar dados para gráficos
	graficos, err := h.getGraficos()
	if err != nil {
		h.sendErrorResponse(w, "Erro ao buscar dados dos gráficos", err)
		return
	}

	// Buscar dados de regiões
	regioes, err := h.getRegioes()
	if err != nil {
		h.sendErrorResponse(w, "Erro ao buscar dados de regiões", err)
		return
	}

	response := DashboardStats{
		Totalizadores: totalizadores,
		Graficos:      graficos,
		Regioes:       regioes,
	}

	h.sendSuccessResponse(w, response)
}

// getTotalizadores - Calcula os totalizadores principais
func (h *DashboardHandler) getTotalizadores() (Totalizadores, error) {
	var totalizadores Totalizadores

	// Total de hospitais
	if err := h.db.Model(&struct {
		ID int `gorm:"primaryKey"`
	}{}).Table("hospitais").Count(&totalizadores.TotalHospitais).Error; err != nil {
		return totalizadores, err
	}

	// Total de médicos
	if err := h.db.Model(&struct {
		UUID string `gorm:"primaryKey"`
	}{}).Table("medicos").Count(&totalizadores.TotalMedicos).Error; err != nil {
		return totalizadores, err
	}

	// Total de pacientes
	if err := h.db.Model(&struct {
		ID string `gorm:"primaryKey"`
	}{}).Table("pacientes").Count(&totalizadores.TotalPacientes).Error; err != nil {
		return totalizadores, err
	}

	// Total de leitos e média de leitos por hospital
	var totalLeitos int
	if err := h.db.Table("hospitais").Select("COALESCE(SUM(leitos_totais), 0)").Scan(&totalLeitos).Error; err != nil {
		return totalizadores, err
	}
	totalizadores.TotalLeitos = totalLeitos

	if totalizadores.TotalHospitais > 0 {
		totalizadores.MediaLeitosHospital = float64(totalLeitos) / float64(totalizadores.TotalHospitais)
	}

	// Percentual de pacientes com convênio
	var totalPacientesComConvenio int
	if err := h.db.Table("pacientes").
		Where("convenio IS NOT NULL AND convenio != '' AND convenio != 'NAO'").
		Count(&totalPacientesComConvenio).Error; err != nil {
		return totalizadores, err
	}

	if totalizadores.TotalPacientes > 0 {
		totalizadores.PercentualConvenio = (float64(totalPacientesComConvenio) / float64(totalizadores.TotalPacientes)) * 100
	}

	return totalizadores, nil
}

// getGraficos - Calcula dados para os gráficos
func (h *DashboardHandler) getGraficos() (Graficos, error) {
	var graficos Graficos

	// Top 10 doenças
	top10Doencas, err := h.getTop10Doencas()
	if err != nil {
		return graficos, err
	}
	graficos.Top10Doencas = top10Doencas

	// Doença mais comum por estado
	doencaPorEstado, err := h.getDoencaPorEstado()
	if err != nil {
		return graficos, err
	}
	graficos.DoencaPorEstado = doencaPorEstado

	// Percentual da população doente
	percentualPopulacao, err := h.getPercentualPopulacaoDoente()
	if err != nil {
		return graficos, err
	}
	graficos.PercentualPopulacao = percentualPopulacao

	return graficos, nil
}

// getTop10Doencas - Busca as 10 doenças mais comuns
func (h *DashboardHandler) getTop10Doencas() ([]DoencaCount, error) {
	var doencas []DoencaCount

	query := `
		SELECT 
			c.codigo,
			c.descricao,
			COUNT(p.id) as quantidade
		FROM cid10s c
		INNER JOIN pacientes p ON p.cid_codigo = c.codigo
		GROUP BY c.codigo, c.descricao
		ORDER BY quantidade DESC
		LIMIT 10
	`

	if err := h.db.Raw(query).Scan(&doencas).Error; err != nil {
		return doencas, err
	}

	return doencas, nil
}

// getDoencaPorEstado - Busca a doença mais comum por estado
func (h *DashboardHandler) getDoencaPorEstado() ([]EstadoDoenca, error) {
	var estadosDoencas []EstadoDoenca

	query := `
		WITH doencas_por_estado AS (
			SELECT 
				e.nome as estado,
				c.codigo as doenca_codigo,
				c.descricao as doenca,
				COUNT(p.id) as quantidade,
				ROW_NUMBER() OVER (PARTITION BY e.codigo ORDER BY COUNT(p.id) DESC) as rn
			FROM estados e
			INNER JOIN municipios m ON m.codigo_uf = e.codigo
			INNER JOIN pacientes p ON p.municipio_codigo = m.codigo
			INNER JOIN cid10s c ON c.codigo = p.cid_codigo
			GROUP BY e.codigo, e.nome, c.codigo, c.descricao
		)
		SELECT estado, doenca, quantidade
		FROM doencas_por_estado
		WHERE rn = 1
		ORDER BY quantidade DESC
	`

	if err := h.db.Raw(query).Scan(&estadosDoencas).Error; err != nil {
		return estadosDoencas, err
	}

	return estadosDoencas, nil
}

// getPercentualPopulacaoDoente - Calcula percentual da população doente
func (h *DashboardHandler) getPercentualPopulacaoDoente() (float64, error) {
	var totalPopulacao int
	var totalPacientes int

	// Total da população dos municípios
	if err := h.db.Table("municipios").Select("COALESCE(SUM(populacao), 0)").Scan(&totalPopulacao).Error; err != nil {
		return 0, err
	}

	// Total de pacientes
	if err := h.db.Table("pacientes").Count(&totalPacientes).Error; err != nil {
		return 0, err
	}

	if totalPopulacao > 0 {
		return (float64(totalPacientes) / float64(totalPopulacao)) * 100, nil
	}

	return 0, nil
}

// getRegioes - Calcula dados de regiões
func (h *DashboardHandler) getRegioes() (Regioes, error) {
	var regioes Regioes

	query := `
		WITH doencas_por_regiao AS (
			SELECT 
				e.regiao,
				COUNT(p.id) as total_pacientes,
				SUM(m.populacao) as total_populacao
			FROM estados e
			INNER JOIN municipios m ON m.codigo_uf = e.codigo
			INNER JOIN pacientes p ON p.municipio_codigo = m.codigo
			GROUP BY e.regiao
		),
		percentuais_regiao AS (
			SELECT 
				regiao,
				total_pacientes,
				total_populacao,
				CASE 
					WHEN total_populacao > 0 THEN (total_pacientes::float / total_populacao::float) * 100
					ELSE 0
				END as percentual
			FROM doencas_por_regiao
		)
		SELECT 
			regiao as regiao_mais_doente,
			percentual as percentual_regiao
		FROM percentuais_regiao
		ORDER BY percentual DESC
		LIMIT 1
	`

	var result struct {
		RegiaoMaisDoente string  `json:"regiao_mais_doente"`
		PercentualRegiao float64 `json:"percentual_regiao"`
	}

	if err := h.db.Raw(query).Scan(&result).Error; err != nil {
		return regioes, err
	}

	regioes.RegiaoMaisDoente = result.RegiaoMaisDoente
	regioes.PercentualRegiao = result.PercentualRegiao

	return regioes, nil
}

// GetHospitaisFiltrados - Endpoint para hospitais com filtros
func (h *DashboardHandler) GetHospitaisFiltrados(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extrair parâmetros de query
	estado := r.URL.Query().Get("estado")
	cidade := r.URL.Query().Get("cidade")
	especialidade := r.URL.Query().Get("especialidade")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	hospitais, total, err := h.getHospitaisComFiltros(estado, cidade, especialidade, limit, offset)
	if err != nil {
		h.sendErrorResponse(w, "Erro ao buscar hospitais", err)
		return
	}

	response := map[string]interface{}{
		"hospitais": hospitais,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"pages":     int(math.Ceil(float64(total) / float64(limit))),
	}

	h.sendSuccessResponse(w, response)
}

// getHospitaisComFiltros - Busca hospitais com filtros aplicados
func (h *DashboardHandler) getHospitaisComFiltros(estado, cidade, especialidade string, limit, offset int) ([]map[string]interface{}, int, error) {
	var hospitais []map[string]interface{}
	var total int

	// Construir query base
	query := `
		SELECT 
			h.uuid,
			h.nome,
			h.bairro,
			h.leitos_totais,
			h.especialidades,
			m.nome as municipio,
			e.nome as estado,
			e.regiao
		FROM hospitais h
		INNER JOIN municipios m ON m.codigo = h.cod_municipio
		INNER JOIN estados e ON e.codigo = m.codigo_uf
		WHERE 1=1
	`

	args := []interface{}{}

	// Aplicar filtros
	if estado != "" {
		query += " AND e.nome ILIKE ?"
		args = append(args, "%"+estado+"%")
	}

	if cidade != "" {
		query += " AND m.nome ILIKE ?"
		args = append(args, "%"+cidade+"%")
	}

	if especialidade != "" {
		query += " AND h.especialidades::text ILIKE ?"
		args = append(args, "%"+especialidade+"%")
	}

	// Contar total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") as count_query"
	if err := h.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return hospitais, 0, err
	}

	// Adicionar paginação
	query += " ORDER BY h.nome LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Executar query
	if err := h.db.Raw(query, args...).Scan(&hospitais).Error; err != nil {
		return hospitais, 0, err
	}

	return hospitais, total, nil
}

// GetHospitalDetalhes - Endpoint para detalhes de um hospital específico
func (h *DashboardHandler) GetHospitalDetalhes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extrair UUID do hospital da URL
	hospitalUUID := r.URL.Query().Get("uuid")
	if hospitalUUID == "" {
		http.Error(w, "UUID do hospital é obrigatório", http.StatusBadRequest)
		return
	}

	hospital, medicos, err := h.getHospitalDetalhesComMedicos(hospitalUUID)
	if err != nil {
		h.sendErrorResponse(w, "Erro ao buscar detalhes do hospital", err)
		return
	}

	response := map[string]interface{}{
		"hospital": hospital,
		"medicos":  medicos,
	}

	h.sendSuccessResponse(w, response)
}

// getHospitalDetalhesComMedicos - Busca detalhes do hospital e seus médicos
func (h *DashboardHandler) getHospitalDetalhesComMedicos(hospitalUUID string) (map[string]interface{}, []map[string]interface{}, error) {
	// Buscar dados do hospital
	var hospital map[string]interface{}
	hospitalQuery := `
		SELECT 
			h.uuid,
			h.nome,
			h.bairro,
			h.leitos_totais,
			h.especialidades,
			m.nome as municipio,
			e.nome as estado,
			e.regiao,
			m.latitude,
			m.longitude
		FROM hospitais h
		INNER JOIN municipios m ON m.codigo = h.cod_municipio
		INNER JOIN estados e ON e.codigo = m.codigo_uf
		WHERE h.uuid = ?
	`

	if err := h.db.Raw(hospitalQuery, hospitalUUID).Scan(&hospital).Error; err != nil {
		return hospital, nil, err
	}

	// Buscar médicos próximos (dentro de 30km)
	var medicos []map[string]interface{}
	medicosQuery := `
		SELECT 
			m.uuid,
			m.nome,
			m.especialidade,
			mun.nome as municipio,
			mun.latitude,
			mun.longitude,
			ST_Distance(
				ST_Point(mun.longitude::float, mun.latitude::float)::geography,
				ST_Point(?, ?)::geography
			) / 1000 as distancia_km
		FROM medicos m
		INNER JOIN municipios mun ON mun.codigo = m.cod_municipio
		INNER JOIN hospitais h ON h.uuid = ?
		INNER JOIN municipios mh ON mh.codigo = h.cod_municipio
		WHERE ST_DWithin(
			ST_Point(mun.longitude::float, mun.latitude::float)::geography,
			ST_Point(mh.longitude::float, mh.latitude::float)::geography,
			30000
		)
		AND (
			m.especialidade = ANY(string_to_array(h.especialidades::text, ','))
			OR m.especialidade = 'Clínica Geral'
		)
		ORDER BY distancia_km
		LIMIT 50
	`

	// Extrair coordenadas do hospital para a query
	hospitalLat := hospital["latitude"].(string)
	hospitalLon := hospital["longitude"].(string)

	if err := h.db.Raw(medicosQuery, hospitalLon, hospitalLat, hospitalUUID).Scan(&medicos).Error; err != nil {
		return hospital, medicos, err
	}

	return hospital, medicos, nil
}

func (h *DashboardHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *DashboardHandler) sendErrorResponse(w http.ResponseWriter, message string, err error) {
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}
