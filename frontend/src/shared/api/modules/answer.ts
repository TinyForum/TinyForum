import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";
import { VoteStatusResponse } from "../types/vote.model";

export const answerApi = {
  // 获取单个答案
  getAnswer: (id: number) =>
    apiClient.get<ApiResponse<Comment>>(`/answers/${id}`),

  // 删除答案
  deleteAnswer: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/answers/${id}`),

  // 答案投票（赞同/反对）
  voteAnswer: (id: number, voteType: "up" | "down") =>
    apiClient.post<ApiResponse<{ vote_count: number; user_vote: number }>>(
      `/answers/${id}/vote`,
      {
        vote_type: voteType,
      },
    ),

  // 取消投票
  removeVote: (id: number) =>
    apiClient.delete<ApiResponse<{ vote_count: number; user_vote: number }>>(
      `/answers/${id}/vote`,
    ),

  // 接受为正确答案
  acceptAnswer: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/answers/${id}/accept`),

  // 取消接受答案
  unacceptAnswer: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/answers/${id}/unaccept`),
  getVoteStatus: (id: number) =>
    apiClient.get<ApiResponse<VoteStatusResponse>>(`/answers/${id}/status`),
};
