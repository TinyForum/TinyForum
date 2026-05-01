import apiClient from "../../client";
import {
  CreateAnnouncementPayload,
  AnnouncementDO,
  UpdateAnnouncementPayload,
} from "../../types/announcement.model";
import { ApiResponse } from "../../types/basic.model";

export const adminAnnouncementApi = {
  /**
   * 创建公告
   * @param data 创建公告的数据
   * @returns 创建的公告
   */
  create: (data: CreateAnnouncementPayload) =>
    apiClient.post<ApiResponse<AnnouncementDO>>("/admin/announcements", data),

  /**
   * 更新公告
   * @param id 公告ID
   * @param data 更新数据
   * @returns 更新的公告
   */
  update: (id: number, data: UpdateAnnouncementPayload) =>
    apiClient.put<ApiResponse<AnnouncementDO>>(
      `/admin/announcements/${id}`,
      data,
    ),

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
    apiClient.put<ApiResponse<null>>(`/admin/announcements/${id}/pin`, {
      pinned,
    }),
};
