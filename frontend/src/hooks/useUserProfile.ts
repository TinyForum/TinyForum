// hooks/useUserProfile.ts
import { useQuery } from "@tanstack/react-query";
import { userApi } from "@/lib/api";

export function useUserProfile() {
  const {
    data: profile,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["current-user-profile"],
    queryFn: () => userApi.getMeRole().then((r) => r.data.data),
    staleTime: 5 * 60 * 1000, // 5 分钟内不重新请求
    retry: 1, // 失败重试 1 次
  });

  return {
    profile,
    isLoading,

    error,
    refetch,
    // 便捷属性
    userId: profile?.user_id,
    role: profile?.role,
    // 角色判断辅助方法
    isAdmin: profile?.role === "admin" || profile?.role === "super_admin",
    isSuperAdmin: profile?.role === "super_admin",
    isModerator: profile?.role === "moderator",
    isReviewer: profile?.role === "reviewer",
    isMember: profile?.role === "member",
  };
}
