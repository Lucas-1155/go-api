package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	logger     *Logger
	dbPostgres *gorm.DB
	rdb        *redis.Client
)

func Init() error {
	var err error

	db, err = InitializeOracle()

	if err != nil {
		return fmt.Errorf("ERRO INITIALIZING ORACLE: %v", err)
	}

	dbPostgres, err = InitializePostgres()

	if err != nil {
		fmt.Printf("AVISO: POSTGRES LOGS OFFLINE: %v\n", err)
	}

	rdb, err = InitializeRedis()

	if err != nil {
		fmt.Printf("ERRO INITIALIZING REDIS: %v\n", err)
	}

	return nil
}

func GetOracle() *gorm.DB {
	return db
}

func GetLogger(p string) *Logger {
	logger = NewLogger(p)
	return logger
}

func GetPostgres() *gorm.DB {
	return dbPostgres
}

func GetRedis() *redis.Client {
	return rdb
}
