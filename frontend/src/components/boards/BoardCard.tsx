// components/boards/BoardCard.tsx
'use client';

import Link from 'next/link';
import { Board } from '@/lib/api';
import {
  ChatBubbleLeftRightIcon,
  FireIcon,
  RectangleGroupIcon,
  ArrowRightIcon,
  DocumentTextIcon,
} from '@heroicons/react/24/outline';

interface BoardCardProps {
  board: Board;
}

export function BoardCard({ board }: BoardCardProps) {
  return (
    <Link
      href={`/boards/${board.slug}`}
      className="group block card bg-base-100 shadow-md border border-base-200 hover:shadow-xl hover:border-primary/20 transition-all duration-300 animate-fade-in"
    >
      <div className="card-body p-5">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 min-w-0">
            {/* 板块图标和名称 */}
            <div className="flex items-center gap-2 mb-2">
              <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-red-100 to-red-200 dark:from-red-900/30 dark:to-red-800/20 flex items-center justify-center group-hover:scale-110 transition-transform">
                {board.icon ? (
                  <span className="text-lg">{board.icon}</span>
                ) : (
                  <RectangleGroupIcon className="w-4 h-4 text-red-500" />
                )}
              </div>
              <h2 className="text-base font-bold text-base-content group-hover:text-primary transition-colors line-clamp-1">
                {board.name}
              </h2>
            </div>
            
            {/* 描述 */}
            {board.description && (
              <p className="text-sm text-base-content/60 line-clamp-2 mb-3">
                {board.description}
              </p>
            )}
            
            {/* 统计信息 */}
            <div className="flex flex-wrap items-center gap-4 text-xs">
              <div className="flex items-center gap-1.5 text-base-content/50">
                <ChatBubbleLeftRightIcon className="w-3.5 h-3.5" />
                <span>{board.post_count || 0} 帖子</span>
              </div>
              <div className="flex items-center gap-1.5 text-base-content/50">
                <DocumentTextIcon className="w-3.5 h-3.5" />
                <span>{board.thread_count || 0} 主题</span>
              </div>
              {board.today_count > 0 && (
                <div className="flex items-center gap-1.5 text-orange-500">
                  <FireIcon className="w-3.5 h-3.5" />
                  <span>今日 {board.today_count}</span>
                </div>
              )}
            </div>
          </div>
          
          {/* 箭头图标 */}
          <div className="shrink-0 w-8 h-8 rounded-full bg-base-200 group-hover:bg-red-50 dark:group-hover:bg-red-900/20 flex items-center justify-center transition-all group-hover:translate-x-1">
            <ArrowRightIcon className="w-4 h-4 text-base-content/40 group-hover:text-primary transition-colors" />
          </div>
        </div>
      </div>
    </Link>
  );
}