package grpc

import (
	basketpb "github.com/ddd-micro/api/proto/basket"
	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Providers contains all gRPC-related dependencies
type Providers struct {
	BasketServer    *BasketServer
	AuthInterceptor *AuthInterceptor
	GRPCServer      *grpc.Server
}

// ProviderSet is the Wire provider set for gRPC
var ProviderSet = wire.NewSet(
	NewBasketServer,
	NewAuthInterceptor,
	NewGRPCServer,
	NewProviders,
)

// NewProviders creates new gRPC providers
func NewProviders(
	basketServer *BasketServer,
	authInterceptor *AuthInterceptor,
	grpcServer *grpc.Server,
) *Providers {
	return &Providers{
		BasketServer:    basketServer,
		AuthInterceptor: authInterceptor,
		GRPCServer:      grpcServer,
	}
}

// NewGRPCServer creates a new gRPC server with interceptors
func NewGRPCServer(
	basketServer *BasketServer,
	authInterceptor *AuthInterceptor,
) *grpc.Server {
	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryAuthInterceptor()),
	)

	// Register basket service
	basketpb.RegisterBasketServiceServer(server, basketServer)

	// Enable reflection for debugging
	reflection.Register(server)

	return server
}
