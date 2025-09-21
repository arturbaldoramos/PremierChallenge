package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"example.com/m/v2/internal/domain"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

type StatsResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type TotalStatsResponse struct {
	TotalHospitais int64 `json:"total_hospitais"`
	TotalMedicos   int64 `json:"total_medicos"`
	TotalLeitos    int64 `json:"total_leitos"`
	TotalPacientes int64 `json:"total_pacientes"`
}

type Stats2Response struct {
	TotalEstados    int64 `json:"total_estados"`
	TotalMunicipios int64 `json:"total_municipios"`
}

type Cid10Response struct {
	Cid10  string `json:"cid10"`
	Doenca string `json:"doenca"`
	Count  int64  `json:"count"`
}

type HospitalAccessResponse struct {
	HospitalUUID   string `json:"hospital_uuid"`
	HospitalNome   string `json:"hospital_nome"`
	Especialidades string `json:"especialidades"`
	Count          int64  `json:"count"`
}

type MedicosDistributionResponse struct {
	ComUmHospital    int64 `json:"com_um_hospital"`
	ComDoisHospitais int64 `json:"com_dois_hospitais"`
	ComTresHospitais int64 `json:"com_tres_hospitais"`
	SemHospitais     int64 `json:"sem_hospitais"`
}

type HospitalsBySpecialtyResponse struct {
	Especialidade string `json:"especialidade"`
	Count         int64  `json:"count"`
}

type HospitalsByMunicipioResponse struct {
	Municipio string `json:"municipio"`
	Count     int64  `json:"count"`
}

type MedicoHospitalAssignment struct {
	MedicoUUID    string  `json:"medico_uuid"`
	MedicoNome    string  `json:"medico_nome"`
	HospitalUUID  string  `json:"hospital_uuid"`
	HospitalNome  string  `json:"hospital_nome"`
	Especialidade string  `json:"especialidade"`
	Distancia     float64 `json:"distancia_km"`
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

// Função auxiliar para calcular distância entre dois pontos (fórmula de Haversine)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Raio da Terra em km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c

	return distance
}

