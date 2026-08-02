import { useCallback } from "react"
import { useQuery } from "@tanstack/react-query"
import { userApi } from "@/shared/api/modules/user"
import { RoleResponse } from "@/shared/api/types/user.model"
import { userKeys } from "./useUserKeys"

// ========== 用户角色 ==========
interface UseUserRoleReturn {
  role: RoleResponse | null
  isLoading: boolean
  error: string | null
  loadRole: () => Promise<void>
  isAdmin: boolean
  isModerator: boolean
  isUser: boolean
  isMember: boolean
  isSuperAdmin: boolean
  isReviewer: boolean
  isSystemMaintainer: boolean
}

export function useUserRole(): UseUserRoleReturn {
  // 查询：拉取当前用户角色（默认自动加载）
  const { data: role, isLoading, error, refetch } = useQuery<RoleResponse>({
    queryKey: userKeys.role(),
    queryFn: async () => {
      const res = await userApi.getMeRole()
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取角色失败")
      }
      if (!res.data.data) {
        throw new Error("获取角色失败")
      }
      return res.data.data
    },
  })

  // 命令式入口：重新拉取角色
  const loadRole = useCallback(async (): Promise<void> => {
    await refetch()
  }, [refetch])

  return {
    role: role ?? null,
    isLoading,
    error: error?.message ?? null,
    loadRole,
    isAdmin: role?.role === "admin",
    isModerator: role?.role === "moderator",
    isUser: role?.role === "user",
    isMember: role?.role === "member",
    isSuperAdmin: role?.role === "super_admin",
    isReviewer: role?.role === "reviewer",
    isSystemMaintainer: role?.role === "system_maintainer",
  }
}
