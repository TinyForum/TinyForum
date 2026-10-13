import apiClient from "../../client";
import { ApiResponse, PageData } from "../../types/basic.model";
import {
  GetUserCommentsRequest,
  UserCommentVO,
} from "../../types/user.model";

export const userCommentApi = {
  // 获取当前用户发表的评论列表
  listMyComments: (params?: GetUserCommentsRequest) =>
    apiClient.get<ApiResponse<PageData<UserCommentVO>>>("/users/me/comments", {
      params,
    }),
};
