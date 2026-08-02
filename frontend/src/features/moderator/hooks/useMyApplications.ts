// hooks/useMyApplications.ts
import { useCallback } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import { moderatorApi } from "@/shared/api/modules/moderator";

// 我的版主申请查询 Hook（分页）
export function useMyApplications(
  page: number,
  pageSize: number,
  enabled: boolean,
) {
  const query = useQuery({
    queryKey: ["my-applications", page],
    queryFn: async () => {
      const res = await moderatorApi.getMyApplications({
        page,
        page_size: pageSize,
      });
      return {
        list: res?.data?.data?.list || [],
        total: res?.data?.data?.total ?? 0,
      };
    },
    enabled,
  });

  // 撤销申请变更
  const cancelMutation = useMutation({
    mutationFn: (applicationId: number) =>
      moderatorApi.cancelApplication(applicationId),
  });

  const cancelApplication = useCallback(
    async (applicationId: number) => {
      await cancelMutation.mutateAsync(applicationId);
    },
    [cancelMutation],
  );

  return {
    data: query.data,
    isLoading: query.isLoading,
    error: query.error,
    refetch: query.refetch,
    cancelApplication,
  };
}
