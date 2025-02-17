package utils

import (
	"itu-minitwit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Retrieves the UserID. Returns a tuple of (int: value, error: err)
func GetUserIDFromName(c *gin.Context, username string) (int, error) {
	db := c.MustGet("DB").(*gorm.DB)
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}
	return int(user.ID), nil
}
