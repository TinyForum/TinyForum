import apiClient from "../../client";
import { ApiResponse } from "../../types/basic.model";
import { UserStatsVO } from "../../types/user.model";

export const userStatsApi = {
  // 获取当前用户的统计数据（所有统计字段）
  getUserStats: () =>
    apiClient.get<ApiResponse<UserStatsVO>>("/users/me/stats"),
};
