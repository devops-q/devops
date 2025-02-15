package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "layout.html", gin.H{
        "Title":   "Sign Up",
        "body":    "register",
        "Error":   "",
        "Username": "",
        "Email":   "",
    })
}

