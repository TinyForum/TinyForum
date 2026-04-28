import apiClient from "../client";
import { ApiResponse } from "../types";

// 设置积分请求参数
// interface SetScoreRequest {
//   operation: "set" | "add" | "subtract";
//   score: number;
//   reason: string;
// }

// 设置积分响应数据
interface SetScoreResponse {
  user_id: number;
  old_score: number;
  new_score: number;
  change: number;
  operation: string;
  operator_id: number;
  reason: string;
  timestamp: number;
}

// 单个用户积分响应
interface UserScoreResponse {
  user_id: number;
  score: number;
}

// 所有用户积分响应（数组）
type AllUserScoreResponse = Array<{
  id: number;
  username: string;
  avatar: string;
  score: number;
}>;

// 查询参数
interface GetUserScoreParams {
  id?: number; // 用户ID，可选
}

export const scoreApi = {
  /**
   * 获取所有用户积分列表
   * @param params - 可选参数，可指定用户ID
   * @returns 用户积分列表
   */
  getAllUserScore: (params?: GetUserScoreParams) =>
    apiClient.get<ApiResponse<AllUserScoreResponse>>("/admin/users/score", {
      params,
    }),

  /**
   * 获取当前用户自己的积分
   * @returns 当前用户积分
   */
  getUserScore: () =>
    apiClient.get<ApiResponse<UserScoreResponse>>("/admin/users/score"),

  /**
   * 设置用户积分为指定值
   * @param userId - 用户ID
   * @param score - 目标积分值
   * @param reason - 操作原因
   * @returns 操作后的积分信息
   */
  setUserScore: (userId: number, score: number, reason: string) =>
    apiClient.put<ApiResponse<SetScoreResponse>>(
      `/admin/users/${userId}/score`,
      {
        operation: "set",
        score: score,
        reason: reason,
      },
    ),

  /**
   * 增加用户积分
   * @param userId - 用户ID
   * @param increment - 增加的积分数量
   * @param reason - 操作原因
   * @returns 操作后的积分信息
   */
  addUserScore: (userId: number, increment: number, reason: string) =>
    apiClient.put<ApiResponse<SetScoreResponse>>(
      `/admin/users/${userId}/score`,
      {
        operation: "add",
        score: increment,
        reason: reason,
      },
    ),

  /**
   * 扣除用户积分
   * @param userId - 用户ID
   * @param decrement - 扣除的积分数量
   * @param reason - 操作原因
   * @returns 操作后的积分信息
   */
  subtractUserScore: (userId: number, decrement: number, reason: string) =>
    apiClient.put<ApiResponse<SetScoreResponse>>(
      `/admin/users/${userId}/score`,
      {
        operation: "subtract",
        score: decrement,
        reason: reason,
      },
    ),
};
