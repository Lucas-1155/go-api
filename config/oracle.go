package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/dzwvip/oracle"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitializeOracle() (*gorm.DB, error) {

	_, filename, _, _ := runtime.Caller(0)

	rootDir := filepath.Join(filepath.Dir(filename), "..")
	envPath := filepath.Join(rootDir, ".env")

	err := godotenv.Load(envPath)

	if err != nil {
		logger.Error("Erro loading .env file")
	}

	logger := GetLogger("oracle")

	dsn := os.Getenv("ORACLE_URL")

	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})

	if err != nil {
		logger.ErrorF("ORACLE CONNECTION ERROR: %v", err)
		return nil, err
	}

	logger.Info("ORACLE CONNECTION SUCCESSFUL")
	return db, nil
}
