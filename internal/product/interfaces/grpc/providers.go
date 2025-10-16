package grpc

import (
	productpb "github.com/ddd-micro/api/proto/product"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// ProviderSet is the Wire provider set for gRPC interface layer
var ProviderSet = wire.NewSet(
	NewProductServer,
	ProvideGRPCServer,
)

// ProvideGRPCServer provides configured gRPC server
func ProvideGRPCServer(productServer *ProductServer) *grpc.Server {
	grpcServer := grpc.NewServer()

	// Register product service
	productpb.RegisterProductServiceServer(grpcServer, productServer)

	return grpcServer
}
