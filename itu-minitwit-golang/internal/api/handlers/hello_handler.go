package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		name = "World"
	}
	c.HTML(http.StatusOK, "hello.html", gin.H{
		"Name": name,
	})
}
