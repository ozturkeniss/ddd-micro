import axios from 'axios';

// Types - Bad formatting test
export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  role: 'user' | 'admin';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface UpdateUserRequest {
  first_name?: string;
  last_name?: string;
}

export interface UpdateUserByAdminRequest {
  first_name?: string;
  last_name?: string;
  role?: 'user' | 'admin';
  is_active?: boolean;
}

export interface AssignRoleRequest {
  role: 'user' | 'admin';
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  token: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface TokenResponse {
  token: string;
}

export interface ListUsersResponse {
  users: User[];
  total: number;
  offset: number;
  limit: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}

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

// User Service Class
export class UserService {
  // ========== PUBLIC ENDPOINTS ==========

  /**
   * Register a new user
   */
  static async register(data: CreateUserRequest): Promise<ApiResponse<User>> {
    const response = await apiClient.post('/users/register', data);
    return response.data;
  }

  /**
   * Login user
   */
  static async login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    const response = await apiClient.post('/users/login', data);

    // Store token and user data in localStorage
    if (response.data.success && response.data.data) {
      localStorage.setItem('auth_token', response.data.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.data.user));
    }

    return response.data;
  }

  /**
   * Refresh JWT token
   */
  static async refreshToken(
    data: RefreshTokenRequest
  ): Promise<ApiResponse<TokenResponse>> {
    const response = await apiClient.post('/users/refresh-token', data);

    // Update stored token
    if (response.data.success && response.data.data) {
      localStorage.setItem('auth_token', response.data.data.token);
    }

    return response.data;
  }

  // ========== AUTHENTICATED USER ENDPOINTS ==========

  /**
   * Get current user profile
   */
  static async getProfile(): Promise<ApiResponse<User>> {
    const response = await apiClient.get('/users/profile');
    return response.data;
  }

  /**
   * Update current user profile
   */
  static async updateProfile(
    data: UpdateUserRequest
  ): Promise<ApiResponse<User>> {
    const response = await apiClient.put('/users/profile', data);
    return response.data;
  }

  /**
   * Change password
   */
  static async changePassword(
    data: ChangePasswordRequest
  ): Promise<ApiResponse<null>> {
    const response = await apiClient.post('/users/change-password', data);
    return response.data;
  }

  // ========== ADMIN ENDPOINTS ==========

  /**
   * List all users (Admin only)
   */
  static async listUsers(params?: {
    offset?: number;
    limit?: number;
    search?: string;
  }): Promise<ApiResponse<ListUsersResponse>> {
    const response = await apiClient.get('/admin/users', { params });
    return response.data;
  }

  /**
   * Get user by ID (Admin only)
   */
  static async getUserById(id: number): Promise<ApiResponse<User>> {
    const response = await apiClient.get(`/admin/users/${id}`);
    return response.data;
  }

  /**
   * Update user by admin (Admin only)
   */
  static async updateUserByAdmin(
    id: number,
    data: UpdateUserByAdminRequest
  ): Promise<ApiResponse<User>> {
    const response = await apiClient.put(`/admin/users/${id}`, data);
    return response.data;
  }

  /**
   * Delete user (Admin only)
   */
  static async deleteUser(id: number): Promise<ApiResponse<null>> {
    const response = await apiClient.delete(`/admin/users/${id}`);
    return response.data;
  }

  /**
   * Assign role to user (Admin only)
   */
  static async assignRole(
    id: number,
    data: AssignRoleRequest
  ): Promise<ApiResponse<User>> {
    const response = await apiClient.post(
      `/admin/users/${id}/assign-role`,
      data
    );
    return response.data;
  }

  // ========== UTILITY METHODS ==========

  /**
   * Logout user (clear local storage)
   */
  static logout(): void {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
  }

  /**
   * Get current user from localStorage
   */
  static getCurrentUser(): User | null {
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
    return user?.role === 'admin';
  }

  /**
   * Get auth token
   */
  static getAuthToken(): string | null {
    return localStorage.getItem('auth_token');
  }
}

export default UserService;
