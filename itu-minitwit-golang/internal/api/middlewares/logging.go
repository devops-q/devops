package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"

	"itu-minitwit/pkg/logger"
)

// SlogMiddleware creates a Gin middleware for logging HTTP requests using slog
func SlogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger.GetLogger().WithService("gin-http").WithGroup("http").Info("request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip", c.ClientIP(),
			"duration", time.Since(start),
			"user_agent", c.Request.UserAgent(),
			"errors", c.Errors.String(),
		)
	}
}
