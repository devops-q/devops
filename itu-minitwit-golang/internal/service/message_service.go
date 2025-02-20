package service

import (
	"gorm.io/gorm"
	"itu-minitwit/internal/api/json_models"
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

func MapMessage(message models.Message) json_models.Message {
	return json_models.Message{
		Content: message.Text,
		PubDate: message.CreatedAt,
		User:    message.Author.Username,
	}
}

func MapMessages(messages []models.Message) []json_models.Message {
	var formattedMessages = make([]json_models.Message, 0)
	for _, message := range messages {
		formattedMessages = append(formattedMessages, MapMessage(message))
	}
	return formattedMessages
}

func CreateMessage(db *gorm.DB, authorID uint, content string) error {
	err := db.Model(&models.Message{}).
		Create(&models.Message{
			AuthorID: authorID,
			Text:     content,
		}).Error

	return err
}
