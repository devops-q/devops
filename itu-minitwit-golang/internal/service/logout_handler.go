package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.AddFlash("You were logged out")
	session.Save()
}
