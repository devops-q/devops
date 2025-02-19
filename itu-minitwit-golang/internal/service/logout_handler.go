package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"itu-minitwit/internal/utils"
)

func LogOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	utils.SetFlashes(c, "You were logged out")
	session.Save()
}
