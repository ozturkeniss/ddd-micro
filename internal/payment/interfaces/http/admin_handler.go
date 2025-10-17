package http

import (
	"net/http"
	"strconv"

	"github.com/ddd-micro/internal/payment/application"
	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-related payment HTTP requests
type AdminHandler struct {
	paymentService *application.PaymentServiceCQRS
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(paymentService *application.PaymentServiceCQRS) *AdminHandler {
	return &AdminHandler{
		paymentService: paymentService,
	}
}

// ListAllPayments lists all payments (admin only)
// @Summary List all payments
// @Description Get a list of all payments in the system (admin only)
// @Tags admin-payments
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param user_id query int false "Filter by user ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.PaymentListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/payments [get]
func (h *AdminHandler) ListAllPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userIDStr := c.Query("user_id")
	status := c.Query("status")

	var userID *uint
	if userIDStr != "" {
		if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	req := dto.AdminListPaymentsRequest{
		Page:   page,
		Limit:  limit,
		UserID: userID,
		Status: status,
	}

	payments, err := h.paymentService.AdminListPayments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentByID gets any payment by ID (admin only)
// @Summary Get payment by ID (admin)
// @Description Get payment details by ID (admin only)
// @Tags admin-payments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/payments/{id} [get]
func (h *AdminHandler) GetPaymentByID(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentService.AdminGetPayment(c.Request.Context(), paymentID)
	if err != nil {
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// UpdatePaymentStatus updates payment status (admin only)
// @Summary Update payment status
// @Description Update payment status (admin only)
// @Tags admin-payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Param request body dto.UpdatePaymentStatusRequest true "Status update request"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/payments/{id}/status [put]
func (h *AdminHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")

	var req dto.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.AdminUpdatePaymentStatus(c.Request.Context(), paymentID, req)
	if err != nil {
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// CreateRefund creates a refund (admin only)
// @Summary Create refund
// @Description Create a refund for a payment (admin only)
// @Tags admin-refunds
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRefundRequest true "Refund creation request"
// @Success 201 {object} dto.RefundResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/refunds [post]
func (h *AdminHandler) CreateRefund(c *gin.Context) {
	var req dto.CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refund, err := h.paymentService.CreateRefund(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, refund)
}

// ListRefunds lists all refunds (admin only)
// @Summary List all refunds
// @Description Get a list of all refunds in the system (admin only)
// @Tags admin-refunds
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param payment_id query string false "Filter by payment ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.RefundListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/refunds [get]
func (h *AdminHandler) ListRefunds(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	paymentID := c.Query("payment_id")
	status := c.Query("status")

	req := dto.AdminListRefundsRequest{
		Page:      page,
		Limit:     limit,
		PaymentID: paymentID,
		Status:    status,
	}

	refunds, err := h.paymentService.AdminListRefunds(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, refunds)
}

// GetRefundByID gets refund by ID (admin only)
// @Summary Get refund by ID (admin)
// @Description Get refund details by ID (admin only)
// @Tags admin-refunds
// @Produce json
// @Security BearerAuth
// @Param id path string true "Refund ID"
// @Success 200 {object} dto.RefundResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/refunds/{id} [get]
func (h *AdminHandler) GetRefundByID(c *gin.Context) {
	refundID := c.Param("id")

	refund, err := h.paymentService.AdminGetRefund(c.Request.Context(), refundID)
	if err != nil {
		if err.Error() == "refund not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Refund not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, refund)
}

// ProcessRefund processes a refund (admin only)
// @Summary Process refund
// @Description Process a pending refund (admin only)
// @Tags admin-refunds
// @Produce json
// @Security BearerAuth
// @Param id path string true "Refund ID"
// @Success 200 {object} dto.RefundResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/refunds/{id}/process [post]
func (h *AdminHandler) ProcessRefund(c *gin.Context) {
	refundID := c.Param("id")

	refund, err := h.paymentService.ProcessRefund(c.Request.Context(), refundID)
	if err != nil {
		if err.Error() == "refund not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Refund not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, refund)
}

// GetPaymentStats gets payment statistics (admin only)
// @Summary Get payment statistics
// @Description Get payment statistics and analytics (admin only)
// @Tags admin-analytics
// @Produce json
// @Security BearerAuth
// @Param period query string false "Time period" Enums(daily, weekly, monthly, yearly) default(monthly)
// @Success 200 {object} dto.PaymentStatsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/analytics/payments [get]
func (h *AdminHandler) GetPaymentStats(c *gin.Context) {
	period := c.DefaultQuery("period", "monthly")

	stats, err := h.paymentService.GetPaymentStats(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
