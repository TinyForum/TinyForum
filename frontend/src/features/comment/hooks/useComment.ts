// hooks/useComments.ts
import { commentApi } from "@/shared/api/modules/comments";
import { PageData } from "@/shared/api/types/basic.model";
import {
  CommentResponse,
  CreateCommentPayload,
} from "@/shared/api/types/comment.model";
import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
  UseMutationOptions,
} from "@tanstack/react-query";

import { AxiosError } from "axios";

// ---------- Query Key 工厂 ----------
export const commentKeys = {
  all: ["comments"] as const,
  lists: () => [...commentKeys.all, "list"] as const,
  list: (postId: number, params?: { page?: number; page_size?: number }) =>
    [...commentKeys.lists(), postId, params] as const,
  details: () => [...commentKeys.all, "detail"] as const,
  detail: (id: number) => [...commentKeys.details(), id] as const,
  trees: () => [...commentKeys.all, "tree"] as const,
  tree: (postId: number) => [...commentKeys.trees(), postId] as const,
};

// ---------- 类型别名 ----------
type CommentsPageData = PageData<CommentResponse>;

// ---------- 1. 获取评论列表（带分页） ----------
export function useCommentsByPost(
  postId: number,
  params?: { page?: number; page_size?: number },
  options?: Omit<
    UseQueryOptions<CommentsPageData, AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<CommentsPageData, AxiosError>({
    queryKey: commentKeys.list(postId, params),
    queryFn: async () => {
      const response = await commentApi.listByPost(postId, params);
      // 防止 response.data.data 为 undefined 或 null
      const data = response.data.data;
      if (!data) {
        return {
          list: [], // 或者 items？根据你的 PageData 定义，可能有 list 字段
          total: 0,
          page: params?.page || 1,
          page_size: params?.page_size || 10,
          totalPages: 0,
          has_more: false,
        } as CommentsPageData;
      }
      return data;
    },
    enabled: !!postId,
    ...options,
  });
}

//  获取评论树
export function useCommentsTree(
  postId: number,
  options?: Omit<
    UseQueryOptions<CommentResponse[], AxiosError>,
    "queryKey" | "queryFn"
  >,
) {
  return useQuery<CommentResponse[], AxiosError>({
    queryKey: commentKeys.tree(postId),
    queryFn: async () => {
      const response = await commentApi.listTree(postId);
      return response.data.data || [];
    },
    enabled: !!postId,
    ...options,
  });
}
// ---------- 2. 创建评论 ----------
export function useCreateComment(
  options?: UseMutationOptions<
    CommentResponse,
    AxiosError,
    CreateCommentPayload
  >,
) {
  const queryClient = useQueryClient();

  return useMutation<CommentResponse, AxiosError, CreateCommentPayload>({
    mutationFn: async (payload) => {
      const response = await commentApi.create(payload);
      const comment = response.data.data;
      if (!comment) {
        throw new Error("创建评论失败：返回数据为空");
      }
      // 类型断言：后端实际返回符合 CommentResponse 结构
      return comment as unknown as CommentResponse;
    },
    onSuccess: (_, variables) => {
      // 使用正确的字段名 post_id
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

// ---------- 3. 删除评论 ----------
export function useDeleteComment(
  options?: UseMutationOptions<void, AxiosError, number>,
) {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (commentId) => {
      await commentApi.delete(commentId);
    },
    onSuccess: () => {
      // 删除后清除所有评论列表缓存（无法获取 postId，全局失效）
      queryClient.invalidateQueries({
        queryKey: commentKeys.lists(),
      });
    },
    ...options,
  });
}

// ---------- （可选）组合 Hook：聚合同一帖子的所有操作 ----------
export function useCommentActions(postId: number) {
  const queryClient = useQueryClient();

  const createMutation = useCreateComment({
    onSuccess: () => {
      // 精确失效该帖子下的评论列表
      queryClient.invalidateQueries({
        queryKey: commentKeys.list(postId),
      });
    },
  });

  const deleteMutation = useDeleteComment({
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: commentKeys.list(postId),
      });
    },
  });

  return {
    // 查询 Hook（使用时调用）
    useList: (params?: { page?: number; page_size?: number }) =>
      useCommentsByPost(postId, params),
    // 操作方法
    createComment: createMutation.mutate,
    createCommentAsync: createMutation.mutateAsync,
    deleteComment: deleteMutation.mutate,
    deleteCommentAsync: deleteMutation.mutateAsync,
    // 状态
    isCreating: createMutation.isPending,
    isDeleting: deleteMutation.isPending,
    createError: createMutation.error,
    deleteError: deleteMutation.error,
  };
}
