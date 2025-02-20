package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hashedPassword string, inputPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
