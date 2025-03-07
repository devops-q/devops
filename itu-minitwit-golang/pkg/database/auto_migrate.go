package database

import (
	"itu-minitwit/internal/models"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Message{}, &models.APIUser{})
	if err != nil {
		log.Printf("Error during auto migration: %v", err)
		panic(err)
	}
}