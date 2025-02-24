package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"itu-minitwit/internal/utils"
	"log"
)

func LogOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	utils.SetFlashes(c, "You were logged out")
	err := session.Save()
	if err != nil {
		log.Printf("Error saving session: %v", err)
		_ = c.AbortWithError(500, err)
	}
}
