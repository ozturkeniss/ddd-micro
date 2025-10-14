package grpc

import (
	"github.com/ddd-micro/api/proto/user"
	"github.com/ddd-micro/internal/user/application"
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
	user.RegisterUserServiceServer(grpcServer, userServer)

	return grpcServer
}

