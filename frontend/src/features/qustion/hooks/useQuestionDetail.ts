// hooks/useQuestionDetail.ts
import { useState, useEffect, useCallback } from "react";
import { toast } from "react-hot-toast";
import { Comment } from "@/shared/api/types/comment.model";
// import { Post } from "@/shared/api/types/post.model";
import { questionApi } from "@/shared/api/modules/questions";
import { QuestionDO } from "@/shared/api/types/question.model";
// import { Creation } from "@/shared/api/types/creation.model";

export function useQuestionDetail(questionId: number) {
  const [question, setQuestion] = useState<QuestionDO | null>(null);
  const [answers, setAnswers] = useState<Comment[]>([]);
  const [answersTotal, setAnswersTotal] = useState<number>(0);
  const [liked, setLiked] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * 加载回答详情
   */
  const loadAnswers = useCallback(async (): Promise<void> => {
    if (!questionId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await questionApi.getAnswers(questionId);
      console.log(response);
      if (response.data.code === 0) {
        const data = response.data.data;
        if (data) {
          // setQuestion(data);
          setAnswers(data.list || []);
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
  }, [questionId]);

  /**
   * 加载问题详情
   */
  const loadQuestion = useCallback(async (): Promise<void> => {
    if (!questionId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await questionApi.getDetail(questionId);

      console.log(response);

      if (response.data.code === 0) {
        const data = response.data.data;
        if (data) {
          setQuestion(data);
          // setAnswers(data.answers || []);
          // setAnswersTotal(data.total || 0);
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

  /**
   * 刷新
   */
  const refresh = useCallback((): void => {
    loadQuestion();
  }, [loadQuestion]);

  useEffect((): void => {
    loadQuestion();
    loadAnswers();
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
