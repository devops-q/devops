package utils

import (
	"fmt"
	"itu-minitwit/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Creates a user and adds the user to the database.
// It returns a boolean value, that determines if the value insertion is done correctly
func CreateUser(c *gin.Context, username string, email string, password string) {
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
	if db.Create(newUser).Error == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to make user!"})

	}

}
