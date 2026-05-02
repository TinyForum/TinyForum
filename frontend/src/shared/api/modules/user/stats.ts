import apiClient from "../../client";
import { ApiResponse, PageData } from "../../types/basic.model";
import { UserDO } from "../../types/user.model";

export const userStatsApi = {
  /**
   * 获取统计数据
   */
  statsCount: () =>
    apiClient.get<ApiResponse<PageData<UserDO>>>("/user/stats/count"),
};
