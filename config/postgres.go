package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializePostgres() (*gorm.DB, error) {
	logger := GetLogger("postgres")
	dsn := os.Getenv("POSTGRES_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.ErrorF("POSTGRES CONNECTION ERROR: %v", err)
		return nil, err
	}

	logger.Info("POSTGRES CONNECTION SUCCESSFUL")
	return db, nil
}
