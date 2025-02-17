package handlers

import (
	"itu-minitwit/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	service.LogOutHandler(c)
	c.Redirect(http.StatusFound, "/public")

}