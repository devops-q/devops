package models

import (
	"gorm.io/gorm"
)

type APIUser struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}
