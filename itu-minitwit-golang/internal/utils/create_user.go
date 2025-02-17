package utils

import (
	"fmt"
	"itu-minitwit/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func CreateUser(c *gin.Context, username string, email string, password string)  {
	db := c.MustGet("DB").(*gorm.DB)

	hashed, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if error != nil {
				fmt.Println("Error hashing password: ", error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}

	var newUser = &models.User{
		Username: username,
		Email:    email,
		PwHash:   string(hashed),
	}
	db.Create(newUser)
	return
}
