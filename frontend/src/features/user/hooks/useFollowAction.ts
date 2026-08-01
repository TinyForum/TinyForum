import { useCallback, useState } from "react"
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
  const queryClient = useQueryClient()
  const [following, setFollowing] = useState<boolean>(false)

  // 变更：关注
  const followMutation = useMutation<boolean, Error, number>({
    mutationFn: async (userId) => {
      const res = await userApi.follow(userId)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "关注失败")
      }
      return true
    },
    onSuccess: () => {
      setFollowing(true)
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success("关注成功")
    },
    onError: (err) => {
      toast.error(err.message || "关注失败")
    },
  })

  // 变更：取消关注
  const unfollowMutation = useMutation<boolean, Error, number>({
    mutationFn: async (userId) => {
      const res = await userApi.unfollow(userId)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "取消关注失败")
      }
      return true
    },
    onSuccess: () => {
      setFollowing(false)
      queryClient.invalidateQueries({ queryKey: userKeys.all })
      toast.success("取消关注成功")
    },
    onError: (err) => {
      toast.error(err.message || "取消关注失败")
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
