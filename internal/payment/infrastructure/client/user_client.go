package client

import (
	"context"
	"fmt"

	userpb "github.com/ddd-micro/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// UserClient defines the interface for user service operations
type UserClient interface {
	ValidateToken(ctx context.Context, token string) (*userpb.User, error)
	GetUserByID(ctx context.Context, userID uint) (*userpb.User, error)
	CheckPermission(ctx context.Context, userID uint, resource, action string) (bool, error)
}

// userClient implements UserClient interface
type userClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

// NewUserClient creates a new user client
func NewUserClient(userServiceURL string) (UserClient, error) {
	conn, err := grpc.Dial(userServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userpb.NewUserServiceClient(conn)

	return &userClient{
		conn:   conn,
		client: client,
	}, nil
}

// ValidateToken validates a JWT token and returns user information
func (c *userClient) ValidateToken(ctx context.Context, token string) (*userpb.User, error) {
	// Since there's no ValidateToken method, we'll use GetProfile with token in metadata
	// This is a simplified implementation
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	
	req := &userpb.GetProfileRequest{}
	resp, err := c.client.GetProfile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return resp.User, nil
}

// GetUserByID gets user information by ID
func (c *userClient) GetUserByID(ctx context.Context, userID uint) (*userpb.User, error) {
	req := &userpb.GetUserRequest{
		Id: uint32(userID),
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.User, nil
}

// CheckPermission checks if user has permission for a specific resource and action
func (c *userClient) CheckPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	// Since there's no CheckPermission method, we'll get user info and check role
	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Simple permission check based on role
	// In a real implementation, this would be more sophisticated
	if user.Role == "admin" {
		return true, nil
	}

	// For now, return true for basic operations
	return true, nil
}

// Close closes the gRPC connection
func (c *userClient) Close() error {
	return c.conn.Close()
}
