package routes

import (
	"inventory-manager-server/handlers"
	"inventory-manager-server/middleware"
	"inventory-manager-server/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up routes
func SetupRoutes() *gin.Engine {
	// Create Gin engine
	router := gin.New()
	// Setup middleware (CORS, logging, etc.)
	// Skip logging for health check endpoint to reduce noise
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	router.Use(gin.Recovery())

	// CORS middleware - permissive configuration
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Set CORS headers for all requests
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check and index routes
	router.GET("/", handlers.Index)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.POST("/testInfra", handlers.TestInfra)

	// Public route
	router.POST("/api/auth/login", handlers.Login)

	// WebSocket route (uses query parameter authentication, not middleware)
	router.GET("/api/ws", func(c *gin.Context) {
		websocket.ServeWS(c.Writer, c.Request)
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	// Authenticated route
	authed := router.Group("/api")
	authed.Use(middleware.AuthMiddleware())

	// User profile route
	userProfile := authed.Group("/profile")
	{
		userProfile.GET("", handlers.GetProfile)
		userProfile.PUT("/password", handlers.ChangePassword)
	}

	// User management route
	userManagement := authed.Group("/manager/users")
	userManagement.Use(middleware.ManagerOnly())
	{
		userManagement.GET("", handlers.ListUsers)
		userManagement.POST("", handlers.CreateUser)
		userManagement.PUT("", handlers.UpdateUser)
		userManagement.DELETE("/:id", handlers.DeleteUser)
	}

	// Store & Staff management route
	storeManagement := authed.Group("/manager/stores")
	storeManagement.Use(middleware.ManagerOnly())
	{
		storeManagement.GET("", handlers.ListStores)
		storeManagement.POST("", handlers.CreateStore)
		storeManagement.DELETE("/:id", handlers.DeleteStore)
		staffManagement := storeManagement.Group("/staff")
		{
			staffManagement.GET("", handlers.ListStoreStaff)
			staffManagement.POST("", handlers.AddStaffToStore)
			staffManagement.DELETE("/:id", handlers.DeleteStaffFromStore)
		}
	}

	// SKU management route
	skuManagement := authed.Group("/manager/skus")
	skuManagement.Use(middleware.ManagerOnly())
	{
		skuManagement.GET("", handlers.ListSKUs)
		skuManagement.GET("/categories", handlers.ListSKUCategories)
		skuManagement.GET("/:id", handlers.GetSKU)
		skuManagement.POST("", handlers.CreateSKU)
		skuManagement.PUT("/:id", handlers.UpdateSKU)
		skuManagement.DELETE("/:id", handlers.DeleteSKU)
	}

	// Inventory management route (all authenticated users can query)
	inventory := authed.Group("/inventory")
	{
		inventory.GET("", handlers.GetInventory)
		inventory.GET("/:id", handlers.GetInventoryByID)
		// Adjust endpoint - staff can only adjust their store's inventory
		inventory.POST("/:id/adjust", handlers.AdjustInventory)
	}

	// Inventory management route (manager only - create, update, delete)
	inventoryManagement := authed.Group("/manager/inventory")
	inventoryManagement.Use(middleware.ManagerOnly())
	{
		inventoryManagement.POST("", handlers.CreateInventory)
		inventoryManagement.PUT("/:id", handlers.UpdateInventory)
		inventoryManagement.DELETE("/:id", handlers.DeleteInventory)
	}

	return router
}
