package handlers

import (
	"itu-minitwit/internal/models"
	"itu-minitwit/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLatest(ctx *gin.Context) {
	log := logger.Init()
	db := ctx.MustGet("DB").(*gorm.DB)
	var latestID models.LatestID
	result := db.Model(&models.LatestID{}).First(&latestID)

	if result.Error != nil {
		log.Error("[GetLatest] Failed to read latest id from DB: %v\n", result.Error)
		ctx.JSON(http.StatusOK, gin.H{"latest": -1})
		return
	}

	// If content is valid and not -1, return it; otherwise, return -1
	if latestID.LatestID != -1 {
		ctx.JSON(http.StatusOK, gin.H{"latest": latestID.LatestID})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"latest": -1})
	}
}
