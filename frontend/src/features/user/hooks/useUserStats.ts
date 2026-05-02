import { ErrorResponse } from "./useUserProfile";
import { useState, useCallback } from "react";
import { userStatsApi } from "@/shared/api/modules/user/stats";
import { UserStatsVO } from "@/shared/api/types/user.model";

type UseUserStatsReturn = UserStatsVO & {
  loadStats: () => Promise<void>;
  isLoading: boolean;
  error: string | null;
};

export function useUserStats(): UseUserStatsReturn {
  // ----- 基础统计字段 -----
  const [totalPost, setTotalPost] = useState<number>(0);
  const [totalComment, setTotalComment] = useState<number>(0);
  const [totalFavorite, setTotalFavorites] = useState<number>(0);
  const [totalLike, setTotalLike] = useState<number>(0);
  const [totalFollower, setTotalFollower] = useState<number>(0);
  const [totalFollowing, setTotalFollowing] = useState<number>(0);
  const [totalReport, setTotalReport] = useState<number>(0);
  const [totalViolation, setTotalViolation] = useState<number>(0);
  const [totalQuestion, setTotalQuestion] = useState<number>(0);
  const [totalAnswer, setTotalAnswer] = useState<number>(0);
  const [totalUpload, setTotalUpload] = useState<number>(0);
  const [totalScore, setTotalScore] = useState<number>(0);

  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadStats = useCallback(async (): Promise<void> => {
    console.log("loadStats");
    setIsLoading(true);
    setError(null);
    try {
      const response = await userStatsApi.getUserStats();
      console.log("response: ", response);
      if (response.status === 200 && response.data.code === 0) {
        const data = response.data.data;
        if (data) {
          setTotalPost(data.total_post ?? 0);
          setTotalComment(data.total_comment ?? 0);
          setTotalFavorites(data.total_favorite ?? 0);
          setTotalLike(data.total_like ?? 0);
          setTotalFollower(data.total_follower ?? 0);
          setTotalFollowing(data.total_following ?? 0);
          setTotalReport(data.total_report ?? 0);
          setTotalViolation(data.total_violation ?? 0);
          setTotalQuestion(data.total_question ?? 0);
          setTotalAnswer(data.total_answer ?? 0);
          setTotalUpload(data.total_upload ?? 0);
          setTotalScore(data.total_score ?? 0);
        }
      } else {
        throw new Error(response.data.message || "获取统计数量失败");
      }
    } catch (err: unknown) {
      const errorObj = err as ErrorResponse;
      const errorMsg =
        errorObj.response?.data?.message ||
        errorObj.message ||
        "获取统计数量失败";
      setError(errorMsg);
    } finally {
      setIsLoading(false);
    }
    console.log("loadStats end");
  }, []);

  return {
    // 返回值与 UserStatsVO 完全对齐（下划线命名）
    total_post: totalPost,
    total_comment: totalComment,
    total_favorite: totalFavorite,
    total_like: totalLike,
    total_follower: totalFollower,
    total_following: totalFollowing,
    total_report: totalReport,
    total_violation: totalViolation,
    total_question: totalQuestion,
    total_answer: totalAnswer,
    total_upload: totalUpload,
    total_score: totalScore,
    isLoading,
    error,
    loadStats,
  };
}
