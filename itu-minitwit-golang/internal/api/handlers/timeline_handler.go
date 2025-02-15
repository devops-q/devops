package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TimelineHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "layout.html", gin.H{
        "Title":   "Sign In",
        "body":    "timeline",
        "Error":   "",
        "Username": "",
		"UserID" : gin.H{
			"Username" : "phbl",
			"UserID": 3,
		},
        "Email":   "",
		"Endpoint" : "/",
    })
}

