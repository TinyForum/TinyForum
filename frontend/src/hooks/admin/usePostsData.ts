import { adminApi } from "@/shared/api";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import toast from "react-hot-toast";
// 帖子数据管理 Hook
export function usePostsData(page: number, keyword: string, enabled: boolean) {
  const queryClient = useQueryClient();
  const t = useTranslations("admin");

  const { data, isLoading } = useQuery({
    queryKey: ["admin-posts", page, keyword],
    queryFn: () =>
      adminApi
        .listPosts({ page, page_size: 20, keyword })
        .then((r) => r.data.data),
    enabled,
  });

  const togglePinMutation = useMutation({
    mutationFn: (id: number) => adminApi.togglePin(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
      toast.success(t("operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // ✅ 修复：直接传递 id，因为 mutationFn 已经期望接收 id
  const handleTogglePin = (id: number) => {
    togglePinMutation.mutate(id);
  };

  return {
    posts: data?.list ?? [],
    total: data?.total ?? 0,
    isLoading,
    togglePin: handleTogglePin,
    isToggling: togglePinMutation.isPending,
  };
}
