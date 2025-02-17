package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash("You were logged out") 
	session.Save()                     
	c.Redirect(http.StatusFound, "/public")

}