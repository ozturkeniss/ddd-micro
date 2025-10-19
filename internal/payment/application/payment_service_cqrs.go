package application

import (
	"context"
	"fmt"

	"github.com/ddd-micro/internal/payment/application/command"
	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/application/query"
	"github.com/ddd-micro/internal/payment/domain"
	"github.com/ddd-micro/internal/payment/infrastructure/client"
	"github.com/ddd-micro/internal/payment/infrastructure/kafka"
	"github.com/ddd-micro/kafka"
)

// PaymentServiceCQRS represents the main payment service using CQRS pattern
type PaymentServiceCQRS struct {
	// Command handlers
	createPaymentHandler       *command.CreatePaymentCommandHandler
	processPaymentHandler      *command.ProcessPaymentCommandHandler
	cancelPaymentHandler       *command.CancelPaymentCommandHandler
	addPaymentMethodHandler    *command.AddPaymentMethodCommandHandler
	updatePaymentMethodHandler *command.UpdatePaymentMethodCommandHandler
	deletePaymentMethodHandler *command.DeletePaymentMethodCommandHandler

	// Query handlers
	getPaymentHandler         *query.GetPaymentQueryHandler
	listPaymentsHandler       *query.ListPaymentsQueryHandler
	getPaymentMethodHandler   *query.GetPaymentMethodQueryHandler
	listPaymentMethodsHandler *query.ListPaymentMethodsQueryHandler

	// Repositories
	paymentRepo       domain.PaymentRepository
	paymentMethodRepo domain.PaymentMethodRepository

	// External service clients
	userClient    client.UserClient
	productClient client.ProductClient
	basketClient  client.BasketClient

	// Kafka event publisher
	eventPublisher *kafka.PaymentEventPublisher
}

// NewPaymentServiceCQRS creates a new PaymentServiceCQRS
func NewPaymentServiceCQRS(
	createPaymentHandler *command.CreatePaymentCommandHandler,
	processPaymentHandler *command.ProcessPaymentCommandHandler,
	cancelPaymentHandler *command.CancelPaymentCommandHandler,
	addPaymentMethodHandler *command.AddPaymentMethodCommandHandler,
	updatePaymentMethodHandler *command.UpdatePaymentMethodCommandHandler,
	deletePaymentMethodHandler *command.DeletePaymentMethodCommandHandler,
	getPaymentHandler *query.GetPaymentQueryHandler,
	listPaymentsHandler *query.ListPaymentsQueryHandler,
	getPaymentMethodHandler *query.GetPaymentMethodQueryHandler,
	listPaymentMethodsHandler *query.ListPaymentMethodsQueryHandler,
	paymentRepo domain.PaymentRepository,
	paymentMethodRepo domain.PaymentMethodRepository,
	userClient client.UserClient,
	productClient client.ProductClient,
	basketClient client.BasketClient,
	eventPublisher *kafka.PaymentEventPublisher,
) *PaymentServiceCQRS {
	return &PaymentServiceCQRS{
		createPaymentHandler:       createPaymentHandler,
		processPaymentHandler:      processPaymentHandler,
		cancelPaymentHandler:       cancelPaymentHandler,
		addPaymentMethodHandler:    addPaymentMethodHandler,
		updatePaymentMethodHandler: updatePaymentMethodHandler,
		deletePaymentMethodHandler: deletePaymentMethodHandler,
		getPaymentHandler:          getPaymentHandler,
		listPaymentsHandler:        listPaymentsHandler,
		getPaymentMethodHandler:    getPaymentMethodHandler,
		listPaymentMethodsHandler:  listPaymentMethodsHandler,
		paymentRepo:                paymentRepo,
		paymentMethodRepo:          paymentMethodRepo,
		userClient:                 userClient,
		productClient:              productClient,
		basketClient:               basketClient,
		eventPublisher:             eventPublisher,
	}
}

// Payment operations

// CreatePayment creates a new payment
func (s *PaymentServiceCQRS) CreatePayment(ctx context.Context, userID uint, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Validate payment based on type
	if req.BasketID != nil {
		// Basket-based payment: validate basket
		basket, err := s.basketClient.ValidateBasket(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("basket validation failed: %w", err)
		}

		// Calculate total amount from basket
		var totalAmount float64
		for _, item := range basket.Items {
			totalAmount += float64(item.Quantity) * item.UnitPrice
		}

		// Validate amount matches basket total
		if req.Amount != totalAmount {
			return nil, fmt.Errorf("payment amount does not match basket total")
		}

		// Note: In a real implementation, you would reserve items here

	} else if req.ProductID != nil && req.Quantity != nil {
		// Direct product payment: validate product
		product, err := s.productClient.GetProduct(ctx, *req.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product validation failed: %w", err)
		}

		// Validate quantity
		if *req.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity")
		}

		// Calculate total amount
		totalAmount := float64(*req.Quantity) * product.Price

		// Validate amount matches product total
		if req.Amount != totalAmount {
			return nil, fmt.Errorf("payment amount does not match product total")
		}

		// Check stock availability
		if product.Stock < int32(*req.Quantity) {
			return nil, fmt.Errorf("insufficient stock")
		}

	} else {
		return nil, fmt.Errorf("either basket_id or product_id with quantity must be provided")
	}

	cmd := command.CreatePaymentCommand{
		UserID:          userID,
		OrderID:         req.OrderID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentMethod:   req.PaymentMethod,
		PaymentMethodID: req.PaymentMethodID,
		ReturnURL:       req.ReturnURL,
		CancelURL:       req.CancelURL,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		BasketID:        req.BasketID,
	}

	return s.createPaymentHandler.Handle(ctx, cmd)
}

