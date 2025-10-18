import axios from 'axios';
import {
  Basket,
  BasketItem,
  CreateBasketRequest,
  AddItemRequest,
  UpdateItemRequest,
  RemoveItemRequest,
  ClearBasketRequest,
  BasketResponse,
  ClearBasketResponse,
  DeleteBasketResponse,
  CleanupExpiredBasketsResponse,
  ApiResponse,
} from '../../types/basket.types';

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

// Basket Service Class
export class BasketService {
  // ========== USER ENDPOINTS ==========

  /**
   * Create a new basket for the current user
   */
  static async createBasket(): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    const data: CreateBasketRequest = {
      user_id: user.id,
    };

    const response = await apiClient.post('/basket', data);
    return response.data;
  }

  /**
   * Get current user's basket
   */
  static async getBasket(): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    const response = await apiClient.get('/basket');
    return response.data;
  }

  /**
   * Add item to basket
   */
  static async addItem(
    productId: number,
    quantity: number,
    unitPrice: number
  ): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    const data: AddItemRequest = {
      user_id: user.id,
      product_id: productId,
      quantity,
      unit_price: unitPrice,
    };

    const response = await apiClient.post('/basket/items', data);
    return response.data;
  }

  /**
   * Update item quantity in basket
   */
  static async updateItem(
    productId: number,
    quantity: number
  ): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    const data: UpdateItemRequest = {
      user_id: user.id,
      quantity,
    };

    const response = await apiClient.put('/basket/items', data, {
      params: { product_id: productId },
    });
    return response.data;
  }

  /**
   * Remove item from basket
   */
  static async removeItem(productId: number): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    const response = await apiClient.delete(`/basket/items/${productId}`);
    return response.data;
  }

  /**
   * Clear all items from basket
   */
  static async clearBasket(): Promise<ApiResponse<ClearBasketResponse>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    const response = await apiClient.delete('/basket/clear');
    return response.data;
  }

  // ========== ADMIN ENDPOINTS ==========

  /**
   * Get user's basket by user ID (Admin only)
   */
  static async getUserBasket(userId: number): Promise<ApiResponse<Basket>> {
    const response = await apiClient.get(`/admin/baskets/${userId}`);
    return response.data;
  }

  /**
   * Delete user's basket by user ID (Admin only)
   */
  static async deleteUserBasket(
    userId: number
  ): Promise<ApiResponse<DeleteBasketResponse>> {
    const response = await apiClient.delete(`/admin/baskets/${userId}`);
    return response.data;
  }

  /**
   * Cleanup expired baskets (Admin only)
   */
  static async cleanupExpiredBaskets(): Promise<
    ApiResponse<CleanupExpiredBasketsResponse>
  > {
    const response = await apiClient.post('/admin/baskets/cleanup');
    return response.data;
  }

  // ========== UTILITY METHODS ==========

  /**
   * Get current user from localStorage
   */
  private static getCurrentUser(): { id: number } | null {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        return JSON.parse(userStr);
      } catch {
        return null;
      }
    }
    return null;
  }

  /**
   * Check if user is authenticated
   */
  static isAuthenticated(): boolean {
    return !!localStorage.getItem('auth_token');
  }

  /**
   * Check if current user is admin
   */
  static isAdmin(): boolean {
    const user = this.getCurrentUser();
    if (!user) {
      return false;
    }

    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const userData = JSON.parse(userStr);
        return userData.role === 'admin';
      } catch {
        return false;
      }
    }
    return false;
  }

  /**
   * Get basket item count
   */
  static async getBasketItemCount(): Promise<number> {
    try {
      const response = await this.getBasket();
      if (response.success && response.data) {
        return response.data.item_count;
      }
      return 0;
    } catch {
      return 0;
    }
  }

  /**
   * Check if basket is empty
   */
  static async isBasketEmpty(): Promise<boolean> {
    try {
      const response = await this.getBasket();
      if (response.success && response.data) {
        return response.data.items.length === 0;
      }
      return true;
    } catch {
      return true;
    }
  }

  /**
   * Get basket total
   */
  static async getBasketTotal(): Promise<number> {
    try {
      const response = await this.getBasket();
      if (response.success && response.data) {
        return response.data.total;
      }
      return 0;
    } catch {
      return 0;
    }
  }

  /**
   * Add multiple items to basket
   */
  static async addMultipleItems(
    items: Array<{
      productId: number;
      quantity: number;
      unitPrice: number;
    }>
  ): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    // Add items one by one
    let lastResponse: ApiResponse<Basket> | null = null;

    for (const item of items) {
      const data: AddItemRequest = {
        user_id: user.id,
        product_id: item.productId,
        quantity: item.quantity,
        unit_price: item.unitPrice,
      };

      const response = await apiClient.post('/basket/items', data);
      lastResponse = response.data;
    }

    return lastResponse || { success: false, message: 'Failed to add items' };
  }

  /**
   * Update multiple items in basket
   */
  static async updateMultipleItems(
    items: Array<{
      productId: number;
      quantity: number;
    }>
  ): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    // Update items one by one
    let lastResponse: ApiResponse<Basket> | null = null;

    for (const item of items) {
      const data: UpdateItemRequest = {
        user_id: user.id,
        quantity: item.quantity,
      };

      const response = await apiClient.put('/basket/items', data, {
        params: { product_id: item.productId },
      });
      lastResponse = response.data;
    }

    return (
      lastResponse || { success: false, message: 'Failed to update items' }
    );
  }

  /**
   * Remove multiple items from basket
   */
  static async removeMultipleItems(
    productIds: number[]
  ): Promise<ApiResponse<Basket>> {
    const user = this.getCurrentUser();
    if (!user) {
      throw new Error('User not authenticated');
    }

    // Remove items one by one
    let lastResponse: ApiResponse<Basket> | null = null;

    for (const productId of productIds) {
      const response = await apiClient.delete(`/basket/items/${productId}`);
      lastResponse = response.data;
    }

    return (
      lastResponse || { success: false, message: 'Failed to remove items' }
    );
  }
}

export default BasketService;
