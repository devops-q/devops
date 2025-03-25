package handlers

import (
	"itu-minitwit/internal/utils"
	"itu-minitwit/pkg/logger"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MessageHandler(ctx *gin.Context) {
	log := logger.Init()
	if sessions.Default(ctx).Get("user_id") == nil {
		log.Error("[MessageHandler] Failed to get user id from session")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in!"})
		return
	}
	text, exists := ctx.GetPostForm("text")

	if exists {
		if created := utils.AddMessage(ctx, text); created {
			utils.SetFlashes(ctx, "Your message was recorded")
			log.Info("[MessageHandler] Successfully added message")
			ctx.Redirect(301, "/")
		}
	} else {
		log.Error("[MessageHandler] Failed to get text from context")
		ctx.JSON(http.StatusFailedDependency, gin.H{"error": "The message could not be created!"})
	}

}
