import axios, { type AxiosInstance, type AxiosRequestConfig } from "axios";
import { getSession } from "next-auth/react";

const apiClient: AxiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1",
  headers: { "Content-Type": "application/json" },
  timeout: 30000,
});

// Request interceptor — inject JWT
apiClient.interceptors.request.use(async (config) => {
  if (typeof window !== "undefined") {
    const session = await getSession();
    if (session?.accessToken) {
      config.headers.Authorization = `Bearer ${session.accessToken}`;
    }
  }
  return config;
});

// Response interceptor — handle errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 && typeof window !== "undefined") {
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export async function get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.get<T>(url, config);
  return res.data;
}

export async function post<R, T = unknown>(url: string, data?: T, config?: AxiosRequestConfig): Promise<R> {
  const res = await apiClient.post<R>(url, data, config);
  return res.data;
}

export async function put<R, T = unknown>(url: string, data?: T, config?: AxiosRequestConfig): Promise<R> {
  const res = await apiClient.put<R>(url, data, config);
  return res.data;
}

export async function del<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.delete<T>(url, config);
  return res.data;
}

export async function postFormData<R>(url: string, formData: FormData): Promise<R> {
  const res = await apiClient.post<R>(url, formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
  return res.data;
}

export default apiClient;
