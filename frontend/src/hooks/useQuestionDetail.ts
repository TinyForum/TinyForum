// hooks/useQuestionDetail.ts
import { useState, useEffect, useCallback } from "react";
import type { Post, Comment, ApiResponse } from "@/lib/api/types";
import { toast } from "react-hot-toast";
import { questionApi } from "@/lib/api";

// interface UseQuestionDetailOptions {
//   answerPage?: number;
//   answerPageSize?: number;
// }

interface QuestionDetailResponse {
  post: Post;
  answers: Comment[];
  total: number;
}

export function useQuestionDetail(
  questionId: number,
  // options: UseQuestionDetailOptions = {},
) {
  // const { answerPage = 1, answerPageSize = 20 } = options;

  const [question, setQuestion] = useState<Post | null>(null);
  const [answers, setAnswers] = useState<Comment[]>([]);
  const [answersTotal, setAnswersTotal] = useState<number>(0);
  const [liked, setLiked] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const loadQuestion = useCallback(async (): Promise<void> => {
    if (!questionId) return;

    setLoading(true);
    setError(null);

    try {
      const response: { data: ApiResponse<QuestionDetailResponse> } =
        await questionApi.getDetail(questionId);

      console.log(response);

      if (response.data.code === 0) {
        const data = response.data.data;
        if (data) {
          setQuestion(data.post);
          setAnswers(data.answers || []);
          setAnswersTotal(data.total || 0);
        }
      } else {
        throw new Error(response.data.message || "加载失败");
      }
    } catch (err: unknown) {
      const errorMsg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message ||
        (err as Error)?.message ||
        "加载问题失败";
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, [questionId]); // 移除 answerPage 和 answerPageSize

  const refresh = useCallback((): void => {
    loadQuestion();
  }, [loadQuestion]);

  useEffect((): void => {
    loadQuestion();
  }, [loadQuestion]);

  return {
    question,
    answers,
    answersTotal,
    liked,
    loading,
    error,
    refresh,
    setLiked,
    setAnswers,
    setAnswersTotal,
  };
}
