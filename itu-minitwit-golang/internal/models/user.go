package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string
	Email     string
	PwHash    string
	Following []*User   `gorm:"many2many:follower"`
	Messages  []Message `gorm:"foreignKey:AuthorID"`
}
