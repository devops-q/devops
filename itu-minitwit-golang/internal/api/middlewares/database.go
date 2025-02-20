package middlewares

import (
	"context"
	"itu-minitwit/pkg/database"
	"time"

	"github.com/gin-gonic/gin"
)

func SetDbMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ctx.Set("DB", database.DB.WithContext(tmt))
		ctx.Next()
	}
}
