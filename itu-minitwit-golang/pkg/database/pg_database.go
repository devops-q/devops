package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"log"
)

var DB *gorm.DB

func InitDb(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	gormConfig := &gorm.Config{}

	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if dbErr != nil {
		log.Fatalf("Could not connect to database: %v", dbErr)
	}
	err := db.AutoMigrate(&models.User{}, &models.Message{}, &models.APIUser{})
	if err != nil {
		log.Printf("Error during auto migration: %v", err)
		panic(err)
	}
	DB = db
}
