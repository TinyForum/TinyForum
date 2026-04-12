'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { questionsApi } from '@/lib/api/questions';
import { useAuthStore } from '@/store/auth';
import { toast } from 'react-hot-toast';
import {
  ArrowUpIcon,
  ArrowDownIcon,
  CheckBadgeIcon,
  UserCircleIcon,
  CalendarIcon,
} from '@heroicons/react/24/outline';
import type { Comment } from '@/types';

interface AnswerCardProps {
  answer: Comment;
  isAccepted: boolean;
  isAuthor: boolean;
  canAccept: boolean;
  onAccept: () => void;
}

export default function AnswerCard({
  answer,
  isAccepted,
  isAuthor,
  canAccept,
  onAccept,
}: AnswerCardProps) {
  const { user, isAuthenticated } = useAuthStore();
  const [userVote, setUserVote] = useState<string>('');
  const [voteCount, setVoteCount] = useState((answer as any).vote_count || 0);

  useEffect(() => {
    if (isAuthenticated) {
      loadVoteStatus();
    }
  }, [isAuthenticated]);

  const loadVoteStatus = async () => {
    try {
      const response = await questionsApi.getVoteStatus(answer.id);
      if (response.data.code === 200) {
        setUserVote(response.data.data.vote_type);
      }
    } catch (error) {
      console.error('Failed to load vote status:', error);
    }
  };

  const handleVote = async (voteType: 'up' | 'down') => {
    if (!isAuthenticated) {
      toast.error('请先登录');
      return;
    }

    // 不能给自己的答案投票
    if (answer.author_id === user?.id) {
      toast.error('不能给自己的答案投票');
      return;
    }

    // 乐观更新
    const wasVoted = userVote === voteType;
    const newVoteType = wasVoted ? '' : voteType;
    const delta = wasVoted
      ? (voteType === 'up' ? -1 : 1)
      : (voteType === 'up' ? 1 : -1);
    
    setUserVote(newVoteType);
    setVoteCount(prev => prev + delta);
    
    try {
      await questionsApi.voteAnswer(answer.id, voteType);
    } catch (error) {
      // 回滚
      setUserVote(userVote);
      setVoteCount(prev => prev - delta);
      toast.error('投票失败');
    }
  };

  return (
    <div className={`bg-white rounded-lg shadow-sm p-6 ${
      isAccepted ? 'border-2 border-green-500' : ''
    }`}>
      <div className="flex gap-4">
        {/* 投票区域 */}
        <div className="flex flex-col items-center gap-1">
          <button
            onClick={() => handleVote('up')}
            className={`p-1 rounded hover:bg-gray-100 transition-colors ${
              userVote === 'up' ? 'text-orange-500' : 'text-gray-400'
            }`}
          >
            <ArrowUpIcon className="w-6 h-6" />
          </button>
          <span className="font-medium text-gray-700">{voteCount}</span>
          <button
            onClick={() => handleVote('down')}
            className={`p-1 rounded hover:bg-gray-100 transition-colors ${
              userVote === 'down' ? 'text-blue-500' : 'text-gray-400'
            }`}
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
                  {answer.author?.username}
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
          {canAccept && (
            <div className="mt-4">
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