// GET /api/v1/stats/totals - Estatísticas gerais
func (h *StatsHandler) GetTotalStats(w http.ResponseWriter, r *http.Request) {
	var stats TotalStatsResponse

	// Contar hospitais
	if err := h.db.Model(&domain.Hospital{}).Count(&stats.TotalHospitais).Error; err != nil {
		http.Error(w, "Erro ao contar hospitais", http.StatusInternalServerError)
		return
	}

	// Contar médicos
	if err := h.db.Model(&domain.Medico{}).Count(&stats.TotalMedicos).Error; err != nil {
		http.Error(w, "Erro ao contar médicos", http.StatusInternalServerError)
		return
	}

	// Contar total de leitos
	if err := h.db.Model(&domain.Hospital{}).Select("COALESCE(SUM(leitos_totais), 0)").Scan(&stats.TotalLeitos).Error; err != nil {
		http.Error(w, "Erro ao contar leitos", http.StatusInternalServerError)
		return
	}

	// Contar pacientes
	if err := h.db.Model(&domain.Paciente{}).Count(&stats.TotalPacientes).Error; err != nil {
		http.Error(w, "Erro ao contar pacientes", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats2 - Estatísticas de estados e municípios
func (h *StatsHandler) GetStats2(w http.ResponseWriter, r *http.Request) {
	var stats Stats2Response

	// Contar estados
	if err := h.db.Model(&domain.Estado{}).Count(&stats.TotalEstados).Error; err != nil {
		http.Error(w, "Erro ao contar estados", http.StatusInternalServerError)
		return
	}

	// Contar municípios
	if err := h.db.Model(&domain.Municipio{}).Count(&stats.TotalMunicipios).Error; err != nil {
		http.Error(w, "Erro ao contar municípios", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/cid10-mais-comuns - Doenças mais comuns (CID10)
func (h *StatsHandler) GetCid10MaisComuns(w http.ResponseWriter, r *http.Request) {
	var results []Cid10Response

	query := `
		SELECT 
			p.cid10,
			COALESCE(c.descricao, 'CID10 não encontrado') as doenca,
			COUNT(p.cid10) as count
		FROM pacientes p
		LEFT JOIN cid10 c ON p.cid10 = c.codigo
		WHERE p.cid10 IS NOT NULL AND p.cid10 != ''
		GROUP BY p.cid10, c.descricao
		ORDER BY count DESC
		LIMIT 5
	`

	if err := h.db.Raw(query).Scan(&results).Error; err != nil {
		http.Error(w, "Erro ao buscar doenças mais comuns", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    results,
		Message: "Top 5 doenças mais comuns (CID10) dos pacientes",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/hospitais-mais-acessados - Hospitais mais acessados pelos pacientes
func (h *StatsHandler) GetHospitaisMaisAcessados(w http.ResponseWriter, r *http.Request) {
	var results []HospitalAccessResponse

	query := `
		SELECT 
			h.uuid as hospital_uuid,
			h.nome as hospital_nome,
			h.especialidades as especialidades,
			COUNT(*) as count
		FROM hospitals h
		JOIN medico_hospital mh ON h.uuid = mh.hospital_uuid
		JOIN medicos m ON mh.medico_uuid = m.uuid
		JOIN pacientes p ON p.cod_municipio = m.cod_municipio
		GROUP BY h.uuid, h.nome, h.especialidades
		ORDER BY count DESC
		LIMIT 20
	`

	if err := h.db.Raw(query).Scan(&results).Error; err != nil {
		http.Error(w, "Erro ao buscar hospitais mais acessados", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    results,
		Message: "Top 20 hospitais mais acessados pelos pacientes",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/debug-cid10 - Debug: Contar doenças por CID10
func (h *StatsHandler) GetDebugCid10(w http.ResponseWriter, r *http.Request) {
	var results []struct {
		Cid10      string `json:"cid10"`
		Quantidade int64  `json:"quantidade"`
	}

	// Query básica para contar por CID10
	query := `
		SELECT cid10, COUNT(*) as quantidade 
		FROM pacientes 
		WHERE cid10 IS NOT NULL AND cid10 != ''
		GROUP BY cid10 
		ORDER BY quantidade DESC
		LIMIT 20
	`

	if err := h.db.Raw(query).Scan(&results).Error; err != nil {
		http.Error(w, "Erro ao executar query de debug", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    results,
		Message: "Debug: Contagem de doenças por CID10 (top 20)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/hospitais-por-especialidade?especialidade=X - Hospitais filtrados por especialidade
func (h *StatsHandler) GetHospitalsBySpecialty(w http.ResponseWriter, r *http.Request) {
	especialidade := r.URL.Query().Get("especialidade")

	if especialidade == "" {
		// Se não especificou, retorna todas as especialidades
		var results []HospitalsBySpecialtyResponse

		query := `
			SELECT 
				TRIM(UNNEST(string_to_array(especialidades, ','))) as especialidade,
				COUNT(*) as count
			FROM hospitals 
			WHERE especialidades IS NOT NULL AND especialidades != ''
			GROUP BY TRIM(UNNEST(string_to_array(especialidades, ',')))
			ORDER BY count DESC
		`

		if err := h.db.Raw(query).Scan(&results).Error; err != nil {
			http.Error(w, "Erro ao buscar hospitais por especialidade", http.StatusInternalServerError)
			return
		}

		response := StatsResponse{
			Success: true,
			Data:    results,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filtrar hospitais que atendem a especialidade específica
	var hospitais []domain.Hospital

	query := h.db.Where("especialidades ILIKE ?", "%"+especialidade+"%")

	if err := query.Find(&hospitais).Error; err != nil {
		http.Error(w, "Erro ao buscar hospitais", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    hospitais,
		Message: "Hospitais que atendem a especialidade: " + especialidade,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/hospitais-por-municipio?municipio=X - Hospitais filtrados por município
func (h *StatsHandler) GetHospitalsByMunicipio(w http.ResponseWriter, r *http.Request) {
	municipio := r.URL.Query().Get("municipio")

	if municipio == "" {
		// Se não especificou, retorna todos os municípios com seus hospitais
		var results []HospitalsByMunicipioResponse

		query := `
			SELECT 
				m.nome as municipio,
				COUNT(h.id) as count
			FROM municipios m
			LEFT JOIN hospitals h ON m.codigo = h.cod_municipio
			GROUP BY m.nome
			HAVING COUNT(h.id) > 0
			ORDER BY count DESC
		`

		if err := h.db.Raw(query).Scan(&results).Error; err != nil {
			http.Error(w, "Erro ao buscar hospitais por município", http.StatusInternalServerError)
			return
		}

		response := StatsResponse{
			Success: true,
			Data:    results,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filtrar hospitais do município específico
	var hospitais []domain.Hospital

	query := `
		SELECT h.* FROM hospitals h
		JOIN municipios m ON h.cod_municipio = m.codigo
		WHERE m.nome ILIKE ?
	`

	if err := h.db.Raw(query, "%"+municipio+"%").Scan(&hospitais).Error; err != nil {
		http.Error(w, "Erro ao buscar hospitais", http.StatusInternalServerError)
		return
	}

	response := StatsResponse{
		Success: true,
		Data:    hospitais,
		Message: "Hospitais do município: " + municipio,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/medicos-distribuicao - Distribuição de médicos por quantidade de hospitais
func (h *StatsHandler) GetMedicosDistribution(w http.ResponseWriter, r *http.Request) {
	// Buscar todos os médicos com suas coordenadas (via município)
	var medicosData []struct {
		UUID          string `json:"uuid"`
		Nome          string `json:"nome"`
		Especialidade string `json:"especialidade"`
		CodMunicipio  string `json:"cod_municipio"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
	}

	queryMedicos := `
		SELECT m.uuid, m.nome, m.especialidade, m.cod_municipio, mu.latitude, mu.longitude
		FROM medicos m
		JOIN municipios mu ON m.cod_municipio = mu.codigo
		WHERE mu.latitude IS NOT NULL AND mu.longitude IS NOT NULL
		AND mu.latitude != '' AND mu.longitude != ''
	`

	if err := h.db.Raw(queryMedicos).Scan(&medicosData).Error; err != nil {
		http.Error(w, "Erro ao buscar médicos com coordenadas", http.StatusInternalServerError)
		return
	}

	// Buscar todos os hospitais com suas coordenadas e especialidades
	var hospitaisData []struct {
		UUID           string `json:"uuid"`
		Nome           string `json:"nome"`
		Especialidades string `json:"especialidades"`
		CodMunicipio   string `json:"cod_municipio"`
		Latitude       string `json:"latitude"`
		Longitude      string `json:"longitude"`
	}

	queryHospitais := `
		SELECT h.uuid, h.nome, h.especialidades, h.cod_municipio, m.latitude, m.longitude
		FROM hospitals h
		JOIN municipios m ON h.cod_municipio = m.codigo
		WHERE m.latitude IS NOT NULL AND m.longitude IS NOT NULL
		AND m.latitude != '' AND m.longitude != ''
		AND h.especialidades IS NOT NULL AND h.especialidades != ''
	`

	if err := h.db.Raw(queryHospitais).Scan(&hospitaisData).Error; err != nil {
		http.Error(w, "Erro ao buscar hospitais com coordenadas", http.StatusInternalServerError)
		return
	}

	// Contar distribuição real baseada na lógica de negócio
	var comUmHospital, comDoisHospitais, comTresHospitais, semHospitais int64

	for _, medico := range medicosData {
		medicoLat, err1 := strconv.ParseFloat(medico.Latitude, 64)
		medicoLon, err2 := strconv.ParseFloat(medico.Longitude, 64)
		if err1 != nil || err2 != nil {
			semHospitais++
			continue
		}

		var hospitaisCompativeis int

		// Verificar cada hospital
		for _, hospital := range hospitaisData {
			// Verificar se o hospital atende a especialidade do médico
			especialidades := strings.Split(hospital.Especialidades, ",")
			especialidadeCompativel := false
			for _, esp := range especialidades {
				if strings.TrimSpace(strings.ToLower(esp)) == strings.TrimSpace(strings.ToLower(medico.Especialidade)) {
					especialidadeCompativel = true
					break
				}
			}

			if !especialidadeCompativel {
				continue
			}

			// Calcular distância
			hospitalLat, err1 := strconv.ParseFloat(hospital.Latitude, 64)
			hospitalLon, err2 := strconv.ParseFloat(hospital.Longitude, 64)
			if err1 != nil || err2 != nil {
				continue
			}

			distancia := calculateDistance(medicoLat, medicoLon, hospitalLat, hospitalLon)

			// Verificar se está dentro do raio de 30km
			if distancia <= 30.0 {
				hospitaisCompativeis++
				// Máximo de 3 hospitais por médico
				if hospitaisCompativeis >= 3 {
					break
				}
			}
		}

		// Classificar médico baseado na quantidade de hospitais compatíveis
		switch hospitaisCompativeis {
		case 0:
			semHospitais++
		case 1:
			comUmHospital++
		case 2:
			comDoisHospitais++
		default: // 3 ou mais (limitado a 3)
			comTresHospitais++
		}
	}

	distribution := MedicosDistributionResponse{
		ComUmHospital:    comUmHospital,
		ComDoisHospitais: comDoisHospitais,
		ComTresHospitais: comTresHospitais,
		SemHospitais:     semHospitais,
	}

	response := StatsResponse{
		Success: true,
		Data:    distribution,
		Message: "Distribuição real baseada na lógica de negócio (raio 30km, mesma especialidade, máx 3 hospitais)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/v1/stats/assign-medicos-hospitais - Implementa a lógica de negócio para atribuir médicos a hospitais
func (h *StatsHandler) AssignMedicosToHospitals(w http.ResponseWriter, r *http.Request) {
	// Buscar todos os médicos
	var medicos []domain.Medico
	if err := h.db.Find(&medicos).Error; err != nil {
		http.Error(w, "Erro ao buscar médicos", http.StatusInternalServerError)
		return
	}

	// Buscar todos os hospitais com suas coordenadas (via município)
	query := `
		SELECT h.*, m.latitude, m.longitude, m.nome as municipio_nome
		FROM hospitals h
		JOIN municipios m ON h.cod_municipio = m.codigo
		WHERE m.latitude IS NOT NULL AND m.longitude IS NOT NULL
		AND m.latitude != '' AND m.longitude != ''
	`

	rows, err := h.db.Raw(query).Rows()
	if err != nil {
		http.Error(w, "Erro ao buscar hospitais", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var hospitais []struct {
		domain.Hospital
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	for rows.Next() {
		var hospital struct {
			domain.Hospital
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		}
		if err := h.db.ScanRows(rows, &hospital); err != nil {
			continue
		}
		hospitais = append(hospitais, hospital)
	}

	var assignments []MedicoHospitalAssignment

	// Para cada médico, encontrar hospitais compatíveis
	for _, medico := range medicos {
		// Buscar coordenadas do município do médico
		var medicoMunicipio domain.Municipio
		if err := h.db.Where("codigo = ?", medico.CodMunicipio).First(&medicoMunicipio).Error; err != nil {
			continue
		}

		if medicoMunicipio.Latitude == "" || medicoMunicipio.Longitude == "" {
			continue
		}

		medicoLat, err1 := strconv.ParseFloat(medicoMunicipio.Latitude, 64)
		medicoLon, err2 := strconv.ParseFloat(medicoMunicipio.Longitude, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		var hospitaisCompativeis []struct {
			Hospital  domain.Hospital
			Distancia float64
		}

		// Verificar cada hospital
		for _, hospital := range hospitais {
			// Verificar se o hospital atende a especialidade do médico
			especialidades := strings.Split(hospital.Especialidades, ",")
			especialidadeCompativel := false
			for _, esp := range especialidades {
				if strings.TrimSpace(strings.ToLower(esp)) == strings.TrimSpace(strings.ToLower(medico.Especialidade)) {
					especialidadeCompativel = true
					break
				}
			}

			if !especialidadeCompativel {
				continue
			}

			// Calcular distância
			hospitalLat, err1 := strconv.ParseFloat(hospital.Latitude, 64)
			hospitalLon, err2 := strconv.ParseFloat(hospital.Longitude, 64)
			if err1 != nil || err2 != nil {
				continue
			}

			distancia := calculateDistance(medicoLat, medicoLon, hospitalLat, hospitalLon)

			// Verificar se está dentro do raio de 30km
			if distancia <= 30.0 {
				hospitaisCompativeis = append(hospitaisCompativeis, struct {
					Hospital  domain.Hospital
					Distancia float64
				}{
					Hospital:  hospital.Hospital,
					Distancia: distancia,
				})
			}
		}

		// Ordenar por distância e pegar até 3 hospitais
		// (implementação simplificada - na prática, usar sort.Slice)
		maxHospitais := 3
		if len(hospitaisCompativeis) > maxHospitais {
			hospitaisCompativeis = hospitaisCompativeis[:maxHospitais]
		}

		// Criar assignments
		for _, hc := range hospitaisCompativeis {
			assignments = append(assignments, MedicoHospitalAssignment{
				MedicoUUID:    medico.UUID.String(),
				MedicoNome:    medico.Nome,
				HospitalUUID:  hc.Hospital.UUID.String(),
				HospitalNome:  hc.Hospital.Nome,
				Especialidade: medico.Especialidade,
				Distancia:     hc.Distancia,
			})
		}
	}

	response := StatsResponse{
		Success: true,
		Data:    assignments,
		Message: "Atribuições médico-hospital baseadas na lógica de negócio (raio 30km, mesma especialidade, máx 3 hospitais)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/stats/especialidades - Lista todas as especialidades disponíveis
func (h *StatsHandler) GetEspecialidades(w http.ResponseWriter, r *http.Request) {
	var especialidades []string

	// Buscar especialidades dos médicos
	query := `
		SELECT DISTINCT especialidade 
		FROM medicos 
		WHERE especialidade IS NOT NULL AND especialidade != ''
		ORDER BY especialidade
	`

	rows, err := h.db.Raw(query).Rows()
	if err != nil {
		http.Error(w, "Erro ao buscar especialidades", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var especialidade string
		if err := rows.Scan(&especialidade); err != nil {
			continue
		}
		especialidades = append(especialidades, especialidade)
	}

	response := StatsResponse{
		Success: true,
		Data:    especialidades,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
