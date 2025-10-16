package grpc

import (
	"context"
	"strings"

	"github.com/ddd-micro/internal/product/application"
	userpb "github.com/ddd-micro/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor handles JWT authentication and authorization
type AuthInterceptor struct {
	userService *application.UserService
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(userService *application.UserService) *AuthInterceptor {
	return &AuthInterceptor{
		userService: userService,
	}
}

// Unary returns a unary server interceptor for authentication
func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if the method requires authentication
		if a.isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract and validate JWT token
		user, err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "user_id", user.Id)
		ctx = context.WithValue(ctx, "user_email", user.Email)
		ctx = context.WithValue(ctx, "user_role", user.Role)
		ctx = context.WithValue(ctx, "user_active", user.IsActive)

		// Check authorization for admin methods
		if a.isAdminMethod(info.FullMethod) {
			if user.Role != "admin" {
				return nil, status.Errorf(codes.PermissionDenied, "admin access required")
			}
		}

		return handler(ctx, req)
	}
}

// authenticate extracts and validates JWT token
func (a *AuthInterceptor) authenticate(ctx context.Context) (*userpb.User, error) {
	// Extract token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header not provided")
	}

	// Extract token from "Bearer <token>" format
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == authHeader[0] {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	// Validate token with user service
	user, err := a.userService.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, status.Errorf(codes.PermissionDenied, "user account is inactive")
	}

	return user, nil
}

// isPublicMethod checks if the method is public (no authentication required)
func (a *AuthInterceptor) isPublicMethod(method string) bool {
	publicMethods := map[string]bool{
		"/product.ProductService/GetProduct":             true,
		"/product.ProductService/GetProductBySKU":        true,
		"/product.ProductService/ListProducts":           true,
		"/product.ProductService/SearchProducts":         true,
		"/product.ProductService/ListProductsByCategory": true,
		"/product.ProductService/IncrementViewCount":     true,
	}
	return publicMethods[method]
}

// isAdminMethod checks if the method requires admin access
func (a *AuthInterceptor) isAdminMethod(method string) bool {
	adminMethods := map[string]bool{
		"/product.ProductService/CreateProduct":     true,
		"/product.ProductService/UpdateProduct":     true,
		"/product.ProductService/DeleteProduct":     true,
		"/product.ProductService/UpdateStock":       true,
		"/product.ProductService/ReduceStock":       true,
		"/product.ProductService/IncreaseStock":     true,
		"/product.ProductService/ActivateProduct":   true,
		"/product.ProductService/DeactivateProduct": true,
		"/product.ProductService/MarkAsFeatured":    true,
		"/product.ProductService/UnmarkAsFeatured":  true,
	}
	return adminMethods[method]
}
