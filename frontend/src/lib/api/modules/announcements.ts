// lib/api/modules/announcements.ts
import apiClient from "../client";
import type { ApiResponse, PageData } from "../types";

// ============ 类型定义 ============

// 公告类型枚举
export type AnnouncementType = "normal" | "important" | "emergency" | "event";

// 公告状态枚举
export type AnnouncementStatus = "draft" | "published" | "archived";

// 公告数据结构
export interface Announcement {
  id: number;
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: AnnouncementType;
  status: AnnouncementStatus;
  is_pinned: boolean;
  is_global: boolean;
  board_id: number | null;
  published_at: string | null;
  expired_at: string | null;
  view_count: number;
  created_by: number;
  updated_by: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  
  // 关联数据
  board?: {
    id: number;
    name: string;
    slug: string;
  } | null;
  creator?: {
    id: number;
    username: string;
    avatar?: string;
  } | null;
}

// ============ 请求参数类型 ============

// 创建公告请求
export interface CreateAnnouncementPayload {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  type?: AnnouncementType;
  is_pinned?: boolean;
  status?: AnnouncementStatus;
  is_global?: boolean;
  board_id?: number | null;
  published_at?: string | null;
  expired_at?: string | null;
}

// 更新公告请求
export interface UpdateAnnouncementPayload {
  title?: string;
  content?: string;
  summary?: string;
  cover?: string;
  type?: AnnouncementType;
  is_pinned?: boolean;
  status?: AnnouncementStatus;
  is_global?: boolean;
  board_id?: number | null;
  published_at?: string | null;
  expired_at?: string | null;
}

// 公告列表查询参数
export interface AnnouncementListParams {
  page?: number;
  page_size?: number;
  board_id?: number;
  type?: AnnouncementType;
  status?: AnnouncementStatus;
  is_pinned?: boolean;
  is_global?: boolean;
  keyword?: string;
  start_time?: string;
  end_time?: string;
}

// 置顶/取消置顶请求
export interface PinAnnouncementPayload {
  pinned: boolean;
}

// ============ API 响应类型 ============

// 分页响应数据
export interface AnnouncementListResponse {
  list: Announcement[];
  total: number;
  page: number;
  page_size: number;
}

// ============ API 方法 ============

export const announcementApi = {
  // ========== 公开接口 ==========
  
  /**
   * 获取公告列表（支持分页和过滤）
   * @param params 查询参数
   * @returns 分页的公告列表
   */
  list: (params: AnnouncementListParams) =>
    apiClient.get<ApiResponse<AnnouncementListResponse>>("/announcements", { params }),

  /**
   * 管理员获取公告
   */
  adminList: (params: AnnouncementListParams) => apiClient.get<ApiResponse<AnnouncementListResponse>>("/admin/announcements", { params }),
  /**
   * 获取置顶公告
   * @param boardId 可选，板块ID
   * @returns 置顶公告列表
   */
  getPinned: (boardId?: number) =>
    apiClient.get<ApiResponse<Announcement[]>>("/announcements/pinned", {
      params: boardId ? { board_id: boardId } : undefined,
    }),

  /**
   * 获取公告详情
   * @param id 公告ID
   * @returns 公告详情
   */
  getById: (id: number) =>
    apiClient.get<ApiResponse<Announcement>>(`/announcements/${id}`),

  // ========== 管理接口（需要管理员权限） ==========

  /**
   * 创建公告
   * @param data 创建公告的数据
   * @returns 创建的公告
   */
  create: (data: CreateAnnouncementPayload) =>
    apiClient.post<ApiResponse<Announcement>>("/admin/announcements", data),

  /**
   * 更新公告
   * @param id 公告ID
   * @param data 更新数据
   * @returns 更新的公告
   */
  update: (id: number, data: UpdateAnnouncementPayload) =>
    apiClient.put<ApiResponse<Announcement>>(`/admin/announcements/${id}`, data),

  /**
   * 删除公告
   * @param id 公告ID
   */
  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/admin/announcements/${id}`),

  /**
   * 发布公告
   * @param id 公告ID
   */
  publish: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/admin/announcements/${id}/publish`),

  /**
   * 归档公告
   * @param id 公告ID
   */
  archive: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/admin/announcements/${id}/archive`),

  /**
   * 置顶/取消置顶公告
   * @param id 公告ID
   * @param pinned 是否置顶
   */
  pin: (id: number, pinned: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/announcements/${id}/pin`, { pinned }),
};

// ============ 辅助函数 ============

/**
 * 获取公告类型的显示文本
 */
export function getAnnouncementTypeText(type: AnnouncementType): string {
  const typeMap: Record<AnnouncementType, string> = {
    normal: "普通公告",
    important: "重要公告",
    emergency: "紧急公告",
    event: "活动公告",
  };
  return typeMap[type] || type;
}

/**
 * 获取公告类型的颜色样式
 */
export function getAnnouncementTypeColor(type: AnnouncementType): string {
  const colorMap: Record<AnnouncementType, string> = {
    normal: "blue",
    important: "orange",
    emergency: "red",
    event: "green",
  };
  return colorMap[type] || "default";
}

/**
 * 获取公告状态的显示文本
 */
export function getAnnouncementStatusText(status: AnnouncementStatus): string {
  const statusMap: Record<AnnouncementStatus, string> = {
    draft: "草稿",
    published: "已发布",
    archived: "已归档",
  };
  return statusMap[status] || status;
}

/**
 * 获取公告状态的颜色样式
 */
export function getAnnouncementStatusColor(status: AnnouncementStatus): string {
  const colorMap: Record<AnnouncementStatus, string> = {
    draft: "gray",
    published: "green",
    archived: "default",
  };
  return colorMap[status] || "default";
}

/**
 * 格式化公告发布时间
 */
export function formatAnnouncementTime(dateStr: string | null): string {
  if (!dateStr) return "未发布";
  const date = new Date(dateStr);
  return date.toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}