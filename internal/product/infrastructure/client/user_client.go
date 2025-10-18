package client

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/ddd-micro/api/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient defines the interface for user service client
type UserClient interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, userID uint) (*userpb.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*userpb.User, error)

	// ValidateToken validates a JWT token
	ValidateToken(ctx context.Context, token string) (*userpb.User, error)

	// Close closes the client connection
	Close() error
}

// userClient implements UserClient interface
type userClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

// NewUserClient creates a new user service client
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

// GetUserByID retrieves a user by ID
func (c *userClient) GetUserByID(ctx context.Context, userID uint) (*userpb.User, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &userpb.GetUserRequest{
		Id: uint32(userID),
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return resp.User, nil
}

// GetUserByEmail retrieves a user by email
func (c *userClient) GetUserByEmail(ctx context.Context, email string) (*userpb.User, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Since there's no direct GetUserByEmail in the proto, we'll need to implement
	// a workaround or add this method to the user service
	// For now, we'll return an error indicating this method needs to be implemented
	return nil, fmt.Errorf("GetUserByEmail not implemented in user service proto")
}

// ValidateToken validates a JWT token
func (c *userClient) ValidateToken(ctx context.Context, token string) (*userpb.User, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create context with token in metadata
	ctx = context.WithValue(ctx, "authorization", "Bearer "+token)

	req := &userpb.GetProfileRequest{}

	resp, err := c.client.GetProfile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return resp.User, nil
}

// Close closes the client connection
func (c *userClient) Close() error {
	return c.conn.Close()
}
