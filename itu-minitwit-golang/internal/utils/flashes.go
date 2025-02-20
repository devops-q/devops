package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFlashes(c *gin.Context) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes()
	if err := session.Save(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	return flashes
}

func SetFlashes(c *gin.Context, flash string) {
	session := sessions.Default(c)
	session.AddFlash(flash)
	if err := session.Save(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

}
