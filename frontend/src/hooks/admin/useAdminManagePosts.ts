// hooks/admin/useAdminManagePosts.ts
import { adminApi } from "@/shared/api";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import toast from "react-hot-toast";

interface UseAdminManagePostsOptions {
  page: number;
  pageSize?: number; // 默认 20
  keyword?: string;
  enabled?: boolean;
}

export function useAdminManagePosts({
  page,
  pageSize = 20,
  keyword = "",
  enabled = true,
}: UseAdminManagePostsOptions) {
  const queryClient = useQueryClient();
  const t = useTranslations("admin");

  // 获取待审核帖子列表
  const { data, isLoading, refetch } = useQuery({
    queryKey: ["admin-pending-posts", page, keyword],
    queryFn: () =>
      adminApi
        .listPendingPosts({ page, page_size: pageSize, keyword })
        .then((res) => res.data.data),
    enabled,
  });

  // 审核通过 Mutation
  const approveMutation = useMutation({
    mutationFn: (params: { id: number; note?: string }) =>
      adminApi.approvePost(params.id, params.note),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("approve_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 审核拒绝 Mutation（可带原因）
  const rejectMutation = useMutation({
    mutationFn: (params: { id: number; reason?: string }) =>
      adminApi.rejectPost(params.id, params.reason),
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
        items.map((item) => adminApi.approvePost(item.id, item.note)),
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
        items.map((item) => adminApi.rejectPost(item.id, item.reason)),
      ),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("batch_reject_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 刷新相关查询
  const invalidateQueries = () => {
    queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
    queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
  };

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
    // 单个操作
    approvePost,
    isApproving: approveMutation.isPending,
    rejectPost,
    isRejecting: rejectMutation.isPending,
    // 批量操作
    batchApprove,
    isBatchApproving: batchApproveMutation.isPending,
    batchReject,
    isBatchRejecting: batchRejectMutation.isPending,
  };
}
