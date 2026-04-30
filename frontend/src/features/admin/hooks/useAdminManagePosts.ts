import { adminPostsApi } from "@/shared/api/modules/admin/post";
import type { ApiResponse, PageData, Post } from "@/shared/api/types"; // 假设类型从这个路径导出
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import toast from "react-hot-toast";

interface UseAdminManagePostsOptions {
  page: number;
  pageSize?: number;
  keyword?: string;
  enabled?: boolean;
}

// 定义 listPendingPosts 的返回类型，根据实际 API 结构调整
// 假设 adminApi.listPendingPosts 返回 Promise<{ data: ApiResponse<PageData<Post>> }>
type ListPendingPostsResponse = { data: ApiResponse<PageData<Post>> };

export function useAdminManagePosts({
  page,
  pageSize = 20,
  keyword = "",
  enabled = true,
}: UseAdminManagePostsOptions) {
  const queryClient = useQueryClient();
  const t = useTranslations("admin");

  // 刷新相关查询的辅助函数（先定义，以便在 mutation 的 onSuccess 中使用）
  const invalidateQueries = () => {
    queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
    queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
  };

  // 获取待审核帖子列表 - 显式指定泛型并修复 res 类型
  const { data, isLoading, refetch } = useQuery<PageData<Post>, Error>({
    queryKey: ["admin-pending-posts", page, keyword],
    queryFn: async () => {
      const response = (await adminPostsApi.listPendingPosts({
        page,
        page_size: pageSize,
        keyword,
      })) as ListPendingPostsResponse;
      // response.data.data 就是 PageData<Post>，如果可能为空则提供默认值
      return (
        response.data.data ?? {
          list: [],
          total: 0,
          page: 1,
          page_size: pageSize,
        }
      );
    },
    enabled,
  });

  // 审核通过 Mutation
  const approveMutation = useMutation({
    mutationFn: (params: { id: number; note?: string }) =>
      adminPostsApi.approvePost(params.id, params.note),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("approve_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 审核拒绝 Mutation（可带原因）
  const rejectMutation = useMutation({
    mutationFn: (params: { id: number; reason?: string }) =>
      adminPostsApi.rejectPost(params.id, params.reason),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("reject_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 批量审核通过
  const batchApproveMutation = useMutation({
    mutationFn: (items: Array<{ id: number; note?: string }>) =>
      Promise.all(
        items.map((item) => adminPostsApi.approvePost(item.id, item.note)),
      ),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("batch_approve_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 批量审核拒绝
  const batchRejectMutation = useMutation({
    mutationFn: (items: Array<{ id: number; reason?: string }>) =>
      Promise.all(
        items.map((item) => adminPostsApi.rejectPost(item.id, item.reason)),
      ),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("batch_reject_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 便捷方法
  const approvePost = (id: number, note?: string) => {
    approveMutation.mutate({ id, note });
  };

  const rejectPost = (id: number, reason?: string) => {
    rejectMutation.mutate({ id, reason });
  };

  const batchApprove = (ids: number[], note?: string) => {
    const items = ids.map((id) => ({ id, note }));
    batchApproveMutation.mutate(items);
  };

  const batchReject = (ids: number[], reason?: string) => {
    const items = ids.map((id) => ({ id, reason }));
    batchRejectMutation.mutate(items);
  };

  return {
    posts: data?.list ?? [],
    total: data?.total ?? 0,
    isLoading,
    refetch,
    approvePost,
    isApproving: approveMutation.isPending,
    rejectPost,
    isRejecting: rejectMutation.isPending,
    batchApprove,
    isBatchApproving: batchApproveMutation.isPending,
    batchReject,
    isBatchRejecting: batchRejectMutation.isPending,
  };
}
