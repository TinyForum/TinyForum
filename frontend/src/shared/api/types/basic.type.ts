// 统一错误类型
export interface ApiError {
  response?: {
    data?: {
      code?: number;
      message?: string;
      errors?: Array<{
        field: string;
        message: string;
      }>;
    };
  };
  message?: string;
}

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data?: T; // 改为可选
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}
