package middleware

import (
	"net/http"
	"strings"

	"inventory-manager-server/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a JWT authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" || token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user information to context
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.Username)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// ManagerOnly is a manager-only middleware
func ManagerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if role is manager
		if role != "manager" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
