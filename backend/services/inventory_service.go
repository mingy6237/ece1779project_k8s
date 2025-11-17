package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"inventory-manager-server/cache"
	"inventory-manager-server/config"
	"inventory-manager-server/database"
	"inventory-manager-server/dto"
	"inventory-manager-server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InventoryService handles inventory operations
type InventoryService struct {
	// Using global database and cache instances
}

// NewInventoryService creates a new inventory service
func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

// GetInventory gets inventory with caching, filtering, sorting, and pagination
// userID and userRole are used to determine allowed stores for staff members
func (s *InventoryService) GetInventory(params dto.InventoryQueryParams, userID *uuid.UUID, userRole string) (*dto.InventoryListResponse, error) {
	// Build cache key
	cacheKey := s.buildCacheKey(params, userID, userRole)

	// Try to get from cache
	if cache.Client != nil {
		cached, err := cache.Get(cacheKey)
		if err == nil && cached != "" {
			var result dto.InventoryListResponse
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return &result, nil
			}
		}
	}

	// Build query
	query := database.DB.Model(&models.Inventory{})

	// Apply filters
	// StoreID: nil means all stores (for manager) or all assigned stores (for staff), non-nil means specific store
	if params.StoreID != nil {
		query = query.Where("store_id = ?", *params.StoreID)
	} else if userRole == "staff" && userID != nil {
		// For staff without specific store_id, query their allowed stores
		var storeUsers []models.StoreUser
		if err := database.DB.Where("user_id = ?", *userID).Find(&storeUsers).Error; err == nil && len(storeUsers) > 0 {
			storeIDs := make([]uuid.UUID, len(storeUsers))
			for i, su := range storeUsers {
				storeIDs[i] = su.StoreID
			}
			query = query.Where("store_id IN ?", storeIDs)
		} else {
			// Staff has no stores, return empty result
			return &dto.InventoryListResponse{
				Items:      []dto.InventoryResponse{},
				Total:      0,
				Page:       params.Page,
				PageSize:   params.PageSize,
				TotalPages: 0,
			}, nil
		}
	}
	// SKUID: nil means all SKUs, non-nil means specific SKU
	// When SKUID is nil, return all SKUs' inventory for the specified store(s)
	if params.SKUID != nil {
		query = query.Where("sku_id = ?", *params.SKUID)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count inventory: %w", err)
	}

	// Apply sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	order := params.Order
	if order == "" {
		order = "desc"
	}
	validSortFields := map[string]bool{
		"quantity":   true,
		"created_at": true,
		"updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	orderBy := sortBy + " " + strings.ToUpper(order)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	var inventories []models.Inventory
	if err := query.Preload("SKU").Preload("Store").
		Offset(offset).Limit(params.PageSize).Order(orderBy).Find(&inventories).Error; err != nil {
		return nil, fmt.Errorf("failed to query inventory: %w", err)
	}

	// Convert to response format
	items := make([]dto.InventoryResponse, len(inventories))
	for i, inv := range inventories {
		items[i] = dto.InventoryResponse{
			ID:        inv.ID,
			SKUID:     inv.SKUID,
			StoreID:   inv.StoreID,
			Quantity:  inv.Quantity,
			Version:   inv.Version,
			CreatedAt: inv.CreatedAt,
			UpdatedAt: inv.UpdatedAt,
		}
		// Include SKU info if loaded
		if inv.SKU.ID != uuid.Nil {
			items[i].SKU = &dto.SKUResponse{
				ID:          inv.SKU.ID,
				Name:        inv.SKU.Name,
				Category:    inv.SKU.Category,
				Description: inv.SKU.Description,
				Price:       inv.SKU.Price,
				Version:     inv.SKU.Version,
				CreatedAt:   inv.SKU.CreatedAt,
				UpdatedAt:   inv.SKU.UpdatedAt,
			}
		}
		// Include Store info if loaded
		if inv.Store.ID != uuid.Nil {
			items[i].Store = &dto.StoreResponse{
				ID:        inv.Store.ID,
				Name:      inv.Store.Name,
				Address:   inv.Store.Address,
				CreatedAt: inv.Store.CreatedAt,
				UpdatedAt: inv.Store.UpdatedAt,
			}
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	result := &dto.InventoryListResponse{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}

	// Cache the result
	if cache.Client != nil {
		if data, err := json.Marshal(result); err == nil {
			cache.Set(cacheKey, string(data), 5*time.Minute)
		}
	}

	return result, nil
}

// GetInventoryByID gets a single inventory by ID
func (s *InventoryService) GetInventoryByID(id uuid.UUID) (*dto.InventoryResponse, error) {
	var inventory models.Inventory
	if err := database.DB.Preload("SKU").Preload("Store").
		First(&inventory, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, fmt.Errorf("failed to query inventory: %w", err)
	}

	response := &dto.InventoryResponse{
		ID:        inventory.ID,
		SKUID:     inventory.SKUID,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity,
		Version:   inventory.Version,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}

	if inventory.SKU.ID != uuid.Nil {
		response.SKU = &dto.SKUResponse{
			ID:          inventory.SKU.ID,
			Name:        inventory.SKU.Name,
			Category:    inventory.SKU.Category,
			Description: inventory.SKU.Description,
			Price:       inventory.SKU.Price,
			Version:     inventory.SKU.Version,
			CreatedAt:   inventory.SKU.CreatedAt,
			UpdatedAt:   inventory.SKU.UpdatedAt,
		}
	}

	if inventory.Store.ID != uuid.Nil {
		response.Store = &dto.StoreResponse{
			ID:        inventory.Store.ID,
			Name:      inventory.Store.Name,
			Address:   inventory.Store.Address,
			CreatedAt: inventory.Store.CreatedAt,
			UpdatedAt: inventory.Store.UpdatedAt,
		}
	}

	return response, nil
}

// CreateInventory creates a new inventory record
func (s *InventoryService) CreateInventory(req dto.CreateInventoryRequest, userID uuid.UUID, userName string) (*dto.InventoryResponse, error) {
	// Check if inventory already exists
	var existing models.Inventory
	if err := database.DB.Where("sku_id = ? AND store_id = ?", req.SKUID, req.StoreID).
		First(&existing).Error; err == nil {
		return nil, fmt.Errorf("inventory already exists for this SKU and store")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing inventory: %w", err)
	}

	// Verify SKU exists
	var sku models.SKU
	if err := database.DB.First(&sku, "id = ?", req.SKUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("SKU not found")
		}
		return nil, fmt.Errorf("failed to query SKU: %w", err)
	}

	// Verify Store exists
	var store models.Store
	if err := database.DB.First(&store, "id = ?", req.StoreID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("store not found")
		}
		return nil, fmt.Errorf("failed to query store: %w", err)
	}

	// Create inventory in transaction
	var inventory models.Inventory
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create inventory
		inventory = models.Inventory{
			SKUID:    req.SKUID,
			StoreID:  req.StoreID,
			Quantity: req.Quantity,
			Version:  1,
		}
		if err := tx.Create(&inventory).Error; err != nil {
			return fmt.Errorf("failed to create inventory: %w", err)
		}

		// Create outbox record
		outbox := models.Outbox{
			OperationType:    "create",
			SenderInstanceID: config.CONFIG.InstanceID,
			InventoryID:      inventory.ID,
			SKUID:            req.SKUID,
			SKUName:          sku.Name,
			StoreID:          req.StoreID,
			StoreName:        store.Name,
			UserID:           userID,
			UserName:         userName,
			DeltaQuantity:    req.Quantity,
			NewQuantity:      req.Quantity,
			Version:          1,
		}
		if err := tx.Create(&outbox).Error; err != nil {
			return fmt.Errorf("failed to create outbox record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Delete cache
	s.invalidateCache(req.StoreID, nil)

	// Return response
	response := &dto.InventoryResponse{
		ID:        inventory.ID,
		SKUID:     inventory.SKUID,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity,
		Version:   inventory.Version,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}

	return response, nil
}

// UpdateInventory updates inventory quantity with optimistic locking
func (s *InventoryService) UpdateInventory(id uuid.UUID, quantity int, userID uuid.UUID, userName string) (*dto.InventoryResponse, error) {
	var inventory models.Inventory
	var sku models.SKU
	var store models.Store

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get inventory with lock
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Preload("SKU").Preload("Store").
			First(&inventory, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("inventory not found")
			}
			return fmt.Errorf("failed to query inventory: %w", err)
		}

		sku = inventory.SKU
		store = inventory.Store

		// Calculate delta
		deltaQuantity := quantity - inventory.Quantity

		// Update inventory
		inventory.Quantity = quantity
		inventory.Version++
		if err := tx.Save(&inventory).Error; err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}

		// Create outbox record
		outbox := models.Outbox{
			OperationType:    "update",
			SenderInstanceID: config.CONFIG.InstanceID,
			InventoryID:      inventory.ID,
			SKUID:            inventory.SKUID,
			SKUName:          sku.Name,
			StoreID:          inventory.StoreID,
			StoreName:        store.Name,
			UserID:           userID,
			UserName:         userName,
			DeltaQuantity:    deltaQuantity,
			NewQuantity:      quantity,
			Version:          inventory.Version,
		}
		if err := tx.Create(&outbox).Error; err != nil {
			return fmt.Errorf("failed to create outbox record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Delete cache
	s.invalidateCache(inventory.StoreID, &inventory.SKUID)

	// Return response
	response := &dto.InventoryResponse{
		ID:        inventory.ID,
		SKUID:     inventory.SKUID,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity,
		Version:   inventory.Version,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}

	return response, nil
}

// AdjustInventory adjusts inventory quantity by delta
func (s *InventoryService) AdjustInventory(id uuid.UUID, deltaQuantity int, userID uuid.UUID, userName string) (*dto.InventoryResponse, error) {
	var inventory models.Inventory
	var sku models.SKU
	var store models.Store

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get inventory with lock
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Preload("SKU").Preload("Store").
			First(&inventory, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("inventory not found")
			}
			return fmt.Errorf("failed to query inventory: %w", err)
		}

		sku = inventory.SKU
		store = inventory.Store

		// Check if adjustment would result in negative quantity
		newQuantity := inventory.Quantity + deltaQuantity
		if newQuantity < 0 {
			return fmt.Errorf("insufficient inventory: current quantity is %d, cannot adjust by %d", inventory.Quantity, deltaQuantity)
		}

		// Update inventory
		inventory.Quantity = newQuantity
		inventory.Version++
		if err := tx.Save(&inventory).Error; err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}

		// Create outbox record
		outbox := models.Outbox{
			OperationType:    "adjust",
			SenderInstanceID: config.CONFIG.InstanceID,
			InventoryID:      inventory.ID,
			SKUID:            inventory.SKUID,
			SKUName:          sku.Name,
			StoreID:          inventory.StoreID,
			StoreName:        store.Name,
			UserID:           userID,
			UserName:         userName,
			DeltaQuantity:    deltaQuantity,
			NewQuantity:      newQuantity,
			Version:          inventory.Version,
		}
		if err := tx.Create(&outbox).Error; err != nil {
			return fmt.Errorf("failed to create outbox record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Delete cache
	s.invalidateCache(inventory.StoreID, &inventory.SKUID)

	// Return response
	response := &dto.InventoryResponse{
		ID:        inventory.ID,
		SKUID:     inventory.SKUID,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity,
		Version:   inventory.Version,
		CreatedAt: inventory.CreatedAt,
		UpdatedAt: inventory.UpdatedAt,
	}

	return response, nil
}

// DeleteInventory deletes an inventory record
func (s *InventoryService) DeleteInventory(id uuid.UUID, userID uuid.UUID, userName string) error {
	var inventory models.Inventory
	var sku models.SKU
	var store models.Store

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get inventory
		if err := tx.Preload("SKU").Preload("Store").
			First(&inventory, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("inventory not found")
			}
			return fmt.Errorf("failed to query inventory: %w", err)
		}

		sku = inventory.SKU
		store = inventory.Store

		// Delete inventory
		if err := tx.Delete(&inventory).Error; err != nil {
			return fmt.Errorf("failed to delete inventory: %w", err)
		}

		// Create outbox record (marking deletion)
		outbox := models.Outbox{
			OperationType:    "delete",
			SenderInstanceID: config.CONFIG.InstanceID,
			InventoryID:      inventory.ID,
			SKUID:            inventory.SKUID,
			SKUName:          sku.Name,
			StoreID:          inventory.StoreID,
			StoreName:        store.Name,
			UserID:           userID,
			UserName:         userName,
			DeltaQuantity:    -inventory.Quantity,
			NewQuantity:      0,
			Version:          inventory.Version + 1,
		}
		if err := tx.Create(&outbox).Error; err != nil {
			return fmt.Errorf("failed to create outbox record: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Delete cache
	s.invalidateCache(inventory.StoreID, &inventory.SKUID)

	return nil
}

// buildCacheKey builds a cache key for inventory query
func (s *InventoryService) buildCacheKey(params dto.InventoryQueryParams, userID *uuid.UUID, userRole string) string {
	parts := []string{"inventory"}
	if params.StoreID != nil {
		parts = append(parts, "store", params.StoreID.String())
	} else if userRole == "staff" && userID != nil {
		// For staff, query allowed stores and include in cache key
		var storeUsers []models.StoreUser
		if err := database.DB.Where("user_id = ?", *userID).Find(&storeUsers).Error; err == nil && len(storeUsers) > 0 {
			storeIDStrs := make([]string, len(storeUsers))
			for i, su := range storeUsers {
				storeIDStrs[i] = su.StoreID.String()
			}
			sort.Strings(storeIDStrs)
			parts = append(parts, "allowed", strings.Join(storeIDStrs, ","))
		} else {
			parts = append(parts, "store", "none")
		}
	} else {
		parts = append(parts, "store", "all")
	}
	if params.SKUID != nil {
		parts = append(parts, "sku", params.SKUID.String())
	}
	parts = append(parts, "page", strconv.Itoa(params.Page))
	parts = append(parts, "size", strconv.Itoa(params.PageSize))
	parts = append(parts, "sort", params.SortBy)
	parts = append(parts, "order", params.Order)
	return strings.Join(parts, ":")
}

// invalidateCache invalidates cache for a store
func (s *InventoryService) invalidateCache(storeID uuid.UUID, skuID *uuid.UUID) {
	if cache.Client == nil {
		return
	}

	ctx := context.Background()

	// Delete all cache keys matching inventory patterns
	// Use SCAN to avoid blocking Redis
	var cursor uint64
	var keys []string

	// Scan for store-specific cache keys
	pattern := fmt.Sprintf("inventory:store:%s:*", storeID.String())
	for {
		var err error
		keys, cursor, err = cache.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		if len(keys) > 0 {
			cache.Client.Del(ctx, keys...)
		}
		if cursor == 0 {
			break
		}
	}

	// Also delete general inventory cache (without store filter)
	cursor = 0
	pattern = "inventory:page:*"
	for {
		var err error
		keys, cursor, err = cache.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		if len(keys) > 0 {
			cache.Client.Del(ctx, keys...)
		}
		if cursor == 0 {
			break
		}
	}

	// Delete inventory cache without store prefix
	cursor = 0
	pattern = "inventory:*"
	for {
		var err error
		keys, cursor, err = cache.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		// Filter out store-specific keys we already deleted
		filteredKeys := []string{}
		for _, key := range keys {
			if !strings.Contains(key, fmt.Sprintf("store:%s", storeID.String())) {
				filteredKeys = append(filteredKeys, key)
			}
		}
		if len(filteredKeys) > 0 {
			cache.Client.Del(ctx, filteredKeys...)
		}
		if cursor == 0 {
			break
		}
	}
}
