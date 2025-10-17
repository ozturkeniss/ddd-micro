// Payment Types
export interface Payment {
  id: string;
  user_id: number;
  order_id: string;
  amount: number;
  currency: string;
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled' | 'refunded';
  payment_method: 'credit_card' | 'debit_card' | 'bank_transfer' | 'paypal' | 'stripe';
  payment_provider: string;
  transaction_id?: string;
  gateway_response?: Record<string, any>;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface PaymentMethod {
  id: string;
  user_id: number;
  type: 'credit_card' | 'debit_card' | 'bank_account';
  provider: string;
  last_four_digits?: string;
  expiry_month?: number;
  expiry_year?: number;
  is_default: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Refund {
  id: string;
  payment_id: string;
  amount: number;
  reason: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

// Request Types
export interface CreatePaymentRequest {
  order_id: string;
  amount: number;
  currency: string;
  payment_method: string;
  payment_method_id?: string;
  return_url?: string;
  cancel_url?: string;
}

export interface ProcessPaymentRequest {
  payment_id: string;
  payment_method_id: string;
  confirmation_data?: Record<string, any>;
}

export interface CreateRefundRequest {
  payment_id: string;
  amount: number;
  reason: string;
}

export interface AddPaymentMethodRequest {
  type: 'credit_card' | 'debit_card' | 'bank_account';
  provider: string;
  token: string;
  is_default?: boolean;
}

export interface UpdatePaymentMethodRequest {
  is_default?: boolean;
  is_active?: boolean;
}

// Response Types
export interface PaymentResponse {
  payment: Payment;
  payment_url?: string;
  client_secret?: string;
}

export interface PaymentMethodResponse {
  payment_method: PaymentMethod;
}

export interface RefundResponse {
  refund: Refund;
}

export interface ListPaymentsResponse {
  payments: Payment[];
  total: number;
  offset: number;
  limit: number;
}

export interface ListPaymentMethodsResponse {
  payment_methods: PaymentMethod[];
  total: number;
}

export interface ListRefundsResponse {
  refunds: Refund[];
  total: number;
  offset: number;
  limit: number;
}

export interface PaymentStatsResponse {
  total_payments: number;
  total_amount: number;
  successful_payments: number;
  failed_payments: number;
  pending_payments: number;
  refunded_amount: number;
  average_payment_amount: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}
