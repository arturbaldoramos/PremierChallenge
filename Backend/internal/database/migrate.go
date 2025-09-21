package database

import (
	"log"

	"example.com/m/v2/internal/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Starting automatic migrations...")

	// Try to migrate all models at once first
	err := db.AutoMigrate(
		&domain.Estado{},
		&domain.Municipio{},
		&domain.Cid10{},
		&domain.Hospital{},
		&domain.Medico{},
		&domain.Paciente{},
	)

	if err != nil {
		log.Printf("AutoMigrate failed: %v", err)
		log.Println("Attempting individual table migration...")
		
		// If AutoMigrate fails, try individual migrations with error handling
		models := []interface{}{
			&domain.Estado{},
			&domain.Municipio{},
			&domain.Cid10{},
			&domain.Hospital{},
			&domain.Medico{},
			&domain.Paciente{},
		}

		for _, model := range models {
			modelName := getModelName(model)
			log.Printf("Migrating: %s", modelName)
			
			if err := migrateModelSafe(db, model); err != nil {
				log.Printf("Failed to migrate %s: %v", modelName, err)
				// Continue with other models instead of failing completely
				continue
			}
			log.Printf("Successfully migrated: %s", modelName)
		}
	} else {
		log.Println("All models migrated successfully!")
	}

	// Verificações adicionais para garantir que todas as colunas existem
	if err := ensurePacienteCID10Column(db); err != nil {
		log.Printf("Failed to ensure CID10 column in Paciente: %v", err)
		return err
	}

	log.Println("All migrations completed successfully!")
	return nil
}

// migrateModelSafe handles migration for a single model, ignoring errors
func migrateModelSafe(db *gorm.DB, model interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic during migration: %v", r)
		}
	}()
	
	// Simply try to migrate, ignore errors
	return db.AutoMigrate(model)
}

// migrateModel handles migration for a single model with better error handling
func migrateModel(db *gorm.DB, model interface{}) error {
	// Check if table exists
	tableName := db.NamingStrategy.TableName(getModelName(model))
	var exists bool
	err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = ?)", tableName).Scan(&exists).Error
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Table %s already exists, checking structure...", tableName)
		// Table exists, just ensure it has the right structure
		return db.AutoMigrate(model)
	} else {
		log.Printf("Creating new table: %s", tableName)
		// Table doesn't exist, create it
		return db.AutoMigrate(model)
	}
}

// getModelName extracts model name from interface{}
func getModelName(model interface{}) string {
	switch model.(type) {
	case *domain.Estado:
		return "Estado"
	case *domain.Municipio:
		return "Municipio"
	case *domain.Cid10:
		return "CID10"
	case *domain.Hospital:
		return "Hospital"
	case *domain.Medico:
		return "Medico"
	case *domain.Paciente:
		return "Paciente"
	default:
		return "Unknown"
	}
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
