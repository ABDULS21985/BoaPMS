export interface BaseAPIResponse<T> {
  isSuccess: boolean;
  message: string;
  data: T;
}

export interface ResponseVm {
  isSuccess: boolean;
  message: string;
  data?: unknown;
}

export interface PaginatedResult<T> {
  items: T[];
  totalCount: number;
  pageNumber: number;
  pageSize: number;
  totalPages: number;
}

export interface EnumItem {
  id: number;
  name: string;
  description?: string;
}

export type EnumList = EnumItem[];
