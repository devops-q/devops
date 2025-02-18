package utils

import (
	"itu-minitwit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindUserWithName(c *gin.Context, name string) (*models.User, error) {
	db := c.MustGet("DB").(*gorm.DB)

	var user models.User
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		return nil, err // Return error if user is not found
	}
	return &user, nil
}
