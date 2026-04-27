// hooks/moderators/useModeratorReports.ts
import { useQuery } from "@tanstack/react-query";
import { moderatorApi } from "@/lib/api/modules/moderator";

export function useModeratorReports(
  boardId: number,
  page: number = 1,
  enabled: boolean = true,
) {
  const query = useQuery({
    queryKey: ["moderator", "reports", boardId, page],
    queryFn: () =>
      moderatorApi
        .getBoardReports(boardId, { page, page_size: 20 })
        .then((r) => r.data.data),
    enabled: !!boardId && enabled,
  });

  return {
    ...query,
    reports: query.data?.list || [],
    total: query.data?.total || 0,
  };
}
