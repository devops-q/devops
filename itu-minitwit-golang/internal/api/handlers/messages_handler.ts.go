package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/internal/api/json_models"
	"itu-minitwit/internal/service"
	"net/http"
	"strconv"
)

func MessagesHandlerAPI(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	nrOfMessagesParam := c.DefaultQuery("no", "100")
	nrOfMessages, err := strconv.Atoi(nrOfMessagesParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, json_models.ErrorResponse{
			Code:         http.StatusBadRequest,
			ErrorMessage: "Invalid number of messages provided in param\"no\"",
		})
		return
	}

	messages, err := service.GetAllMessagesWithAuthors(db, nrOfMessages)

	if err != nil {
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error fetching messages",
		})
		return
	}

	var formattedMessages = make([]json_models.Message, 0)
	for _, message := range messages {
		formattedMessages = append(formattedMessages, json_models.Message{
			Content: message.Text,
			PubDate: message.CreatedAt,
			User:    message.Author.Username,
		})
	}

	c.JSON(http.StatusOK, formattedMessages)
}
