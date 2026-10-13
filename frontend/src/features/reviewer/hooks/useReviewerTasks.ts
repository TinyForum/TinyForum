import { useCallback } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { reviewerApi, ReviewTaskPayload } from "@/shared/api/modules/reviewer";
import { AuditTaskVO } from "@/shared/api/types/reviewer.model";
import { PageData } from "@/shared/api/types/basic.model";
import { reviewerKeys } from "./useReviewerKeys";

interface UseReviewerTasksParams {
  page?: number;
  page_size?: number;
}

type UseReviewerTasksReturn = PageData<AuditTaskVO> & {
  isLoading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  approveTask: (taskId: number, payload?: ReviewTaskPayload) => void;
  rejectTask: (taskId: number, payload?: ReviewTaskPayload) => void;
  isReviewing: boolean;
  reviewingTaskId: number | null;
};

// 获取待审核任务列表 + 通过/拒绝操作
export function useReviewerTasks(
  params: UseReviewerTasksParams = {},
): UseReviewerTasksReturn {
  const t = useTranslations("Reviewer");
  const queryClient = useQueryClient();
  const queryParams = { page: 1, page_size: 20, ...params };

  const { data, isFetching, error, refetch } = useQuery<
    PageData<AuditTaskVO>
  >({
    queryKey: reviewerKeys.taskList(queryParams),
    queryFn: async () => {
      const res = await reviewerApi.listPendingTasks(queryParams);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("load_failed"));
      }
      if (!res.data.data) {
        throw new Error(t("load_failed"));
      }
      return res.data.data;
    },
  });

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: reviewerKeys.tasks() });

  const approveMutation = useMutation({
    mutationFn: ({
      taskId,
      payload,
    }: {
      taskId: number;
      payload?: ReviewTaskPayload;
    }) => reviewerApi.approveTask(taskId, payload),
    onSuccess: invalidate,
  });

  const rejectMutation = useMutation({
    mutationFn: ({
      taskId,
      payload,
    }: {
      taskId: number;
      payload?: ReviewTaskPayload;
    }) => reviewerApi.rejectTask(taskId, payload),
    onSuccess: invalidate,
  });

  const refresh = useCallback(async () => {
    await refetch();
  }, [refetch]);

  const reviewVariables = approveMutation.isPending
    ? approveMutation.variables
    : rejectMutation.isPending
      ? rejectMutation.variables
      : undefined;

  return {
    list: data?.list ?? [],
    total: data?.total ?? 0,
    page: data?.page ?? queryParams.page,
    page_size: data?.page_size ?? queryParams.page_size,
    has_more: data?.has_more ?? false,
    isLoading: isFetching,
    error: error?.message ?? null,
    refresh,
    approveTask: (taskId, payload) =>
      approveMutation.mutate({ taskId, payload }),
    rejectTask: (taskId, payload) => rejectMutation.mutate({ taskId, payload }),
    isReviewing: approveMutation.isPending || rejectMutation.isPending,
    reviewingTaskId: reviewVariables?.taskId ?? null,
  };
}
