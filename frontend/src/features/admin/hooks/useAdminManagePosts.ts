import { Post } from "@/shared/api";
import { adminPostsApi } from "@/shared/api/modules/admin/post";
import { PageData } from "@/shared/api/types/basic.model";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import toast from "react-hot-toast";

interface UseAdminManagePostsOptions {
  page: number;
  pageSize?: number;
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

  const invalidateQueries = () => {
    queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
    queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
  };

  // 获取待审核帖子列表
  const { data, isLoading, refetch } = useQuery<PageData<Post>>({
    queryKey: ["admin-pending-posts", page, keyword],
    queryFn: async (): Promise<PageData<Post>> => {
      const response = await adminPostsApi.listPendingPosts({
        page,
        page_size: pageSize,
        keyword,
      });
      const pageData = response.data?.data;
      if (pageData && typeof pageData === "object" && "list" in pageData) {
        return pageData as PageData<Post>;
      }
      // 降级默认值（后端返回异常时使用）
      return {
        list: [],
        total: 0,
        page: page,
        page_size: pageSize,
        has_more: false,
      };
    },
    enabled,
  });

  // 以下 mutations 保持不变...
  const approveMutation = useMutation({
    mutationFn: (params: { id: number; note?: string }) =>
      adminPostsApi.approvePost(params.id, params.note),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("approve_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  const rejectMutation = useMutation({
    mutationFn: (params: { id: number; reason?: string }) =>
      adminPostsApi.rejectPost(params.id, params.reason),
    onSuccess: () => {
      invalidateQueries();
      toast.success(t("reject_success"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

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
