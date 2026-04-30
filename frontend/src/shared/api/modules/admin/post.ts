import apiClient from "../../client";
import { ApiResponse, PageData, Post } from "../../types";

export const adminPostsApi = {
  // ── 帖子管理 ──────────────────────────────────────────────────────────────
  /** 获取所有帖子列表（分页，可选关键词） */
  listPosts: (params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }) => apiClient.get<ApiResponse<PageData<Post>>>("/admin/posts", { params }),

  /** 获取待审核帖子列表 */
  listPendingPosts: (params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/admin/posts/pending", {
      params,
    }),

  /** 置顶/取消置顶帖子 */
  togglePin: (id: number) =>
    apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/pin`),

  /** 审核通过（可用于帖子、评论等） */
  approvePost: (id: number, note?: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/audit/tasks/${id}/approve`, {
      note,
    }),

  /** 审核拒绝（可附带原因） */
  rejectPost: (id: number, reason?: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/audit/tasks/${id}/reject`, {
      reason,
    }),
};
