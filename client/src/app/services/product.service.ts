import axios from 'axios';
import {
  Product,
  ProductVariant,
  Category,
  CreateProductRequest,
  UpdateProductRequest,
  UpdateStockRequest,
  SearchProductsRequest,
  ViewProductRequest,
  ListProductsResponse,
  ProductResponse,
  SearchProductsResponse,
  CategoryResponse,
  StockUpdateResponse,
  ProductActivationResponse,
  ProductFeaturedResponse,
  ApiResponse,
} from '../../types/product.types';

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081';
const API_VERSION = '/api/v1';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}${API_VERSION}`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(config => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
apiClient.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Product Service Class
export class ProductService {
  // ========== PUBLIC ENDPOINTS ==========

  /**
   * Get all products with pagination and filters
   */
  static async getProducts(params?: {
    offset?: number;
    limit?: number;
    category?: string;
    brand?: string;
    min_price?: number;
    max_price?: number;
    is_featured?: boolean;
    sort_by?: 'name' | 'price' | 'created_at' | 'updated_at';
    sort_order?: 'asc' | 'desc';
  }): Promise<ApiResponse<ListProductsResponse>> {
    const response = await apiClient.get('/products', { params });
    return response.data;
  }

  /**
   * Search products
   */
  static async searchProducts(
    searchParams: SearchProductsRequest
  ): Promise<ApiResponse<SearchProductsResponse>> {
    const response = await apiClient.get('/products/search', {
      params: searchParams,
    });
    return response.data;
  }

  /**
   * Get products by category
   */
  static async getProductsByCategory(
    category: string,
    params?: {
      offset?: number;
      limit?: number;
      sort_by?: 'name' | 'price' | 'created_at' | 'updated_at';
      sort_order?: 'asc' | 'desc';
    }
  ): Promise<ApiResponse<ListProductsResponse>> {
    const response = await apiClient.get(`/products/category/${category}`, {
      params,
    });
    return response.data;
  }

  /**
   * Get product by ID
   */
  static async getProductById(
    id: number
  ): Promise<ApiResponse<ProductResponse>> {
    const response = await apiClient.get(`/products/${id}`);
    return response.data;
  }

  /**
   * Record product view
   */
  static async viewProduct(
    id: number,
    data?: ViewProductRequest
  ): Promise<ApiResponse<null>> {
    const response = await apiClient.post(`/products/${id}/view`, data);
    return response.data;
  }

  // ========== ADMIN ENDPOINTS ==========

  /**
   * Create new product (Admin only)
   */
  static async createProduct(
    data: CreateProductRequest
  ): Promise<ApiResponse<Product>> {
    const response = await apiClient.post('/admin/products', data);
    return response.data;
  }

  /**
   * Update product (Admin only)
   */
  static async updateProduct(
    id: number,
    data: UpdateProductRequest
  ): Promise<ApiResponse<Product>> {
    const response = await apiClient.put(`/admin/products/${id}`, data);
    return response.data;
  }

  /**
   * Delete product (Admin only)
   */
  static async deleteProduct(id: number): Promise<ApiResponse<null>> {
    const response = await apiClient.delete(`/admin/products/${id}`);
    return response.data;
  }

  /**
   * Update product stock (Admin only)
   */
  static async updateStock(
    id: number,
    data: UpdateStockRequest
  ): Promise<ApiResponse<StockUpdateResponse>> {
    const response = await apiClient.put(`/admin/products/${id}/stock`, data);
    return response.data;
  }

  /**
   * Activate product (Admin only)
   */
  static async activateProduct(
    id: number
  ): Promise<ApiResponse<ProductActivationResponse>> {
    const response = await apiClient.post(`/admin/products/${id}/activate`);
    return response.data;
  }

  /**
   * Set product as featured (Admin only)
   */
  static async setFeatured(
    id: number
  ): Promise<ApiResponse<ProductFeaturedResponse>> {
    const response = await apiClient.post(`/admin/products/${id}/featured`);
    return response.data;
  }

  // ========== UTILITY METHODS ==========

  /**
   * Get product categories
   */
  static async getCategories(): Promise<ApiResponse<CategoryResponse>> {
    // This would need to be implemented in the backend
    // For now, return a mock response
    return {
      success: true,
      message: 'Categories retrieved successfully',
      data: {
        categories: [],
        total: 0,
      },
    };
  }

  /**
   * Get featured products
   */
  static async getFeaturedProducts(
    limit: number = 10
  ): Promise<ApiResponse<ListProductsResponse>> {
    return this.getProducts({ is_featured: true, limit });
  }

  /**
   * Get products by brand
   */
  static async getProductsByBrand(
    brand: string,
    params?: {
      offset?: number;
      limit?: number;
      sort_by?: 'name' | 'price' | 'created_at' | 'updated_at';
      sort_order?: 'asc' | 'desc';
    }
  ): Promise<ApiResponse<ListProductsResponse>> {
    return this.getProducts({ ...params, brand });
  }

  /**
   * Get products in price range
   */
  static async getProductsInPriceRange(
    minPrice: number,
    maxPrice: number,
    params?: {
      offset?: number;
      limit?: number;
      sort_by?: 'name' | 'price' | 'created_at' | 'updated_at';
      sort_order?: 'asc' | 'desc';
    }
  ): Promise<ApiResponse<ListProductsResponse>> {
    return this.getProducts({
      ...params,
      min_price: minPrice,
      max_price: maxPrice,
    });
  }

  /**
   * Get low stock products (Admin only)
   */
  static async getLowStockProducts(
    threshold: number = 10,
    params?: {
      offset?: number;
      limit?: number;
    }
  ): Promise<ApiResponse<ListProductsResponse>> {
    // This would need to be implemented in the backend
    // For now, return a mock response
    return {
      success: true,
      message: 'Low stock products retrieved successfully',
      data: {
        products: [],
        total: 0,
        offset: params?.offset || 0,
        limit: params?.limit || 20,
        filters: {
          categories: [],
          brands: [],
          price_range: { min: 0, max: 0 },
        },
      },
    };
  }

  /**
   * Get product statistics (Admin only)
   */
  static async getProductStats(): Promise<
    ApiResponse<{
      total_products: number;
      active_products: number;
      inactive_products: number;
      featured_products: number;
      low_stock_products: number;
      total_categories: number;
      total_brands: number;
    }>
  > {
    // This would need to be implemented in the backend
    // For now, return a mock response
    return {
      success: true,
      message: 'Product statistics retrieved successfully',
      data: {
        total_products: 0,
        active_products: 0,
        inactive_products: 0,
        featured_products: 0,
        low_stock_products: 0,
        total_categories: 0,
        total_brands: 0,
      },
    };
  }
}

export default ProductService;
