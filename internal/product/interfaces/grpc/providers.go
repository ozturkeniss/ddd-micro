package grpc

import (
	productpb "github.com/ddd-micro/api/proto/product"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// ProviderSet is the Wire provider set for gRPC interface layer
var ProviderSet = wire.NewSet(
	NewProductServer,
	NewAuthInterceptor,
	ProvideGRPCServer,
)

// ProvideGRPCServer provides configured gRPC server
func ProvideGRPCServer(productServer *ProductServer, authInterceptor *AuthInterceptor) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// Register product service
	productpb.RegisterProductServiceServer(grpcServer, productServer)

	return grpcServer
}
