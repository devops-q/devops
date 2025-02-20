package utils

import (
	"itu-minitwit/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserFomContext(ctx *gin.Context) *models.User {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	if userID == nil {
		return nil
	}
	return ctx.MustGet("user").(*models.User)
}
