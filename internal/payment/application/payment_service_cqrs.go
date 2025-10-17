package application

import (
	"context"
	"time"

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
	setDefaultPaymentMethodHandler *command.SetDefaultPaymentMethodCommandHandler
	createRefundHandler       *command.CreateRefundCommandHandler
	processRefundHandler      *command.ProcessRefundCommandHandler
	
	// Query handlers
	getPaymentHandler         *query.GetPaymentQueryHandler
	getUserPaymentsHandler    *query.GetUserPaymentsQueryHandler
	getPaymentMethodsHandler  *query.GetPaymentMethodsQueryHandler
	getRefundsHandler         *query.GetRefundsQueryHandler
	getPaymentStatsHandler    *query.GetPaymentStatsQueryHandler
	getRefundStatsHandler     *query.GetRefundStatsQueryHandler
	
	// Repositories
	paymentRepo       domain.PaymentRepository
	paymentMethodRepo domain.PaymentMethodRepository
	refundRepo        domain.RefundRepository
}

// NewPaymentServiceCQRS creates a new PaymentServiceCQRS
func NewPaymentServiceCQRS(
	paymentRepo domain.PaymentRepository,
	paymentMethodRepo domain.PaymentMethodRepository,
	refundRepo domain.RefundRepository,
	paymentGateway domain.PaymentGateway,
) *PaymentServiceCQRS {
	return &PaymentServiceCQRS{
		createPaymentHandler:      command.NewCreatePaymentCommandHandler(paymentRepo, paymentGateway),
		processPaymentHandler:     command.NewProcessPaymentCommandHandler(paymentRepo, paymentGateway),
		cancelPaymentHandler:      command.NewCancelPaymentCommandHandler(paymentRepo, paymentGateway),
		addPaymentMethodHandler:   command.NewAddPaymentMethodCommandHandler(paymentMethodRepo, paymentGateway),
		updatePaymentMethodHandler: command.NewUpdatePaymentMethodCommandHandler(paymentMethodRepo),
		deletePaymentMethodHandler: command.NewDeletePaymentMethodCommandHandler(paymentMethodRepo),
		setDefaultPaymentMethodHandler: command.NewSetDefaultPaymentMethodCommandHandler(paymentMethodRepo),
		createRefundHandler:       command.NewCreateRefundCommandHandler(paymentRepo, refundRepo, paymentGateway),
		processRefundHandler:      command.NewProcessRefundCommandHandler(refundRepo, paymentGateway),
		getPaymentHandler:         query.NewGetPaymentQueryHandler(paymentRepo),
		getUserPaymentsHandler:    query.NewGetUserPaymentsQueryHandler(paymentRepo),
		getPaymentMethodsHandler:  query.NewGetPaymentMethodsQueryHandler(paymentMethodRepo),
		getRefundsHandler:         query.NewGetRefundsQueryHandler(refundRepo),
		getPaymentStatsHandler:    query.NewGetPaymentStatsQueryHandler(paymentRepo),
		getRefundStatsHandler:     query.NewGetRefundStatsQueryHandler(refundRepo),
		paymentRepo:               paymentRepo,
		paymentMethodRepo:         paymentMethodRepo,
		refundRepo:                refundRepo,
	}
}

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

// ProcessPayment processes a payment
func (s *PaymentServiceCQRS) ProcessPayment(ctx context.Context, paymentID string, req dto.ProcessPaymentRequest) (*dto.PaymentResponse, error) {
	cmd := command.ProcessPaymentCommand{
		PaymentID:       paymentID,
		PaymentMethodID: req.PaymentMethodID,
		ConfirmationData: req.ConfirmationData,
	}
	
	return s.processPaymentHandler.Handle(ctx, cmd)
}

// GetPayment retrieves a payment by ID
func (s *PaymentServiceCQRS) GetPayment(ctx context.Context, paymentID string) (*dto.PaymentResponse, error) {
	query := query.GetPaymentQuery{
		PaymentID: paymentID,
	}
	
	return s.getPaymentHandler.Handle(ctx, query)
}

// GetUserPayments retrieves payments for a user
func (s *PaymentServiceCQRS) GetUserPayments(ctx context.Context, userID uint, offset, limit int) (*dto.ListPaymentsResponse, error) {
	query := query.GetUserPaymentsQuery{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	}
	
	return s.getUserPaymentsHandler.Handle(ctx, query)
}

