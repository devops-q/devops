package handlers

import (
	"fmt"
	"itu-minitwit/internal/api/json_models"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"itu-minitwit/pkg/logger"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	log := logger.Init()
	var db = c.MustGet("DB").(*gorm.DB)
	var user *models.User = utils.GetUserFomContext(c)

	if user != nil {
		log.Error("[RegisterHandler] User already exists")
		c.Redirect(http.StatusFound, "/")
	}

	var err string
	username := ""
	email := ""

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		password2 := c.PostForm("password2")
		success, registerErrorMessage := service.RegisterUser(db, username, email, password, password2)
		err = registerErrorMessage
		if success {

			utils.SetFlashes(c, "You were successfully registered and can login now")
			c.Redirect(http.StatusFound, "/login")
			return

		} else {
			log.Error("[RegisterHandler] Error registering user: %v", err)
		}
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign Up",
		"body":     "register",
		"Error":    err,
		"Username": username,
		"Email":    email,
		"Endpoint": "/register",
		"Flashes":  utils.GetFlashes(c),
	})

}

func RegisterHandlerAPI(c *gin.Context) {
	var db = c.MustGet("DB").(*gorm.DB)
	var body json_models.RegisterUserBody

	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error binding json", err)
		_ = c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, err := service.RegisterUser(db, body.Username, body.Email, body.Pwd, body.Pwd)
	if success {
		utils.SetFlashes(c, "You were successfully registered and can login now")
		c.JSON(http.StatusNoContent, nil)
		return
	} else {
		c.JSON(http.StatusBadRequest, json_models.ErrorResponse{Code: http.StatusBadRequest, ErrorMessage: err})
	}
}
