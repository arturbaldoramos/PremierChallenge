package services

import (
	"os"
	"strconv"
	"strings"
	"example.com/m/v2/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"runtime"
	"sync"
)

// getEnvAsInt retrieves an environment variable as an integer, with a default fallback
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

type DataService struct {
	db *gorm.DB
}

type BatchResult struct {
	Inserted int
	Updated  int
	Err      error
}

func NewDataService(db *gorm.DB) *DataService {
	return &DataService{db: db}
}

// Estado operations
func (s *DataService) EstadoExists(codigo string) (bool, error) {
	var count int64
	err := s.db.Model(&domain.Estado{}).Where("codigo = ?", codigo).Count(&count).Error
	return count > 0, err
}

func (s *DataService) GetEstadoByCodigo(codigo string) (*domain.Estado, error) {
	var estado domain.Estado
	err := s.db.Where("codigo = ?", codigo).First(&estado).Error
	if err != nil {
		return nil, err
	}
	return &estado, nil
}

func (s *DataService) BulkInsertEstados(estados []domain.Estado) error {
	return s.db.CreateInBatches(estados, 100).Error
}

func (s *DataService) UpsertEstados(estados []domain.Estado) (int, int, error) {
	if len(estados) == 0 {
		return 0, 0, nil
	}

	// Usar PostgreSQL UPSERT com ON CONFLICT para melhor performance
	query := `
		INSERT INTO estados (codigo, unidade_federativa, nome, regiao, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (codigo) DO UPDATE SET
			unidade_federativa = EXCLUDED.unidade_federativa,
			nome = EXCLUDED.nome,
			regiao = EXCLUDED.regiao,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude
	`

	// Contar registros existentes antes da operação
	var existingCodes []string
	for _, e := range estados {
		existingCodes = append(existingCodes, e.Codigo)
	}

	var existingCount int64
	s.db.Model(&domain.Estado{}).Where("codigo IN ?", existingCodes).Count(&existingCount)

	// Executar upsert em transação
	tx := s.db.Begin()
	for _, estado := range estados {
		result := tx.Exec(query,
			estado.Codigo,
			estado.UnidadeFederativa,
			estado.Nome,
			estado.Regiao,
			estado.Latitude,
			estado.Longitude,
		)
		if result.Error != nil {
			tx.Rollback()
			return 0, 0, result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	// Calcular inserções e atualizações
	inserted := len(estados) - int(existingCount)
	if inserted < 0 {
		inserted = 0
	}
	updated := int(existingCount)

	return inserted, updated, nil
}

// Municipio operations
func (s *DataService) MunicipioExists(codigo string) (bool, error) {
	var count int64
	err := s.db.Model(&domain.Municipio{}).Where("codigo = ?", codigo).Count(&count).Error
	return count > 0, err
}

func (s *DataService) BulkInsertMunicipios(municipios []domain.Municipio) error {
	return s.db.CreateInBatches(municipios, 100).Error
}

func (s *DataService) UpsertMunicipios(municipios []domain.Municipio) (int, int, error) {
	if len(municipios) == 0 {
		return 0, 0, nil
	}

	// Usar PostgreSQL UPSERT com ON CONFLICT para melhor performance
	query := `
		INSERT INTO municipios (codigo, nome, latitude, longitude, capital, codigo_uf, siafi_id, ddd, fuso_hora, populacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (codigo) DO UPDATE SET
			nome = EXCLUDED.nome,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			capital = EXCLUDED.capital,
			codigo_uf = EXCLUDED.codigo_uf,
			siafi_id = EXCLUDED.siafi_id,
			ddd = EXCLUDED.ddd,
			fuso_hora = EXCLUDED.fuso_hora,
			populacao = EXCLUDED.populacao
	`

	var inserted, updated int

	// Processar em batches para evitar sobrecarga de memória
	batchSize := 500
	for i := 0; i < len(municipios); i += batchSize {
		end := i + batchSize
		if end > len(municipios) {
			end = len(municipios)
		}

		batch := municipios[i:end]

		// Contar registros existentes antes da operação
		var existingCodes []string
		for _, m := range batch {
			existingCodes = append(existingCodes, m.Codigo)
		}

		var existingCount int64
		s.db.Model(&domain.Municipio{}).Where("codigo IN ?", existingCodes).Count(&existingCount)

		// Executar upsert em batch
		tx := s.db.Begin()
		for _, municipio := range batch {
			result := tx.Exec(query,
				municipio.Codigo,
				municipio.Nome,
				municipio.Latitude,
				municipio.Longitude,
				municipio.Capital,
				municipio.CodigoUF,
				municipio.SiafiId,
				municipio.DDD,
				municipio.FusoHora,
				municipio.Populacao,
			)
			if result.Error != nil {
				tx.Rollback()
				return inserted, updated, result.Error
			}
		}

		if err := tx.Commit().Error; err != nil {
			return inserted, updated, err
		}

		// Calcular inserções e atualizações aproximadas
		batchInserted := len(batch) - int(existingCount)
		if batchInserted < 0 {
			batchInserted = 0
		}
		batchUpdated := int(existingCount)

		inserted += batchInserted
		updated += batchUpdated
	}

	return inserted, updated, nil
}

// UpsertMunicipiosConcurrent - Versão concorrente otimizada para grandes volumes
func (s *DataService) UpsertMunicipiosConcurrent(municipios []domain.Municipio) (int, int, error) {
	if len(municipios) == 0 {
		return 0, 0, nil
	}

	// Número de workers baseado no número de CPU cores
	numWorkers := runtime.NumCPU()
	if numWorkers > 8 {
		numWorkers = 8 // Limita para não sobrecarregar o banco
	}

	batchSize := 200
	batches := make([][]domain.Municipio, 0)

	// Dividir em batches
	for i := 0; i < len(municipios); i += batchSize {
		end := i + batchSize
		if end > len(municipios) {
			end = len(municipios)
		}
		batches = append(batches, municipios[i:end])
	}

	// Channels para coordenação
	batchChan := make(chan []domain.Municipio, len(batches))
	resultChan := make(chan BatchResult, len(batches))

	// Enviar batches para o channel
	for _, batch := range batches {
		batchChan <- batch
	}
	close(batchChan)

	// Worker function
	worker := func() {
		for batch := range batchChan {
			inserted, updated, err := s.processMunicipioBatch(batch)
			resultChan <- BatchResult{Inserted: inserted, Updated: updated, Err: err}
		}
	}

	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	// Aguardar conclusão
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	var totalInserted, totalUpdated int
	for result := range resultChan {
		if result.Err != nil {
			return totalInserted, totalUpdated, result.Err
		}
		totalInserted += result.Inserted
		totalUpdated += result.Updated
	}

	return totalInserted, totalUpdated, nil
}

func (s *DataService) processMunicipioBatch(batch []domain.Municipio) (int, int, error) {
	query := `
		INSERT INTO municipios (codigo, nome, latitude, longitude, capital, codigo_uf, siafi_id, ddd, fuso_hora, populacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (codigo) DO UPDATE SET
			nome = EXCLUDED.nome,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			capital = EXCLUDED.capital,
			codigo_uf = EXCLUDED.codigo_uf,
			siafi_id = EXCLUDED.siafi_id,
			ddd = EXCLUDED.ddd,
			fuso_hora = EXCLUDED.fuso_hora,
			populacao = EXCLUDED.populacao
	`

	// Contar existentes
	var existingCodes []string
	for _, m := range batch {
		existingCodes = append(existingCodes, m.Codigo)
	}

	var existingCount int64
	if err := s.db.Model(&domain.Municipio{}).Where("codigo IN ?", existingCodes).Count(&existingCount).Error; err != nil {
		return 0, 0, err
	}

	// Executar upsert
	tx := s.db.Begin()
	for _, municipio := range batch {
		if err := tx.Exec(query,
			municipio.Codigo,
			municipio.Nome,
			municipio.Latitude,
			municipio.Longitude,
			municipio.Capital,
			municipio.CodigoUF,
			municipio.SiafiId,
			municipio.DDD,
			municipio.FusoHora,
			municipio.Populacao,
		).Error; err != nil {
			tx.Rollback()
			return 0, 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	inserted := len(batch) - int(existingCount)
	if inserted < 0 {
		inserted = 0
	}
	updated := int(existingCount)

	return inserted, updated, nil
}

// Hospital operations
func (s *DataService) HospitalExists(uuid uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&domain.Hospital{}).Where("uuid = ?", uuid).Count(&count).Error
	return count > 0, err
}

func (s *DataService) BulkInsertHospitais(hospitais []domain.Hospital) error {
	return s.db.CreateInBatches(hospitais, 100).Error
}

func (s *DataService) UpsertHospitais(hospitais []domain.Hospital) (int, int, error) {
	var inserted, updated int

	for _, hospital := range hospitais {
		exists, err := s.HospitalExists(hospital.UUID)
		if err != nil {
			return inserted, updated, err
		}

		if exists {
			err = s.db.Model(&domain.Hospital{}).Where("uuid = ?", hospital.UUID).Updates(hospital).Error
			if err != nil {
				return inserted, updated, err
			}
			updated++
		} else {
			err = s.db.Create(&hospital).Error
			if err != nil {
				return inserted, updated, err
			}
			inserted++
		}
	}

	return inserted, updated, nil
}

// Paciente operations
func (s *DataService) PacienteExists(cpf string) (bool, error) {
	var count int64
	err := s.db.Model(&domain.Paciente{}).Where("cpf = ?", cpf).Count(&count).Error
	return count > 0, err
}

func (s *DataService) BulkInsertPacientes(pacientes []domain.Paciente) error {
	return s.db.CreateInBatches(pacientes, 100).Error
}

func (s *DataService) UpsertPacientes(pacientes []domain.Paciente) (int, int, error) {
	var inserted, updated int

	for _, paciente := range pacientes {
		exists, err := s.PacienteExists(paciente.CPF)
		if err != nil {
			return inserted, updated, err
		}

		if exists {
			err = s.db.Model(&domain.Paciente{}).Where("cpf = ?", paciente.CPF).Updates(paciente).Error
			if err != nil {
				return inserted, updated, err
			}
			updated++
		} else {
			err = s.db.Create(&paciente).Error
			if err != nil {
				return inserted, updated, err
			}
			inserted++
		}
	}

	return inserted, updated, nil
}

// Medico operations
func (s *DataService) MedicoExists(uuid uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&domain.Medico{}).Where("uuid = ?", uuid).Count(&count).Error
	return count > 0, err
}

func (s *DataService) BulkInsertMedicos(medicos []domain.Medico) error {
	return s.db.CreateInBatches(medicos, 100).Error
}

func (s *DataService) UpsertMedicos(medicos []domain.Medico) (int, int, error) {
	if len(medicos) == 0 {
		return 0, 0, nil
	}

	var inserted, updated int

	// Processar em batches para otimizar performance
	batchSize := getEnvAsInt("BATCH_SIZE_MEDICOS", 1000)
	for i := 0; i < len(medicos); i += batchSize {
		end := i + batchSize
		if end > len(medicos) {
			end = len(medicos)
		}

		batch := medicos[i:end]
		batchInserted, batchUpdated, err := s.processMedicoBulk(batch)
		if err != nil {
			return inserted, updated, err
		}

		inserted += batchInserted
		updated += batchUpdated
	}

	return inserted, updated, nil
}

// processMedicoBulk executa bulk insert para um batch de médicos
func (s *DataService) processMedicoBulk(batch []domain.Medico) (int, int, error) {
	if len(batch) == 0 {
		return 0, 0, nil
	}

	// Contar registros existentes antes da operação
	var existingUUIDs []string
	for _, m := range batch {
		existingUUIDs = append(existingUUIDs, m.UUID.String())
	}

	var existingCount int64
	s.db.Model(&domain.Medico{}).Where("uuid::text IN ?", existingUUIDs).Count(&existingCount)

	// Construir query de bulk insert
	placeholders := make([]string, len(batch))
	values := make([]interface{}, 0, len(batch)*4)

	for i, medico := range batch {
		placeholders[i] = "(?, ?, ?, ?)"
		values = append(values, medico.UUID, medico.Nome, medico.Especialidade, medico.CodMunicipio)
	}

	query := `
		INSERT INTO medicos (uuid, nome, especialidade, cod_municipio)
		VALUES ` + strings.Join(placeholders, ", ") + `
		ON CONFLICT (uuid) DO UPDATE SET
			nome = EXCLUDED.nome,
			especialidade = EXCLUDED.especialidade,
			cod_municipio = EXCLUDED.cod_municipio
	`

	// Executar bulk insert
	result := s.db.Exec(query, values...)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	// Calcular inserções e atualizações aproximadas
	batchInserted := len(batch) - int(existingCount)
	if batchInserted < 0 {
		batchInserted = 0
	}
	batchUpdated := int(existingCount)

	return batchInserted, batchUpdated, nil
}

// UpsertMedicosConcurrent - Versão concorrente otimizada para grandes volumes
func (s *DataService) UpsertMedicosConcurrent(medicos []domain.Medico) (int, int, error) {
	if len(medicos) == 0 {
		return 0, 0, nil
	}

	// Número de workers baseado no número de CPU cores
	numWorkers := runtime.NumCPU()
	if numWorkers > 6 {
		numWorkers = 6 // Limita para não sobrecarregar o banco
	}

	batchSize := getEnvAsInt("BATCH_SIZE_MEDICOS", 300)
	batches := make([][]domain.Medico, 0)

	// Dividir em batches
	for i := 0; i < len(medicos); i += batchSize {
		end := i + batchSize
		if end > len(medicos) {
			end = len(medicos)
		}
		batches = append(batches, medicos[i:end])
	}

	// Channels para coordenação
	batchChan := make(chan []domain.Medico, len(batches))
	resultChan := make(chan BatchResult, len(batches))

	// Enviar batches para o channel
	for _, batch := range batches {
		batchChan <- batch
	}
	close(batchChan)

	// Worker function
	worker := func() {
		for batch := range batchChan {
			inserted, updated, err := s.processMedicoBatch(batch)
			resultChan <- BatchResult{Inserted: inserted, Updated: updated, Err: err}
		}
	}

	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	// Aguardar conclusão
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	var totalInserted, totalUpdated int
	for result := range resultChan {
		if result.Err != nil {
			return totalInserted, totalUpdated, result.Err
		}
		totalInserted += result.Inserted
		totalUpdated += result.Updated
	}

	return totalInserted, totalUpdated, nil
}

func (s *DataService) processMedicoBatch(batch []domain.Medico) (int, int, error) {
	query := `
		INSERT INTO medicos (uuid, nome, especialidade, cod_municipio)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (uuid) DO UPDATE SET
			nome = EXCLUDED.nome,
			especialidade = EXCLUDED.especialidade,
			cod_municipio = EXCLUDED.cod_municipio
	`

	// Contar existentes
	var existingUUIDs []string
	for _, m := range batch {
		existingUUIDs = append(existingUUIDs, m.UUID.String())
	}

	var existingCount int64
	if err := s.db.Model(&domain.Medico{}).Where("uuid::text IN ?", existingUUIDs).Count(&existingCount).Error; err != nil {
		return 0, 0, err
	}

	// Executar upsert
	tx := s.db.Begin()
	for _, medico := range batch {
		if err := tx.Exec(query,
			medico.UUID,
			medico.Nome,
			medico.Especialidade,
			medico.CodMunicipio,
		).Error; err != nil {
			tx.Rollback()
			return 0, 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	inserted := len(batch) - int(existingCount)
	if inserted < 0 {
		inserted = 0
	}
	updated := int(existingCount)

	return inserted, updated, nil
}

// CID10 operations
func (s *DataService) CID10Exists(codigo string) (bool, error) {
	var count int64
	err := s.db.Model(&domain.Cid10{}).Where("codigo = ?", codigo).Count(&count).Error
	return count > 0, err
}

func (s *DataService) BulkInsertCID10(cid10s []domain.Cid10) error {
	return s.db.CreateInBatches(cid10s, 100).Error
}

func (s *DataService) UpsertCID10(cid10s []domain.Cid10) (int, int, error) {
	var inserted, updated int

	for _, cid10 := range cid10s {
		exists, err := s.CID10Exists(cid10.Codigo)
		if err != nil {
			return inserted, updated, err
		}

		if exists {
			err = s.db.Model(&domain.Cid10{}).Where("codigo = ?", cid10.Codigo).Updates(cid10).Error
			if err != nil {
				return inserted, updated, err
			}
			updated++
		} else {
			err = s.db.Create(&cid10).Error
			if err != nil {
				return inserted, updated, err
			}
			inserted++
		}
	}

	return inserted, updated, nil
}