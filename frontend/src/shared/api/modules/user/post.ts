import apiClient from "../../client";
import { ApiResponse, PageData } from "../../types/basic.model";
import { GetUserPostsRequest, UserPostsVO } from "../../types/user.model";

export const userPostApi = {
  // 获取当前用户的帖子列表（支持卡片视图可选）
  getUserPosts: (params?: GetUserPostsRequest) =>
    apiClient.get<ApiResponse<PageData<UserPostsVO>>>("/users/me/posts", {
      params,
    }),
  // 获取当前用户的单个帖子（支持卡片视图可选）
  getUserPostById: (postId: number, view?: "card" | "full") =>
    apiClient.get<ApiResponse<UserPostsVO>>(`/users/me/posts/${postId}`, {
      params: { view },
    }),
};
