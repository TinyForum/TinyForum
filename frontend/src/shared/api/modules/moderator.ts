// 版主
import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";
import {
  ApplyModeratorForm,
  Moderator,
  AddModeratorRequest,
  UpdatePermissionsRequest,
  BanUserRequest,
  ModeratorBoard,
  ModeratorPost,
  ModeratorReport,
  BannedUser,
  ModeratorApplication,
} from "../types/moderator.model";

// MARK: API
export const moderatorApi = {
  // ── 版主申请 ──────────────────────────────────────────────────────────────

  /**
   * 申请成为版主
   * @param boardId 板块ID
   * @param data 申请表单
   */
  applyModerator: (boardId: number, data: ApplyModeratorForm) =>
    apiClient.post<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/moderators/apply`,
      data,
    ),
  /**
   * 查看申请状态 (传递申请 ID)
   */
  getMyApplications: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<
      ApiResponse<{
        list: ModeratorApplication[];
        total: number;
        page: number;
        page_size: number;
      }>
    >("/boards/moderators/apply-status", { params }),
  /**
   * 撤销版主申请
   * @param applicationId 申请ID
   */
  cancelApplication: (applicationId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(
      `/boards/applications/${applicationId}`,
    ),

  // ── 版主管理 ──────────────────────────────────────────────────────────────

  /**
   * 获取板块版主列表
   * @param boardId 板块ID
   */
  getModerators: (boardId: number) =>
    apiClient.get<ApiResponse<Moderator[]>>(`/boards/${boardId}/moderators`),

  /**
   * 任命版主（管理员/有权限的版主）
   * @param boardId 板块ID
   * @param data 版主信息
   */
  addModerator: (boardId: number, data: AddModeratorRequest) =>
    apiClient.post<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/moderators`,
      data,
    ),

  /**
   * 移除版主
   * @param boardId 板块ID
   * @param userId 用户ID
   */
  removeModerator: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/moderators/${userId}`,
    ),

  /**
   * 更新版主权限（管理员）
   * @param boardId 板块ID
   * @param userId 用户ID
   * @param data 权限配置
   */
  updateModeratorPermissions: (
    boardId: number,
    userId: number,
    data: UpdatePermissionsRequest,
  ) =>
    apiClient.put<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/moderators/${userId}/permissions`,
      data,
    ),

  // ── 禁言管理 ──────────────────────────────────────────────────────────────

  /**
   * 禁言用户
   * @param boardId 板块ID
   * @param data 禁言信息
   */
  banUser: (boardId: number, data: BanUserRequest) =>
    apiClient.post<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/bans`,
      data,
    ),

  /**
   * 解除禁言
   * @param boardId 板块ID
   * @param userId 用户ID
   */
  unbanUser: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/bans/${userId}`,
    ),

  // ── 帖子管理（版主） ───────────────────────────────────────────────────────

  /**
   * 删除帖子（版主/管理员）
   * @param boardId 板块ID
   * @param postId 帖子ID
   */
  deletePost: (boardId: number, postId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/posts/${postId}`,
    ),

  /**
   * 置顶/取消置顶帖子（版主/管理员）
   * @param boardId 板块ID
   * @param postId 帖子ID
   * @param pinInBoard 是否置顶
   */
  pinPost: (boardId: number, postId: number, pinInBoard: boolean) =>
    apiClient.put<ApiResponse<{ message: string }>>(
      `/boards/${boardId}/posts/${postId}/pin`,
      { pin_in_board: pinInBoard },
    ),
  // 获取当前用户管理的板块
  getMyModeratorBoards: () =>
    apiClient.get<ApiResponse<ModeratorBoard>>("/boards/moderators/managed"),

  // 获取板块帖子（版主视角）
  getBoardPosts: (
    boardId: number,
    params?: {
      page?: number;
      page_size?: number;
      keyword?: string;
      status?: string;
    },
  ) =>
    apiClient.get<ApiResponse<{ list: ModeratorPost[]; total: number }>>(
      `/moderator/boards/${boardId}/posts`,
      { params },
    ),

  // 获取板块举报
  getBoardReports: (
    boardId: number,
    params?: { page?: number; page_size?: number; status?: string },
  ) =>
    apiClient.get<ApiResponse<{ list: ModeratorReport[]; total: number }>>(
      `/moderator/boards/${boardId}/reports`,
      { params },
    ),

  // 获取板块禁言用户
  getBoardBannedUsers: (
    boardId: number,
    params?: { page?: number; page_size?: number },
  ) =>
    apiClient.get<ApiResponse<{ list: BannedUser[]; total: number }>>(
      `/moderator/boards/${boardId}/bans`,
      { params },
    ),
};
