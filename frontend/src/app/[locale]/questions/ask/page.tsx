// app/[locale]/questions/ask/page.tsx
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { 
  ArrowLeftIcon, 
  SparklesIcon,
  ExclamationTriangleIcon 
} from '@heroicons/react/24/outline';
import { useAuthStore } from '@/store/auth';
import { useBoard } from '@/hooks/useBoard';
import { useQuestionForm } from '@/hooks/useQuestionForm';
import { useTags } from '@/hooks/useTags';
import { QuestionForm } from '@/components/question/QuestionForm';
import { BoardSelector } from '@/components/question/BoardSelector';
import { TagSelector } from '@/components/question/TagSelector';
import { RewardScoreInput } from '@/components/question/RewardScoreInput';
import { FormActions } from '@/components/question/FormActions';
import { useTranslations } from 'next-intl';

export default function AskQuestionPage() {
  const t = useTranslations("Questions");
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
      }
    }
  }, [boards, boardsLoading, getDefaultBoard, selectedBoardId]);

  // 检查登录状态
  useEffect(() => {
    if (typeof window !== 'undefined' && !isAuthenticated) {
      router.push('/auth/login?redirect=/questions/ask');
    }
  }, [isAuthenticated, router]);

  // 显示加载状态
  if (boardsLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="card bg-base-100 shadow-md border border-base-200">
            <div className="card-body p-8">
              <div className="flex flex-col items-center justify-center py-12">
                <span className="loading loading-spinner loading-lg text-primary mb-4"></span>
                <p className="text-base-content/60">{t('loading')}</p>
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
      <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
        <div className="max-w-4xl mx-auto px-4 py-8">
          <div className="card bg-base-100 shadow-md border border-base-200">
            <div className="card-body p-8">
              <div className="text-center">
                <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-100 dark:bg-red-900/20 flex items-center justify-center">
                  <ExclamationTriangleIcon className="w-10 h-10 text-red-500" />
                </div>
                <h3 className="text-lg font-semibold text-base-content mb-2">{t('load_boards_failed')}</h3>
                <p className="text-base-content/60 mb-4">{boardsError}</p>
                <button
                  onClick={() => window.location.reload()}
                  className="btn btn-primary"
                >
                  {t('reload')}
                </button>
              </div>
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
    <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Hero Section */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-red-500 to-red-600 shadow-lg mb-4">
            <SparklesIcon className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-red-600 to-red-500 bg-clip-text text-transparent">
            {t('ask_question_title')}
          </h1>
          <p className="text-base-content/60 mt-2">
            {t('ask_question_description')}
          </p>
        </div>

        {/* 返回按钮 */}
        <div className="mb-4">
          <Link
            href="/questions"
            className="inline-flex items-center gap-2 text-base-content/60 hover:text-primary transition-colors"
          >
            <ArrowLeftIcon className="w-4 h-4" />
            {t('back_to_questions')}
          </Link>
        </div>

        {/* 表单卡片 */}
        <div className="card bg-base-100 shadow-md border border-base-200">
          <div className="card-body p-0">
            <div className="p-6 border-b border-base-200">
              <h2 className="text-xl font-bold text-base-content">{t('publish_question')}</h2>
              <p className="text-sm text-base-content/60 mt-1">
                {t('publish_question_description')}
              </p>
            </div>

            <form onSubmit={form.handleSubmit(onSubmit)} className="p-6 space-y-6">
              {/* 板块选择 */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-base-content">
                  {t('select_board')} <span className="text-red-500">*</span>
                </label>
                <BoardSelector
                  boards={boards}
                  selectedBoardId={selectedBoardId}
                  onBoardChange={setSelectedBoardId}
                  loading={boardsLoading}
                />
              </div>

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

        {/* 提示信息 */}
        <div className="mt-6 text-center text-sm text-base-content/40">
          <p>{t('agree_to_rules')}</p>
        </div>
      </div>
    </div>
  );
}