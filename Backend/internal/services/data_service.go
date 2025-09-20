package services

import (
	"example.com/m/v2/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DataService struct {
	db *gorm.DB
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
	var inserted, updated int

	for _, estado := range estados {
		exists, err := s.EstadoExists(estado.Codigo)
		if err != nil {
			return inserted, updated, err
		}

		if exists {
			err = s.db.Model(&domain.Estado{}).Where("codigo = ?", estado.Codigo).Updates(estado).Error
			if err != nil {
				return inserted, updated, err
			}
			updated++
		} else {
			err = s.db.Create(&estado).Error
			if err != nil {
				return inserted, updated, err
			}
			inserted++
		}
	}

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
	var inserted, updated int

	for _, municipio := range municipios {
		exists, err := s.MunicipioExists(municipio.Codigo)
		if err != nil {
			return inserted, updated, err
		}

		if exists {
			err = s.db.Model(&domain.Municipio{}).Where("codigo = ?", municipio.Codigo).Updates(municipio).Error
			if err != nil {
				return inserted, updated, err
			}
			updated++
		} else {
			err = s.db.Create(&municipio).Error
			if err != nil {
				return inserted, updated, err
			}
			inserted++
		}
	}

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
	var inserted, updated int

	for _, medico := range medicos {
		exists, err := s.MedicoExists(medico.UUID)
		if err != nil {
			return inserted, updated, err
		}

		if exists {
			err = s.db.Model(&domain.Medico{}).Where("uuid = ?", medico.UUID).Updates(medico).Error
			if err != nil {
				return inserted, updated, err
			}
			updated++
		} else {
			err = s.db.Create(&medico).Error
			if err != nil {
				return inserted, updated, err
			}
			inserted++
		}
	}

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