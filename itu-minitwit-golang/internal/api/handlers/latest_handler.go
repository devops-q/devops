package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLatest(ctx *gin.Context) {
	// Attempt to read the file
	content, err := os.ReadFile("/Users/philipblomholt/devops/itu-minitwit-golang/internal/api/handlers/latest_processed_sim_action_id.txt")
	if err != nil {
		log.Printf("Failed to read file: %v\n", err)
		ctx.JSON(http.StatusOK, gin.H{"latest": -1})
		return
	}

	// Attempt to convert file content to integer
	contentInt, err := strconv.Atoi(string(content))
	if err != nil {
		log.Printf("Failed to parse file content to integer: %v\n", err)
		ctx.JSON(http.StatusOK, gin.H{"latest": -1})
		return
	}

	// If content is valid and not -1, return it; otherwise, return -1
	if contentInt != -1 {
		ctx.JSON(http.StatusOK, gin.H{"latest": contentInt})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"latest": -1})
	}
}