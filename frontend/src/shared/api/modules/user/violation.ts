import apiClient from "../../client";
import { ApiResponse, PageRequest } from "../../types/basic.model";
import { UserStatsVO } from "../../types/user.model";

/** 用户违规 */
export const userViolationApi = {
  // 获取当前用户的违规列表
  listUserViolations: () =>
    apiClient.get<ApiResponse<ViolationVO[]>>("/users/me/violations"),

  /** 用户申诉 */
  appeal: (id: string, reason: string) =>
    apiClient.post<ApiResponse<null>>(`/users/me/violations/${id}/appeal`, {
      reason,
    }),

  // 获取单条违规详情
  getViolationDetail: (id: string) =>
    apiClient.get<ApiResponse<ViolationVO>>(`/users/me/violations/${id}`),
};

export interface ViolationVO {
  id: string;
  created_at: string;
  updated_at: string;
  user_id: string;
  operator_id: string;
  violation_type: string;
  reason: string;
  source: string;
  status: string;
  punish_type: string;
  punish_expire_at: string;
  appeal_status: string;
  appeal_reason: string;
  appeal_time: string;
  appeal_result: string;
}
