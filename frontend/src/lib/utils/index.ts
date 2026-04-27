import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { formatDistanceToNow, format } from "date-fns";
import { zhCN } from "date-fns/locale";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function timeAgo(dateStr: string): string {
  // 添加参数验证
  if (!dateStr) {
    return "刚刚";
  }

  const date = new Date(dateStr);

  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    console.warn("Invalid date string in timeAgo:", dateStr);
    return "刚刚";
  }

  return formatDistanceToNow(date, { addSuffix: true, locale: zhCN });
}

export function formatDate(dateStr: string): string {
  // 添加参数验证
  if (!dateStr) {
    return "";
  }

  const date = new Date(dateStr);

  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    console.warn("Invalid date string in formatDate:", dateStr);
    return "";
  }

  return format(date, "yyyy-MM-dd HH:mm", { locale: zhCN });
}

export function truncate(str: string, maxLength: number): string {
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength) + "...";
}

export function getErrorMessage(error: unknown): string {
  if (error && typeof error === "object" && "response" in error) {
    const axiosError = error as { response?: { data?: { message?: string } } };
    return axiosError.response?.data?.message || "操作失败，请稍后重试";
  }
  if (error instanceof Error) return error.message;
  return "操作失败，请稍后重试";
}
