// lib/api/modules/announcements.ts
import apiClient from "../client";
import type { ApiResponse, PageData, Announcement } from "../types";
// TODO: 添加公告模块的接口定义

export interface CreateAnnouncementPayload {
  title: string;
  content: string;
  type?: "system" | "feature" | "maintenance" | "policy";
  is_pinned?: boolean;
}

export interface UpdateAnnouncementPayload {
  title?: string;
  content?: string;
  type?: "system" | "feature" | "maintenance" | "policy";
  is_pinned?: boolean;
}

export const announcementApi = {
  // 获取公告列表
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Announcement>>>("/announcements", { params }),

  // 获取公告详情
  getById: (id: number) =>
    apiClient.get<ApiResponse<Announcement>>(`/announcements/${id}`),

  // 创建公告（管理员）
  create: (data: CreateAnnouncementPayload) =>
    apiClient.post<ApiResponse<Announcement>>("/admin/announcements", data),

  // 更新公告（管理员）
  update: (id: number, data: UpdateAnnouncementPayload) =>
    apiClient.put<ApiResponse<Announcement>>(`/admin/announcements/${id}`, data),

  // 删除公告（管理员）
  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/admin/announcements/${id}`),
};