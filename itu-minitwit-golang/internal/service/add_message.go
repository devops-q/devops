package service

import (
	"fmt"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MessageHandler(ctx *gin.Context) {
	db := ctx.MustGet("DB").(*gorm.DB)

	// Check if the user is logged in
	if sessions.Default(ctx).Get("user_id") == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in!"})
		return
	}

	// Get text input from form
	text, exists := ctx.GetPostForm("text")

	text = "THIS IS A TEST PLEASE WORK!"
	if exists {
		user := utils.GetUserFomContext(ctx)

		newMessage := &models.Message{
			AuthorID: user.ID,
			Author:   user,
			Text:     text,
			Flagged:  false,
		}
		fmt.Println(db.Create(&newMessage))

		ctx.JSON(http.StatusCreated, &newMessage)
	}
	ctx.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign In",
		"body":     "timeline",
		"Error":    "",
		"Username": "Tester",
		"UserID" : 3,
		"Email":    "",
		"Endpoint": "/",
	})
}
