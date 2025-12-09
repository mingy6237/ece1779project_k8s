package handlers

import (
	"net/http"
	"time"

	"inventory-manager-server/config"

	"github.com/gin-gonic/gin"
)

// Index handles the home page
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"instance_id": config.CONFIG.InstanceID,
		"message":     "Inventory Manager Server API",
		"status":      "running",
		"timestamp":   time.Now().Format(time.RFC3339),
		"version":     "1.0.5",
	})
}
