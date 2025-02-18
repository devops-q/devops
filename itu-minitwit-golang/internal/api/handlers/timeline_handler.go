package handlers

import (
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TimelineHandler(c *gin.Context) {

	
	var user *models.User = utils.GetUserFomContext(c)

	if user != nil {
    c.HTML(http.StatusOK, "layout.html", gin.H{
        "Title":   "Sign In",
        "body":    "timeline",
        "Error":   "",
        "Username": user.Username,
		"UserID" : user.ID,
        "Email":   user.Email,
		"Endpoint" : c.FullPath(),
    })
} else {
	// This is when the user is not logged in. 
	c.HTML(http.StatusOK, "layout.html", gin.H{
        "Title":   "Sign In",
        "body":    "timeline",
        "Error":   "",
		"Endpoint" : c.FullPath(),
    })
}
}

