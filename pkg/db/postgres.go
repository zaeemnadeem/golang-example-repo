package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect connects to the Postgres database using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Basic logging, can integrate with zap later
	})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Ping the DB (in GORM this often means getting the underlying *sql.DB and checking it)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
