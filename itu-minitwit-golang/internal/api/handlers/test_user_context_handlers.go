package handlers

import (
	"itu-minitwit/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ForceSetUserId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	session := sessions.Default(ctx)
	session.Set("user_id", id)
}

func GetUserInSession(ctx *gin.Context) {
	user := utils.GetUserFomContext(ctx)

	if user == nil {
		ctx.JSON(http.StatusNotFound, nil)
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}
