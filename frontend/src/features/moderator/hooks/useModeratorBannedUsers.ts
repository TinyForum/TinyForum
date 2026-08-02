// hooks/moderators/useModeratorBannedUsers.ts
import { useQuery } from "@tanstack/react-query";
import { moderatorApi } from "@/shared/api/modules/moderator";
import { moderatorKeys } from "./useModeratorKeys";

export function useModeratorBannedUsers(
  boardId: number,
  page: number = 1,
  enabled: boolean = true,
) {
  const query = useQuery({
    queryKey: moderatorKeys.bannedUsers(boardId, page),
    queryFn: () =>
      moderatorApi
        .getBoardBannedUsers(boardId, { page, page_size: 20 })
        .then((r) => r.data.data),
    enabled: !!boardId && enabled,
  });

  return {
    ...query,
    users: query.data?.list || [],
    total: query.data?.total || 0,
  };
}
