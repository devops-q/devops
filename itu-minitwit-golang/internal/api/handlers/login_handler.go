package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "layout.html", gin.H{
        "Title":   "Sign In",
        "body":    "login",
        "Error":   "",
        "Username": "",
        "Email":   "",
    })
}

