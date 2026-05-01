"use client";

import { format, parseISO, isValid } from "date-fns";
import { zhCN } from "date-fns/locale";

// 日期格式化
export const formatDate = (dateStr: string | null): string => {
  if (!dateStr) return "未发布";
  const date = parseISO(dateStr);
  if (!isValid(date)) return "无效日期";
  return format(date, "yyyy-MM-dd HH:mm", { locale: zhCN });
};
