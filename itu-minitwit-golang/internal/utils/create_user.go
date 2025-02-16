package utils

import (
	"itu-minitwit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func CreateUser(c *gin.Context, username string, email string, password string)  {
	db := c.MustGet("DB").(*gorm.DB)


	var newUser = &models.User{
		Username: username,
		Email:    email,
		PwHash:   password,
	}


	db.Create(newUser)
	return
}