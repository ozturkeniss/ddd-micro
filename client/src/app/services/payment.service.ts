import axios from 'axios';
import {
  Payment,
  PaymentMethod,
  Refund,
  CreatePaymentRequest,
  ProcessPaymentRequest,
  CreateRefundRequest,
  AddPaymentMethodRequest,
  UpdatePaymentMethodRequest,
  PaymentResponse,
  PaymentMethodResponse,
  RefundResponse,
  ListPaymentsResponse,
  ListPaymentMethodsResponse,
  ListRefundsResponse,
  PaymentStatsResponse,
  ApiResponse,
} from '../../types/payment.types';

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

// Payment Service Class
export class PaymentService {
  // ========== PAYMENT ENDPOINTS ==========

  /**
   * Create a new payment
   */
  static async createPayment(
    data: CreatePaymentRequest
  ): Promise<ApiResponse<PaymentResponse>> {
    const response = await apiClient.post('/payments', data);
    return response.data;
  }

  /**
   * Process payment
   */
  static async processPayment(
    data: ProcessPaymentRequest
  ): Promise<ApiResponse<PaymentResponse>> {
    const response = await apiClient.post('/payments/process', data);
    return response.data;
  }

  /**
   * Get payment by ID
   */
  static async getPayment(id: string): Promise<ApiResponse<Payment>> {
    const response = await apiClient.get(`/payments/${id}`);
    return response.data;
  }

  /**
   * Get user's payments
   */
  static async getUserPayments(params?: {
    offset?: number;
    limit?: number;
    status?: string;
    payment_method?: string;
  }): Promise<ApiResponse<ListPaymentsResponse>> {
    const response = await apiClient.get('/payments', { params });
    return response.data;
  }

  /**
   * Cancel payment
   */
  static async cancelPayment(id: string): Promise<ApiResponse<Payment>> {
    const response = await apiClient.post(`/payments/${id}/cancel`);
    return response.data;
  }

  // ========== REFUND ENDPOINTS ==========

  /**
   * Create refund
   */
  static async createRefund(
    data: CreateRefundRequest
  ): Promise<ApiResponse<RefundResponse>> {
    const response = await apiClient.post('/refunds', data);
    return response.data;
  }

  /**
   * Get refund by ID
   */
  static async getRefund(id: string): Promise<ApiResponse<Refund>> {
    const response = await apiClient.get(`/refunds/${id}`);
    return response.data;
  }

  /**
   * Get user's refunds
   */
  static async getUserRefunds(params?: {
    offset?: number;
    limit?: number;
    status?: string;
  }): Promise<ApiResponse<ListRefundsResponse>> {
    const response = await apiClient.get('/refunds', { params });
    return response.data;
  }

  // ========== PAYMENT METHOD ENDPOINTS ==========

  /**
   * Add payment method
   */
  static async addPaymentMethod(
    data: AddPaymentMethodRequest
  ): Promise<ApiResponse<PaymentMethodResponse>> {
    const response = await apiClient.post('/payment-methods', data);
    return response.data;
  }

  /**
   * Get user's payment methods
   */
  static async getPaymentMethods(): Promise<
    ApiResponse<ListPaymentMethodsResponse>
  > {
    const response = await apiClient.get('/payment-methods');
    return response.data;
  }

  /**
   * Update payment method
   */
  static async updatePaymentMethod(
    id: string,
    data: UpdatePaymentMethodRequest
  ): Promise<ApiResponse<PaymentMethodResponse>> {
    const response = await apiClient.put(`/payment-methods/${id}`, data);
    return response.data;
  }

  /**
   * Delete payment method
   */
  static async deletePaymentMethod(id: string): Promise<ApiResponse<null>> {
    const response = await apiClient.delete(`/payment-methods/${id}`);
    return response.data;
  }

  /**
   * Set default payment method
   */
  static async setDefaultPaymentMethod(
    id: string
  ): Promise<ApiResponse<PaymentMethodResponse>> {
    const response = await apiClient.post(`/payment-methods/${id}/set-default`);
    return response.data;
  }

  // ========== ADMIN ENDPOINTS ==========

  /**
   * Get all payments (Admin only)
   */
  static async getAllPayments(params?: {
    offset?: number;
    limit?: number;
    status?: string;
    user_id?: number;
    payment_method?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<ApiResponse<ListPaymentsResponse>> {
    const response = await apiClient.get('/admin/payments', { params });
    return response.data;
  }

  /**
   * Get all refunds (Admin only)
   */
  static async getAllRefunds(params?: {
    offset?: number;
    limit?: number;
    status?: string;
    user_id?: number;
    start_date?: string;
    end_date?: string;
  }): Promise<ApiResponse<ListRefundsResponse>> {
    const response = await apiClient.get('/admin/refunds', { params });
    return response.data;
  }

  /**
   * Get payment statistics (Admin only)
   */
  static async getPaymentStats(params?: {
    start_date?: string;
    end_date?: string;
  }): Promise<ApiResponse<PaymentStatsResponse>> {
    const response = await apiClient.get('/admin/payments/stats', { params });
    return response.data;
  }

  /**
   * Process refund (Admin only)
   */
  static async processRefund(id: string): Promise<ApiResponse<Refund>> {
    const response = await apiClient.post(`/admin/refunds/${id}/process`);
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
   * Format currency amount
   */
  static formatCurrency(amount: number, currency: string = 'USD'): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount);
  }

  /**
   * Get payment status color
   */
  static getPaymentStatusColor(status: string): string {
    const statusColors: Record<string, string> = {
      pending: 'yellow',
      processing: 'blue',
      completed: 'green',
      failed: 'red',
      cancelled: 'gray',
      refunded: 'purple',
    };
    return statusColors[status] || 'gray';
  }

  /**
   * Get payment method icon
   */
  static getPaymentMethodIcon(type: string): string {
    const icons: Record<string, string> = {
      credit_card: '💳',
      debit_card: '💳',
      bank_transfer: '🏦',
      paypal: '🅿️',
      stripe: '💳',
    };
    return icons[type] || '💳';
  }
}

export default PaymentService;
