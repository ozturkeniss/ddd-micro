package http

import (
	"net/http"
	"strconv"

	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/gin-gonic/gin"
)

// BasketHandler handles HTTP requests for basket operations
type BasketHandler struct {
	basketService *application.BasketServiceCQRS
}

// NewBasketHandler creates a new basket handler
func NewBasketHandler(basketService *application.BasketServiceCQRS) *BasketHandler {
	return &BasketHandler{
		basketService: basketService,
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
	var req dto.CreateBasketRequest
	
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

	basket, err := h.basketService.CreateBasket(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

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
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	basket, err := h.basketService.GetBasket(c.Request.Context(), userID.(uint))
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
	var req dto.AddItemRequest
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

	basket, err := h.basketService.AddItem(c.Request.Context(), req.UserID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

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

	basket, err := h.basketService.UpdateItem(c.Request.Context(), userID.(uint), uint(productID), req)
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

	basket, err := h.basketService.RemoveItem(c.Request.Context(), userID.(uint), uint(productID))
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

	basket, err := h.basketService.ClearBasket(c.Request.Context(), userID.(uint))
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

	basket, err := h.basketService.GetBasket(c.Request.Context(), uint(userID))
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