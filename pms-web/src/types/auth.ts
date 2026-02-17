export interface LoginRequest {
  username: string;
  password: string;
}

export interface AuthenticateResponse {
  userId: string;
  username: string;
  firstName: string;
  lastName: string;
  email: string;
  roles: string[];
  permissions: string[];
  organizationalUnit?: string;
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}

export interface TokenResponse {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface CurrentUserData {
  userId: string;
  username: string;
  email: string;
  name: string;
  firstName: string;
  lastName: string;
  roles: string[];
  permissions: string[];
  organizationalUnit?: string;
}

export interface CookieData {
  id: string;
  email: string;
  role: string;
  name: string;
  token: string;
  phone: string;
  needPasswordReset: boolean;
  expiry: string;
}

export interface ForgotPasswordRequest {
  email: string;
  clientHost?: string;
}

export interface ResetPasswordRequest {
  email: string;
  password: string;
  confirmPassword: string;
  code: string;
}
