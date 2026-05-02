import apiClient from "../../client";
import { ApiResponse } from "../../types/basic.model";
import { ModeratorApplication, ReviewApplicationRequest } from "../moderator";

export const adminModeratorApi = {
  /**
   * 获取版主申请列表（管理员）s
   * @param params 查询参数
   */
  listApplications: (params?: {
    board_id?: number;
    status?: "pending" | "approved" | "rejected";
    page?: number;
    page_size?: number;
  }) =>
    apiClient.get<
      ApiResponse<{
        list: ModeratorApplication[];
        total: number;
        page: number;
        page_size: number;
      }>
    >("/admin/boards/applications", { params }),

  /**
   * 审批版主申请（管理员）
   * @param applicationId 申请ID
   * @param data 审批信息
   */
  reviewApplication: (applicationId: number, data: ReviewApplicationRequest) =>
    apiClient.post<ApiResponse<{ message: string }>>(
      `/admin/boards/applications/${applicationId}/review`,
      data,
    ),
};
