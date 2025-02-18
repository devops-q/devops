package service

import (
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
)

func GetMessagesByAuthor(db *gorm.DB, userID uint, perPage int) ([]models.Message, error) {
	var messages []models.Message
	err := db.Model(&models.Message{}).
		Preload("Author").
		Where("author_id = ? AND flagged = ?", userID, false).
		Order("created_at desc").
		Limit(perPage).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func GetAllMessagesWithAuthors(db *gorm.DB, limit int) ([]models.Message, error) {
	var messages []models.Message

	err := db.Model(&models.Message{}).
		Preload("Author").
		Joins("JOIN users ON messages.author_id = users.id").
		Where("messages.flagged = ?", false).
		Order("messages.created_at desc").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
