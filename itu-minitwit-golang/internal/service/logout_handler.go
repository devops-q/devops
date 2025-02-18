package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)


func LogOutHandler(c *gin.Context) {
session := sessions.Default(c)
session.Clear()
session.AddFlash("You were logged out") 
session.Save()
}