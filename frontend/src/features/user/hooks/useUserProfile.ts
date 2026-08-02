// hooks/useUserProfile.ts
import { useState, useCallback } from "react"
import { toast } from "react-hot-toast"
import {
  useQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query"
import { userApi } from "@/shared/api/modules/user"
import { UpdateProfilePayload } from "@/shared/api/types/user.model"
import { ProfileResponse, UserDO } from "@/shared/api/types/user.model.do"
import { userKeys } from "./useUserKeys"

export interface ErrorResponse {
  response?: { data?: { message?: string } }
  message?: string
}

// ========== 用户资料 ==========
interface UseProfileReturn {
  user: UserDO | null
  profile: ProfileResponse | null
  loading: boolean
  error: string | null
  loadProfile: (id: number) => Promise<void>
  updateProfile: (data: UpdateProfilePayload) => Promise<boolean>
}

export function useUserProfile(): UseProfileReturn {
  const queryClient = useQueryClient()
  // 当前查询的用户 ID（null 时不发起请求）
  const [userId, setUserId] = useState<number | null>(null)
  // user 由资料更新成功后写入（用于登录态展示）
  const [user, setUser] = useState<UserDO | null>(null)

  // 查询：按 userId 拉取用户公开资料
  const profileQuery = useQuery<ProfileResponse>({
    queryKey: userKeys.detail(userId ?? 0),
    queryFn: async () => {
      if (userId === null) {
        throw new Error("缺少用户 ID")
      }
      const res = await userApi.getProfile(userId)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取用户信息失败")
      }
      if (!res.data.data) {
        throw new Error("获取用户信息失败")
      }
      return res.data.data
    },
    enabled: userId !== null,
  })

  // 命令式入口：仅更新本地 userId，由 useQuery 自动触发请求
  const loadProfile = useCallback(async (id: number): Promise<void> => {
    setUserId(id)
  }, [])

  // 变更：更新当前用户资料
  const updateProfileMutation = useMutation<UserDO, Error, UpdateProfilePayload>({
    mutationFn: async (data) => {
      const res = await userApi.updateProfile(data)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "更新失败")
      }
      if (!res.data.data) {
        throw new Error("更新失败")
      }
      return res.data.data
    },
    onSuccess: (updatedUser) => {
      setUser(updatedUser)
      queryClient.invalidateQueries({ queryKey: userKeys.details() })
      toast.success("资料更新成功")
    },
    onError: (err) => {
      toast.error(err.message || "更新失败")
    },
  })

  const updateProfile = useCallback(
    async (data: UpdateProfilePayload): Promise<boolean> => {
      try {
        await updateProfileMutation.mutateAsync(data)
        return true
      } catch {
        return false
      }
    },
    [updateProfileMutation],
  )

  return {
    user,
    profile: profileQuery.data ?? null,
    loading: profileQuery.isLoading,
    error: profileQuery.error?.message ?? null,
    loadProfile,
    updateProfile,
  }
}
