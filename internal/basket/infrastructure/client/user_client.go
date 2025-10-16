package client

import (
	"context"
	"fmt"

	userpb "github.com/ddd-micro/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient interface for user service operations
type UserClient interface {
	GetUser(ctx context.Context, userID uint) (*userpb.User, error)
	ValidateToken(ctx context.Context, token string) (*userpb.User, error)
	Close() error
}

// userClient implements UserClient interface
type userClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

// NewUserClient creates a new user service gRPC client
func NewUserClient(userServiceURL string) (UserClient, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(userServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	// Create client
	client := userpb.NewUserServiceClient(conn)

	return &userClient{
		conn:   conn,
		client: client,
	}, nil
}

// GetUser retrieves user information by ID
func (c *userClient) GetUser(ctx context.Context, userID uint) (*userpb.User, error) {
	req := &userpb.GetUserRequest{
		Id: uint32(userID),
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.User, nil
}

// ValidateToken validates a JWT token and returns user information
func (c *userClient) ValidateToken(ctx context.Context, token string) (*userpb.User, error) {
	// For now, we'll use a simple approach - in a real implementation,
	// you might want to add a dedicated ValidateToken RPC method
	// or implement token validation logic here
	return nil, fmt.Errorf("ValidateToken not implemented yet")
}

// Close closes the gRPC connection
func (c *userClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// HealthCheck performs a health check on the user service
func (c *userClient) HealthCheck(ctx context.Context) error {
	// For now, just try to get a user with ID 1 to check if service is alive
	// In a real implementation, you might want to add a dedicated health check endpoint
	req := &userpb.GetUserRequest{Id: 1}
	_, err := c.client.GetUser(ctx, req)
	if err != nil {
		return fmt.Errorf("user service health check failed: %w", err)
	}
	return nil
}
