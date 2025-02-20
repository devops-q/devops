package middlewares

import (
	"itu-minitwit/internal/utils"

	"github.com/gin-gonic/gin"
)

func UpdateLatestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		utils.UpdateLatest(ctx.DefaultQuery("latest", "-1"))
	}
}
