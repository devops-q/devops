package utils

import (
	"itu-minitwit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func UserExists(c *gin.Context, username string) bool {
    db := c.MustGet("DB").(*gorm.DB)
    var user models.User
    if err := db.Where("username = ?", username).First(&user).Error; err == nil {
        return true
    }
    return false
}
