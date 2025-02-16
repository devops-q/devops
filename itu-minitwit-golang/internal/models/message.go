package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	AuthorID uint   `gorm:"not null"`
	Author   *User  `gorm:"foreignkey:author_id;association_foreignkey:id"`
	Text     string `gorm:"not null"`
	Flagged  bool   `gorm:"default:false"`
}
