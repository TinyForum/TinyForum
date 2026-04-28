// hooks/useAnswerVote.ts
import { useState, useEffect, useCallback } from "react";
import { toast } from "react-hot-toast";
import { answerApi } from "@/lib/api/modules/answer";
import { ApiResponse } from "@/lib/api/types";

type VoteType = "up" | "down" | "";

interface VoteStatusResponse {
  user_vote: number; // 1: up, -1: down, 0: no vote
  up_count: number;
  down_count: number;
}

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

export function useAnswerVote(answerId: number, currentUserId?: number) {
  const [userVote, setUserVote] = useState<VoteType>("");
  const [voteCount, setVoteCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);
  const [voting, setVoting] = useState<boolean>(false);

  const loadVoteStatus = useCallback(async (): Promise<void> => {
    if (!currentUserId) {
      setLoading(false);
      return;
    }

    setLoading(true);
    try {
      const response: { data: ApiResponse<VoteStatusResponse> } =
        await answerApi.getVoteStatus(answerId);

      if (response.data.code === 0) {
        const data = response.data.data;
        if (data) {
          const userVoteValue = data.user_vote;

          // 手动计算净得票数 = 赞同数 - 反对数
          const upCount = data.up_count || 0;
          const downCount = data.down_count || 0;
          const netVotes = upCount - downCount;

          // 转换 user_vote: 1 -> 'up', -1 -> 'down', 0 -> ''
          const newUserVote: VoteType =
            userVoteValue === 1 ? "up" : userVoteValue === -1 ? "down" : "";
          setUserVote(newUserVote);
          setVoteCount(netVotes);

          console.log(
            `Answer ${answerId} - user_vote: ${userVoteValue}, up: ${upCount}, down: ${downCount}, net: ${netVotes}`,
          );
        }
      }
    } catch (err: unknown) {
      console.error("Failed to load vote status:", err);
    } finally {
      setLoading(false);
    }
  }, [answerId, currentUserId]);

  const handleVote = useCallback(
    async (voteType: "up" | "down"): Promise<boolean> => {
      if (!currentUserId) {
        toast.error("请先登录");
        return false;
      }
      if (voting || loading) return false;

      setVoting(true);
      try {
        const isCurrentlyVoted = userVote === voteType;
        if (isCurrentlyVoted) {
          await answerApi.removeVote(answerId);
        } else {
          await answerApi.voteAnswer(answerId, voteType);
        }
        await loadVoteStatus(); // 用服务端真实数据更新 UI
        toast.success(isCurrentlyVoted ? "已取消投票" : "投票成功");
        return true;
      } catch (err: unknown) {
        const error = err as ErrorResponse;
        toast.error(error.response?.data?.message || "投票失败");
        return false;
      } finally {
        setVoting(false);
      }
    },
    [answerId, currentUserId, userVote, voting, loading, loadVoteStatus],
  );

  useEffect((): void => {
    loadVoteStatus();
  }, [loadVoteStatus]);

  return {
    userVote,
    voteCount,
    loading: loading || voting,
    handleVote,
  };
}
