import { post } from "@/lib/api-client";
import type { AuthenticateResponse, LoginRequest, TokenResponse, RefreshTokenRequest } from "@/types/auth";

export function login(data: LoginRequest) {
  return post<AuthenticateResponse>("/auth/login", data);
}

export function refreshToken(data: RefreshTokenRequest) {
  return post<TokenResponse>("/auth/refresh-token", data);
}

export function validateToken() {
  return post<{ valid: boolean }>("/auth/validate");
}
