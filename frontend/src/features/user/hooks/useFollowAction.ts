import { useCallback, useState } from "react"
import { useTranslations } from "next-intl"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import toast from "react-hot-toast"
import { userApi } from "@/shared/api/modules/user"
import { userKeys } from "./useUserKeys"

// ========== 关注/取消关注 ==========
interface UseFollowReturn {
  following: boolean
  loading: boolean
  follow: (userId: number) => Promise<boolean>
  unfollow: (userId: number) => Promise<boolean>
  checkFollowStatus: (userId: number) => Promise<boolean>
}

export function useFollowAction(): UseFollowReturn {
  const t = useTranslations("Profile")
  const queryClient = useQueryClient()
  const [following, setFollowing] = useState<boolean>(false)

  // 变更：关注
  const followMutation = useMutation<boolean, Error, number>({
    mutationFn: async (userId) => {
      const res = await userApi.follow(userId)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("follow_failed"))
      }
      return true
    },
    onSuccess: () => {
      setFollowing(true)
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success(t("follow_success"))
    },
    onError: (err) => {
      toast.error(err.message || t("follow_failed"))
    },
  })

  // 变更：取消关注
  const unfollowMutation = useMutation<boolean, Error, number>({
    mutationFn: async (userId) => {
      const res = await userApi.unfollow(userId)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || t("unfollow_failed"))
      }
      return true
    },
    onSuccess: () => {
      setFollowing(false)
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success(t("unfollow_success"))
    },
    onError: (err) => {
      toast.error(err.message || t("unfollow_failed"))
    },
  })

  const follow = useCallback(
    async (userId: number): Promise<boolean> => {
      try {
        return await followMutation.mutateAsync(userId)
      } catch {
        return false
      }
    },
    [followMutation],
  )

  const unfollow = useCallback(
    async (userId: number): Promise<boolean> => {
      try {
        return await unfollowMutation.mutateAsync(userId)
      } catch {
        return false
      }
    },
    [unfollowMutation],
  )

  const checkFollowStatus = useCallback(
    async (userId: number): Promise<boolean> => {
      // 需要真实接口时可调用 userApi.getFollowing 检查
      console.log("check follow status", userId, following)
      return following
    },
    [following],
  )

  return {
    following,
    loading: followMutation.isPending || unfollowMutation.isPending,
    follow,
    unfollow,
    checkFollowStatus,
  }
}
