import apiClient from "../client";
import {
  AnnouncementListParams,
  AnnouncementListResponse,
  AnnouncementDO,
} from "../types/announcement.model";
import { ApiResponse } from "../types/basic.model";

// ============ API 方法 ============

export const announcementApi = {
  // ========== 公开接口 ==========

  /**
   * 获取公告列表（支持分页和过滤）
   * @param params 查询参数
   * @returns 分页的公告列表
   */
  list: (params: AnnouncementListParams) =>
    apiClient.get<ApiResponse<AnnouncementListResponse>>("/announcements", {
      params,
    }),

  /**
   * 管理员获取公告
   */
  adminList: (params: AnnouncementListParams) =>
    apiClient.get<ApiResponse<AnnouncementListResponse>>(
      "/admin/announcements",
      { params },
    ),
  /**
   * 获取置顶公告
   * @param boardId 可选，板块ID
   * @returns 置顶公告列表
   */
  getPinned: (boardId?: number) =>
    apiClient.get<ApiResponse<AnnouncementDO[]>>("/announcements/pinned", {
      params: boardId ? { board_id: boardId } : undefined,
    }),

  /**
   * 获取公告详情
   * @param id 公告ID
   * @returns 公告详情
   */
  getById: (id: number) =>
    apiClient.get<ApiResponse<AnnouncementDO>>(`/announcements/${id}`),
};
