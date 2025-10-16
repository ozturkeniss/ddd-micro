package grpc

import (
	userpb "github.com/ddd-micro/api/proto/user"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// ProviderSet is the Wire provider set for gRPC interface layer
var ProviderSet = wire.NewSet(
	NewUserServer,
	NewAuthInterceptor,
	ProvideGRPCServer,
)

// ProvideGRPCServer provides configured gRPC server
func ProvideGRPCServer(userServer *UserServer, authInterceptor *AuthInterceptor) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// Register user service
	userpb.RegisterUserServiceServer(grpcServer, userServer)

	return grpcServer
}
