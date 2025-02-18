package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RetrieveFlashes(c *gin.Context) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save() 

	var flashMSG []interface{}
	if len(flashes) > 0 {
		flashMSG = flashes 
	}
	return flashMSG
}