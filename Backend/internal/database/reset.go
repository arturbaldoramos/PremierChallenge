package database

import (
	"log"

	"example.com/m/v2/internal/domain"
	"gorm.io/gorm"
)

func ResetDatabase(db *gorm.DB) error {
	log.Println("Resetting database...")

	// Drop tables if they exist
	tables := []interface{}{
		&domain.Paciente{},
		&domain.Medico{},
		&domain.Hospital{},
		&domain.Cid10{},
		&domain.Municipio{},
		&domain.Estado{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %T: %v", table, err)
		}
	}

	log.Println("Database reset completed!")
	return nil
}