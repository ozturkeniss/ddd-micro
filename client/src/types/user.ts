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
