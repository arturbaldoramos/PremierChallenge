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

	log.Println("All migrations completed successfully!")
	return nil
}
