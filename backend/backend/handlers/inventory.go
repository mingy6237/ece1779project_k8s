package handlers

import (
	"net/http"
	"strings"

	"inventory-manager-server/database"
	"inventory-manager-server/dto"
	"inventory-manager-server/models"
	"inventory-manager-server/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var inventoryService = services.NewInventoryService()

// GetInventory queries inventory with filtering, sorting, and pagination
func GetInventory(c *gin.Context) {
	// Parse query parameters
	var params dto.InventoryQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters", "errors": err.Error()})
		return
	}

	// Parse UUID query parameters
	var storeID *uuid.UUID
	var skuID *uuid.UUID

	if params.StoreID != nil && *params.StoreID != "" {
		parsedStoreID, err := uuid.Parse(*params.StoreID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters", "errors": "Invalid store_id format"})
			return
		}
		storeID = &parsedStoreID
	}

	if params.SKUID != nil && *params.SKUID != "" {
		parsedSKUID, err := uuid.Parse(*params.SKUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters", "errors": "Invalid sku_id format"})
			return
		}
		skuID = &parsedSKUID
	}

	// Get user role and ID from context
	userRole, _ := c.Get("userRole")
	userIDRaw, _ := c.Get("userID")
	userIDUUID := userIDRaw.(uuid.UUID)

	// For staff, verify access if store_id is specified
	if userRole == "staff" && storeID != nil {
		// Verify staff belongs to the specified store
		var storeUser models.StoreUser
		if err := database.DB.Where("user_id = ? AND store_id = ?", userIDUUID, *storeID).
			First(&storeUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"message": "You can only access inventory for your assigned stores"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
	}

	// Query inventory - service will handle store filtering based on user role
	result, err := inventoryService.GetInventory(params, storeID, skuID, &userIDUUID, userRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query inventory", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetInventoryByID gets a single inventory by ID
func GetInventoryByID(c *gin.Context) {
	// Get inventory ID from path parameter
	inventoryIDStr := c.Param("id")
	inventoryID, err := uuid.Parse(inventoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inventory ID"})
		return
	}

	// Check permissions: staff can only view their store's inventory. Do this before querying inventory.
	userRole, _ := c.Get("userRole")
	if userRole == "staff" {
		userID, _ := c.Get("userID")
		userIDUUID := userID.(uuid.UUID)

		// Need to know the storeID for this inventory: do a pre-query for permission before proceeding.
		var inventory models.Inventory
		// Select only store_id to minimize query
		if err := database.DB.Select("store_id").Where("id = ?", inventoryID).First(&inventory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "Inventory not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
			return
		}

		var storeUser models.StoreUser
		if err := database.DB.Where("user_id = ? AND store_id = ?", userIDUUID, inventory.StoreID).
			First(&storeUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"message": "You can only access inventory for your assigned stores"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
	}

	// Now query the inventory fully
	inventoryRes, err := inventoryService.GetInventoryByID(inventoryID)
	if err != nil {
		if err.Error() == "inventory not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Inventory not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventoryRes)
}

// CreateInventory creates a new inventory record (manager only)
func CreateInventory(c *gin.Context) {
	// Parse request body
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")
	userIDUUID := userID.(uuid.UUID)

	// Create inventory
	inventory, err := inventoryService.CreateInventory(req, userIDUUID, userName.(string))
	if err != nil {
		if err.Error() == "inventory already exists for this SKU and store" {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		if err.Error() == "SKU not found" || err.Error() == "store not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create inventory", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

// UpdateInventory updates inventory quantity (manager only)
func UpdateInventory(c *gin.Context) {
	// Get inventory ID from path parameter
	inventoryIDStr := c.Param("id")
	inventoryID, err := uuid.Parse(inventoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inventory ID"})
		return
	}

	// Parse request body
	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")
	userIDUUID := userID.(uuid.UUID)

	// Update inventory
	inventory, err := inventoryService.UpdateInventory(inventoryID, req.Quantity, userIDUUID, userName.(string))
	if err != nil {
		if err.Error() == "inventory not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Inventory not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update inventory", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// DeleteInventory deletes an inventory record (manager only)
func DeleteInventory(c *gin.Context) {
	// Get inventory ID from path parameter
	inventoryIDStr := c.Param("id")
	inventoryID, err := uuid.Parse(inventoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inventory ID"})
		return
	}

	// Get user info from context
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")
	userIDUUID := userID.(uuid.UUID)

	// Delete inventory
	if err := inventoryService.DeleteInventory(inventoryID, userIDUUID, userName.(string)); err != nil {
		if err.Error() == "inventory not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Inventory not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete inventory", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory deleted successfully"})
}

// AdjustInventory adjusts inventory quantity by delta
// Staff can only adjust their store's inventory, manager can adjust any
func AdjustInventory(c *gin.Context) {
	// Get inventory ID from path parameter
	inventoryIDStr := c.Param("id")
	inventoryID, err := uuid.Parse(inventoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inventory ID"})
		return
	}

	// Parse request body
	var req dto.AdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Get user info from context
	userID, _ := c.Get("userID")
	userName, _ := c.Get("userName")
	userRole, _ := c.Get("userRole")
	userIDUUID := userID.(uuid.UUID)

	// Check permissions: staff can only adjust their store's inventory
	if userRole == "staff" {
		// Get the inventory to check store
		var inventory models.Inventory
		if err := database.DB.First(&inventory, "id = ?", inventoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"message": "Inventory not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}

		// Verify staff belongs to this store
		var storeUser models.StoreUser
		if err := database.DB.Where("user_id = ? AND store_id = ?", userIDUUID, inventory.StoreID).
			First(&storeUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"message": "You can only adjust inventory for your assigned stores"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
			return
		}
	}

	// Adjust inventory
	inventory, err := inventoryService.AdjustInventory(inventoryID, req.DeltaQuantity, userIDUUID, userName.(string))
	if err != nil {
		if err.Error() == "inventory not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Inventory not found"})
			return
		}
		if strings.Contains(err.Error(), "insufficient inventory") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to adjust inventory", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}
