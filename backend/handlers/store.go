package handlers

import (
	"net/http"

	"inventory-manager-server/database"
	"inventory-manager-server/dto"
	"inventory-manager-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateStore creates a store (manager only)
func CreateStore(c *gin.Context) {

	// Parse request body
	var req dto.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Create store
	store := models.Store{
		Name:    req.Name,
		Address: req.Address,
	}

	if err := database.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create store"})
		return
	}

	// Return store information
	storeResponse := dto.StoreResponse{
		ID:        store.ID,
		Name:      store.Name,
		Address:   store.Address,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}

	c.JSON(http.StatusCreated, storeResponse)
}

// DeleteStore deletes a store (manager only)
func DeleteStore(c *gin.Context) {
	// Get store ID from path parameter
	storeIDStr := c.Param("id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid store ID"})
		return
	}

	// Get store to delete
	var store models.Store
	if err := database.DB.First(&store, "id = ?", storeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Store not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Check if store has inventory
	var inventoryCount int64
	if err := database.DB.Model(&models.Inventory{}).Where("store_id = ?", storeID).Count(&inventoryCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	if inventoryCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Cannot delete store with active inventory"})
		return
	}

	// Check if store has associated users
	var storeUserCount int64
	if err := database.DB.Model(&models.StoreUser{}).Where("store_id = ?", storeID).Count(&storeUserCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	if storeUserCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Cannot delete store with associated users"})
		return
	}

	// Delete store
	if err := database.DB.Delete(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store deleted successfully"})
}

// ListStores lists all stores
func ListStores(c *gin.Context) {
	// Query all stores
	var stores []models.Store
	if err := database.DB.Order("created_at DESC").Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Convert to response format
	storeResponses := make([]dto.StoreResponse, len(stores))
	for i, store := range stores {
		storeResponses[i] = dto.StoreResponse{
			ID:        store.ID,
			Name:      store.Name,
			Address:   store.Address,
			CreatedAt: store.CreatedAt,
			UpdatedAt: store.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": storeResponses,
	})
}

// ListStoreStaff lists all staff members for a store
func ListStoreStaff(c *gin.Context) {
	// Get store ID from query parameter
	storeIDStr := c.Query("store_id")
	if storeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "store_id query parameter is required"})
		return
	}

	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid store ID"})
		return
	}

	// Verify store exists
	var store models.Store
	if err := database.DB.First(&store, "id = ?", storeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Store not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Query store staff with user information
	var storeStaff []models.StoreUser
	if err := database.DB.Preload("User").
		Where("store_id = ?", storeID).
		Order("created_at DESC").
		Find(&storeStaff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Convert to response format - just user info, no repeated store info
	staffResponses := make([]dto.UserResponse, len(storeStaff))
	for i, staff := range storeStaff {
		staffResponses[i] = dto.UserResponse{
			ID:        staff.User.ID,
			Username:  staff.User.Username,
			Email:     staff.User.Email,
			Role:      staff.User.Role,
			CreatedAt: staff.User.CreatedAt,
			UpdatedAt: staff.User.UpdatedAt,
		}
	}

	response := dto.StoreStaffListResponse{
		StoreID: storeID,
		Staff:   staffResponses,
	}

	c.JSON(http.StatusOK, response)
}

// AddStaffToStore adds a staff member to a store (manager only)
func AddStaffToStore(c *gin.Context) {
	// Parse request body
	var req dto.AddStaffToStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Verify store exists
	var store models.Store
	if err := database.DB.First(&store, "id = ?", req.StoreID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Store not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Verify user exists
	var user models.User
	if err := database.DB.First(&user, "id = ?", req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Check if association already exists
	var existingStoreUser models.StoreUser
	if err := database.DB.Where("store_id = ? AND user_id = ?", req.StoreID, req.UserID).
		First(&existingStoreUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "User is already associated with this store"})
		return
	}

	// Create store-user association
	storeUser := models.StoreUser{
		StoreID: req.StoreID,
		UserID:  req.UserID,
	}

	if err := database.DB.Create(&storeUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add staff to store"})
		return
	}

	// Load user and store for response
	if err := database.DB.Preload("User").Preload("Store").First(&storeUser, storeUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve created association"})
		return
	}

	// Return staff information
	staffResponse := dto.StoreStaffResponse{
		ID:      storeUser.ID,
		StoreID: storeUser.StoreID,
		UserID:  storeUser.UserID,
		User: dto.UserResponse{
			ID:        storeUser.User.ID,
			Username:  storeUser.User.Username,
			Email:     storeUser.User.Email,
			Role:      storeUser.User.Role,
			CreatedAt: storeUser.User.CreatedAt,
			UpdatedAt: storeUser.User.UpdatedAt,
		},
		Store: dto.StoreResponse{
			ID:        storeUser.Store.ID,
			Name:      storeUser.Store.Name,
			Address:   storeUser.Store.Address,
			CreatedAt: storeUser.Store.CreatedAt,
			UpdatedAt: storeUser.Store.UpdatedAt,
		},
		CreatedAt: storeUser.CreatedAt,
		UpdatedAt: storeUser.UpdatedAt,
	}

	c.JSON(http.StatusCreated, staffResponse)
}

// DeleteStaffFromStore removes a staff member from a store (manager only)
func DeleteStaffFromStore(c *gin.Context) {
	// Get store staff ID from path parameter
	staffIDStr := c.Param("id")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid staff ID"})
		return
	}

	// Get store-user association to delete
	var storeUser models.StoreUser
	if err := database.DB.First(&storeUser, "id = ?", staffID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Store staff association not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Delete store-user association
	if err := database.DB.Delete(&storeUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove staff from store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff removed from store successfully"})
}
