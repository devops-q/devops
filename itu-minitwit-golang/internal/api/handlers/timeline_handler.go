package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"log"
	"net/http"
)

func TimelineHandler(c *gin.Context) {
	log.Printf("We got a visitor from: %s\n", c.ClientIP())

	db := c.MustGet("DB").(*gorm.DB)
	cfg := c.MustGet("Config").(*config.Config)

	value, userLoggedIn := c.Get("user")
	var user *models.User

	if !userLoggedIn {
		c.Redirect(http.StatusFound, "/public")
	}

	user = value.(*models.User)

	var messages []models.Message
	db.
		Model(&models.Message{}).Preload("Author").
		Joins("JOIN users ON messages.author_id = users.id").
		Where("messages.flagged = ? AND (messages.author_id = ? OR messages.author_id IN (SELECT following_id FROM follower WHERE user_id = ?))",
			false, user.ID, user.ID).
		Order("messages.created_at desc").Limit(cfg.PerPage).
		Find(&messages)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "My Timeline",
		"body":     "timeline",
		"User":     user,
		"Messages": messages,
		"Endpoint": "/",
	})
}
