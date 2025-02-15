package main

import (
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/internal/api"
	"itu-minitwit/internal/models"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	initDb(cfg)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	api.SetupRoutes(r, cfg)

	log.Printf("Server starting on port %d", cfg.Port)
	r.Run(fmt.Sprintf("localhost:%d", cfg.Port))
}

func initDb(cfg *config.Config) {
	if _, err := os.Stat(cfg.DBPath); err != nil {
		fmt.Print("Creating new db file")
		path := strings.Split(cfg.DBPath, "/")
		if len(path) > 0 {
			dirPath := strings.Join(path[:len(path)-1], "/")
			os.MkdirAll(dirPath, 0755)
		}
		os.Create(cfg.DBPath)
	}

	db, dbErr := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("Could not connect to database: %v", dbErr)
	}
	db.AutoMigrate(&models.User{}, &models.Message{})
}
