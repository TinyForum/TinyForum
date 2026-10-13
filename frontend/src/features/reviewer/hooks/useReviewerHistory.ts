import { useCallback } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { reviewerApi } from "@/shared/api/modules/reviewer";
import { AuditTaskVO } from "@/shared/api/types/reviewer.model";
import { PageData } from "@/shared/api/types/basic.model";
import { reviewerKeys } from "./useReviewerKeys";

interface UseReviewerHistoryParams {
  page?: number;
  page_size?: number;
}

type UseReviewerHistoryReturn = PageData<AuditTaskVO> & {
  isLoading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
};

// 获取审核历史（已处理任务）
export function useReviewerHistory(
  params: UseReviewerHistoryParams = {},
): UseReviewerHistoryReturn {
  const t = useTranslations("Reviewer");
  const queryParams = { page: 1, page_size: 20, ...params };

  const { data, isFetching, error, refetch } = useQuery<
    PageData<AuditTaskVO>
  >({
    queryKey: reviewerKeys.historyList(queryParams),
    queryFn: async () => {
      const res = await reviewerApi.listHistory(queryParams);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("load_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("load_failed"));
      }
      return res.data.data;
    },
  });

  const refresh = useCallback(async () => {
    await refetch();
  }, [refetch]);

  return {
    list: data?.list ?? [],
    total: data?.total ?? 0,
    page: data?.page ?? queryParams.page,
    page_size: data?.page_size ?? queryParams.page_size,
    has_more: data?.has_more ?? false,
    isLoading: isFetching,
    error: error?.message ?? null,
    refresh,
  };
}
