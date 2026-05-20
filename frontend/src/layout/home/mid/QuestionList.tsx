"use client";

import { useAuthStore } from "@/store/auth";
import { PostListSkeleton } from "@/shared/ui/common/PostListSkeleton";
import { EmptyPostList } from "@/shared/ui/common/EmptyPostList";
import { Pagination } from "@/shared/ui/common/Pagination";
import { QuestionSimple } from "@/shared/api/types/question.model";
import Link from "next/link";
import Avatar from "@/features/user/components/Avatar";
import { Tag } from "lucide-react";

// 简单的问答卡片组件
function QuestionCard({ question }: { question: QuestionSimple }) {
  const { isAuthenticated } = useAuthStore();

  return (
    <div className="bg-base-100 rounded-lg shadow-sm p-5 hover:shadow-md transition-shadow">
      <Link href={`/question/${question.id}`} className="block">
        <h2 className="text-xl font-semibold mb-2 hover:text-primary line-clamp-2">
          {question.title}
        </h2>
      </Link>
      <p className="text-base-content/70 text-sm mb-3 line-clamp-2">
        {question.summary}
      </p>
      <div className="flex flex-wrap items-center gap-3 text-sm text-base-content/60">
        {/* 作者信息 */}
        <div className="flex items-center gap-1.5">
          <Avatar avatarUrl={question.author.avatar_url} size="sm" />
          <span>{question.author.username}</span>
        </div>
        {/* 统计数据 */}
        <div className="flex items-center gap-3">
          <span>👁️ {question.view_count}</span>
          <span>💬 {question.answer_count}</span>
          <span className="text-yellow-500">⭐ {question.reward_score}</span>
        </div>
        {/* 标签 */}
        {question.tags.length > 0 && (
          <div className="flex gap-1">
            {question.tags.map((tag) => (
              <Tag key={tag.id} name={tag.name} size="sm" />
            ))}
          </div>
        )}
        {/* 已采纳标记 */}
        {question.accepted_answer_id && (
          <span className="badge badge-success badge-sm">已采纳</span>
        )}
      </div>
    </div>
  );
}

interface QuestionListProps {
  questions: QuestionSimple[];
  isLoading: boolean;
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}

export default function QuestionList({
  questions,
  isLoading,
  totalPages,
  currentPage,
  onPageChange,
}: QuestionListProps) {
  const { isAuthenticated } = useAuthStore();

  if (isLoading) {
    return <PostListSkeleton />;
  }

  if (questions.length === 0) {
    return <EmptyPostList isAuthenticated={isAuthenticated} />;
  }

  return (
    <>
      <div className="space-y-3 z-0">
        {questions.map((question) => (
          <QuestionCard key={question.id} question={question} />
        ))}
      </div>

      {totalPages > 1 && (
        <Pagination
          totalPages={totalPages}
          currentPage={currentPage}
          onPageChange={onPageChange}
        />
      )}
    </>
  );
}
