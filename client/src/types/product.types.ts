// Product Types
export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  category: string;
  brand: string;
  sku: string;
  stock_quantity: number;
  is_active: boolean;
  is_featured: boolean;
  images: string[];
  specifications: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface ProductVariant {
  id: number;
  product_id: number;
  name: string;
  value: string;
  price_adjustment: number;
  stock_quantity: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  description: string;
  parent_id?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Request Types
export interface CreateProductRequest {
  name: string;
  description: string;
  price: number;
  category: string;
  brand: string;
  sku: string;
  stock_quantity: number;
  images?: string[];
  specifications?: Record<string, any>;
  variants?: Omit<ProductVariant, 'id' | 'product_id' | 'created_at' | 'updated_at'>[];
}

export interface UpdateProductRequest {
  name?: string;
  description?: string;
  price?: number;
  category?: string;
  brand?: string;
  sku?: string;
  stock_quantity?: number;
  images?: string[];
  specifications?: Record<string, any>;
  is_active?: boolean;
  is_featured?: boolean;
}

export interface UpdateStockRequest {
  stock_quantity: number;
}

export interface SearchProductsRequest {
  query?: string;
  category?: string;
  brand?: string;
  min_price?: number;
  max_price?: number;
  is_featured?: boolean;
  is_active?: boolean;
  offset?: number;
  limit?: number;
  sort_by?: 'name' | 'price' | 'created_at' | 'updated_at';
  sort_order?: 'asc' | 'desc';
}

export interface ViewProductRequest {
  user_id?: number;
  ip_address?: string;
  user_agent?: string;
}

// Response Types
export interface ListProductsResponse {
  products: Product[];
  total: number;
  offset: number;
  limit: number;
  filters: {
    categories: string[];
    brands: string[];
    price_range: {
      min: number;
      max: number;
    };
  };
}

export interface ProductResponse {
  product: Product;
  variants: ProductVariant[];
  related_products: Product[];
  view_count: number;
}

export interface SearchProductsResponse {
  products: Product[];
  total: number;
  offset: number;
  limit: number;
  search_query: string;
  filters_applied: Record<string, any>;
}

export interface CategoryResponse {
  categories: Category[];
  total: number;
}

export interface StockUpdateResponse {
  product_id: number;
  old_stock: number;
  new_stock: number;
  updated_at: string;
}

export interface ProductActivationResponse {
  product_id: number;
  is_active: boolean;
  message: string;
}

export interface ProductFeaturedResponse {
  product_id: number;
  is_featured: boolean;
  message: string;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}
