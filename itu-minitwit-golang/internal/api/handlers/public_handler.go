package handlers

import (
	"gorm.io/gorm"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublicHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	cfg := c.MustGet("Config").(*config.Config)

	value, userLoggedIn := c.Get("user")
	var user *models.User

	if userLoggedIn {
		user = value.(*models.User)
	} else {
		user = nil
	}

	var messages []models.Message
	db.Model(&models.Message{}).
		Preload("Author").
		Where(map[string]interface{}{"flagged": false}).
		Order("created_at desc").Limit(cfg.PerPage).
		Find(&messages)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Public Timeline",
		"Endpoint": "/public",
		"body":     "timeline",
		"User":     user,
		"Messages": messages,
	})
}
