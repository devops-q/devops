package database

import (
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"itu-minitwit/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	gormLogger := NewGormLogger()

	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if dbErr != nil {
		gormLogger.logger.Fatalf("Could not connect to database: %v", dbErr)
	}
	err := db.AutoMigrate(&models.User{}, &models.Message{}, &models.APIUser{}, &models.LatestID{})
	if err != nil {
		gormLogger.logger.Error("Error during auto migration: %v", err)
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
		logger.GetLogger().WithService("gorm").Info("Initial API user created successfully", "apiUser", cfg.InitialApiUser)
	}

	return nil
}
