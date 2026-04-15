import apiClient from "../client";
import { ApiResponse } from "../types";

export interface ApplyModeratorForm {
  reason: string;
  req_delete_post?: boolean;
  req_pin_post?: boolean;
  req_edit_any_post?: boolean;
  req_manage_moderator?: boolean;
  req_ban_user?: boolean;
}

export interface ReviewApplicationRequest {
  approve: boolean;
  review_note?: string;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface AddModeratorRequest {
  user_id: number;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface UpdatePermissionsRequest {
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface BanUserRequest {
  user_id: number;
  reason: string;
  expires_at?: string;
}

export interface ModeratorApplication {
  id: number;
  user_id: number;
  username: string;
  board_id: number;
  board_name: string;
  reason: string;
  status: "pending" | "approved" | "rejected" | "canceled";
  review_note: string;
  reviewed_by?: number;
  reviewed_at?: string;
  created_at: string;
  updated_at: string;
  // 添加缺失的字段
  board?: {
    id: number;
    name: string;
    slug: string;
  };
  req_delete_post: boolean;
  req_pin_post: boolean;
  req_edit_any_post: boolean;
  req_manage_moderator: boolean;
  req_ban_user: boolean;
}

export interface Moderator {
  id: number;
  user_id: number;
  created_at: string;
  board_id: number;
  permissions: {
    can_delete_post: boolean;
    can_pin_post: boolean;
    can_edit_any_post: boolean;
    can_manage_moderator: boolean;
    can_ban_user: boolean;
  };
  user?: {
    id: number;
    username: string;
    avatar?: string;
  };
  board?: {
    id: number;
    name: string;
  };
}

export interface BanRecord {
  id: number;
  user_id: number;
  board_id: number;
  moderator_id: number;
  reason: string;
  expires_at?: string;
  is_active: boolean;
  created_at: string;
}

// 申请状态
type ApplicationStatus = 'pending' | 'approved' | 'rejected' | 'canceled';

// 申请状态详情
interface ApplicationStatusDetailResponse {
  has_application: boolean;
  application_id?: number;
  status?: ApplicationStatus;
  reason?: string;
  created_at?: string;
  review_note?: string;
  reviewer_id?: number;
  reviewed_at?: string | null;
  can_cancel: boolean;
  can_resubmit: boolean;
  requested_perms?: {
    delete_post: boolean;
    pin_post: boolean;
    edit_any_post: boolean;
    manage_moderator: boolean;
    ban_user: boolean;
  };
  can_apply: boolean;
}

export const moderatorApi = {
  // ── 版主申请 ──────────────────────────────────────────────────────────────

  /**
   * 申请成为版主
   * @param boardId 板块ID
   * @param data 申请表单
   */
  applyModerator: (boardId: number, data: ApplyModeratorForm) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/boards/${boardId}/moderators/apply`, data),
/**
 * 查看申请状态 (传递申请 ID)
 */
getMyApplications: (params?: { page?: number; page_size?: number }) =>
  apiClient.get<ApiResponse<{ list: ModeratorApplication[]; total: number; page: number; page_size: number }>>(
    "/boards/moderators/apply",
    { params }
  ),
  /**
   * 撤销版主申请
   * @param applicationId 申请ID
   */
  cancelApplication: (applicationId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(`/boards/applications/${applicationId}`),

  /**
   * 获取版主申请列表（管理员）
   * @param params 查询参数
   */
  listApplications: (params?: {
    board_id?: number;
    status?: "pending" | "approved" | "rejected";
    page?: number;
    page_size?: number;
  }) =>
    apiClient.get<ApiResponse<{ list: ModeratorApplication[]; total: number; page: number; page_size: number }>>(
      "/admin/boards/applications",
      { params }
    ),

  /**
   * 审批版主申请（管理员）
   * @param applicationId 申请ID
   * @param data 审批信息
   */
  reviewApplication: (applicationId: number, data: ReviewApplicationRequest) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/admin/boards/applications/${applicationId}/review`, data),

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
    apiClient.post<ApiResponse<{ message: string }>>(`/boards/${boardId}/moderators`, data),

  /**
   * 移除版主
   * @param boardId 板块ID
   * @param userId 用户ID
   */
  removeModerator: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(`/boards/${boardId}/moderators/${userId}`),

  /**
   * 更新版主权限（管理员）
   * @param boardId 板块ID
   * @param userId 用户ID
   * @param data 权限配置
   */
  updateModeratorPermissions: (boardId: number, userId: number, data: UpdatePermissionsRequest) =>
    apiClient.put<ApiResponse<{ message: string }>>(`/boards/${boardId}/moderators/${userId}/permissions`, data),

  // ── 禁言管理 ──────────────────────────────────────────────────────────────

  /**
   * 禁言用户
   * @param boardId 板块ID
   * @param data 禁言信息
   */
  banUser: (boardId: number, data: BanUserRequest) =>
    apiClient.post<ApiResponse<{ message: string }>>(`/boards/${boardId}/bans`, data),

  /**
   * 解除禁言
   * @param boardId 板块ID
   * @param userId 用户ID
   */
  unbanUser: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(`/boards/${boardId}/bans/${userId}`),

  // ── 帖子管理（版主） ───────────────────────────────────────────────────────

  /**
   * 删除帖子（版主/管理员）
   * @param boardId 板块ID
   * @param postId 帖子ID
   */
  deletePost: (boardId: number, postId: number) =>
    apiClient.delete<ApiResponse<{ message: string }>>(`/boards/${boardId}/posts/${postId}`),

  /**
   * 置顶/取消置顶帖子（版主/管理员）
   * @param boardId 板块ID
   * @param postId 帖子ID
   * @param pinInBoard 是否置顶
   */
  pinPost: (boardId: number, postId: number, pinInBoard: boolean) =>
    apiClient.put<ApiResponse<{ message: string }>>(`/boards/${boardId}/posts/${postId}/pin`, { pin_in_board: pinInBoard }),
};