package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/internal/api/json_models"
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"net/http"
	"strconv"
)

func MessagesHandlerAPI(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	utils.UpdateLatest(c.DefaultQuery("latest", "-1"))

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

	var formattedMessages = service.MapMessages(messages)

	c.JSON(http.StatusOK, formattedMessages)
}

func MessagesPerUserHandlerAPI(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	utils.UpdateLatest(c.DefaultQuery("latest", "-1"))

	nrOfMessagesParam := c.DefaultQuery("no", "100")
	nrOfMessages, err := strconv.Atoi(nrOfMessagesParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, json_models.ErrorResponse{
			Code:         http.StatusBadRequest,
			ErrorMessage: "Invalid number of messages provided in param\"no\"",
		})
		return
	}

	username := c.Param("username")

	userId, err := service.GetUserIdByUsername(db, username)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error fetching user",
		})
		return
	}

	messages, err := service.GetMessagesByAuthor(db, uint(userId), nrOfMessages)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error fetching messages",
		})
		return
	}

	var formattedMessages = service.MapMessages(messages)

	c.JSON(http.StatusOK, formattedMessages)
}

func MessagesCreateHandlerAPI(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	username := c.Param("username")

	userId, err := service.GetUserIdByUsername(db, username)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error fetching user",
		})
		return
	}

	var body json_models.CreateMessageBody

	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error binding json", err)
		_ = c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = service.CreateMessage(db, uint(userId), body.Content)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error creating message",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
