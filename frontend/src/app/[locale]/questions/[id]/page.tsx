// app/[locale]/questions/[id]/page.tsx
"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import {
  ArrowLeftIcon,
  ChatBubbleLeftRightIcon,
  UserGroupIcon,
} from "@heroicons/react/24/outline";
import { useAuthStore } from "@/store/auth";
import { useQuestionDetail } from "@/hooks/useQuestionDetail";
import { AnswerCard } from "@/components/question/AnswerCard";
import { toast } from "react-hot-toast";
import { postApi } from "@/shared/api";
import { AnswerForm } from "@/components/question/AnswerForm";
import { QuestionHeader } from "@/components/question/QuestionHeader";
import { answerApi } from "@/shared/api/modules/answer";
import { useTranslations } from "next-intl";

// 错误响应类型
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

// 加载骨架屏组件
function LoadingSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="animate-pulse">
          {/* 返回按钮骨架 */}
          <div className="h-5 w-24 bg-base-200 rounded mb-6" />

          {/* 问题卡片骨架 */}
          <div className="bg-base-100 rounded-2xl shadow-sm p-6 mb-6">
            <div className="h-7 bg-base-200 rounded-lg w-3/4 mb-4" />
            <div className="flex items-center gap-4 mb-4">
              <div className="h-4 w-24 bg-base-200 rounded" />
              <div className="h-4 w-32 bg-base-200 rounded" />
            </div>
            <div className="space-y-2">
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-full" />
              <div className="h-4 bg-base-200 rounded w-2/3" />
            </div>
          </div>

          {/* 回答区域骨架 */}
          <div className="space-y-3">
            <div className="h-6 w-32 bg-base-200 rounded" />
            {[1, 2].map((i) => (
              <div key={i} className="bg-base-100 rounded-xl p-5">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-10 h-10 bg-base-200 rounded-full" />
                  <div className="flex-1">
                    <div className="h-4 bg-base-200 rounded w-32 mb-2" />
                    <div className="h-3 bg-base-200 rounded w-24" />
                  </div>
                </div>
                <div className="h-4 bg-base-200 rounded w-full mb-2" />
                <div className="h-4 bg-base-200 rounded w-5/6" />
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

// 错误状态组件
function ErrorState({ message }: { message: string }) {
  const t = useTranslations("Questions");
  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100 flex items-center justify-center">
      <div className="text-center">
        <div className="text-6xl mb-5">😕</div>
        <h3 className="text-xl font-semibold text-base-content mb-2">
          {t("load_failed")}
        </h3>
        <p className="text-base-content/60 mb-6">{message}</p>
        <Link href="/questions" className="btn btn-primary btn-sm gap-2">
          <ArrowLeftIcon className="w-4 h-4" />
          {t("back_to_questions")}
        </Link>
      </div>
    </div>
  );
}

// 分页组件
function Pagination({
  currentPage,
  totalPages,
  onPageChange,
}: {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  const t = useTranslations("Questions");

  const getPageNumbers = () => {
    const pages: number[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
    } else {
      if (currentPage <= 3) {
        for (let i = 1; i <= maxVisible; i++) pages.push(i);
      } else if (currentPage >= totalPages - 2) {
        for (let i = totalPages - maxVisible + 1; i <= totalPages; i++)
          pages.push(i);
      } else {
        for (let i = currentPage - 2; i <= currentPage + 2; i++) pages.push(i);
      }
    }
    return pages;
  };

  return (
    <div className="flex justify-center items-center gap-2 mt-8">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="btn btn-ghost btn-sm gap-1"
      >
        <ArrowLeftIcon className="w-4 h-4" />
        {t("prev_page")}
      </button>

      <div className="flex gap-1.5 mx-2">
        {getPageNumbers().map((pageNum) => (
          <button
            key={pageNum}
            onClick={() => onPageChange(pageNum)}
            className={`btn btn-sm min-w-[2.5rem] ${
              currentPage === pageNum ? "btn-primary" : "btn-ghost"
            }`}
          >
            {pageNum}
          </button>
        ))}
      </div>

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages}
        className="btn btn-ghost btn-sm gap-1"
      >
        {t("next_page")}
        <ArrowLeftIcon className="w-4 h-4 rotate-180" />
      </button>
    </div>
  );
}

export default function QuestionDetailPage() {
  const t = useTranslations("Questions");
  const params = useParams();
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const [answerPage, setAnswerPage] = useState(1);
  const pageSize = 20;

  const questionId = Number(params.id);

  const { question, answers, answersTotal, liked, loading, refresh, setLiked } =
    useQuestionDetail(questionId);

  const handleAcceptAnswer = async (answerId: number) => {
    if (!isAuthenticated) {
      toast.error(t("please_login"));
      router.push("/login");
      return;
    }

    try {
      const response = await answerApi.acceptAnswer(answerId);
      if (response.data.code === 0) {
        toast.success(t("answer_accepted"));
        refresh();
      } else {
        toast.error(response.data.message || t("operation_failed"));
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("operation_failed"));
    }
  };

  const handleLike = async () => {
    if (!isAuthenticated) {
      toast.error(t("please_login"));
      router.push("/login");
      return;
    }

    try {
      if (liked) {
        await postApi.unlike(questionId);
        setLiked(false);
      } else {
        await postApi.like(questionId);
        setLiked(true);
      }
      refresh();
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(error.response?.data?.message || t("operation_failed"));
    }
  };

  const onAnswerCreated = () => {
    setAnswerPage(1);
    refresh();
  };

  if (loading) {
    return <LoadingSkeleton />;
  }

  if (!question) {
    return <ErrorState message={t("question_not_found")} />;
  }

  const isAuthor = user?.id === question.author_id;
  const hasAccepted = answers.some((a) => a.is_accepted);
  const totalPages = Math.ceil(answersTotal / pageSize);

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-4xl mx-auto px-4 py-8 md:py-12">
        {/* 返回按钮 - 优化样式 */}
        <div className="flex items-center justify-between mb-6">
          <Link
            href="/questions"
            className="inline-flex items-center gap-2 text-base-content/60 hover:text-primary transition-all duration-200 group"
          >
            <ArrowLeftIcon className="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
            <span className="text-sm font-medium">{t("back_to_list")}</span>
          </Link>

          {/* 统计信息 */}
          <div className="flex items-center gap-4 text-sm text-base-content/50">
            <div className="flex items-center gap-1.5">
              <ChatBubbleLeftRightIcon className="w-4 h-4" />
              <span>{t("answers_count", { count: answersTotal })}</span>
            </div>
          </div>
        </div>

        {/* 问题头部 */}
        <div className="mb-8">
          <QuestionHeader
            question={question}
            answersCount={answersTotal}
            liked={liked}
            likesCount={question.like_count || 0}
            hasAccepted={hasAccepted}
            rewardScore={0}
            onLike={handleLike}
          />
        </div>

        {/* 答案列表区域 */}
        <div className="mb-8">
          <div className="flex items-center gap-2.5 mb-5">
            <div className="w-1 h-6 bg-primary rounded-full" />
            <h2 className="text-lg font-semibold text-base-content">
              {t("all_answers")}
            </h2>
            <span className="text-sm text-base-content/40">
              ({answersTotal})
            </span>
            <div className="flex-1 h-px bg-gradient-to-r from-base-200 to-transparent" />
          </div>

          {answers.length === 0 ? (
            <div className="bg-base-100 rounded-2xl shadow-sm p-12 text-center border border-base-200">
              <div className="text-5xl mb-4 opacity-50">💬</div>
              <p className="text-base-content/60">{t("no_answers")}</p>
              <p className="text-sm text-base-content/40 mt-1">
                {t("be_first_to_answer")}
              </p>
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
            <Pagination
              currentPage={answerPage}
              totalPages={totalPages}
              onPageChange={setAnswerPage}
            />
          )}
        </div>

        {/* 回答表单区域 */}
        <div className="mt-8">
          <div className="flex items-center gap-2.5 mb-5">
            <div className="w-1 h-6 bg-secondary rounded-full" />
            <h2 className="text-lg font-semibold text-base-content">
              {t("post_answer")}
            </h2>
            <div className="flex-1 h-px bg-gradient-to-r from-base-200 to-transparent" />
          </div>

          {isAuthenticated ? (
            <div className="bg-base-100 rounded-2xl shadow-sm border border-base-200 overflow-hidden">
              <AnswerForm
                questionId={questionId}
                onSuccess={onAnswerCreated}
                hasAccepted={hasAccepted}
                rewardScore={0}
              />
            </div>
          ) : (
            <div className="bg-base-100 rounded-2xl shadow-sm p-8 text-center border border-base-200">
              <div className="text-5xl mb-4 opacity-50">🔒</div>
              <p className="text-base-content/60 mb-4">
                {t("login_to_answer")}
              </p>
              <Link
                href={`/login?redirect=/questions/${questionId}`}
                className="btn btn-primary btn-md gap-2"
              >
                <UserGroupIcon className="w-4 h-4" />
                {t("login_now")}
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
