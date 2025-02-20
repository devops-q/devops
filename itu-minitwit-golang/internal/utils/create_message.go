package utils

import (
	"itu-minitwit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func AddMessage(ctx *gin.Context, text string) (created bool) {

		user := GetUserFomContext(ctx)
		db := ctx.MustGet("DB").(*gorm.DB)


	newMessage := &models.Message{
		AuthorID: user.ID,
		Author:   user,
		Text:     text,
		Flagged:  false,
	}

	if (db.Create(&newMessage) != nil) {
		return true
	} else {
		return false
	}


}