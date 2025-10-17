package grpc

import (
	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
)

// Providers contains all gRPC-related dependencies
type Providers struct {
	BasketServer    *BasketServer
	AuthInterceptor *AuthInterceptor
}

// NewProviders creates new gRPC providers
func NewProviders(
	basketService *application.BasketServiceCQRS,
	userClient *client.UserClient,
) *Providers {
	basketServer := NewBasketServer(basketService)
	authInterceptor := NewAuthInterceptor(userClient)

	return &Providers{
		BasketServer:    basketServer,
		AuthInterceptor: authInterceptor,
	}
}
