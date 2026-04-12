// components/question/QuestionHeader.tsx
'use client';

import Link from 'next/link';
import { 
  EyeIcon, 
  ChatBubbleLeftRightIcon, 
  CheckBadgeIcon,
  UserCircleIcon,
  CalendarIcon,
  HeartIcon,
  ShareIcon,
  FlagIcon,
  TagIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import type { Post, Tag } from '@/lib/api/types';

interface QuestionHeaderProps {
  question: Post;
  answersCount: number;
  liked: boolean;
  likesCount: number;
  hasAccepted: boolean;
  rewardScore?: number;
  onLike: () => void;
  onShare?: () => void;
  onReport?: () => void;
}

export function QuestionHeader({ 
  question, 
  answersCount, 
  liked, 
  likesCount, 
  hasAccepted,
  rewardScore,
  onLike,
  onShare,
  onReport,
}: QuestionHeaderProps) {
  return (
    <div className="bg-white rounded-lg shadow-sm mb-6">
      <div className="p-6">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">{question.title}</h1>
        
        {/* 元信息 */}
        <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500 mb-4 pb-4 border-b">
          <div className="flex items-center gap-1">
            <UserCircleIcon className="w-4 h-4" />
            <Link href={`/users/${question.author_id}`} className="hover:text-indigo-600">
              {question.author?.username || `用户${question.author_id}`}
            </Link>
          </div>
          <div className="flex items-center gap-1">
            <CalendarIcon className="w-4 h-4" />
            {new Date(question.created_at).toLocaleDateString()}
          </div>
          <div className="flex items-center gap-1">
            <EyeIcon className="w-4 h-4" />
            {question.view_count || 0} 浏览
          </div>
          <div className="flex items-center gap-1">
            <ChatBubbleLeftRightIcon className="w-4 h-4" />
            {answersCount} 回答
          </div>
          <button
            onClick={onLike}
            className="flex items-center gap-1 hover:text-red-500 transition-colors"
          >
            {liked ? (
              <HeartSolidIcon className="w-4 h-4 text-red-500" />
            ) : (
              <HeartIcon className="w-4 h-4" />
            )}
            {likesCount} 点赞
          </button>
          {rewardScore && rewardScore > 0 && (
            <div className="flex items-center gap-1 text-orange-500 bg-orange-50 px-2 py-1 rounded">
              💰 {rewardScore} 积分悬赏
            </div>
          )}
          {hasAccepted && (
            <div className="flex items-center gap-1 text-green-500 bg-green-50 px-2 py-1 rounded">
              <CheckBadgeIcon className="w-4 h-4" />
              已解决
            </div>
          )}
        </div>

        {/* 标签 */}
        {question.tags && question.tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-4">
            {question.tags.map((tag: Tag) => (
              <Link
                key={tag.id}
                href={`/questions?tag_id=${tag.id}`}
                className="px-2 py-1 bg-gray-100 text-gray-600 text-sm rounded-md hover:bg-gray-200 transition-colors"
                style={{ borderLeft: `3px solid ${tag.color || '#6366f1'}` }}
              >
                <TagIcon className="w-3 h-3 inline mr-1" />
                {tag.name}
              </Link>
            ))}
          </div>
        )}

        {/* 内容 */}
        <div
          className="prose max-w-none"
          dangerouslySetInnerHTML={{ __html: question.content }}
        />
      </div>

      {/* 操作按钮 */}
      <div className="px-6 pb-4 flex gap-2 border-t pt-4">
        <button 
          onClick={onShare}
          className="flex items-center gap-1 px-3 py-1 text-gray-500 hover:text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <ShareIcon className="w-4 h-4" />
          分享
        </button>
        <button 
          onClick={onReport}
          className="flex items-center gap-1 px-3 py-1 text-gray-500 hover:text-red-500 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <FlagIcon className="w-4 h-4" />
          举报
        </button>
      </div>
    </div>
  );
}