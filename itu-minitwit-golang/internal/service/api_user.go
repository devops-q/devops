package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
	"itu-minitwit/pkg/logger"
)

func GetApiUsers(db *gorm.DB) (gin.Accounts, error) {
	log := logger.Init()
	var users []models.APIUser

	if err := db.Find(&users).Error; err != nil {
		log.Error("[GetApiUsers] Error: %v", err)
		return nil, err
	}

	accounts := gin.Accounts{}

	for _, user := range users {
		accounts[user.Username] = user.Password
	}
	log.Info("[GetApiUsers] success")
	return accounts, nil
}

func CreateApiUser(db *gorm.DB, username, password string) (bool, error) {
	log := logger.Init()
	// Check if the user already exists
	var existingUser models.APIUser
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		log.Error("[CreateApiUser] Username already exists")
		// User already exists
		return false, nil
	}

	// If no existing user, create a new one
	if err := db.Create(&models.APIUser{Username: username, Password: password}).Error; err != nil {
		log.Error("[CreateApiUser] Error: %v", err)
		return false, err
	}

	// Successfully created the user
	log.Info("[CreateApiUser] success")
	return true, nil
}
