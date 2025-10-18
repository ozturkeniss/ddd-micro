package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
)

// BasketHandler handles HTTP requests for basket operations
type BasketHandler struct {
	basketService *application.BasketServiceCQRS
	metrics       *monitoring.PrometheusMetrics
}

// NewBasketHandler creates a new basket handler
func NewBasketHandler(basketService *application.BasketServiceCQRS, metrics *monitoring.PrometheusMetrics) *BasketHandler {
	return &BasketHandler{
		basketService: basketService,
		metrics:       metrics,
	}
}

// CreateBasket creates a new basket for a user
// @Summary Create a new basket
// @Description Creates a new basket for the authenticated user
// @Tags basket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.BasketResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/basket [post]
func (h *BasketHandler) CreateBasket(c *gin.Context) {
	// Start tracing span
	span, _ := monitoring.StartSpanFromGinContext(c, "basket.create")
	defer span.Finish()

	var req dto.CreateBasketRequest

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	req.UserID = userID.(uint)

	start := time.Now()
	basket, err := h.basketService.CreateBasket(c.Request.Context(), req)
	duration := time.Since(start)

	// Record Redis operation duration
	h.metrics.RecordRedisOperationDuration("create_basket", duration)

	if err != nil {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	// Record successful basket creation
	h.metrics.RecordBasketCreation()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"user.id":      userID.(uint),
		"basket.id":    basket.ID,
		"operation":    "create_basket",
		"success":      true,
	})

	c.JSON(http.StatusCreated, basket)
}

// GetBasket retrieves the user's basket
// @Summary Get user's basket
// @Description Retrieves the basket for the authenticated user
// @Tags basket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.BasketResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/basket [get]
func (h *BasketHandler) GetBasket(c *gin.Context) {
	// Start tracing span
	span, _ := monitoring.StartSpanFromGinContext(c, "basket.get")
	defer span.Finish()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	start := time.Now()
	basket, err := h.basketService.GetBasketHTTP(c.Request.Context(), userID.(uint))
	duration := time.Since(start)

	// Record Redis operation duration
	h.metrics.RecordRedisOperationDuration("get_basket", duration)

	if err != nil {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		if err.Error() == "basket not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Basket not found for user",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	// Record successful basket retrieval
	h.metrics.RecordBasketRetrieval()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"user.id":      userID.(uint),
		"basket.id":    basket.ID,
		"operation":    "get_basket",
		"success":      true,
	})

	c.JSON(http.StatusOK, basket)
}

// AddItem adds an item to the user's basket
// @Summary Add item to basket
// @Description Adds a product to the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddItemRequest true "Add item request"
// @Success 200 {object} dto.BasketResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/basket/items [post]
func (h *BasketHandler) AddItem(c *gin.Context) {
	// Start tracing span
	span, _ := monitoring.StartSpanFromGinContext(c, "basket.add_item")
	defer span.Finish()

	var req dto.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	req.UserID = userID.(uint)

	start := time.Now()
	basket, err := h.basketService.AddItemHTTP(c.Request.Context(), req.UserID, req)
	duration := time.Since(start)

	// Record Redis operation duration
	h.metrics.RecordRedisOperationDuration("add_item", duration)

	if err != nil {
		monitoring.LogSpanEvent(span, "User ID not found in context")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	// Record successful item addition
	h.metrics.RecordItemAddition()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"user.id":      userID.(uint),
		"basket.id":    basket.ID,
		"product.id":   req.ProductID,
		"quantity":     req.Quantity,
		"operation":    "add_item",
		"success":      true,
	})

	c.JSON(http.StatusOK, basket)
}

// UpdateItem updates an item in the user's basket
// @Summary Update basket item
// @Description Updates the quantity of an item in the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Param request body dto.UpdateItemRequest true "Update item request"
// @Success 200 {object} dto.BasketResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/basket/items/{product_id} [put]
func (h *BasketHandler) UpdateItem(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	req.UserID = userID.(uint)

	basket, err := h.basketService.UpdateItemHTTP(c.Request.Context(), userID.(uint), uint(productID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, basket)
}

// RemoveItem removes an item from the user's basket
// @Summary Remove basket item
// @Description Removes an item from the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Success 200 {object} dto.BasketResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/basket/items/{product_id} [delete]
func (h *BasketHandler) RemoveItem(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	basket, err := h.basketService.RemoveItemHTTP(c.Request.Context(), userID.(uint), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, basket)
}

// ClearBasket clears all items from the user's basket
// @Summary Clear basket
// @Description Removes all items from the user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/basket/clear [delete]
func (h *BasketHandler) ClearBasket(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	basket, err := h.basketService.ClearBasketHTTP(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Basket cleared successfully",
		Data:    basket,
	})
}

// AdminGetBasket retrieves any user's basket (admin only)
// @Summary Get user basket (Admin)
// @Description Retrieves the basket for any user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} dto.BasketResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/baskets/{user_id} [get]
func (h *BasketHandler) AdminGetBasket(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid user ID",
		})
		return
	}

	basket, err := h.basketService.GetBasketHTTP(c.Request.Context(), uint(userID))
	if err != nil {
		if err.Error() == "basket not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Basket not found for user",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, basket)
}

// AdminDeleteBasket deletes a user's basket (admin only)
// @Summary Delete user basket (Admin)
// @Description Deletes the basket for any user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/baskets/{user_id} [delete]
func (h *BasketHandler) AdminDeleteBasket(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid user ID",
		})
		return
	}

	err = h.basketService.DeleteBasket(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Basket deleted successfully",
	})
}

// AdminCleanupExpiredBaskets cleans up expired baskets (admin only)
// @Summary Cleanup expired baskets (Admin)
// @Description Removes all expired baskets from the system (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/baskets/cleanup [post]
func (h *BasketHandler) AdminCleanupExpiredBaskets(c *gin.Context) {
	err := h.basketService.CleanupExpiredBaskets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Expired baskets cleaned up successfully",
	})
}
