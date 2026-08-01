// hooks/useComments.ts
import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
  UseMutationOptions,
} from "@tanstack/react-query";

import { AxiosError } from "axios";
import { commentKeys } from "./useCommentKeys";
import { commentApi } from "@/shared/api/modules/comments";
import { PageData } from "@/shared/api/types/basic.model";
import {
  CommentResponse,
  CreateCommentPayload,
  CreateCommentResponse,
} from "@/shared/api/types/comment.model";

// ---------- 查询：获取某帖子的评论列表（分页） ----------
export function useCommentsByPost(
  postId: number,
  params?: { page?: number; page_size?: number },
  options?: Omit<
    UseQueryOptions<PageData<CommentResponse>, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<PageData<CommentResponse>, AxiosError>({
    queryKey: commentKeys.list(postId, params),
    queryFn: async () => {
      const res = await commentApi.listByPost(postId, params);
      // 假设 res 为 AxiosResponse<ApiResponse<PageData<CommentResponse>>>
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取评论列表失败");
      }
      if (!res.data.data) {
        throw new Error("评论数据为空");
      }
      return res.data.data;
    },
    enabled: !!postId,
    ...options,
  });
}

// ---------- 变更：创建评论 ----------
export function useCreateComment(
  options?: UseMutationOptions<
    CreateCommentResponse,
    AxiosError,
    CreateCommentPayload
  >,
) {
  const queryClient = useQueryClient();

  return useMutation<CreateCommentResponse, AxiosError, CreateCommentPayload>({
    mutationFn: async (payload) => {
      const res = await commentApi.create(payload);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "创建评论失败");
      }
      if (!res.data.data) {
        throw new Error("创建评论返回数据为空");
      }
      return res.data.data;
    },
    onSuccess: (_, variables) => {
      // 使用 post_id 字段
      if (variables.post_id) {
        queryClient.invalidateQueries({
          queryKey: commentKeys.list(variables.post_id),
        });
      } else {
        queryClient.invalidateQueries({
          queryKey: commentKeys.lists(),
        });
      }
    },
    ...options,
  });
}

// ---------- 变更：删除评论 ----------
export function useDeleteComment(
  options?: UseMutationOptions<void, AxiosError, number>,
) {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (commentId) => {
      const res = await commentApi.delete(commentId);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "删除评论失败");
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: commentKeys.lists(),
      });
    },
    ...options,
  });
}

// ---------- 组合 Hook（可选） ----------
export function useCommentActions(postId: number) {
  const queryClient = useQueryClient();

  const createMutation = useCreateComment({
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commentKeys.list(postId) });
    },
  });

  const deleteMutation = useDeleteComment({
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: commentKeys.list(postId) });
    },
  });

  return {
    useList: (params?: { page?: number; page_size?: number }) =>
      useCommentsByPost(postId, params),
    createComment: createMutation.mutate,
    createCommentAsync: createMutation.mutateAsync,
    deleteComment: deleteMutation.mutate,
    deleteCommentAsync: deleteMutation.mutateAsync,
    isCreating: createMutation.isPending,
    isDeleting: deleteMutation.isPending,
    createError: createMutation.error,
    deleteError: deleteMutation.error,
  };
}
