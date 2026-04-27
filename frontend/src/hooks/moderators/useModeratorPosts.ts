// hooks/moderators/useModeratorPosts.ts
import { useQuery } from "@tanstack/react-query";
import { moderatorApi } from "@/lib/api/modules/moderator";

export function useModeratorPosts(
  boardId: number,
  page: number = 1,
  keyword: string = "",
  enabled: boolean = true,
) {
  const query = useQuery({
    queryKey: ["moderator", "posts", boardId, page, keyword],
    queryFn: () =>
      moderatorApi
        .getBoardPosts(boardId, { page, page_size: 20, keyword })
        .then((r) => r.data.data),
    enabled: !!boardId && enabled,
  });

  return {
    ...query,
    posts: query.data?.list || [],
    total: query.data?.total || 0,
  };
}
