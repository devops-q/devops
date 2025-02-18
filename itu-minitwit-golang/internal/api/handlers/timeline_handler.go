package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/config"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/utils"
	"log"
	"net/http"
)

func PublicTimelineHandler(c *gin.Context) {
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

	flashes := utils.GetFlashes(c)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Public Timeline",
		"Endpoint": "/public",
		"body":     "timeline",
		"User":     user,
		"Messages": messages,
		"Flashes":  flashes,
	})
}

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

	flashes := utils.GetFlashes(c)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "My Timeline",
		"body":     "timeline",
		"User":     user,
		"Messages": messages,
		"Endpoint": "/",
		"Flashes":  flashes,
	})
	return
}

func UserTimelineHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	cfg := c.MustGet("Config").(*config.Config)
	username := c.Param("username")

	// Get profile user
	var profileUser models.User
	if err := db.Where("username = ?", username).First(&profileUser).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Check if current user is following profile user
	value, userLoggedIn := c.Get("user")
	var currentUser *models.User
	followed := false

	if userLoggedIn {
		currentUser = value.(*models.User)
		var count int64
		db.Table("follower").
			Where("user_id = ? AND following_id = ?", currentUser.ID, profileUser.ID).
			Count(&count)
		followed = count > 0
	}

	// Get messages for profile user
	var messages []models.Message
	db.Model(&models.Message{}).
		Preload("Author").
		Where("author_id = ? AND flagged = ?", profileUser.ID, false).
		Order("created_at desc").
		Limit(cfg.PerPage).
		Find(&messages)

	flashes := utils.GetFlashes(c)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":       profileUser.Username + "'s Timeline",
		"body":        "timeline",
		"User":        currentUser,
		"Messages":    messages,
		"ProfileUser": profileUser,
		"Followed":    followed,
		"Endpoint":    "/" + username,
		"Flashes":     flashes,
	})

	return
}
