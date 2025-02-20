package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
)

func GetApiUsers(db *gorm.DB) (gin.Accounts, error) {
	var users []models.APIUser

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	accounts := gin.Accounts{}

	for _, user := range users {
		accounts[user.Username] = user.Password
	}

	return accounts, nil
}

func CreateApiUser(db *gorm.DB, username, password string) (bool, error) {
	if err := db.Create(&models.APIUser{Username: username, Password: password}).Error; err != nil {
		return false, err
	}
	return true, nil
}
