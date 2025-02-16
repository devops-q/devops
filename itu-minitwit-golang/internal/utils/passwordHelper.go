package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hashedPassword string, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil 
}
