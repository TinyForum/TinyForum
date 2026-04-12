// app/questions/ask/page.tsx
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { useAuthStore } from '@/store/auth';
import { useBoard } from '@/hooks/useBoard';
import { useQuestionForm } from '@/hooks/useQuestionForm';
import { useTags } from '@/hooks/useTags';
import { QuestionForm } from '@/components/question/QuestionForm';
import { BoardSelector } from '@/components/question/BoardSelector';
import { TagSelector } from '@/components/question/TagSelector';
import { RewardScoreInput } from '@/components/question/RewardScoreInput';
import { FormActions } from '@/components/question/FormActions';

export default function AskQuestionPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();
  const { boards, loading: boardsLoading, getDefaultBoard, error: boardsError } = useBoard({ 
    autoLoad: true,
    pageSize: 100 
  });
  const { tags, selectedTags, loading: tagsLoading, toggleTag } = useTags();
  const { form, content, setContent, loading, handleSubmit } = useQuestionForm({
    onSuccess: (questionId) => router.push(`/questions/${questionId}`),
  });
  
  const [selectedBoardId, setSelectedBoardId] = useState<number>(0);

  const { register, watch, formState: { errors } } = form;
  const rewardScore = watch('reward_score');

  // 设置默认板块（只在加载完成后设置一次）
  useEffect(() => {
    if (!boardsLoading && boards.length > 0 && selectedBoardId === 0) {
      const defaultBoard = getDefaultBoard();
      if (defaultBoard) {
        setSelectedBoardId(defaultBoard.id);
        console.log('设置默认板块:', defaultBoard.id, defaultBoard.name);
      }
    }
  }, [boards, boardsLoading, getDefaultBoard, selectedBoardId]);

  // 检查登录状态
  useEffect(() => {
    if (typeof window !== 'undefined' && !isAuthenticated) {
      router.push('/login?redirect=/questions/ask');
    }
  }, [isAuthenticated, router]);

  // 显示加载状态
  if (boardsLoading) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-4xl mx-auto px-4">
          <div className="bg-white rounded-lg shadow-sm p-8">
            <div className="flex items-center justify-center">
              <div className="text-center">
                <div className="w-12 h-12 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
                <p className="text-gray-600">加载中...</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // 显示错误状态
  if (boardsError) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-4xl mx-auto px-4">
          <div className="bg-white rounded-lg shadow-sm p-8">
            <div className="text-center">
              <div className="text-red-500 text-lg mb-4">⚠️ 加载板块失败</div>
              <p className="text-gray-600 mb-4">{boardsError}</p>
              <button
                onClick={() => window.location.reload()}
                className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
              >
                重新加载
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  const onSubmit = async () => {
    if (selectedBoardId === 0) {
      return;
    }
    await handleSubmit(selectedTags, selectedBoardId);
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* 返回按钮 */}
        <Link
          href="/questions"
          className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeftIcon className="w-4 h-4" />
          返回问答列表
        </Link>

        <div className="bg-white rounded-lg shadow-sm">
          <div className="p-6 border-b">
            <h1 className="text-2xl font-bold text-gray-900">提问</h1>
            <p className="text-gray-500 mt-1">详细描述你的问题，获得更精准的回答</p>
          </div>

          <form onSubmit={form.handleSubmit(onSubmit)} className="p-6 space-y-6">
            {/* 板块选择 */}
            <BoardSelector
              boards={boards}
              selectedBoardId={selectedBoardId}
              onBoardChange={setSelectedBoardId}
              loading={boardsLoading}
            />

            {/* 问题表单 */}
            <QuestionForm
              register={register}
              errors={errors}
              content={content}
              setContent={setContent}
            />

            {/* 标签选择 */}
            <TagSelector
              tags={tags}
              selectedTags={selectedTags}
              onToggleTag={toggleTag}
              loading={tagsLoading}
            />

            {/* 悬赏积分 */}
            <RewardScoreInput
              register={register}
              rewardScore={rewardScore}
              userScore={user?.score}
            />

            {/* 表单按钮 */}
            <FormActions loading={loading} />
          </form>
        </div>
      </div>
    </div>
  );
}