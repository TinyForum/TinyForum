// hooks/useFollowStatus.ts
import { useQuery } from "@tanstack/react-query";
import { useAuthStore } from "@/store/auth";
import { userApi } from "@/shared/api/modules/user";
import { UserDO } from "@/shared/api/types/user.model.do";
import { userKeys } from "./useUserKeys";

// 当前用户是否已关注目标用户 Hook
export function useFollowStatus(userId: number) {
  const { user: currentUser } = useAuthStore();

  const { data: isFollowing, refetch: refetchFollowStatus } = useQuery({
    queryKey: userKeys.followStatus(userId, currentUser?.id),
    queryFn: async () => {
      if (!currentUser || currentUser.id === userId) return false;
      const res = await userApi.getFollowing(currentUser.id, {
        page: 1,
        page_size: 100,
      });
      const data = res.data.data;
      return data?.list?.some((u: UserDO) => u.id === userId) ?? false;
    },
    enabled: !!currentUser && currentUser.id !== userId,
  });

  return { isFollowing, refetchFollowStatus };
}
