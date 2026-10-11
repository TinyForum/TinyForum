import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";

// ---------- 请求参数 ----------
export interface CreateReportRequest {
  target_id: number; // 被举报对象ID
  target_type: "post" | "comment" | "user"; // 被举报对象类型
  type?: string; // 举报类型：spam/offensive/illegal/misinformation/privacy/other
  reason: string; // 举报理由（必填，2-500 字符）
  content_snapshot?: string; // 被举报内容快照
  is_anonymous?: boolean; // 是否匿名举报
}

// ---------- 响应体 ----------
export interface ReportDO {
  id: number;
  reporter_id: number;
  target_id: number;
  target_type: string;
  type: string;
  reason: string;
  status: string;
  handler_id: number | null;
  handle_note: string;
  handle_at: string | null;
  reporter_ip: string;
  is_anonymous: boolean;
  priority: number;
  content_snapshot: string;
  created_at: string;
  updated_at: string;
}

export const reportApi = {
  /** 提交举报 */
  create: (data: CreateReportRequest) =>
    apiClient.post<ApiResponse<ReportDO>>("/reports", data),
};