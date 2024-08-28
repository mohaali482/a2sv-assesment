package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func LogController(c *gin.Context) {
	logs, err := os.ReadFile("log")
	if err != nil {
		log.Println("Error reading log file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve logs"})
		return
	}

	// Return the logs as a response
	c.Data(http.StatusOK, "text/plain", logs)
}
