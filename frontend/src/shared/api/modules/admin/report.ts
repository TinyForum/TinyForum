import apiClient from "../../client";
import { ApiResponse, PageData } from "../../types/basic.model";

export const adminReportsApi = {
  /**
   * 获取举报列表（支持分页、筛选、排序）
   */
  listReports: (params: ListReportsParams) =>
    apiClient.get<ApiResponse<PageData<ReportResponse>>>("/admin/reports", {
      params,
    }),

  /**
   * 获取单个举报详情
   */
  getReport: (id: number) =>
    apiClient.get<ApiResponse<ReportResponse>>(`/admin/reports/${id}`),

  /**
   * 处理举报（后端自动获取当前管理员ID，无需前端传 handler_id）
   */
  handleReport: (id: number, data: HandleReportRequest) =>
    apiClient.post<ApiResponse<ReportResponse>>(
      `/admin/reports/${id}/handle`,
      data,
    ),
};

// ---------- 请求参数 ----------
export interface ListReportsParams {
  page?: number; // 页码，默认1
  pageSize?: number; // 每页条数，默认10
  keyword?: string; // 全局搜索：理由(reason)、内容快照(content_snapshot)、举报人IP(reporter_ip)
  reporterId?: number; // 按举报人ID筛选（仅管理员可用）
  targetId?: number; // 按被举报对象ID筛选
  targetType?: string; // 按对象类型筛选：post, comment, user
  type?: string; // 按举报类型筛选：spam, offensive, illegal, misinformation, privacy, other
  status?: string; // 按状态筛选：pending, resolved, rejected
  priority?: 1 | 2 | 3; // 按优先级筛选：1高,2中,3低
  isAnonymous?: boolean; // 是否匿名举报
  handlerId?: number; // 按处理人ID筛选（仅管理员）
  sortBy?: "created_at" | "updated_at" | "priority" | "status" | "handle_at";
  order?: "ASC" | "DESC";
}

export interface HandleReportRequest {
  status: "resolved" | "rejected"; // 处理结果：通过或驳回（注意数据库用的 resolved/rejected）
  handle_note: string; // 处理备注（必填）
  reject_reason?: string; // 驳回原因（当 status=rejected 时建议提供，可存入 handle_note 或单独字段）
  // 注意：不传 handler_id，由后端从登录上下文获取
}

// ---------- 响应体（VO，经过脱敏/隐藏处理）----------
export interface ReportResponse {
  id: number;
  created_at: string;
  updated_at: string;

  // 业务核心字段
  target_id: number;
  target_type: string; // post/comment/user
  type: string; // 举报类型（枚举值）
  reason: string; // 举报理由（用户填写）
  status: string; // pending/resolved/rejected
  handle_note: string; // 处理备注
  handle_at: string | null; // 处理时间，未处理为 null
  is_anonymous: boolean; // 是否匿名举报
  priority: number; // 1高 2中 3低

  // 内容快照（必要字段，已脱敏处理敏感信息）
  content_snapshot: string; // 被举报内容快照（后端脱敏手机号、身份证等）

  // 举报人信息（匿名时为空）
  reporter: UserPublicVO | null;

  // 处理人信息（未处理时为空）
  handler: UserPublicVO | null;

  // 可选：被举报对象的简要信息（减少前端额外请求）
  target_info?: TargetBrief;
}

// 用户公开信息（不包含手机号、邮箱、真实姓名、IP等）
export interface UserPublicVO {
  id: number;
  nickname: string; // 昵称（推荐，不要返回真名）
  avatar: string; // 头像URL
}

// 被举报对象的摘要信息
export interface TargetBrief {
  id: number;
  title?: string; // 帖子标题（若target_type=post）
  content_preview?: string; // 内容预览（截取前100字）
  author?: UserPublicVO; // 作者信息
}
