package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"inventory-manager-server/database"
	"inventory-manager-server/dto"
	"inventory-manager-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ListSKUCategories lists all SKU categories
func ListSKUCategories(c *gin.Context) {
	var categories []string
	if err := database.DB.Model(&models.SKU{}).Distinct("category").Pluck("category", &categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// ListSKUs lists all SKUs with pagination, search, and filtering
func ListSKUs(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	pageSize := 20
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build query
	query := database.DB.Model(&models.SKU{})

	// Search by name, category or description
	if search := c.Query("search"); search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Apply sorting
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	validSortFields := map[string]bool{
		"name":       true,
		"category":   true,
		"price":      true,
		"created_at": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	orderBy := sortBy + " " + strings.ToUpper(order)

	// Get SKUs
	var skus []models.SKU
	if err := query.Offset(offset).Limit(pageSize).Order(orderBy).Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Convert to response format
	skuResponses := make([]dto.SKUResponse, len(skus))
	for i, sku := range skus {
		skuResponses[i] = dto.SKUResponse{
			ID:          sku.ID,
			Name:        sku.Name,
			Category:    sku.Category,
			Description: sku.Description,
			Price:       sku.Price,
			Version:     sku.Version,
			CreatedAt:   sku.CreatedAt,
			UpdatedAt:   sku.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"items":       skuResponses,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetSKU gets a single SKU by ID
func GetSKU(c *gin.Context) {
	// Get SKU ID from path parameter
	skuIDStr := c.Param("id")
	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid SKU ID"})
		return
	}

	// Get SKU
	var sku models.SKU
	if err := database.DB.First(&sku, "id = ?", skuID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "SKU not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Convert to response format
	skuResponse := dto.SKUResponse{
		ID:          sku.ID,
		Name:        sku.Name,
		Category:    sku.Category,
		Description: sku.Description,
		Price:       sku.Price,
		Version:     sku.Version,
		CreatedAt:   sku.CreatedAt,
		UpdatedAt:   sku.UpdatedAt,
	}

	c.JSON(http.StatusOK, skuResponse)
}

// CreateSKU creates a SKU (manager only)
func CreateSKU(c *gin.Context) {
	// Parse request body
	var req dto.CreateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Create SKU
	sku := models.SKU{
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := database.DB.Create(&sku).Error; err != nil {
		// Check for unique constraint violation using PostgreSQL error code
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// 23505 is the error code for unique_violation in PostgreSQL
			if strings.Contains(pgErr.ConstraintName, "name") {
				c.JSON(http.StatusBadRequest, gin.H{"message": "SKU name already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create SKU"})
		return
	}

	// Return SKU information
	skuResponse := dto.SKUResponse{
		ID:          sku.ID,
		Name:        sku.Name,
		Category:    sku.Category,
		Description: sku.Description,
		Price:       sku.Price,
		Version:     sku.Version,
		CreatedAt:   sku.CreatedAt,
		UpdatedAt:   sku.UpdatedAt,
	}

	c.JSON(http.StatusCreated, skuResponse)
}

// UpdateSKU updates SKU information (manager only)
func UpdateSKU(c *gin.Context) {
	// Get SKU ID from path parameter
	skuIDStr := c.Param("id")
	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid SKU ID"})
		return
	}

	// Parse request body
	var req dto.UpdateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "errors": err.Error()})
		return
	}

	// Get SKU
	var sku models.SKU
	if err := database.DB.First(&sku, "id = ?", skuID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "SKU not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		sku.Name = req.Name
	}
	if req.Category != "" {
		sku.Category = req.Category
	}
	if req.Description != "" {
		sku.Description = req.Description
	}
	if req.Price > 0 {
		sku.Price = req.Price
	}

	// Increment version on update
	sku.Version++

	if err := database.DB.Save(&sku).Error; err != nil {
		// Check for unique constraint violation using PostgreSQL error code
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// 23505 is the error code for unique_violation in PostgreSQL
			if strings.Contains(pgErr.ConstraintName, "name") {
				c.JSON(http.StatusBadRequest, gin.H{"message": "SKU name already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update SKU"})
		return
	}

	// Return updated SKU
	skuResponse := dto.SKUResponse{
		ID:          sku.ID,
		Name:        sku.Name,
		Category:    sku.Category,
		Description: sku.Description,
		Price:       sku.Price,
		Version:     sku.Version,
		CreatedAt:   sku.CreatedAt,
		UpdatedAt:   sku.UpdatedAt,
	}

	c.JSON(http.StatusOK, skuResponse)
}

// DeleteSKU deletes a SKU (manager only)
func DeleteSKU(c *gin.Context) {
	// Get SKU ID from path parameter
	skuIDStr := c.Param("id")
	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid SKU ID"})
		return
	}

	// Get SKU to delete
	var sku models.SKU
	if err := database.DB.First(&sku, "id = ?", skuID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "SKU not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Check if SKU has inventory
	var inventoryCount int64
	if err := database.DB.Model(&models.Inventory{}).Where("sku_id = ?", skuID).Count(&inventoryCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	if inventoryCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Cannot delete SKU with active inventory"})
		return
	}

	// Delete SKU
	if err := database.DB.Delete(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete SKU"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SKU deleted successfully"})
}
