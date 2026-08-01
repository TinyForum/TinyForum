// hooks/useAnswer.ts
import { answerApi } from "@/shared/api/modules/answer";
import { CommentResponse } from "@/shared/api/types/comment.model";
import { VoteStatusResponse } from "@/shared/api/types/vote.model";
import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
  UseMutationOptions,
} from "@tanstack/react-query";
import { AxiosError } from "axios";
import { answerKeys } from "./useAnswerKey";

// 辅助：从 ResponseResultVO 中提取 data（假设所有接口返回 { code, message, data }）
function unwrap<T>(res: { code: number; message?: string; data?: T }): T {
  if (res.code !== 0) {
    throw new Error(res.message || "请求失败");
  }
  if (!res.data) {
    throw new Error("返回数据为空");
  }
  return res.data;
}

// ---------- 查询：获取单个答案 ----------
export function useAnswer(
  answerId: number,
  options?: Omit<
    UseQueryOptions<CommentResponse, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<CommentResponse, AxiosError>({
    queryKey: answerKeys.detail(answerId),
    queryFn: async () => {
      const res = await answerApi.getAnswer(answerId);
      return unwrap<CommentResponse>(res.data);
    },
    enabled: !!answerId,
    ...options,
  });
}

// ---------- 查询：获取投票状态 ----------
export function useVoteStatus(
  answerId: number,
  options?: Omit<
    UseQueryOptions<VoteStatusResponse, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<VoteStatusResponse, AxiosError>({
    queryKey: answerKeys.voteStatus(answerId),
    queryFn: async () => {
      const res = await answerApi.getVoteStatus(answerId);
      return unwrap<VoteStatusResponse>(res.data);
    },
    enabled: !!answerId,
    ...options,
  });
}

// ---------- 变更：删除答案 ----------
export function useDeleteAnswer(
  options?: UseMutationOptions<void, AxiosError, number>,
) {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (answerId) => {
      const res = await answerApi.deleteAnswer(answerId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "删除答案失败");
      }
    },
    onSuccess: (_, answerId) => {
      // 删除后使该答案的详情和投票状态失效
      queryClient.invalidateQueries({ queryKey: answerKeys.detail(answerId) });
      queryClient.invalidateQueries({
        queryKey: answerKeys.voteStatus(answerId),
      });
      // 同时使所有答案列表失效（如果页面有列表缓存）
      queryClient.invalidateQueries({ queryKey: answerKeys.lists() });
    },
    ...options,
  });
}

// ---------- 变更：投票（赞同/反对） ----------
export function useVoteAnswer(
  options?: UseMutationOptions<
    { vote_count: number; user_vote: number },
    AxiosError,
    { answerId: number; voteType: "up" | "down" }
  >,
) {
  const queryClient = useQueryClient();

  return useMutation<
    { vote_count: number; user_vote: number },
    AxiosError,
    { answerId: number; voteType: "up" | "down" }
  >({
    mutationFn: async ({ answerId, voteType }) => {
      const res = await answerApi.voteAnswer(answerId, voteType);
      return unwrap<{ vote_count: number; user_vote: number }>(res.data);
    },
    onSuccess: (data, { answerId }) => {
      // 更新该答案的详情缓存（可能包含投票数）和投票状态
      queryClient.invalidateQueries({ queryKey: answerKeys.detail(answerId) });
      queryClient.invalidateQueries({
        queryKey: answerKeys.voteStatus(answerId),
      });
    },
    ...options,
  });
}

// ---------- 变更：取消投票 ----------
export function useRemoveVote(
  options?: UseMutationOptions<
    { vote_count: number; user_vote: number },
    AxiosError,
    number
  >,
) {
  const queryClient = useQueryClient();

  return useMutation<
    { vote_count: number; user_vote: number },
    AxiosError,
    number
  >({
    mutationFn: async (answerId) => {
      const res = await answerApi.removeVote(answerId);
      return unwrap<{ vote_count: number; user_vote: number }>(res.data);
    },
    onSuccess: (data, answerId) => {
      queryClient.invalidateQueries({ queryKey: answerKeys.detail(answerId) });
      queryClient.invalidateQueries({
        queryKey: answerKeys.voteStatus(answerId),
      });
    },
    ...options,
  });
}

// ---------- 变更：采纳为正确答案 ----------
export function useAcceptAnswer(
  options?: UseMutationOptions<void, AxiosError, number>,
) {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (answerId) => {
      const res = await answerApi.acceptAnswer(answerId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "采纳失败");
      }
    },
    onSuccess: (_, answerId) => {
      // 采纳后影响该答案本身、问题详情（可能存储 acceptedAnswerId）以及回答列表
      // 这里我们使相关缓存失效，具体哪些需要失效取决于应用数据结构
      queryClient.invalidateQueries({ queryKey: answerKeys.detail(answerId) });
      // 同时使问题详情和回答列表失效（问题详情可能包含 acceptedAnswerId）
      // 我们可以利用 React Query 的全局 invalidate，但更好的方式是在问题模块中也使问题详情失效
      // 这里我们只处理答案自己的缓存，问题相关的由外层调用时一同处理
      // 但为了安全，我们可以使所有问题详情失效（不精确）
      // 或者暴露一个回调让调用方处理
    },
    ...options,
  });
}

// ---------- 变更：取消采纳 ----------
export function useUnacceptAnswer(
  options?: UseMutationOptions<void, AxiosError, number>,
) {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (answerId) => {
      const res = await answerApi.unacceptAnswer(answerId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "取消采纳失败");
      }
    },
    onSuccess: (_, answerId) => {
      queryClient.invalidateQueries({ queryKey: answerKeys.detail(answerId) });
    },
    ...options,
  });
}
