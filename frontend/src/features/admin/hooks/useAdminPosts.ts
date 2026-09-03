// hooks/admin/useAdminPosts.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { Post } from "@/shared/api/types/post.model";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";
import { adminPostsApi } from "@/shared/api/modules/admin/post";
import { PageData } from "@/shared/api/types/basic.model";

const adminPostsKeys = {
  all: ["admin", "posts"] as const,
  lists: () => [...adminPostsKeys.all, "list"] as const,
  list: (params: object) => [...adminPostsKeys.lists(), params] as const,
  pending: () => [...adminPostsKeys.all, "pending"] as const,
};

// ========== 获取所有帖子列表 ==========
export function useAdminGetPosts(params?: {
  page?: number;
  page_size?: number;
  keyword?: string;
}) {
  const t = useTranslations("Admin");
  return useQuery({
    queryKey: adminPostsKeys.list(params || {}),
    queryFn: async () => {
      const res = await adminPostsApi.listPosts(params);
      if (res.data.code !== 0)
        throw new Error(res.data.message || t("get_post_list_failed"));
      return res.data.data as PageData<Post>;
    },
  });
}

// ========== 获取待审核帖子列表 ==========
export function useAdminPendingPosts(params?: {
  page?: number;
  page_size?: number;
  keyword?: string;
}) {
  const t = useTranslations("Admin");
  return useQuery({
    queryKey: adminPostsKeys.pending(),
    queryFn: async () => {
      const res = await adminPostsApi.listPendingPosts(params);
      if (res.data.code !== 0)
        throw new Error(res.data.message || t("get_pending_posts_failed"));
      return res.data.data as PageData<Post>;
    },
  });
}

// ========== 置顶/取消置顶 ==========
export function useAdminTogglePin() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: (id: number) => adminPostsApi.togglePin(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminPostsKeys.all });
      toast.success(t("pin_toggled"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("operation_failed"));
    },
  });
}

// ========== 审核通过 ==========
export function useAdminApprovePost() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: ({ id, note }: { id: number; note?: string }) =>
      adminPostsApi.approvePost(id, note),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminPostsKeys.all });
      queryClient.invalidateQueries({ queryKey: adminPostsKeys.pending() });
      toast.success(t("approve_success"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("approve_failed_short"));
    },
  });
}

// ========== 审核拒绝 ==========
export function useAdminRejectPost() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason?: string }) =>
      adminPostsApi.rejectPost(id, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminPostsKeys.all });
      queryClient.invalidateQueries({ queryKey: adminPostsKeys.pending() });
      toast.success(t("rejected"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("operation_failed"));
    },
  });
}
