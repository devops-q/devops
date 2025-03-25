package utils

import (
	"fmt"
	"itu-minitwit/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateLatest(ctx *gin.Context, parsedCommandId string) {
	if parsedCommandIdInt, err := strconv.Atoi(parsedCommandId); err != nil {
		fmt.Println("Couldn't convert value to Integer")
		return
	} else {
		if parsedCommandIdInt != -1 {
			db := ctx.MustGet("DB").(*gorm.DB)
			var latestID *models.LatestID
			result := db.Model(&models.LatestID{}).First(&latestID)
			if result.Error == nil {
				db.Delete(&models.LatestID{}, latestID.LatestID)
			}
			db.Model(&models.LatestID{}).Create(&models.LatestID{
				LatestID: parsedCommandIdInt,
			})
		}
	}
}
