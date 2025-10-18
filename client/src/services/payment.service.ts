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
} from '@/types/payment.types';

const API_BASE_URL = process.env.NEXT_PUBLIC_PAYMENT_API_URL || 'http://localhost:8084/api/v1';

class PaymentService {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = localStorage.getItem('token');
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Payment Methods
  async createPayment(data: CreatePaymentRequest): Promise<PaymentResponse> {
    return this.request<PaymentResponse>('/payments', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getPayment(paymentId: string): Promise<Payment> {
    return this.request<Payment>(`/payments/${paymentId}`);
  }

  async listPayments(params?: {
    page?: number;
    limit?: number;
    status?: string;
  }): Promise<ListPaymentsResponse> {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('limit', params.limit.toString());
    if (params?.status) searchParams.append('status', params.status);

    const queryString = searchParams.toString();
    return this.request<ListPaymentsResponse>(`/payments${queryString ? `?${queryString}` : ''}`);
  }

  async processPayment(paymentId: string, data: ProcessPaymentRequest): Promise<Payment> {
    return this.request<Payment>(`/payments/${paymentId}/process`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async cancelPayment(paymentId: string): Promise<Payment> {
    return this.request<Payment>(`/payments/${paymentId}/cancel`, {
      method: 'POST',
    });
  }

  // Payment Methods Management
  async getPaymentMethods(): Promise<ListPaymentMethodsResponse> {
    return this.request<ListPaymentMethodsResponse>('/payment-methods');
  }

  async addPaymentMethod(data: AddPaymentMethodRequest): Promise<PaymentMethodResponse> {
    return this.request<PaymentMethodResponse>('/payment-methods', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updatePaymentMethod(paymentMethodId: string, data: UpdatePaymentMethodRequest): Promise<PaymentMethodResponse> {
    return this.request<PaymentMethodResponse>(`/payment-methods/${paymentMethodId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deletePaymentMethod(paymentMethodId: string): Promise<void> {
    await this.request<void>(`/payment-methods/${paymentMethodId}`, {
      method: 'DELETE',
    });
  }

  async setDefaultPaymentMethod(paymentMethodId: string): Promise<PaymentMethodResponse> {
    return this.request<PaymentMethodResponse>(`/payment-methods/${paymentMethodId}/set-default`, {
      method: 'POST',
    });
  }

  // Admin Methods
  async adminListPayments(params?: {
    page?: number;
    limit?: number;
    user_id?: number;
    status?: string;
  }): Promise<ListPaymentsResponse> {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('limit', params.limit.toString());
    if (params?.user_id) searchParams.append('user_id', params.user_id.toString());
    if (params?.status) searchParams.append('status', params.status);

    const queryString = searchParams.toString();
    return this.request<ListPaymentsResponse>(`/admin/payments${queryString ? `?${queryString}` : ''}`);
  }

  async adminGetPayment(paymentId: string): Promise<Payment> {
    return this.request<Payment>(`/admin/payments/${paymentId}`);
  }

  async adminUpdatePaymentStatus(paymentId: string, status: string, reason?: string): Promise<Payment> {
    return this.request<Payment>(`/admin/payments/${paymentId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status, reason }),
    });
  }

  async adminListRefunds(params?: {
    page?: number;
    limit?: number;
    payment_id?: string;
    status?: string;
  }): Promise<ListRefundsResponse> {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('limit', params.limit.toString());
    if (params?.payment_id) searchParams.append('payment_id', params.payment_id);
    if (params?.status) searchParams.append('status', params.status);

    const queryString = searchParams.toString();
    return this.request<ListRefundsResponse>(`/admin/refunds${queryString ? `?${queryString}` : ''}`);
  }

  async adminCreateRefund(data: CreateRefundRequest): Promise<RefundResponse> {
    return this.request<RefundResponse>('/admin/refunds', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async adminGetRefund(refundId: string): Promise<Refund> {
    return this.request<Refund>(`/admin/refunds/${refundId}`);
  }

  async adminProcessRefund(refundId: string): Promise<Refund> {
    return this.request<Refund>(`/admin/refunds/${refundId}/process`, {
      method: 'POST',
    });
  }

  async getPaymentStats(period: 'daily' | 'weekly' | 'monthly' | 'yearly' = 'monthly'): Promise<PaymentStatsResponse> {
    return this.request<PaymentStatsResponse>(`/admin/analytics/payments?period=${period}`);
  }
}

export const paymentService = new PaymentService();
export default paymentService;
