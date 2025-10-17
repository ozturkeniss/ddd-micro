package client

import (
	"context"
	"fmt"

	userpb "github.com/ddd-micro/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	req := &userpb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.client.ValidateToken(ctx, req)
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
	req := &userpb.CheckPermissionRequest{
		UserId:   uint32(userID),
		Resource: resource,
		Action:   action,
	}

	resp, err := c.client.CheckPermission(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return resp.HasPermission, nil
}

// Close closes the gRPC connection
func (c *userClient) Close() error {
	return c.conn.Close()
}
