import { useState, useCallback } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { userFavoriteApi } from "@/shared/api/modules/user/favorite";
import {
  UserFavoriteVO,
  ListFavoritesRequest,
} from "@/shared/api/types/user.model";
import { PageData } from "@/shared/api/types/basic.model";
import { userKeys } from "./useUserKeys";

type UserFavoritesPageData = PageData<UserFavoriteVO>;

type UseMeFavoritesReturn = UserFavoritesPageData & {
  isLoading: boolean;
  error: string | null;
  loadFavorites: (params?: Partial<ListFavoritesRequest>) => Promise<void>;
  loadMore: () => Promise<void>;
  refresh: () => Promise<void>;
  removeFavorite: (targetId: number) => void;
  isRemoving: boolean;
};

const DEFAULT_PARAMS: ListFavoritesRequest = {
  page: 1,
  page_size: 10,
  target_type: "creation",
};

// 获取当前用户的收藏列表（分页）+ 取消收藏
export function useMeFavorites(
  initialParams?: Partial<ListFavoritesRequest>,
): UseMeFavoritesReturn {
  const t = useTranslations("User");
  const queryClient = useQueryClient();
  const [params, setParams] = useState<ListFavoritesRequest>({
    ...DEFAULT_PARAMS,
    ...initialParams,
  });

  const { data, isFetching, error, refetch } = useQuery<
    PageData<UserFavoriteVO>
  >({
    queryKey: userKeys.favorites(params),
    queryFn: async () => {
      const res = await userFavoriteApi.listMyFavorites(params);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("favorites_load_failed"));
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

  const removeMutation = useMutation({
    mutationFn: (targetId: number) =>
      userFavoriteApi.removeFavorite(targetId, params.target_type),
    onSuccess: () => {
      // 取消收藏后使列表失效，触发重新拉取
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });

  const loadFavorites = useCallback(
    async (next?: Partial<ListFavoritesRequest>): Promise<void> => {
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
    loadFavorites,
    loadMore,
    refresh,
    removeFavorite: removeMutation.mutate,
    isRemoving: removeMutation.isPending,
  };
}
