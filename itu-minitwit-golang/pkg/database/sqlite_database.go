package database

import (
	"fmt"
	"gorm.io/gorm/logger"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"log"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb(cfg *config.Config) {
	if _, err := os.Stat(cfg.DBPath); err != nil {
		fmt.Print("Creating new db file")
		path := strings.Split(cfg.DBPath, "/")
		if len(path) > 0 {
			dirPath := strings.Join(path[:len(path)-1], "/")
			os.MkdirAll(dirPath, 0755)
		}
		os.Create(cfg.DBPath)
	}

	gormConfig := &gorm.Config{}

	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, dbErr := gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
	if dbErr != nil {
		log.Fatalf("Could not connect to database: %v", dbErr)
	}
	db.AutoMigrate(&models.User{}, &models.Message{}, &models.APIUser{})
	DB = db
}
