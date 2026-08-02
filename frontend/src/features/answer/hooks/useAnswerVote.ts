// hooks/useAnswerVote.ts
import { useCallback } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'
import { answerApi } from '@/shared/api/modules/answer'
import { VoteStatusResponse } from '@/shared/api/types/vote.model'
import { answerKeys } from './useAnswerKey'

type VoteType = 'up' | 'down' | ''

export function useAnswerVote(answerId: number, currentUserId?: number) {
  const queryClient = useQueryClient()

  // 查询：获取投票状态（未登录时不发起请求）
  const query = useQuery<VoteStatusResponse>({
    queryKey: answerKeys.voteStatus(answerId),
    queryFn: async () => {
      const res = await answerApi.getVoteStatus(answerId)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '获取投票状态失败')
      }
      if (!res.data.data) {
        throw new Error('投票状态数据为空')
      }
      return res.data.data
    },
    enabled: !!currentUserId && !!answerId,
  })

  // 转换 user_vote: 1 -> 'up', -1 -> 'down', 0 -> ''
  const userVote: VoteType =
    query.data?.user_vote === 1
      ? 'up'
      : query.data?.user_vote === -1
        ? 'down'
        : ''
  // 手动计算净得票数 = 赞同数 - 反对数
  const voteCount = (query.data?.up_count ?? 0) - (query.data?.down_count ?? 0)

  // 变更：投票/取消投票（当前已投相同票则取消，否则投票）
  const voteMutation = useMutation({
    mutationFn: async (voteType: 'up' | 'down') => {
      if (userVote === voteType) {
        await answerApi.removeVote(answerId)
      } else {
        await answerApi.voteAnswer(answerId, voteType)
      }
    },
    onSuccess: () => {
      // 投票后使该答案的投票状态失效，重新拉取服务端真实数据
      queryClient.invalidateQueries({
        queryKey: answerKeys.voteStatus(answerId),
      })
    },
  })

  const handleVote = useCallback(
    async (voteType: 'up' | 'down'): Promise<boolean> => {
      if (!currentUserId) {
        toast.error('请先登录')
        return false
      }
      try {
        await voteMutation.mutateAsync(voteType)
        toast.success(userVote === voteType ? '已取消投票' : '投票成功')
        return true
      } catch {
        toast.error('投票失败')
        return false
      }
    },
    [currentUserId, userVote, voteMutation],
  )

  return {
    userVote,
    voteCount,
    loading: query.isLoading || voteMutation.isPending,
    handleVote,
  }
}
