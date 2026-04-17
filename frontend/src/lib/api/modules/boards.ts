/**
 * api/modules/boards.ts
 * 板块 + 版主管理 + 禁言管理 + 帖子管理
 */

import { UserRoleType } from "@/type/roles.types";
import apiClient from "../client";
import type {
  ApiResponse,
  PageData,
  Board,
  Moderator,
  ModeratorApplication,
  Post,
  // UserRole,
} from "../types";

// ─── Payloads ─────────────────────────────────────────────────────────────────

export interface CreateBoardPayload {
  name: string;
  slug: string;
  description?: string;
  icon?: string;
  cover?: string;
  parent_id?: number;
  sort_order?: number;
  view_role?: UserRoleType;
  post_role?: UserRoleType;
  reply_role?: UserRoleType;
}

export type UpdateBoardPayload = Partial<CreateBoardPayload>;

export interface AddModeratorPayload {
  user_id: number;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface BanUserPayload {
  user_id: number;
  reason: string;
  expires_at?: string;
}

export interface PinPostPayload {
  pin_in_board: boolean;
}

// ─── API ──────────────────────────────────────────────────────────────────────

export const boardApi = {
  // ── 读取 ──────────────────────────────────────────────────────────────────────
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Board>>>("/boards", { params }),

  getTree: () =>
    apiClient.get<ApiResponse<Board[]>>("/boards/tree"),

  getById: (id: number | string) =>
    apiClient.get<ApiResponse<Board>>(`/boards/${id}`),

  getBySlug: (slug: string) =>
    apiClient.get<ApiResponse<Board>>(`/boards/slug/${slug}`),

  getPostsBySlug: (slug: string, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Post>>>(`/boards/slug/${slug}/posts`, {
      params,
    }),

  // ── 增删改（Admin） ───────────────────────────────────────────────────────────
  create: (data: CreateBoardPayload) =>
    apiClient.post<ApiResponse<Board>>("/boards", data),

  update: (id: number, data: UpdateBoardPayload) =>
    apiClient.put<ApiResponse<Board>>(`/boards/${id}`, data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/boards/${id}`),

  // ── 版主管理 ──────────────────────────────────────────────────────────────────
  getModerators: (boardId: number) =>
    apiClient.get<ApiResponse<Moderator[]>>(`/boards/${boardId}/moderators`),

  addModerator: (boardId: number, data: AddModeratorPayload) =>
    apiClient.post<ApiResponse<null>>(`/boards/${boardId}/moderators`, data),

  removeModerator: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<null>>(
      `/boards/${boardId}/moderators/${userId}`
    ),

  // ── 申请版主 ──────────────────────────────────────────────────────────────────
  // applyForModerator: (boardId: number, reason: string) =>
  //   apiClient.post<ApiResponse<null>>(
  //     `/boards/${boardId}/moderator-apply`,
  //     { reason }
  //   ),

  // checkApplicationStatus: (boardId: number) =>
  //   apiClient.get<
  //     ApiResponse<{
  //       has_applied: boolean;
  //       application?: ModeratorApplication;
  //     }>
  //   >(`/boards/${boardId}/application-status`),

  // getMyApplications: (params?: { page?: number; page_size?: number }) =>
  //   apiClient.get<ApiResponse<PageData<ModeratorApplication>>>(
  //     "/boards/my-applications",
  //     { params }
  //   ),

  // checkModeratorStatus: (boardId: number) =>
  //   apiClient.get<
  //     ApiResponse<{ is_moderator: boolean; moderator?: Moderator }>
  //   >(`/boards/${boardId}/moderator-status`),

  // ── 禁言管理 ──────────────────────────────────────────────────────────────────
  banUser: (boardId: number, data: BanUserPayload) =>
    apiClient.post<ApiResponse<null>>(`/boards/${boardId}/bans`, data),

  unbanUser: (boardId: number, userId: number) =>
    apiClient.delete<ApiResponse<null>>(`/boards/${boardId}/bans/${userId}`),

  // ── 帖子管理（版主） ──────────────────────────────────────────────────────────
  deletePost: (boardId: number, postId: number) =>
    apiClient.delete<ApiResponse<null>>(`/boards/${boardId}/posts/${postId}`),

  pinPost: (boardId: number, postId: number, data: PinPostPayload) =>
    apiClient.put<ApiResponse<null>>(
      `/boards/${boardId}/posts/${postId}/pin`,
      data
    ),
};