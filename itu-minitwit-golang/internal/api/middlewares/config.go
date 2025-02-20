package middlewares

import (
	"github.com/gin-gonic/gin"
	"itu-minitwit/config"
)

func SetConfigMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("Config", cfg)
		ctx.Next()
	}
}
