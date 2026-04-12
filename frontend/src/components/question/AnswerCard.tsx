// components/question/AnswerCard.tsx
'use client';

import Link from 'next/link';
import { ArrowUpIcon, ArrowDownIcon } from 'lucide-react';
import { UserCircleIcon, CalendarIcon, CheckBadgeIcon } from '@heroicons/react/24/outline';
import { useAnswerVote } from '@/hooks/useAnswerVote';
import type { Comment } from '@/lib/api/types';

interface AnswerCardProps {
  answer: Comment;
  isAccepted: boolean;
  isAuthor: boolean;
  canAccept: boolean;
  onAccept: () => void;
  currentUserId?: number;
}

export function AnswerCard({ 
  answer, 
  isAccepted, 
  canAccept, 
  onAccept,
  currentUserId,
}: AnswerCardProps) {
  const { userVote, voteCount, loading: voteLoading, handleVote } = useAnswerVote(
    answer.id,
    currentUserId
  );

  const handleUpVote = () => handleVote('up');
  const handleDownVote = () => handleVote('down');

  return (
    <div className={`bg-white rounded-lg shadow-sm p-6 transition-all ${
      isAccepted ? 'border-2 border-green-500 shadow-md' : ''
    }`}>
      <div className="flex gap-4">
        {/* 投票区域 */}
        <div className="flex flex-col items-center gap-1">
          <button
            onClick={handleUpVote}
            disabled={voteLoading}
            className={`p-1 rounded hover:bg-gray-100 transition-colors ${
              userVote === 'up' ? 'text-orange-500' : 'text-gray-400'
            } disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            <ArrowUpIcon className="w-6 h-6" />
          </button>
          <span className="font-medium text-gray-700">{voteCount}</span>
          <button
            onClick={handleDownVote}
            disabled={voteLoading}
            className={`p-1 rounded hover:bg-gray-100 transition-colors ${
              userVote === 'down' ? 'text-blue-500' : 'text-gray-400'
            } disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            <ArrowDownIcon className="w-6 h-6" />
          </button>
        </div>

        {/* 内容区域 */}
        <div className="flex-1">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-3 text-sm text-gray-500">
              <div className="flex items-center gap-1">
                <UserCircleIcon className="w-4 h-4" />
                <Link href={`/users/${answer.author_id}`} className="hover:text-indigo-600">
                  {answer.author?.username || `用户${answer.author_id}`}
                </Link>
              </div>
              <div className="flex items-center gap-1">
                <CalendarIcon className="w-4 h-4" />
                {new Date(answer.created_at).toLocaleDateString()}
              </div>
            </div>
            {isAccepted && (
              <div className="flex items-center gap-1 text-green-500 bg-green-50 px-2 py-1 rounded">
                <CheckBadgeIcon className="w-4 h-4" />
                <span className="text-sm font-medium">已采纳</span>
              </div>
            )}
          </div>

          <div
            className="prose max-w-none text-gray-700"
            dangerouslySetInnerHTML={{ __html: answer.content }}
          />

          {/* 采纳按钮 */}
          {canAccept && !isAccepted && (
            <div className="mt-4 pt-3 border-t">
              <button
                onClick={onAccept}
                className="flex items-center gap-1 text-green-600 hover:text-green-700 text-sm font-medium transition-colors"
              >
                <CheckBadgeIcon className="w-4 h-4" />
                采纳为答案
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}