// CancelPayment cancels a payment
func (s *PaymentServiceCQRS) CancelPayment(ctx context.Context, paymentID string) (*dto.PaymentResponse, error) {
	cmd := command.CancelPaymentCommand{
		PaymentID: paymentID,
	}
	
	return s.cancelPaymentHandler.Handle(ctx, cmd)
}

// AddPaymentMethod adds a payment method for a user
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

// GetPaymentMethods retrieves payment methods for a user
func (s *PaymentServiceCQRS) GetPaymentMethods(ctx context.Context, userID uint) (*dto.ListPaymentMethodsResponse, error) {
	query := query.GetPaymentMethodsQuery{
		UserID: userID,
	}
	
	return s.getPaymentMethodsHandler.Handle(ctx, query)
}

// UpdatePaymentMethod updates a payment method
func (s *PaymentServiceCQRS) UpdatePaymentMethod(ctx context.Context, paymentMethodID string, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	cmd := command.UpdatePaymentMethodCommand{
		PaymentMethodID: paymentMethodID,
		IsDefault:       req.IsDefault,
		IsActive:        req.IsActive,
	}
	
	return s.updatePaymentMethodHandler.Handle(ctx, cmd)
}

// DeletePaymentMethod deletes a payment method
func (s *PaymentServiceCQRS) DeletePaymentMethod(ctx context.Context, paymentMethodID string) error {
	cmd := command.DeletePaymentMethodCommand{
		PaymentMethodID: paymentMethodID,
	}
	
	return s.deletePaymentMethodHandler.Handle(ctx, cmd)
}

// SetDefaultPaymentMethod sets a payment method as default
func (s *PaymentServiceCQRS) SetDefaultPaymentMethod(ctx context.Context, userID uint, paymentMethodID string) (*dto.PaymentMethodResponse, error) {
	cmd := command.SetDefaultPaymentMethodCommand{
		UserID:          userID,
		PaymentMethodID: paymentMethodID,
	}
	
	return s.setDefaultPaymentMethodHandler.Handle(ctx, cmd)
}

// CreateRefund creates a refund for a payment
func (s *PaymentServiceCQRS) CreateRefund(ctx context.Context, req dto.CreateRefundRequest) (*dto.RefundResponse, error) {
	cmd := command.CreateRefundCommand{
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
	}
	
	return s.createRefundHandler.Handle(ctx, cmd)
}

// GetRefunds retrieves refunds for a payment
func (s *PaymentServiceCQRS) GetRefunds(ctx context.Context, paymentID string) (*dto.ListRefundsResponse, error) {
	query := query.GetRefundsQuery{
		PaymentID: paymentID,
	}
	
	return s.getRefundsHandler.Handle(ctx, query)
}

// ProcessRefund processes a refund
func (s *PaymentServiceCQRS) ProcessRefund(ctx context.Context, refundID string) (*dto.RefundResponse, error) {
	cmd := command.ProcessRefundCommand{
		RefundID: refundID,
	}
	
	return s.processRefundHandler.Handle(ctx, cmd)
}

// GetPaymentStats retrieves payment statistics
func (s *PaymentServiceCQRS) GetPaymentStats(ctx context.Context, userID *uint, startDate, endDate *string) (*dto.PaymentStatsResponse, error) {
	query := query.GetPaymentStatsQuery{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	
	return s.getPaymentStatsHandler.Handle(ctx, query)
}

// GetRefundStats retrieves refund statistics
func (s *PaymentServiceCQRS) GetRefundStats(ctx context.Context, userID *uint, startDate, endDate *string) (*dto.RefundStatsResponse, error) {
	query := query.GetRefundStatsQuery{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}
	
	return s.getRefundStatsHandler.Handle(ctx, query)
}

// ProcessWebhook processes a webhook from payment gateway
func (s *PaymentServiceCQRS) ProcessWebhook(ctx context.Context, provider string, payload []byte, signature string) (*dto.WebhookResponse, error) {
	// This would be implemented to handle webhook events
	// For now, return a success response
	return &dto.WebhookResponse{
		Success: true,
		Message: "Webhook processed successfully",
	}, nil
}

// CleanupExpiredPayments removes expired payments
func (s *PaymentServiceCQRS) CleanupExpiredPayments(ctx context.Context) (int, error) {
	return s.paymentRepo.CleanupExpiredPayments(ctx)
}
