package main

import (
	"flag"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"itu-minitwit/pkg/database"
	"log"

	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	oldDbPath := flag.String("path", "", "Old DB path")

	var newDB *gorm.DB
	if cfg.DBEngine == "sqlite" {
		newDB = database.InitSQLiteDb(cfg)
	} else if cfg.DBEngine == "postgresql" {
		newDB = database.InitPostgreSQLDb(cfg)
	}

	cfg.DBPath = *oldDbPath
	oldDB := database.InitSQLiteDb(cfg)

	var users []*models.User
	oldDB.Preload("User.Following").Find(&users)
	newDB.CreateInBatches(users, 1000)

	var messages []*models.Message
	oldDB.Find(&messages)
	newDB.CreateInBatches(messages, 1000)
}
