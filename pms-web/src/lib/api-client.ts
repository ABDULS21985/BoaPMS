import axios, { type AxiosInstance, type AxiosRequestConfig } from "axios";
import { getSession } from "next-auth/react";

// Convert snake_case keys to camelCase recursively.
// The Go API returns snake_case JSON; the frontend types expect camelCase.
function toCamelCase(str: string): string {
  return str.replace(/_([a-z0-9])/g, (_, c) => c.toUpperCase());
}

function camelCaseKeys(obj: unknown): unknown {
  if (Array.isArray(obj)) return obj.map(camelCaseKeys);
  if (obj !== null && typeof obj === "object" && !(obj instanceof Date)) {
    return Object.fromEntries(
      Object.entries(obj as Record<string, unknown>).map(([k, v]) => [
        toCamelCase(k),
        camelCaseKeys(v),
      ])
    );
  }
  return obj;
}

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

// Response interceptor — convert snake_case to camelCase and handle errors
apiClient.interceptors.response.use(
  (response) => {
    if (response.data) {
      response.data = camelCaseKeys(response.data);
    }
    return response;
  },
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
