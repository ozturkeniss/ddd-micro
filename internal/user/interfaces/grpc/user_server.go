package grpc

import (
	"context"

	userpb "github.com/ddd-micro/api/proto/user"
	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserServer implements the gRPC UserService
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userService *application.UserServiceCQRS
}

// NewUserServer creates a new gRPC user server
func NewUserServer(userService *application.UserServiceCQRS) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

// Register handles user registration
func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	appReq := application.CreateUserRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	userResp, err := s.userService.CreateUser(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userpb.RegisterResponse{
		User: toProtoUser(userResp),
	}, nil
}

// Login handles user authentication
func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	appReq := application.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	loginResp, err := s.userService.Login(ctx, appReq)
	if err != nil {
		if err == application.ErrInvalidCredentials {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		if err == application.ErrUserInactive {
			return nil, status.Errorf(codes.PermissionDenied, "user account is inactive")
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &userpb.LoginResponse{
		User:  toProtoUser(&loginResp.User),
		Token: loginResp.Token,
	}, nil
}

// RefreshToken handles token refresh
func (s *UserServer) RefreshToken(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.RefreshTokenResponse, error) {
	newToken, err := s.userService.RefreshToken(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token")
	}

	return &userpb.RefreshTokenResponse{
		Token: newToken,
	}, nil
}

// GetProfile retrieves the authenticated user's profile
func (s *UserServer) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.UserResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userResp, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userpb.UserResponse{
		User: toProtoUser(userResp),
	}, nil
}

// UpdateProfile updates the authenticated user's profile
func (s *UserServer) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.UserResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	appReq := application.UpdateUserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	userResp, err := s.userService.UpdateUser(ctx, userID, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &userpb.UserResponse{
		User: toProtoUser(userResp),
	}, nil
}

// ChangePassword handles password change
func (s *UserServer) ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.userService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to change password: %v", err)
	}

	return &userpb.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}

// GetUser retrieves a user by ID (admin only)
func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	userResp, err := s.userService.GetUserByID(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userpb.UserResponse{
		User: toProtoUser(userResp),
	}, nil
}

// ListUsers retrieves all users with pagination (admin only)
func (s *UserServer) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	listResp, err := s.userService.ListUsers(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	users := make([]*userpb.User, len(listResp.Users))
	for i, u := range listResp.Users {
		users[i] = toProtoUser(&u)
	}

	return &userpb.ListUsersResponse{
		Users:  users,
		Total:  int32(listResp.Total),
		Offset: int32(listResp.Offset),
		Limit:  int32(listResp.Limit),
	}, nil
}

// UpdateUser updates a user (admin only)
func (s *UserServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	isActive := req.IsActive
	appReq := application.UpdateUserByAdminRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      domain.Role(req.Role),
		IsActive:  isActive,
	}

	userResp, err := s.userService.UpdateUserByAdmin(ctx, uint(req.Id), appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &userpb.UserResponse{
		User: toProtoUser(userResp),
	}, nil
}

// DeleteUser deletes a user (admin only)
func (s *UserServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	err := s.userService.DeleteUser(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &userpb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

// AssignRole assigns a role to a user (admin only)
func (s *UserServer) AssignRole(ctx context.Context, req *userpb.AssignRoleRequest) (*userpb.UserResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	userResp, err := s.userService.AssignRole(ctx, uint(req.Id), domain.Role(req.Role))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign role: %v", err)
	}

	return &userpb.UserResponse{
		User: toProtoUser(userResp),
	}, nil
}

// Helper functions

func toProtoUser(u *application.UserResponse) *userpb.User {
	return &userpb.User{
		Id:        uint32(u.ID),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func getUserIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value("user_id").(uint)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	return userID, nil
}

func requireAdmin(ctx context.Context) error {
	role, ok := ctx.Value("user_role").(string)
	if !ok || role != string(domain.RoleAdmin) {
		return status.Errorf(codes.PermissionDenied, "admin access required")
	}
	return nil
}
