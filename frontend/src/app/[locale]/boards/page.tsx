'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { boardApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import type { Board } from '@/types';
import {
  MagnifyingGlassIcon,
  FolderPlusIcon,
  ChatBubbleLeftRightIcon,
  FireIcon,
  ChevronRightIcon,
  RectangleGroupIcon,
  PlusIcon,
} from '@heroicons/react/24/outline';

// ─── BoardCard ─────────────────────────────────────────────────────────────────

function BoardCard({ board }: { board: Board }) {
  return (
    <Link
      href={`/boards/${board.slug}`}
      className="group flex items-start gap-4 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 hover:border-blue-200 dark:hover:border-blue-700 hover:shadow-md transition-all duration-200 p-5"
    >
      {/* Icon */}
      <div className="shrink-0 w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center text-lg">
        {board.icon || <RectangleGroupIcon className="w-5 h-5 text-blue-400" />}
      </div>

      {/* Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors truncate">
            {board.name}
          </h3>
          {board.today_count > 0 && (
            <span className="shrink-0 inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-orange-50 dark:bg-orange-900/30 text-orange-500">
              <FireIcon className="w-3 h-3" />
              今日 +{board.today_count}
            </span>
          )}
        </div>
        {board.description && (
          <p className="mt-1 text-sm text-gray-400 dark:text-gray-500 line-clamp-1">
            {board.description}
          </p>
        )}
        <div className="flex items-center gap-3 mt-2 text-xs text-gray-400">
          <span className="flex items-center gap-1">
            <ChatBubbleLeftRightIcon className="w-3.5 h-3.5" />
            {board.post_count} 帖子
          </span>
          <span>{board.thread_count} 主题</span>
        </div>

        {/* Sub-boards */}
        {board.children && board.children.length > 0 && (
          <div className="flex flex-wrap gap-1.5 mt-3">
            {board.children.map(child => (
              <Link
                key={child.id}
                href={`/boards/${child.slug}`}
                onClick={e => e.stopPropagation()}
                className="text-xs px-2.5 py-1 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
              >
                {child.name}
              </Link>
            ))}
          </div>
        )}
      </div>

      <ChevronRightIcon className="shrink-0 w-4 h-4 text-gray-300 group-hover:text-blue-400 transition-colors mt-1" />
    </Link>
  );
}

// ─── BoardsPage ────────────────────────────────────────────────────────────────

export default function BoardsPage() {
  const { user } = useAuthStore();
  const [boards, setBoards] = useState<Board[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await boardApi.getTree();
      setBoards(res.data.data);
    } catch (err) {
      console.error('Failed to load boards:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  // Client-side search filter (name + description)
  const filtered = search.trim()
    ? boards.filter(b => {
        const q = search.toLowerCase();
        return (
          b.name.toLowerCase().includes(q) ||
          b.description?.toLowerCase().includes(q) ||
          b.children?.some(c => c.name.toLowerCase().includes(q))
        );
      })
    : boards;

  // Stats totals
  const totalPosts   = boards.reduce((s, b) => s + b.post_count, 0);
  const totalToday   = boards.reduce((s, b) => s + b.today_count, 0);
  const totalBoards  = boards.reduce((s, b) => s + 1 + (b.children?.length ?? 0), 0);

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Header */}
      <div className="flex items-start justify-between gap-4 mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">板块</h1>
          {!loading && (
            <p className="text-sm text-gray-400 mt-1">
              共 {totalBoards} 个板块 · {totalPosts} 篇帖子
              {totalToday > 0 && (
                <span className="ml-2 text-orange-400">今日新增 {totalToday} 篇</span>
              )}
            </p>
          )}
        </div>

        {user && (
          <Link
            href="/boards/new"
            className="shrink-0 flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl font-medium transition-colors"
          >
            <PlusIcon className="w-4 h-4" />
            创建板块
          </Link>
        )}
      </div>

      {/* Search */}
      <div className="relative mb-6">
        <MagnifyingGlassIcon className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          type="text"
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="搜索板块..."
          className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition text-sm"
        />
        {search && (
          <button
            onClick={() => setSearch('')}
            className="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-300 hover:text-gray-500 transition-colors text-lg leading-none"
          >
            ×
          </button>
        )}
      </div>

      {/* Board list */}
      {loading ? (
        <div className="space-y-3">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="h-20 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 animate-pulse" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-20 bg-white dark:bg-gray-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
          {search ? (
            <>
              <MagnifyingGlassIcon className="w-10 h-10 text-gray-200 dark:text-gray-600 mx-auto mb-3" />
              <p className="text-gray-400 mb-1">没有匹配 "{search}" 的板块</p>
              <button onClick={() => setSearch('')} className="text-sm text-blue-500 hover:text-blue-600 transition-colors">
                清除搜索
              </button>
            </>
          ) : (
            <>
              <FolderPlusIcon className="w-10 h-10 text-gray-200 dark:text-gray-600 mx-auto mb-3" />
              <p className="text-gray-400 mb-4">还没有任何板块</p>
              {user ? (
                <Link
                  href="/boards/new"
                  className="inline-flex items-center gap-1.5 px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl font-medium transition-colors"
                >
                  <PlusIcon className="w-4 h-4" />
                  创建第一个板块
                </Link>
              ) : (
                <Link
                  href="/auth/login"
                  className="inline-block px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl font-medium transition-colors"
                >
                  登录后创建
                </Link>
              )}
            </>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map(board => <BoardCard key={board.id} board={board} />)}
        </div>
      )}
    </div>
  );
}