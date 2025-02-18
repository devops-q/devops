package service

import (
	"itu-minitwit/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleRegister(c *gin.Context, username string, email string, password string, password2 string) (message string) {

	if username == "" {
		return "You have to enter a username"
	} else if email == "" || !strings.Contains(email, "@") {
		return "You have to enter a valid email address"
	} else if password == "" {
		return "You have to enter a password"
	} else if password2 != password {
		return "The two passwords do not match"
	} else if utils.UserExists(c, username) {
		return "The username is already taken"
	}

	db := c.MustGet("DB").(*gorm.DB)

	if utils.UserExists(c, username) {
		return "The username is already taken"
	}

	userMade, err := utils.CreateUser(db, username, email, password)
	if err != nil {
		return "Failed to create user: " + err.Error()
	}

	if userMade {
		return ""
	}

	return "Unexpected error occurred"
}
