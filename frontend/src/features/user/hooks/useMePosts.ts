import { userPostApi } from "@/shared/api/modules/user/post";
import { useState, useCallback } from "react";
import {
  UserPostsVO,
  GetUserPostsRequest,
} from "@/shared/api/types/user.model";
import { PageData } from "@/shared/api/types/basic.model";

// 分页数据核心字段（复用后端 PageData 结构）
type UserPostsPageData = PageData<UserPostsVO>;

// Hook 返回类型：分页数据 + 控制字段
type UseMePostsReturn = UserPostsPageData & {
  isLoading: boolean;
  error: string | null;
  loadPosts: (params?: Partial<GetUserPostsRequest>) => Promise<void>;
  loadMore: () => Promise<void>;
  refresh: () => Promise<void>;
};

// 默认请求参数（避免 undefined 传递给 API）
const DEFAULT_PARAMS: GetUserPostsRequest = {
  page: 1,
  page_size: 10,
};

export function useMePosts(
  initialParams?: Partial<GetUserPostsRequest>,
): UseMePostsReturn {
  // 分页数据状态
  const [list, setList] = useState<UserPostsVO[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [hasMore, setHasMore] = useState<boolean>(false);
  // 请求状态
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  // 保存当前完整查询参数（合并默认值）
  const [currentParams, setCurrentParams] = useState<GetUserPostsRequest>({
    ...DEFAULT_PARAMS,
    ...initialParams,
  });

  // 核心请求方法
  const loadPosts = useCallback(
    async (params?: Partial<GetUserPostsRequest>) => {
      // 合并参数：当前参数 + 传入的部分参数
      const finalParams: GetUserPostsRequest = {
        ...currentParams,
        ...params,
        // 确保 page 和 page_size 为数字（如果传入 undefined，则回退到 currentParams 的值）
        page: params?.page ?? currentParams.page,
        page_size: params?.page_size ?? currentParams.page_size,
      };

      setIsLoading(true);
      setError(null);
      try {
        const response = await userPostApi.getUserPosts(finalParams);
        if (response.status === 200 && response.data.code === 0) {
          // 关键修复：处理 response.data.data 可能为 undefined 的情况
          const data: PageData<UserPostsVO> | undefined = response.data.data;
          if (!data) {
            throw new Error("返回数据为空");
          }
          setList(data.list);
          setTotal(data.total);
          setPage(data.page);
          setPageSize(data.page_size);
          setHasMore(data.has_more);
          // 更新当前参数（用于后续刷新或加载更多）
          setCurrentParams(finalParams);
        } else {
          throw new Error(response.data.message || "获取帖子列表失败");
        }
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : String(err);
        setError(errorMessage);

        // 出错时清空列表，防止展示旧数据
        setList([]);
        setTotal(0);
        setHasMore(false);
      } finally {
        setIsLoading(false);
      }
    },
    [currentParams],
  );

  // 加载下一页（基于当前 page+1）
  const loadMore = useCallback(async () => {
    if (isLoading || !hasMore) return;
    await loadPosts({ page: page + 1 });
  }, [isLoading, hasMore, page, loadPosts]);

  // 刷新（重新加载当前参数的第一页）
  const refresh = useCallback(async () => {
    await loadPosts({ page: 1 });
  }, [loadPosts]);

  return {
    list,
    total,
    page,
    page_size: pageSize,
    has_more: hasMore,
    isLoading,
    error,
    loadPosts,
    loadMore,
    refresh,
  };
}
