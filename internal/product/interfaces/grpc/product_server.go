package grpc

import (
	"context"

	productpb "github.com/ddd-micro/api/proto/product"
	"github.com/ddd-micro/internal/product/application"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductServer implements the gRPC ProductService
type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	productService *application.ProductServiceCQRS
}

// NewProductServer creates a new gRPC product server
func NewProductServer(productService *application.ProductServiceCQRS) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

// CreateProduct handles product creation
func (s *ProductServer) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.ProductResponse, error) {
	appReq := application.CreateProductRequest{
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		ComparePrice:     req.ComparePrice,
		CostPrice:        req.CostPrice,
		Stock:            int(req.Stock),
		MinStock:         int(req.MinStock),
		MaxStock:         int(req.MaxStock),
		Category:         req.Category,
		SubCategory:      req.SubCategory,
		Brand:            req.Brand,
		SKU:              req.Sku,
		Barcode:          req.Barcode,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		Color:            req.Color,
		Size:             req.Size,
		Material:         req.Material,
		Tags:             req.Tags,
		Images:           req.Images,
		IsDigital:        req.IsDigital,
		IsFeatured:       req.IsFeatured,
		IsOnSale:         req.IsOnSale,
		SortOrder:        int(req.SortOrder),
	}

	productResp, err := s.productService.CreateProduct(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// GetProduct handles product retrieval by ID
func (s *ProductServer) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.ProductResponse, error) {
	productResp, err := s.productService.GetProductByID(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// GetProductBySKU handles product retrieval by SKU
func (s *ProductServer) GetProductBySKU(ctx context.Context, req *productpb.GetProductBySKURequest) (*productpb.ProductResponse, error) {
	productResp, err := s.productService.GetProductBySKU(ctx, req.Sku)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// UpdateProduct handles product updates
func (s *ProductServer) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.ProductResponse, error) {
	appReq := application.UpdateProductRequest{
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		ComparePrice:     req.ComparePrice,
		CostPrice:        req.CostPrice,
		Stock:            req.Stock,
		MinStock:         req.MinStock,
		MaxStock:         req.MaxStock,
		Category:         req.Category,
		SubCategory:      req.SubCategory,
		Brand:            req.Brand,
		Barcode:          req.Barcode,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		Color:            req.Color,
		Size:             req.Size,
		Material:         req.Material,
		Tags:             req.Tags,
		Images:           req.Images,
		IsActive:         req.IsActive,
		IsDigital:        req.IsDigital,
		IsFeatured:       req.IsFeatured,
		IsOnSale:         req.IsOnSale,
		SortOrder:        req.SortOrder,
	}

	productResp, err := s.productService.UpdateProduct(ctx, uint(req.Id), appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// DeleteProduct handles product deletion
func (s *ProductServer) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error) {
	err := s.productService.DeleteProduct(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &productpb.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

// ListProducts handles product listing
func (s *ProductServer) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	listResp, err := s.productService.ListProducts(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	products := make([]*productpb.Product, len(listResp.Products))
	for i, p := range listResp.Products {
		products[i] = toProtoProduct(&p)
	}

	return &productpb.ListProductsResponse{
		Products: products,
		Total:    int32(listResp.Total),
		Offset:   int32(listResp.Offset),
		Limit:    int32(listResp.Limit),
	}, nil
}

// SearchProducts handles product search
func (s *ProductServer) SearchProducts(ctx context.Context, req *productpb.SearchProductsRequest) (*productpb.ListProductsResponse, error) {
	searchResp, err := s.productService.SearchProducts(ctx, req.Query, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search products: %v", err)
	}

	products := make([]*productpb.Product, len(searchResp.Products))
	for i, p := range searchResp.Products {
		products[i] = toProtoProduct(&p)
	}

	return &productpb.ListProductsResponse{
		Products: products,
		Total:    int32(searchResp.Total),
		Offset:   int32(searchResp.Offset),
		Limit:    int32(searchResp.Limit),
	}, nil
}

// ListProductsByCategory handles product listing by category
func (s *ProductServer) ListProductsByCategory(ctx context.Context, req *productpb.ListProductsByCategoryRequest) (*productpb.ListProductsResponse, error) {
	listResp, err := s.productService.ListProductsByCategory(ctx, req.Category, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products by category: %v", err)
	}

	products := make([]*productpb.Product, len(listResp.Products))
	for i, p := range listResp.Products {
		products[i] = toProtoProduct(&p)
	}

	return &productpb.ListProductsResponse{
		Products: products,
		Total:    int32(listResp.Total),
		Offset:    int32(listResp.Offset),
		Limit:     int32(listResp.Limit),
	}, nil
}

// UpdateStock handles stock updates
func (s *ProductServer) UpdateStock(ctx context.Context, req *productpb.UpdateStockRequest) (*productpb.UpdateStockResponse, error) {
	err := s.productService.UpdateStock(ctx, uint(req.ProductId), int(req.Stock))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
	}

	return &productpb.UpdateStockResponse{
		Message: "Stock updated successfully",
	}, nil
}

// ReduceStock handles stock reduction
func (s *ProductServer) ReduceStock(ctx context.Context, req *productpb.ReduceStockRequest) (*productpb.ReduceStockResponse, error) {
	err := s.productService.ReduceStock(ctx, uint(req.ProductId), int(req.Amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reduce stock: %v", err)
	}

	return &productpb.ReduceStockResponse{
		Message: "Stock reduced successfully",
	}, nil
}

// IncreaseStock handles stock increase
func (s *ProductServer) IncreaseStock(ctx context.Context, req *productpb.IncreaseStockRequest) (*productpb.IncreaseStockResponse, error) {
	err := s.productService.IncreaseStock(ctx, uint(req.ProductId), int(req.Amount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to increase stock: %v", err)
	}

	return &productpb.IncreaseStockResponse{
		Message: "Stock increased successfully",
	}, nil
}

// ActivateProduct handles product activation
func (s *ProductServer) ActivateProduct(ctx context.Context, req *productpb.ActivateProductRequest) (*productpb.ProductResponse, error) {
	productResp, err := s.productService.ActivateProduct(ctx, uint(req.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to activate product: %v", err)
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// DeactivateProduct handles product deactivation
func (s *ProductServer) DeactivateProduct(ctx context.Context, req *productpb.DeactivateProductRequest) (*productpb.ProductResponse, error) {
	productResp, err := s.productService.DeactivateProduct(ctx, uint(req.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to deactivate product: %v", err)
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// MarkAsFeatured handles marking product as featured
func (s *ProductServer) MarkAsFeatured(ctx context.Context, req *productpb.MarkAsFeaturedRequest) (*productpb.ProductResponse, error) {
	productResp, err := s.productService.MarkAsFeatured(ctx, uint(req.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark product as featured: %v", err)
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// UnmarkAsFeatured handles unmarking product as featured
func (s *ProductServer) UnmarkAsFeatured(ctx context.Context, req *productpb.UnmarkAsFeaturedRequest) (*productpb.ProductResponse, error) {
	productResp, err := s.productService.UnmarkAsFeatured(ctx, uint(req.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmark product as featured: %v", err)
	}

	return &productpb.ProductResponse{
		Product: toProtoProduct(productResp),
	}, nil
}

// IncrementViewCount handles view count increment
func (s *ProductServer) IncrementViewCount(ctx context.Context, req *productpb.IncrementViewCountRequest) (*productpb.IncrementViewCountResponse, error) {
	err := s.productService.IncrementViewCount(ctx, uint(req.ProductId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to increment view count: %v", err)
	}

	return &productpb.IncrementViewCountResponse{
		Message: "View count incremented successfully",
	}, nil
}

// Helper function to convert application.ProductResponse to proto.Product
func toProtoProduct(p *application.ProductResponse) *productpb.Product {
	return &productpb.Product{
		Id:               uint32(p.ID),
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Price:            p.Price,
		ComparePrice:     p.ComparePrice,
		CostPrice:        p.CostPrice,
		Stock:            int32(p.Stock),
		MinStock:         int32(p.MinStock),
		MaxStock:         int32(p.MaxStock),
		Category:         p.Category,
		SubCategory:      p.SubCategory,
		Brand:            p.Brand,
		Sku:              p.SKU,
		Barcode:          p.Barcode,
		Weight:           p.Weight,
		Dimensions:       p.Dimensions,
		Color:            p.Color,
		Size:             p.Size,
		Material:         p.Material,
		Tags:             p.Tags,
		Images:           p.Images,
		IsActive:         p.IsActive,
		IsDigital:        p.IsDigital,
		IsFeatured:       p.IsFeatured,
		IsOnSale:         p.IsOnSale,
		SortOrder:        int32(p.SortOrder),
		ViewCount:        int32(p.ViewCount),
		CreatedAt:        timestamppb.New(p.CreatedAt),
		UpdatedAt:        timestamppb.New(p.UpdatedAt),
	}
}
