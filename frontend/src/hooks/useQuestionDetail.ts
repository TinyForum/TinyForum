// hooks/useQuestionDetail.ts
import { useState, useEffect, useCallback } from 'react';
import type { Post, Comment } from '@/lib/api/types';
import { toast } from 'react-hot-toast';
import { postApi, questionApi } from '@/lib/api';

interface UseQuestionDetailOptions {
  answerPage?: number;
  answerPageSize?: number;
}

export function useQuestionDetail(questionId: number, options: UseQuestionDetailOptions = {}) {
  const { answerPage = 1, answerPageSize = 20 } = options;
  
  const [question, setQuestion] = useState<Post | null>(null);
  const [answers, setAnswers] = useState<Comment[]>([]);
  const [answersTotal, setAnswersTotal] = useState(0);
  const [liked, setLiked] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadQuestion = useCallback(async () => {
    if (!questionId) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await questionApi.getDetail(questionId);
      
      console.log(response);
      if (response.data.code === 0 || response.data.code === 0) {
        const data = response.data.data;
        setQuestion(data.post);
        setAnswers(data.answers || []);
        setAnswersTotal(data.total || 0);
        setLiked(data.liked || false);
        setLoading(false);
      } else {
        throw new Error(response.data.message || '加载失败');
      }
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.message || '加载问题失败';
      setError(errorMsg);
      toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  }, [questionId, answerPage, answerPageSize]);

  const refresh = useCallback(() => {
    loadQuestion();
  }, [loadQuestion]);

  useEffect(() => {
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