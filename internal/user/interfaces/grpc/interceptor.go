package grpc

import (
	"context"
	"strings"

	"github.com/ddd-micro/internal/user/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a gRPC interceptor for JWT authentication
type AuthInterceptor struct {
	userService *application.UserServiceCQRS
	// Methods that don't require authentication
	publicMethods map[string]bool
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(userService *application.UserServiceCQRS) *AuthInterceptor {
	publicMethods := map[string]bool{
		"/user.UserService/Register":     true,
		"/user.UserService/Login":        true,
		"/user.UserService/RefreshToken": true,
	}

	return &AuthInterceptor{
		userService:   userService,
		publicMethods: publicMethods,
	}
}

// Unary returns a server interceptor function for unary RPCs
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if method is public
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract and validate token
		newCtx, err := i.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// authorize validates the JWT token and adds user info to context
func (i *AuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(accessToken, "Bearer ") {
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	}

	// Validate token
	claims, err := i.userService.ValidateToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token: %v", err)
	}

	// Add user info to context
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_email", claims.Email)
	ctx = context.WithValue(ctx, "user_role", string(claims.Role))

	return ctx, nil
}
