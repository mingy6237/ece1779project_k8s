package handlers

import (
	"inventory-manager-server/cache"
	"inventory-manager-server/config"
	"inventory-manager-server/database"
	"inventory-manager-server/dto"
	"inventory-manager-server/kafka"
	"inventory-manager-server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TestInfra(c *gin.Context) {
	userCount := int64(0)

	// Test database
	db := database.GetDB()
	db.Model(&models.User{}).Count(&userCount)
	newUser := models.User{
		Username:     "test" + strconv.Itoa(int(userCount)),
		PasswordHash: "test1234",
		Email:        "test" + strconv.Itoa(int(userCount)) + "@test.com",
		Role:         "staff",
	}
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
	users := []dto.UserPreviewResponse{}
	db.Model(&models.User{}).Scan(&users)

	// Test cache
	if err := cache.Set("user_count", userCount, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set user count", "details": err.Error()})
		return
	}
	redisUserCount, err := cache.Get("user_count")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count", "details": err.Error()})
		return
	}

	// Test Kafka
	kafka.PublishMessage(config.CONFIG.KafkaTopic, newUser)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Infrastructure test completed",
		"users":          users,
		"count":          len(users),
		"redisUserCount": redisUserCount,
		"status": gin.H{
			"database": "connected",
			"redis":    "connected",
			"kafka":    "connected",
		},
	})
}
