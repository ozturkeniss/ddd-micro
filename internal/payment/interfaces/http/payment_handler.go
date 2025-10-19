package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ddd-micro/internal/payment/application"
	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *application.PaymentServiceCQRS
	metrics        *monitoring.PrometheusMetrics
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *application.PaymentServiceCQRS, metrics *monitoring.PrometheusMetrics) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		metrics:        metrics,
	}
}

// CreatePayment creates a new payment
// @Summary Create a new payment
// @Description Create a new payment for the authenticated user
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePaymentRequest true "Payment creation request"
// @Success 201 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	// Start tracing span
	span, _ := monitoring.StartSpanFromGinContext(c, "payment.create")
	defer span.Finish()

	userID := c.GetUint("user_id")

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		monitoring.LogSpanError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	payment, err := h.paymentService.CreatePayment(c.Request.Context(), userID, req)
	duration := time.Since(start)

	// Record database query duration
	h.metrics.RecordDatabaseQueryDuration("create_payment", "payments", duration)

	if err != nil {
		monitoring.LogSpanError(span, err)
		h.metrics.RecordPaymentFailure()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful payment creation
	h.metrics.RecordPaymentCreation()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"payment.id":    payment.ID,
		"payment.amount": payment.Amount,
		"success":       true,
	})

	c.JSON(http.StatusCreated, payment)
}

// GetPayment gets a payment by ID
// @Summary Get payment by ID
// @Description Get payment details by ID for the authenticated user
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID := c.GetUint("user_id")
	paymentID := c.Param("id")

	payment, err := h.paymentService.GetPayment(c.Request.Context(), userID, paymentID)
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

// ProcessPayment processes a payment
// @Summary Process a payment
// @Description Process a payment with payment method confirmation
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Param request body dto.ProcessPaymentRequest true "Payment processing request"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id}/process [post]
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	// Start tracing span
	span, _ := monitoring.StartSpanFromGinContext(c, "payment.process")
	defer span.Finish()

	userID := c.GetUint("user_id")
	paymentID := c.Param("id")

	var req dto.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		monitoring.LogSpanError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	payment, err := h.paymentService.ProcessPayment(c.Request.Context(), userID, paymentID, req)
	duration := time.Since(start)

	// Record payment processing duration
	h.metrics.RecordPaymentProcessingDuration("credit_card", "processing", duration)

	if err != nil {
		monitoring.LogSpanError(span, err)
		h.metrics.RecordPaymentFailure()
		if err.Error() == "payment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful payment processing
	h.metrics.RecordPaymentCompletion()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"payment.id":    payment.ID,
		"payment.status": payment.Status,
		"success":       true,
	})

	c.JSON(http.StatusOK, payment)
}

// CancelPayment cancels a payment
// @Summary Cancel a payment
// @Description Cancel a pending payment
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id}/cancel [post]
func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	userID := c.GetUint("user_id")
	paymentID := c.Param("id")

	payment, err := h.paymentService.CancelPayment(c.Request.Context(), userID, paymentID)
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

// ListPayments lists user's payments
// @Summary List user payments
// @Description Get a list of payments for the authenticated user
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.PaymentListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments [get]
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	req := dto.ListPaymentsRequest{
		UserID: userID,
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	payments, err := h.paymentService.ListPayments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetPaymentMethods gets user's payment methods
// @Summary Get user payment methods
// @Description Get a list of payment methods for the authenticated user
// @Tags payment-methods
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PaymentMethodListResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods [get]
func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	userID := c.GetUint("user_id")

	paymentMethods, err := h.paymentService.GetPaymentMethods(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethods)
}

// AddPaymentMethod adds a new payment method
// @Summary Add payment method
// @Description Add a new payment method for the authenticated user
// @Tags payment-methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddPaymentMethodRequest true "Payment method request"
// @Success 201 {object} dto.PaymentMethodResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods [post]
func (h *PaymentHandler) AddPaymentMethod(c *gin.Context) {
	// Start tracing span
	span, _ := monitoring.StartSpanFromGinContext(c, "payment_method.add")
	defer span.Finish()

	userID := c.GetUint("user_id")

	var req dto.AddPaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		monitoring.LogSpanError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	paymentMethod, err := h.paymentService.AddPaymentMethod(c.Request.Context(), userID, req)
	duration := time.Since(start)

	// Record database query duration
	h.metrics.RecordDatabaseQueryDuration("create_payment_method", "payment_methods", duration)

	if err != nil {
		monitoring.LogSpanError(span, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record successful payment method addition
	h.metrics.RecordPaymentMethodAddition()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"payment_method.id":   paymentMethod.ID,
		"payment_method.type": paymentMethod.Type,
		"success":             true,
	})

	c.JSON(http.StatusCreated, paymentMethod)
}

// UpdatePaymentMethod updates a payment method
// @Summary Update payment method
// @Description Update an existing payment method
// @Tags payment-methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment Method ID"
// @Param request body dto.UpdatePaymentMethodRequest true "Payment method update request"
// @Success 200 {object} dto.PaymentMethodResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods/{id} [put]
func (h *PaymentHandler) UpdatePaymentMethod(c *gin.Context) {
	userID := c.GetUint("user_id")
	paymentMethodID := c.Param("id")

	var req dto.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentMethod, err := h.paymentService.UpdatePaymentMethod(c.Request.Context(), userID, paymentMethodID, req)
	if err != nil {
		if err.Error() == "payment method not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

// DeletePaymentMethod deletes a payment method
// @Summary Delete payment method
// @Description Delete an existing payment method
// @Tags payment-methods
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment Method ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods/{id} [delete]
func (h *PaymentHandler) DeletePaymentMethod(c *gin.Context) {
	userID := c.GetUint("user_id")
	paymentMethodID := c.Param("id")

	err := h.paymentService.DeletePaymentMethod(c.Request.Context(), userID, paymentMethodID)
	if err != nil {
		if err.Error() == "payment method not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetDefaultPaymentMethod sets a payment method as default
// @Summary Set default payment method
// @Description Set a payment method as the default for the user
// @Tags payment-methods
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment Method ID"
// @Success 200 {object} dto.PaymentMethodResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment-methods/{id}/set-default [post]
func (h *PaymentHandler) SetDefaultPaymentMethod(c *gin.Context) {
	userID := c.GetUint("user_id")
	paymentMethodID := c.Param("id")

	paymentMethod, err := h.paymentService.SetDefaultPaymentMethod(c.Request.Context(), userID, paymentMethodID)
	if err != nil {
		if err.Error() == "payment method not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}
