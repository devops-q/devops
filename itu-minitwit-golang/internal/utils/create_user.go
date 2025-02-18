package utils

import (
	"itu-minitwit/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Creates a user and adds the user to the database.
// It returns a boolean value, that determines if the value insertion is done correctly
func CreateUser(db *gorm.DB, username string, email string, password string) (bool, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err 
	}

	newUser := &models.User{
		Username: username,
		Email:    email,
		PwHash:   string(hashed),
	}

	if err := db.Create(newUser).Error; err != nil {
		return false, err 
	}

	return true, nil 
}

