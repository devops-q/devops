package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	PwHash    string    `gorm:"not null"`
	Following []User    `gorm:"many2many:follower"`
	Messages  []Message `gorm:"foreignKey:AuthorID"`
}
