// Basket Types
export interface Basket {
  id: string;
  user_id: number;
  items: BasketItem[];
  total: number;
  item_count: number;
  created_at: string;
  updated_at: string;
  expires_at: string;
  is_expired: boolean;
}

export interface BasketItem {
  id: number;
  product_id: number;
  quantity: number;
  unit_price: number;
  total_price: number;
  created_at: string;
  updated_at: string;
}

// Request Types
export interface CreateBasketRequest {
  user_id: number;
}

export interface AddItemRequest {
  user_id: number;
  product_id: number;
  quantity: number;
  unit_price: number;
}

export interface UpdateItemRequest {
  user_id: number;
  quantity: number;
}

export interface RemoveItemRequest {
  user_id: number;
  product_id: number;
}

export interface ClearBasketRequest {
  user_id: number;
}

// Response Types
export interface BasketResponse {
  basket: Basket;
  message?: string;
}

export interface ClearBasketResponse {
  success: boolean;
  message: string;
}

export interface DeleteBasketResponse {
  success: boolean;
  message: string;
}

export interface CleanupExpiredBasketsResponse {
  success: boolean;
  message: string;
  cleaned_count: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}