// GetPayment gets a payment by ID
func (s *PaymentServiceCQRS) GetPayment(ctx context.Context, userID uint, paymentID string) (*dto.PaymentResponse, error) {
	query := query.GetPaymentQuery{
		UserID:    userID,
		PaymentID: paymentID,
	}

	return s.getPaymentHandler.Handle(ctx, query)
}

// ProcessPayment processes a payment
func (s *PaymentServiceCQRS) ProcessPayment(ctx context.Context, userID uint, paymentID string, req dto.ProcessPaymentRequest) (*dto.PaymentResponse, error) {
	// Get payment to check type
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the payment
	if payment.UserID != userID {
		return nil, domain.ErrPaymentNotFound
	}

	cmd := command.ProcessPaymentCommand{
		PaymentID:        paymentID,
		PaymentMethodID:  req.PaymentMethodID,
		ConfirmationData: req.ConfirmationData,
	}

	// Process payment
	paymentResp, err := s.processPaymentHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// If payment is successful, publish events for stock update and basket clearing
	if paymentResp.Status == "completed" {
		// Convert payment items for Kafka events
		var items []kafka.PaymentItem
		if payment.ProductID != nil && payment.Quantity != nil {
			// Direct product purchase
			items = []kafka.PaymentItem{
				{
					ProductID:  *payment.ProductID,
					Quantity:   *payment.Quantity,
					UnitPrice:  payment.Amount / float64(*payment.Quantity),
					TotalPrice: payment.Amount,
				},
			}
		} else if payment.BasketID != nil {
			// Basket-based purchase - get items from basket
			basket, err := s.basketClient.GetBasket(ctx, userID)
			if err == nil {
				for _, item := range basket.Items {
					items = append(items, kafka.PaymentItem{
						ProductID:  uint(item.ProductId),
						Quantity:   int(item.Quantity),
						UnitPrice:  item.UnitPrice,
						TotalPrice: float64(item.Quantity) * item.UnitPrice,
					})
				}
			}
		}

		// Publish payment completed event
		if err := s.eventPublisher.PublishPaymentCompleted(ctx, paymentID, userID, payment.OrderID,
			payment.Amount, payment.Currency, string(payment.PaymentMethod), items, payment.BasketID); err != nil {
			// Log error but don't fail the payment
			// In production, you might want to implement compensation logic
		}
	}

	return paymentResp, nil
}

// CancelPayment cancels a payment
func (s *PaymentServiceCQRS) CancelPayment(ctx context.Context, userID uint, paymentID string) (*dto.PaymentResponse, error) {
	// Get payment to check type
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the payment
	if payment.UserID != userID {
		return nil, domain.ErrPaymentNotFound
	}

	cmd := command.CancelPaymentCommand{
		PaymentID: paymentID,
	}

	// Cancel payment
	paymentResp, err := s.cancelPaymentHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Release reservations if payment was cancelled
	if paymentResp.Status == "cancelled" {
		// Note: In a real implementation, you would release reservations here
	}

	return paymentResp, nil
}

// ListPayments lists user's payments
func (s *PaymentServiceCQRS) ListPayments(ctx context.Context, req dto.ListPaymentsRequest) (*dto.PaymentListResponse, error) {
	query := query.ListPaymentsQuery{
		UserID: req.UserID,
		Page:   req.Page,
		Limit:  req.Limit,
		Status: req.Status,
	}

	return s.listPaymentsHandler.Handle(ctx, query)
}

// Payment method operations

// GetPaymentMethods gets user's payment methods
func (s *PaymentServiceCQRS) GetPaymentMethods(ctx context.Context, userID uint) (*dto.PaymentMethodListResponse, error) {
	query := query.ListPaymentMethodsQuery{
		UserID: userID,
	}

	return s.listPaymentMethodsHandler.Handle(ctx, query)
}

// AddPaymentMethod adds a new payment method
func (s *PaymentServiceCQRS) AddPaymentMethod(ctx context.Context, userID uint, req dto.AddPaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	cmd := command.AddPaymentMethodCommand{
		UserID:    userID,
		Type:      req.Type,
		Provider:  req.Provider,
		Token:     req.Token,
		IsDefault: req.IsDefault,
	}

	return s.addPaymentMethodHandler.Handle(ctx, cmd)
}

