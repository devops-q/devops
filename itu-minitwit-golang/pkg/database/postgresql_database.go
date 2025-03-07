package database

import (
	"fmt"
	"itu-minitwit/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgreSQLDb(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=gorm port=%d sslmode=disable TimeZone=GMT", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBPort)

	gormConfig := &gorm.Config{}

	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	AutoMigrate(db)
	DB = db
	return DB
}
