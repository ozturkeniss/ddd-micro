package grpc

import (
	"context"

	basketpb "github.com/ddd-micro/api/proto/basket"
	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/application/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// BasketServer implements the gRPC BasketService
type BasketServer struct {
	basketpb.UnimplementedBasketServiceServer
	basketService *application.BasketServiceCQRS
}

// NewBasketServer creates a new gRPC basket server
func NewBasketServer(basketService *application.BasketServiceCQRS) *BasketServer {
	return &BasketServer{
		basketService: basketService,
	}
}

// CreateBasket creates a new basket for a user
func (s *BasketServer) CreateBasket(ctx context.Context, req *basketpb.CreateBasketRequest) (*basketpb.BasketResponse, error) {
	appReq := dto.CreateBasketRequest{
		UserID: uint(req.UserId),
	}

	basketResp, err := s.basketService.CreateBasket(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create basket: %v", err)
	}

	return toProtoBasket(basketResp), nil
}

// GetBasket retrieves a user's basket
func (s *BasketServer) GetBasket(ctx context.Context, req *basketpb.GetBasketRequest) (*basketpb.BasketResponse, error) {
	appReq := dto.GetBasketRequest{
		UserID: uint(req.UserId),
	}

	basketResp, err := s.basketService.GetBasket(ctx, appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get basket: %v", err)
	}

	return toProtoBasket(basketResp), nil
}

// AddItem adds an item to the basket
func (s *BasketServer) AddItem(ctx context.Context, req *basketpb.AddItemRequest) (*basketpb.BasketResponse, error) {
	appReq := dto.AddItemRequest{
		UserID:    uint(req.UserId),
		ProductID: uint(req.ProductId),
		Quantity:  int(req.Quantity),
		UnitPrice: req.UnitPrice,
	}

	basketResp, err := s.basketService.AddItem(ctx, appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		if err == application.ErrInvalidQuantity {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
		}
		if err == application.ErrInvalidPrice {
			return nil, status.Errorf(codes.InvalidArgument, "invalid price")
		}
		return nil, status.Errorf(codes.Internal, "failed to add item: %v", err)
	}

	return toProtoBasket(basketResp), nil
}

// UpdateItem updates the quantity of an item in the basket
func (s *BasketServer) UpdateItem(ctx context.Context, req *basketpb.UpdateItemRequest) (*basketpb.BasketResponse, error) {
	appReq := dto.UpdateItemRequest{
		UserID:   uint(req.UserId),
		Quantity: int(req.Quantity),
	}

	basketResp, err := s.basketService.UpdateItem(ctx, uint(req.ProductId), appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		if err == application.ErrItemNotFound {
			return nil, status.Errorf(codes.NotFound, "item not found in basket")
		}
		if err == application.ErrInvalidQuantity {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
		}
		return nil, status.Errorf(codes.Internal, "failed to update item: %v", err)
	}

	return toProtoBasket(basketResp), nil
}

// RemoveItem removes an item from the basket
func (s *BasketServer) RemoveItem(ctx context.Context, req *basketpb.RemoveItemRequest) (*basketpb.BasketResponse, error) {
	appReq := dto.RemoveItemRequest{
		UserID:    uint(req.UserId),
		ProductID: uint(req.ProductId),
	}

	basketResp, err := s.basketService.RemoveItem(ctx, appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		if err == application.ErrItemNotFound {
			return nil, status.Errorf(codes.NotFound, "item not found in basket")
		}
		return nil, status.Errorf(codes.Internal, "failed to remove item: %v", err)
	}

	return toProtoBasket(basketResp), nil
}

// ClearBasket clears all items from the basket
func (s *BasketServer) ClearBasket(ctx context.Context, req *basketpb.ClearBasketRequest) (*basketpb.ClearBasketResponse, error) {
	appReq := dto.ClearBasketRequest{
		UserID: uint(req.UserId),
	}

	err := s.basketService.ClearBasket(ctx, appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to clear basket: %v", err)
	}

	return &basketpb.ClearBasketResponse{
		Success: true,
		Message: "Basket cleared successfully",
	}, nil
}

// GetUserBasket retrieves a specific user's basket (admin only)
func (s *BasketServer) GetUserBasket(ctx context.Context, req *basketpb.GetUserBasketRequest) (*basketpb.BasketResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	appReq := dto.GetBasketRequest{
		UserID: uint(req.UserId),
	}

	basketResp, err := s.basketService.GetBasket(ctx, appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user basket: %v", err)
	}

	return toProtoBasket(basketResp), nil
}

// DeleteUserBasket deletes a specific user's basket (admin only)
func (s *BasketServer) DeleteUserBasket(ctx context.Context, req *basketpb.DeleteUserBasketRequest) (*basketpb.DeleteUserBasketResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	appReq := dto.ClearBasketRequest{
		UserID: uint(req.UserId),
	}

	err := s.basketService.ClearBasket(ctx, appReq)
	if err != nil {
		if err == application.ErrBasketNotFound {
			return nil, status.Errorf(codes.NotFound, "basket not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user basket: %v", err)
	}

	return &basketpb.DeleteUserBasketResponse{
		Success: true,
		Message: "User basket deleted successfully",
	}, nil
}

// CleanupExpiredBaskets removes all expired baskets (admin only)
func (s *BasketServer) CleanupExpiredBaskets(ctx context.Context, req *basketpb.CleanupExpiredBasketsRequest) (*basketpb.CleanupExpiredBasketsResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	cleanedCount, err := s.basketService.AdminCleanupExpiredBaskets(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cleanup expired baskets: %v", err)
	}

	return &basketpb.CleanupExpiredBasketsResponse{
		Success:      true,
		Message:      "Expired baskets cleaned up successfully",
		CleanedCount: int32(cleanedCount),
	}, nil
}

// Helper functions

func toProtoBasket(basket *dto.BasketResponse) *basketpb.BasketResponse {
	items := make([]*basketpb.BasketItem, len(basket.Items))
	for i, item := range basket.Items {
		items[i] = &basketpb.BasketItem{
			Id:         uint32(item.ID),
			ProductId:  uint32(item.ProductID),
			Quantity:   int32(item.Quantity),
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
			CreatedAt:  timestamppb.New(item.CreatedAt),
			UpdatedAt:  timestamppb.New(item.UpdatedAt),
		}
	}

	return &basketpb.BasketResponse{
		Id:        basket.ID,
		UserId:    uint32(basket.UserID),
		Items:     items,
		Total:     basket.Total,
		ItemCount: int32(basket.ItemCount),
		CreatedAt: timestamppb.New(basket.CreatedAt),
		UpdatedAt: timestamppb.New(basket.UpdatedAt),
		ExpiresAt: timestamppb.New(basket.ExpiresAt),
		IsExpired: basket.IsExpired,
	}
}

func requireAdmin(ctx context.Context) error {
	role, ok := ctx.Value("user_role").(string)
	if !ok || role != "admin" {
		return status.Errorf(codes.PermissionDenied, "admin access required")
	}
	return nil
}
