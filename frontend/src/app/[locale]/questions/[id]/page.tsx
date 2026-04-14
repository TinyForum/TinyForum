// app/questions/[id]/page.tsx
'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { useAuthStore } from '@/store/auth';
import { useQuestionDetail } from '@/hooks/useQuestionDetail';
import { AnswerCard } from '@/components/question/AnswerCard';
import { toast } from 'react-hot-toast';
import {  questionApi } from '@/lib/api';
import { AnswerForm } from '@/components/question/AnswerForm';
import { QuestionHeader } from '@/components/question/QuestionHeader';
import { answerApi } from '@/lib/api/modules/answer';

export default function QuestionDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const [answerPage, setAnswerPage] = useState(1);
  const pageSize = 20;

  const questionId = Number(params.id);
  
  const { 
    question, 
    answers, 
    answersTotal, 
    liked, 
    loading, 
    refresh,
    setLiked,
    setAnswers,
    setAnswersTotal,
  } = useQuestionDetail(questionId, { answerPage, answerPageSize: pageSize });

  const handleAcceptAnswer = async (answerId: number) => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      router.push('/login');
      return;
    }

    try {
      const response = await answerApi.acceptAnswer(answerId);
      if (response.data.code === 200) {
        toast.success('已采纳答案');
        refresh();
      } else {
        toast.error(response.data.message || '操作失败');
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleLike = async () => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      router.push('/login');
      return;
    }

    try {
      if (liked) {
        await answerApi.removeVote(questionId);
        setLiked(false);
      } else {
        await answerApi.voteAnswer(questionId,voteType);
        setLiked(true);
      }
      refresh();
    } catch (error: any) {
      toast.error(error.response?.data?.message || '操作失败');
    }
  };

  const onAnswerCreated = () => {
    setAnswerPage(1);
    refresh();
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-500">加载中...</div>
      </div>
    );
  }

  if (!question) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-500 mb-4">问题不存在</p>
          <Link href="/questions" className="text-indigo-600 hover:underline">
            返回问答列表
          </Link>
        </div>
      </div>
    );
  }

  const isAuthor = user?.id === question.author_id;
  const hasAccepted = answers.some(a => a.is_accepted);
  const totalPages = Math.ceil(answersTotal / pageSize);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 返回按钮 */}
        <div className="flex items-center justify-between mb-4">
          <Link
            href="/questions"
            className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            返回列表
          </Link>
        </div>

        {/* 问题头部 */}
        <QuestionHeader
          question={question}
          answersCount={answersTotal}
          liked={liked}
          likesCount={question.like_count || 0}
          hasAccepted={hasAccepted}
          rewardScore={0} // 需要从 API 获取悬赏积分
          onLike={handleLike}
        />

        {/* 答案列表 */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900">
              {answersTotal} 个回答
            </h2>
          </div>
          
          {answers.length === 0 ? (
            <div className="bg-white rounded-lg shadow-sm p-8 text-center text-gray-500">
              <p>暂无回答</p>
              <p className="text-sm mt-1">成为第一个回答的人吧！</p>
            </div>
          ) : (
            <div className="space-y-4">
              {answers.map((answer) => (
                <AnswerCard
                  key={answer.id}
                  answer={answer}
                  isAccepted={answer.is_accepted || false}
                  isAuthor={isAuthor}
                  canAccept={!hasAccepted && isAuthor}
                  onAccept={() => handleAcceptAnswer(answer.id)}
                  currentUserId={user?.id}
                />
              ))}
            </div>
          )}

          {/* 分页 */}
          {totalPages > 1 && (
            <div className="flex justify-center gap-2 mt-6">
              <button
                onClick={() => setAnswerPage(p => Math.max(1, p - 1))}
                disabled={answerPage === 1}
                className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
              >
                上一页
              </button>
              <span className="px-3 py-1 text-gray-600">
                第 {answerPage} / {totalPages} 页
              </span>
              <button
                onClick={() => setAnswerPage(p => Math.min(totalPages, p + 1))}
                disabled={answerPage >= totalPages}
                className="px-3 py-1 border rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
              >
                下一页
              </button>
            </div>
          )}
        </div>

        {/* 回答表单 */}
        {isAuthenticated ? (
          <AnswerForm questionId={questionId} onSuccess={onAnswerCreated} />
        ) : (
          <div className="bg-white rounded-lg shadow-sm p-6 text-center">
            <p className="text-gray-500 mb-3">登录后回答这个问题</p>
            <Link
              href={`/login?redirect=/questions/${questionId}`}
              className="inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
            >
              登录
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}