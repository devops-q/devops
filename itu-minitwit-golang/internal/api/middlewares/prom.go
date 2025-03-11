package middlewares

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware(r *gin.Engine) *ginprom.Prometheus {
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	return p
}