// UpdatePaymentMethod updates a payment method
func (s *PaymentServiceCQRS) UpdatePaymentMethod(ctx context.Context, userID uint, paymentMethodID string, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	cmd := command.UpdatePaymentMethodCommand{
		PaymentMethodID: paymentMethodID,
		IsDefault:       req.IsDefault,
		IsActive:        req.IsActive,
	}

	return s.updatePaymentMethodHandler.Handle(ctx, cmd)
}

// DeletePaymentMethod deletes a payment method
func (s *PaymentServiceCQRS) DeletePaymentMethod(ctx context.Context, userID uint, paymentMethodID string) error {
	cmd := command.DeletePaymentMethodCommand{
		PaymentMethodID: paymentMethodID,
	}

	return s.deletePaymentMethodHandler.Handle(ctx, cmd)
}

// SetDefaultPaymentMethod sets a payment method as default
func (s *PaymentServiceCQRS) SetDefaultPaymentMethod(ctx context.Context, userID uint, paymentMethodID string) (*dto.PaymentMethodResponse, error) {
	// First, get the payment method to verify ownership
	paymentMethod, err := s.paymentMethodRepo.GetByID(ctx, paymentMethodID)
	if err != nil {
		return nil, err
	}

	if paymentMethod.UserID != userID {
		return nil, domain.ErrPaymentMethodNotFound
	}

	// Update all payment methods for this user to not be default
	if err := s.paymentMethodRepo.SetAllNonDefault(ctx, userID); err != nil {
		return nil, err
	}

	// Set this payment method as default
	paymentMethod.IsDefault = true
	if err := s.paymentMethodRepo.Update(ctx, paymentMethod); err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.PaymentMethodResponse{
		ID:             paymentMethod.ID,
		UserID:         paymentMethod.UserID,
		Type:           paymentMethod.Type,
		Provider:       paymentMethod.Provider,
		LastFourDigits: paymentMethod.LastFourDigits,
		ExpiryMonth:    paymentMethod.ExpiryMonth,
		ExpiryYear:     paymentMethod.ExpiryYear,
		IsDefault:      paymentMethod.IsDefault,
		IsActive:       paymentMethod.IsActive,
		CreatedAt:      paymentMethod.CreatedAt,
		UpdatedAt:      paymentMethod.UpdatedAt,
	}, nil
}

// Admin operations (placeholder implementations)

// AdminListPayments lists all payments (admin only)
func (s *PaymentServiceCQRS) AdminListPayments(ctx context.Context, req dto.AdminListPaymentsRequest) (*dto.PaymentListResponse, error) {
	// TODO: Implement admin list payments
	return &dto.PaymentListResponse{
		Payments:   []dto.PaymentResponse{},
		Total:      0,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: 0,
		HasNext:    false,
		HasPrev:    false,
	}, nil
}

// AdminGetPayment gets any payment by ID (admin only)
func (s *PaymentServiceCQRS) AdminGetPayment(ctx context.Context, paymentID string) (*dto.PaymentResponse, error) {
	// TODO: Implement admin get payment
	return nil, domain.ErrPaymentNotFound
}

// AdminUpdatePaymentStatus updates payment status (admin only)
func (s *PaymentServiceCQRS) AdminUpdatePaymentStatus(ctx context.Context, paymentID string, req dto.UpdatePaymentStatusRequest) (*dto.PaymentResponse, error) {
	// TODO: Implement admin update payment status
	return nil, domain.ErrPaymentNotFound
}

// CreateRefund creates a refund (admin only)
func (s *PaymentServiceCQRS) CreateRefund(ctx context.Context, req dto.CreateRefundRequest) (*dto.RefundResponse, error) {
	// TODO: Implement create refund
	return nil, domain.ErrRefundNotFound
}

// AdminListRefunds lists all refunds (admin only)
func (s *PaymentServiceCQRS) AdminListRefunds(ctx context.Context, req dto.AdminListRefundsRequest) (*dto.RefundListResponse, error) {
	// TODO: Implement admin list refunds
	return &dto.RefundListResponse{
		Refunds: []dto.RefundResponse{},
		Total:   0,
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}

// AdminGetRefund gets refund by ID (admin only)
func (s *PaymentServiceCQRS) AdminGetRefund(ctx context.Context, refundID string) (*dto.RefundResponse, error) {
	// TODO: Implement admin get refund
	return nil, domain.ErrRefundNotFound
}

// ProcessRefund processes a refund (admin only)
func (s *PaymentServiceCQRS) ProcessRefund(ctx context.Context, refundID string) (*dto.RefundResponse, error) {
	// TODO: Implement process refund
	return nil, domain.ErrRefundNotFound
}

// GetPaymentStats gets payment statistics (admin only)
func (s *PaymentServiceCQRS) GetPaymentStats(ctx context.Context, period string) (*dto.PaymentStatsResponse, error) {
	// TODO: Implement get payment stats
	return &dto.PaymentStatsResponse{
		TotalPayments:      0,
		TotalAmount:        0,
		SuccessfulPayments: 0,
		FailedPayments:     0,
		PendingPayments:    0,
		RefundedAmount:     0,
		AverageAmount:      0,
	}, nil
}
