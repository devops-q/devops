package service

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
	"strings"
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

func UserExists(db *gorm.DB, username string) bool {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err == nil {
		return true
	}
	return false
}

func RegisterUser(db *gorm.DB, username string, email string, password string, password2 string) (success bool, message string) {
	if username == "" {
		return false, "You have to enter a username"
	} else if email == "" || !strings.Contains(email, "@") {
		return false, "You have to enter a valid email address"
	} else if password == "" {
		return false, "You have to enter a password"
	} else if password2 != password {
		return false, "The two passwords do not match"
	} else if UserExists(db, username) {
		return false, "The username is already taken"
	}

	userMade, err := CreateUser(db, username, email, password)
	if err != nil {
		return false, "Failed to create user: " + err.Error()
	}

	if userMade {
		return true, ""
	}

	return false, "Unexpected error occurred"
}
