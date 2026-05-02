import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import toast from "react-hot-toast";
import { adminPostsApi } from "@/shared/api/modules/admin/post";
import { Post } from "@/shared/api";
import { ApiResponse, PageData } from "@/shared/api/types/basic.model";

type AdminApiListResponse = { data: ApiResponse<PageData<Post>> };
interface UsePostsDataReturn {
  posts: Post[];
  total: number;
  isLoading: boolean;
  togglePin: (id: number) => void;
  isToggling: boolean;
}

/**
 * 管理员帖子数据管理 Hook
 */
export function usePostsData(
  page: number,
  keyword: string,
  enabled: boolean,
): UsePostsDataReturn {
  const queryClient = useQueryClient();
  const t = useTranslations("admin");
  const queryKey = ["admin-posts", page, keyword] as const;

  // 显式指定泛型：查询返回的类型是 PageData<Post>，错误类型为 Error
  const { data, isLoading } = useQuery<PageData<Post>, Error>({
    queryKey,
    queryFn: async () => {
      const response = (await adminPostsApi.listPosts({
        page,
        page_size: 20,
        keyword,
      })) as AdminApiListResponse;
      // 确保返回有效的 PageData 对象，即使接口异常也不应返回 undefined
      // 如果后端可能返回空，这里提供默认结构
      return (
        response.data.data ?? {
          list: [],
          total: 0,
          page,
          page_size: 20,
          has_more: false,
        }
      );
    },
    enabled,
  });

  const togglePinMutation = useMutation({
    mutationFn: async (id: number) => {
      await adminPostsApi.togglePin(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-posts"] });
      toast.success(t("operation_successful"));
    },
    onError: (error: unknown) => {
      console.error("Toggle pin failed:", error);
      toast.error(t("operation_failed"));
    },
  });

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
