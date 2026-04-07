import apiClient from './api-client';
import type {
  ApiResponse, AuthResult, User, Post, Comment, Tag,
  Notification, PageData, PostType,
} from '@/types';

// ─── Auth ────────────────────────────────────────────────────────────────────

export const authApi = {
  register: (data: { username: string; email: string; password: string }) =>
    apiClient.post<ApiResponse<AuthResult>>('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    apiClient.post<ApiResponse<AuthResult>>('/auth/login', data),

  me: () => apiClient.get<ApiResponse<User>>('/auth/me'),
};

// ─── Posts ───────────────────────────────────────────────────────────────────

export const postApi = {
  list: (params: {
    page?: number;
    page_size?: number;
    keyword?: string;
    sort_by?: string;
    type?: PostType;
    author_id?: number;
    tag_id?: number;
  }) => apiClient.get<ApiResponse<PageData<Post>>>('/posts', { params }),

  getById: (id: number) =>
    apiClient.get<ApiResponse<{ post: Post; liked: boolean }>>(`/posts/${id}`),

  create: (data: {
    title: string;
    content: string;
    summary?: string;
    cover?: string;
    type?: PostType;
    tag_ids?: number[];
  }) => apiClient.post<ApiResponse<Post>>('/posts', data),

  update: (id: number, data: {
    title?: string;
    content?: string;
    summary?: string;
    cover?: string;
    tag_ids?: number[];
  }) => apiClient.put<ApiResponse<Post>>(`/posts/${id}`, data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}`),

  like: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/posts/${id}/like`),

  unlike: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}/like`),
};

// ─── Comments ────────────────────────────────────────────────────────────────

export const commentApi = {
  listByPost: (postId: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Comment>>>(`/comments/post/${postId}`, { params }),

  create: (data: { post_id: number; content: string; parent_id?: number }) =>
    apiClient.post<ApiResponse<Comment>>('/comments', data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/comments/${id}`),
};

// ─── Tags ────────────────────────────────────────────────────────────────────

export const tagApi = {
  list: () => apiClient.get<ApiResponse<Tag[]>>('/tags'),

  create: (data: { name: string; description?: string; color?: string }) =>
    apiClient.post<ApiResponse<Tag>>('/tags', data),

  update: (id: number, data: { name?: string; description?: string; color?: string }) =>
    apiClient.put<ApiResponse<Tag>>(`/tags/${id}`, data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/tags/${id}`),
};

// ─── Users ───────────────────────────────────────────────────────────────────

export const userApi = {
  getProfile: (id: number) =>
    apiClient.get<ApiResponse<User>>(`/users/${id}`),

  updateProfile: (data: { bio?: string; avatar?: string }) =>
    apiClient.put<ApiResponse<User>>('/users/profile', data),

  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/users/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/users/${id}/follow`),

  leaderboard: (limit?: number) =>
    apiClient.get<ApiResponse<User[]>>('/users/leaderboard', { params: { limit } }),
};

// ─── Notifications ───────────────────────────────────────────────────────────

export const notificationApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Notification>>>('/notifications', { params }),

  unreadCount: () =>
    apiClient.get<ApiResponse<{ count: number }>>('/notifications/unread-count'),

  markAllRead: () =>
    apiClient.post<ApiResponse<null>>('/notifications/read-all'),
};

// ─── Admin ───────────────────────────────────────────────────────────────────

export const adminApi = {
  listUsers: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<User>>>('/admin/users', { params }),

  setUserActive: (id: number, active: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, { active }),

  listPosts: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<Post>>>('/admin/posts', { params }),

  togglePin: (id: number) =>
    apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/pin`),
};
