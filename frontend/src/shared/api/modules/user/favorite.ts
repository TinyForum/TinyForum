import apiClient from "../../client";
import { ApiResponse, PageData } from "../../types/basic.model";
import {
  AddFavoritePayload,
  ListFavoritesRequest,
  UserFavoriteVO,
} from "../../types/user.model";

export const userFavoriteApi = {
  // 我的收藏列表
  listMyFavorites: (params?: ListFavoritesRequest) =>
    apiClient.get<ApiResponse<PageData<UserFavoriteVO>>>(
      "/users/me/favorites",
      { params },
    ),
  // 收藏内容
  addFavorite: (data: AddFavoritePayload) =>
    apiClient.post<ApiResponse<{ favorited: boolean }>>(
      "/users/me/favorites",
      data,
    ),
  // 取消收藏
  removeFavorite: (targetId: number, targetType = "creation") =>
    apiClient.delete<ApiResponse<{ favorited: boolean }>>(
      `/users/me/favorites/${targetId}`,
      { params: { target_type: targetType } },
    ),
};
