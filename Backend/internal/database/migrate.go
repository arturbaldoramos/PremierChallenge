package database

import (
	"log"

	"example.com/m/v2/internal/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Starting automatic migrations...")

	// Migrate models one by one with better error handling
	if err := db.AutoMigrate(&domain.Estado{}); err != nil {
		log.Printf("Failed to migrate Estado: %v", err)
		return err
	}
	log.Println("Successfully migrated: Estado")

	if err := db.AutoMigrate(&domain.Municipio{}); err != nil {
		log.Printf("Failed to migrate Municipio: %v", err)
		return err
	}
	log.Println("Successfully migrated: Municipio")

	if err := db.AutoMigrate(&domain.Cid10{}); err != nil {
		log.Printf("Failed to migrate CID10: %v", err)
		return err
	}
	log.Println("Successfully migrated: CID10")

	if err := db.AutoMigrate(&domain.Hospital{}); err != nil {
		log.Printf("Failed to migrate Hospital: %v", err)
		return err
	}
	log.Println("Successfully migrated: Hospital")

	if err := db.AutoMigrate(&domain.Medico{}); err != nil {
		log.Printf("Failed to migrate Medico: %v", err)
		return err
	}
	log.Println("Successfully migrated: Medico")

	if err := db.AutoMigrate(&domain.Paciente{}); err != nil {
		log.Printf("Failed to migrate Paciente: %v", err)
		return err
	}
	log.Println("Successfully migrated: Paciente")

	// Verificações adicionais para garantir que todas as colunas existem
	if err := ensurePacienteCID10Column(db); err != nil {
		log.Printf("Failed to ensure CID10 column in Paciente: %v", err)
		return err
	}

	log.Println("All migrations completed successfully!")
	return nil
}

// ensurePacienteCID10Column garante que a coluna CID10 existe na tabela pacientes e remove duplicatas
func ensurePacienteCID10Column(db *gorm.DB) error {
	log.Println("Ensuring CID10 column exists in pacientes table...")

	// Verificar se a coluna c_id10 duplicada existe
	var hasOldColumn bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'pacientes' AND column_name = 'c_id10')").Scan(&hasOldColumn).Error
	if err != nil {
		return err
	}

	if hasOldColumn {
		log.Println("Found duplicate column 'c_id10', removing it...")
		err = db.Exec("ALTER TABLE pacientes DROP COLUMN IF EXISTS c_id10").Error
		if err != nil {
			log.Printf("Warning: Failed to drop duplicate column c_id10: %v", err)
		} else {
			log.Println("Successfully removed duplicate column c_id10")
		}
	}

	// Verificar se a coluna cid10 correta existe
	var hasColumn bool
	err = db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'pacientes' AND column_name = 'cid10')").Scan(&hasColumn).Error
	if err != nil {
		return err
	}

	if !hasColumn {
		log.Println("CID10 column not found, adding it...")
		// Adicionar a coluna CID10
		err = db.Exec("ALTER TABLE pacientes ADD COLUMN cid10 VARCHAR(10)").Error
		if err != nil {
			return err
		}
		log.Println("Successfully added CID10 column to pacientes table")
	} else {
		log.Println("CID10 column already exists in pacientes table")
	}

	return nil
}
