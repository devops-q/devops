package middlewares

import (
	"itu-minitwit/internal/models"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUserContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("user_id")
		if userID != nil {
			db := ctx.MustGet("DB").(*gorm.DB)
			var user *models.User
			db.Limit(1).Find(&user, userID)
			ctx.Set("user", user)
		}
		ctx.Next()
		err := session.Save()
		if err != nil {
			log.Printf("Error saving session: %v", err)
			_ = ctx.AbortWithError(500, err)
		}
	}
}
