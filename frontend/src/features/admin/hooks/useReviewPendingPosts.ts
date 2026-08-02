// hooks/useReviewPendingPosts.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { adminPostsApi } from "@/shared/api/modules/admin/post";
import { Post } from "@/shared/api/types/post.model";
import { adminKeys } from "./useAdminKeys";

// 扩展 Post 类型以包含风险信息
export interface PostWithRisk extends Post {
  risk_score?: number;
  risk_reason?: string;
  risk_logs?: Array<{
    id: number;
    level: string;
    rule: string;
    matched_content: string;
    created_at: string;
  }>;
}

// 后端返回的分页数据结构
export interface PendingPostsResponse {
  list: PostWithRisk[];
  total: number;
}

// 待审核帖子 Hook：列表查询 + 审核通过/拒绝
export function useReviewPendingPosts(page: number, keyword: string) {
  const queryClient = useQueryClient();

  // 查询待审核列表
  const { data, isLoading, refetch } = useQuery<PendingPostsResponse>({
    queryKey: adminKeys.pendingPosts(page, keyword),
    queryFn: async () => {
      const res = await adminPostsApi.listPendingPosts({
        page,
        page_size: 20,
        keyword,
      });
      // 确保返回的数据结构正确，并处理可能的 undefined
      const d = res.data.data;
      return {
        list: d?.list || [],
        total: d?.total || 0,
      };
    },
    placeholderData: (prev) => prev,
  });

  // 审核通过 mutation
  const approveMutation = useMutation({
    mutationFn: ({ id, note }: { id: number; note?: string }) =>
      adminPostsApi.approvePost(id, note),
    onSuccess: () => {
      toast.success("帖子已通过审核");
      queryClient.invalidateQueries({ queryKey: adminKeys.pendingPosts() });
    },
    onError: () => toast.error("审核通过失败"),
  });

  // 审核拒绝 mutation
  const rejectMutation = useMutation({
    mutationFn: ({ id, reason }: { id: number; reason?: string }) =>
      adminPostsApi.rejectPost(id, reason),
    onSuccess: () => {
      toast.success("帖子已拒绝");
      queryClient.invalidateQueries({ queryKey: adminKeys.pendingPosts() });
    },
    onError: () => toast.error("拒绝失败"),
  });

  return {
    posts: data?.list ?? [],
    total: data?.total ?? 0,
    isLoading,
    refetch,
    approveMutation,
    rejectMutation,
  };
}
