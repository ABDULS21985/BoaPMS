import { get, post } from "@/lib/api-client";
import type { AuthenticateResponse, LoginRequest, TokenResponse, RefreshTokenRequest } from "@/types/auth";

export function login(data: LoginRequest) {
  return post<AuthenticateResponse>("/auth/login", data);
}

export function refreshToken(data: RefreshTokenRequest) {
  return post<TokenResponse>("/auth/refresh", { refresh_token: data.refreshToken });
}

export function validateToken() {
  return get<{ valid: boolean }>("/auth/validate");
}
