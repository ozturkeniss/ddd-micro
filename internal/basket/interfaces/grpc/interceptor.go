package grpc

import (
	"context"
	"strings"

	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor handles authentication for gRPC requests
type AuthInterceptor struct {
	userClient *client.UserClient
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(userClient *client.UserClient) *AuthInterceptor {
	return &AuthInterceptor{
		userClient: userClient,
	}
}

// UnaryAuthInterceptor returns a unary server interceptor for authentication
func (a *AuthInterceptor) UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for certain methods
		if a.shouldSkipAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization header not provided")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
		}

		// Validate token with user service
		user, err := (*a.userClient).ValidateToken(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "user_id", uint(user.Id))
		ctx = context.WithValue(ctx, "user_role", string(user.Role))
		ctx = context.WithValue(ctx, "user_email", user.Email)

		return handler(ctx, req)
	}
}

// shouldSkipAuth determines if authentication should be skipped for a method
func (a *AuthInterceptor) shouldSkipAuth(method string) bool {
	// Add methods that don't require authentication
	// For now, all basket methods require authentication
	return false
}

// AdminAuthInterceptor returns a unary server interceptor for admin authentication
func (a *AuthInterceptor) AdminAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// First apply regular auth
		ctx, err := a.authenticate(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		// Check if user is admin
		role, ok := ctx.Value("user_role").(string)
		if !ok || role != "admin" {
			return nil, status.Errorf(codes.PermissionDenied, "admin access required")
		}

		return handler(ctx, req)
	}
}

// authenticate handles the authentication logic
func (a *AuthInterceptor) authenticate(ctx context.Context, method string) (context.Context, error) {
	// Skip authentication for certain methods
	if a.shouldSkipAuth(method) {
		return ctx, nil
	}

	// Extract token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header not provided")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == authHeader[0] {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	// Validate token with user service
	user, err := (*a.userClient).ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Add user info to context
	ctx = context.WithValue(ctx, "user_id", uint(user.Id))
	ctx = context.WithValue(ctx, "user_role", string(user.Role))
	ctx = context.WithValue(ctx, "user_email", user.Email)

	return ctx, nil
}
