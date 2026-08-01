// hooks/usePosts.ts
import { postApi } from "@/shared/api/modules/posts";
import { PageData } from "@/shared/api/types/basic.model";
import {
  PostListParams,
  Post,
  PostDetailResponse,
  CreatePostPayload,
  UpdatePostPayload,
} from "@/shared/api/types/post.model";
import {
  useQuery,
  useMutation,
  useQueryClient,
  UseQueryOptions,
  UseMutationOptions,
} from "@tanstack/react-query";

// ---------- Query Keys ----------
export const postKeys = {
  all: ["posts"] as const,
  lists: () => [...postKeys.all, "list"] as const,
  list: (params?: PostListParams) => [...postKeys.lists(), params] as const,
  details: () => [...postKeys.all, "detail"] as const,
  detail: (id: number) => [...postKeys.details(), id] as const,
};

// ---------- Helper: 确保 API 返回有效数据 ----------
function assertData<T>(data: T | undefined | null): T {
  if (data === undefined || data === null) {
    throw new Error("API 返回的数据为空");
  }
  return data;
}

// ---------- Helper: 链式调用用户传入的 onSuccess ----------
function invokeOnSuccess<TData, TVariables>(
  onSuccess: UseMutationOptions<TData, Error, TVariables>["onSuccess"],
  data: TData,
  variables: TVariables,
) {
  // 仅取前三个参数与旧实现兼容，忽略第四个 mutation context 参数
  const fn = onSuccess as unknown as
    | ((data: TData, variables: TVariables, context: unknown) => void)
    | undefined;
  fn?.(data, variables, undefined);
}

// ---------- Hooks ----------

/**
 * 获取帖子列表（分页 + 筛选）
 */
export const usePosts = (
  params?: PostListParams,
  options?: Omit<UseQueryOptions<PageData<Post>>, "queryKey" | "queryFn">,
) => {
  return useQuery({
    queryKey: postKeys.list(params),
    queryFn: async () => {
      const response = await postApi.list(params);
      return assertData(response.data.data);
    },
    ...options,
  });
};

/**
 * 获取单个帖子详情
 */
export const usePost = (
  id: number,
  options?: Omit<UseQueryOptions<PostDetailResponse>, "queryKey" | "queryFn">,
) => {
  return useQuery({
    queryKey: postKeys.detail(id),
    queryFn: async () => {
      const response = await postApi.getById(id);
      return assertData(response.data.data);
    },
    enabled: !!id,
    ...options,
  });
};

/**
 * 创建帖子
 */
export const useCreatePost = (
  options?: UseMutationOptions<Post, Error, CreatePostPayload>,
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CreatePostPayload) => {
      const response = await postApi.create(payload);
      return assertData(response.data.data);
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      // 安全调用用户传入的 onSuccess
      invokeOnSuccess(options?.onSuccess, data, variables);
    },
    ...options,
  });
};

/**
 * 更新帖子
 */
export const useUpdatePost = (
  options?: UseMutationOptions<
    Post,
    Error,
    { id: number; data: UpdatePostPayload }
  >,
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }) => {
      const response = await postApi.update(id, data);
      return assertData(response.data.data);
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({
        queryKey: postKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      invokeOnSuccess(options?.onSuccess, data, variables);
    },
    ...options,
  });
};

/**
 * 删除帖子
 */
export const useDeletePost = (
  options?: UseMutationOptions<null, Error, number>,
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await postApi.delete(id);
      // delete 返回 null，确保不是 undefined
      return response.data.data !== undefined ? response.data.data : null;
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      queryClient.removeQueries({ queryKey: postKeys.detail(variables) });
      invokeOnSuccess(options?.onSuccess, data, variables);
    },
    ...options,
  });
};

/**
 * 点赞帖子
 */
export const useLikePost = (
  options?: UseMutationOptions<null, Error, number>,
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await postApi.like(id);
      return response.data.data !== undefined ? response.data.data : null;
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: postKeys.detail(variables) });
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      invokeOnSuccess(options?.onSuccess, data, variables);
    },
    ...options,
  });
};

/**
 * 取消点赞
 */
export const useUnlikePost = (
  options?: UseMutationOptions<null, Error, number>,
) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await postApi.unlike(id);
      return response.data.data !== undefined ? response.data.data : null;
    },
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: postKeys.detail(variables) });
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      invokeOnSuccess(options?.onSuccess, data, variables);
    },
    ...options,
  });
};
