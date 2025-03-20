package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"log"
)

var DB *gorm.DB

func InitDb(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

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

func InitApiUserIfNotExists(cfg *config.Config) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if cfg.InitialApiUser == "" || cfg.InitialApiPassword == "" {
		return nil
	}

	success, err := service.CreateApiUser(DB, cfg.InitialApiUser, cfg.InitialApiPassword)

	if err != nil {
		return fmt.Errorf("failed to create initial API user: %v", err)
	}

	if success {
		log.Printf("Initial API user %s created successfully", cfg.InitialApiUser)
	}

	return nil
}
