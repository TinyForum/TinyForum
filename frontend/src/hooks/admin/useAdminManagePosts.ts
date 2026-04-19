// 管理员审核用户的帖子
import { adminApi } from "@/lib/api";
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

  // 审核帖子（通过/拒绝 - 根据后端 API 扩展）
  // 当前 reviewPosts 只传 id，若需传递状态可修改
  const reviewMutation = useMutation({
    mutationFn: (id: number) => adminApi.reviewPosts(id,{status}),
    onSuccess: () => {
      // 刷新待审核列表和普通帖子列表（因为帖子状态会变化）
      queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
      queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
      toast.success(t("operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 如果需要支持批量审核，可添加批量方法
  const batchReviewMutation = useMutation({
    mutationFn: (ids: number[]) =>
      Promise.all(ids.map((id) => adminApi.reviewPosts(id,{status}))),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
      queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
      toast.success(t("batch_operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  const handleReview = (id: number) => {
    reviewMutation.mutate(id);
  };

  const handleBatchReview = (ids: number[]) => {
    batchReviewMutation.mutate(ids);
  };

  return {
    posts: data?.list ?? [],
    total: data?.total ?? 0,
    isLoading,
    refetch,
    reviewPost: handleReview,
    isReviewing: reviewMutation.isPending,
    batchReviewPosts: handleBatchReview,
    isBatchReviewing: batchReviewMutation.isPending,
  };
}