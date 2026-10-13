import { useState, useCallback } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { userCommentApi } from "@/shared/api/modules/user/comment";
import {
  UserCommentVO,
  GetUserCommentsRequest,
} from "@/shared/api/types/user.model";
import { PageData } from "@/shared/api/types/basic.model";
import { userKeys } from "./useUserKeys";

type UserCommentsPageData = PageData<UserCommentVO>;

type UseMeCommentsReturn = UserCommentsPageData & {
  isLoading: boolean;
  error: string | null;
  loadComments: (params?: Partial<GetUserCommentsRequest>) => Promise<void>;
  loadMore: () => Promise<void>;
  refresh: () => Promise<void>;
};

const DEFAULT_PARAMS: GetUserCommentsRequest = {
  page: 1,
  page_size: 10,
};

// 获取当前用户的评论列表（分页）
export function useMeComments(
  initialParams?: Partial<GetUserCommentsRequest>,
): UseMeCommentsReturn {
  const t = useTranslations("User");
  const [params, setParams] = useState<GetUserCommentsRequest>({
    ...DEFAULT_PARAMS,
    ...initialParams,
  });

  const { data, isFetching, error, refetch } = useQuery<
    PageData<UserCommentVO>
  >({
    queryKey: userKeys.comments(params),
    queryFn: async () => {
      const res = await userCommentApi.listMyComments(params);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("comments_load_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("empty_response"));
      }
      return res.data.data;
    },
  });

  const list = data?.list ?? [];
  const total = data?.total ?? 0;
  const page = data?.page ?? params.page;
  const pageSize = data?.page_size ?? params.page_size;
  const hasMore = data?.has_more ?? false;

  const loadComments = useCallback(
    async (next?: Partial<GetUserCommentsRequest>): Promise<void> => {
      setParams((prev) => ({
        ...prev,
        ...next,
        page: next?.page ?? prev.page,
        page_size: next?.page_size ?? prev.page_size,
      }));
    },
    [],
  );

  const loadMore = useCallback(async (): Promise<void> => {
    if (isFetching || !hasMore) return;
    setParams((prev) => ({ ...prev, page: prev.page + 1 }));
  }, [isFetching, hasMore]);

  const refresh = useCallback(async (): Promise<void> => {
    setParams((prev) => ({ ...prev, page: 1 }));
    await refetch();
  }, [refetch]);

  return {
    list,
    total,
    page,
    page_size: pageSize,
    has_more: hasMore,
    isLoading: isFetching,
    error: error?.message ?? null,
    loadComments,
    loadMore,
    refresh,
  };
}
