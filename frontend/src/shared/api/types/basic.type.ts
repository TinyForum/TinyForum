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
