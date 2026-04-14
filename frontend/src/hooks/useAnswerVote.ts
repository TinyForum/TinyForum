// hooks/useAnswerVote.ts
import { useState, useEffect, useCallback } from 'react';
import { commentApi } from '@/lib/api';
import { toast } from 'react-hot-toast';
import { answerApi } from '@/lib/api/modules/answer';

export function useAnswerVote(answerId: number, currentUserId?: number) {
  const [userVote, setUserVote] = useState<string>('');
  const [voteCount, setVoteCount] = useState<number>(0);
  const [loading, setLoading] = useState(true);
  const [voting, setVoting] = useState(false);

  const loadVoteStatus = useCallback(async () => {
    if (!currentUserId) {
      setLoading(false);
      return;
    }
    
    setLoading(true);
    try {
      const response = await commentApi.getVoteStatus(answerId);
      if (response.data.code === 200) {
        const { vote_type, vote_count } = response.data.data;
        setUserVote(vote_type || '');
        setVoteCount(vote_count || 0);
      }
    } catch (error) {
      console.error('Failed to load vote status:', error);
    } finally {
      setLoading(false);
    }
  }, [answerId, currentUserId]);

  const handleVote = useCallback(async (voteType: 'up' | 'down') => {
    if (!currentUserId) {
      toast.error('请先登录');
      return false;
    }

    if (voting || loading) return false;

    const wasVoted = userVote === voteType;
    const newVoteType = wasVoted ? '' : voteType;
    const delta = wasVoted
      ? (voteType === 'up' ? -1 : 1)
      : (voteType === 'up' ? 1 : -1);
    
    // 乐观更新
    setVoting(true);
    const oldVoteType = userVote;
    const oldVoteCount = voteCount;
    setUserVote(newVoteType);
    setVoteCount(prev => prev + delta);
    
    try {
      await answerApi.voteAnswer(answerId, voteType);
      await loadVoteStatus(); // 重新获取最新状态
      toast.success(wasVoted ? '已取消投票' : '投票成功');
      return true;
    } catch (error: any) {
      // 回滚
      setUserVote(oldVoteType);
      setVoteCount(oldVoteCount);
      toast.error(error.response?.data?.message || '投票失败');
      return false;
    } finally {
      setVoting(false);
    }
  }, [answerId, currentUserId, userVote, voteCount, voting, loading, loadVoteStatus]);

  useEffect(() => {
    loadVoteStatus();
  }, [loadVoteStatus]);

  return {
    userVote,
    voteCount,
    loading: loading || voting,
    handleVote,
  };
}