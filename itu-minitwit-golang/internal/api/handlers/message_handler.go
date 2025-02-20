package handlers

import (
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// This method adds a new message to the DB, and then redirects the user to their own timeline if sucessfull.
func MessageHandler(ctx *gin.Context) {

	if sessions.Default(ctx).Get("user_id") == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in!"})
		return
	}
	text, exists := ctx.GetPostForm("text")

	if exists {
		if created := utils.AddMessage(ctx, text); created {
			utils.SetFlashes(ctx, "Your message was recorded")
			ctx.Redirect(301, "/")
		}
	} else {
		ctx.JSON(http.StatusFailedDependency, gin.H{"error": "The message could not be created!"})
	}

	return
}
