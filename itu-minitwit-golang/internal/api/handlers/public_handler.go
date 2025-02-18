package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublicHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	cfg := c.MustGet("Config").(*config.Config)

	var messages []models.Message
	db.Model(&models.Message{}).Preload("Author").Find(&messages).Where("flagged = ?", false).Order("created_at desc").Limit(cfg.PerPage)
	fmt.Printf("Messages: %v\n", messages)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Public Timeline",
		"body":     "timeline",
		"Error":    "",
		"Username": "",
		"Email":    "",
		"Messages": messages,
	})
}
