package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func RetrieveFlashes(c *gin.Context) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes()
	err := session.Save()
	if err != nil {
		log.Printf("Error saving session: %v", err)
		_ = c.AbortWithError(500, err)
	}
	var flashMSG []interface{}
	if len(flashes) > 0 {
		flashMSG = flashes
	}
	return flashMSG
}
