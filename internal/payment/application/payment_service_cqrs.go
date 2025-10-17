package application

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/command"
	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/application/query"
	"github.com/ddd-micro/internal/payment/domain"
)

// PaymentServiceCQRS represents the main payment service using CQRS pattern
type PaymentServiceCQRS struct {
	// Command handlers
	createPaymentHandler      *command.CreatePaymentCommandHandler
	processPaymentHandler     *command.ProcessPaymentCommandHandler
	cancelPaymentHandler      *command.CancelPaymentCommandHandler
	addPaymentMethodHandler   *command.AddPaymentMethodCommandHandler
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
) *PaymentServiceCQRS {
	return &PaymentServiceCQRS{
		createPaymentHandler:      createPaymentHandler,
		processPaymentHandler:     processPaymentHandler,
		cancelPaymentHandler:      cancelPaymentHandler,
		addPaymentMethodHandler:   addPaymentMethodHandler,
		updatePaymentMethodHandler: updatePaymentMethodHandler,
		deletePaymentMethodHandler: deletePaymentMethodHandler,
		getPaymentHandler:         getPaymentHandler,
		listPaymentsHandler:       listPaymentsHandler,
		getPaymentMethodHandler:   getPaymentMethodHandler,
		listPaymentMethodsHandler: listPaymentMethodsHandler,
		paymentRepo:               paymentRepo,
		paymentMethodRepo:         paymentMethodRepo,
	}
}

// Payment operations

// CreatePayment creates a new payment
func (s *PaymentServiceCQRS) CreatePayment(ctx context.Context, userID uint, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	cmd := command.CreatePaymentCommand{
		UserID:          userID,
		OrderID:         req.OrderID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentMethod:   req.PaymentMethod,
		PaymentMethodID: req.PaymentMethodID,
		ReturnURL:       req.ReturnURL,
		CancelURL:       req.CancelURL,
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
	cmd := command.ProcessPaymentCommand{
		PaymentID:       paymentID,
		PaymentMethodID: req.PaymentMethodID,
		ConfirmationData: req.ConfirmationData,
	}

	return s.processPaymentHandler.Handle(ctx, cmd)
}

// CancelPayment cancels a payment
func (s *PaymentServiceCQRS) CancelPayment(ctx context.Context, userID uint, paymentID string) (*dto.PaymentResponse, error) {
	cmd := command.CancelPaymentCommand{
		PaymentID: paymentID,
	}

	return s.cancelPaymentHandler.Handle(ctx, cmd)
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