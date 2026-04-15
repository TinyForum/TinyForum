// src/app/[locale]/boards/[slug]/page.tsx
'use client';

import { useState, useEffect, useCallback } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { Board, boardApi, Post, postApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
// import type { Board, Post } from '@/types';
import {
  ArrowLeftIcon,
  ChatBubbleLeftRightIcon,
  HeartIcon,
  EyeIcon,
  PlusIcon,
  ShieldCheckIcon,
  FolderPlusIcon,
  ExclamationCircleIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import { BoardPostCard } from '@/components/boards/BoardPostCard';
import { CreateBoardInline } from '@/components/boards/CreateBoardInline';

export default function BoardDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuthStore();
  const slug = params.slug as string;

  const [board, setBoard] = useState<Board | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loadingBoard, setLoadingBoard] = useState(true);
  const [loadingPosts, setLoadingPosts] = useState(false);
  const [notFound, setNotFound] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const PAGE_SIZE = 20;

  const loadBoard = useCallback(async () => {
    setLoadingBoard(true);
    setNotFound(false);
    try {
      const res = await boardApi.getBySlug(slug);

      console.log("data: ",res.data.data.slug)
    setBoard(res.data.data);
    } catch (error: any) {
      if (error.response?.status === 404) {
        setNotFound(true);
      } else {
        console.error('Failed to load board:', error);
        setNotFound(true);
      }
    } finally {
      setLoadingBoard(false);
    }
  }, [slug]);

  const loadPosts = useCallback(async () => {
    if (!board) return;
    setLoadingPosts(true);
    try {
      const res = await boardApi.getPostsBySlug(slug, { page, page_size: PAGE_SIZE });
      setPosts(res.data.data.list);
      setTotal(res.data.data.total);
    } catch (error) {
      console.error('Failed to load posts:', error);
    } finally {
      setLoadingPosts(false);
    }
  }, [slug, board, page]);

  useEffect(() => { loadBoard(); }, [loadBoard]);
  useEffect(() => { loadPosts(); }, [loadPosts]);

  const handleBoardCreated = (newBoard: Board) => {
    setBoard(newBoard);
    setNotFound(false);
    setPosts([]);
    setTotal(0);
  };

  const totalPages = Math.ceil(total / PAGE_SIZE);

  // Loading skeleton
  if (loadingBoard) {
    return (
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-pulse">
        <div className="h-4 w-32 bg-gray-200 dark:bg-gray-700 rounded mb-6" />
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 mb-6">
          <div className="h-7 w-48 bg-gray-200 dark:bg-gray-700 rounded mb-3" />
          <div className="h-4 w-72 bg-gray-100 dark:bg-gray-700 rounded" />
        </div>
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="bg-white dark:bg-gray-800 rounded-xl p-5 h-24 border border-gray-100 dark:border-gray-700" />
          ))}
        </div>
      </div>
    );
  }

  // Board not found → show create form
  if (notFound) {
    return <CreateBoardInline slug={slug} onCreated={handleBoardCreated} />;
  }

  if (!board) return null;

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Breadcrumb */}
      <nav className="mb-6 flex items-center gap-1.5 text-sm text-gray-400">
        <Link href="/" className="hover:text-blue-500 transition-colors">首页</Link>
        <span>/</span>
        <Link href="/boards" className="hover:text-blue-500 transition-colors">板块</Link>
        <span>/</span>
        <span className="text-gray-700 dark:text-gray-200 font-medium">{board.name}</span>
      </nav>

      {/* Board header */}
      <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm p-6 mb-6">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1 min-w-0">
            <h1 className="text-xl font-bold text-gray-900 dark:text-white truncate">{board.name}</h1>
            {board.description && (
              <p className="mt-1.5 text-sm text-gray-500 dark:text-gray-400 leading-relaxed">{board.description}</p>
            )}
            <div className="flex items-center gap-4 mt-3 text-xs text-gray-400">
              <span>帖子 <strong className="text-gray-600 dark:text-gray-300">{board.post_count}</strong></span>
              <span>主题 <strong className="text-gray-600 dark:text-gray-300">{board.thread_count}</strong></span>
              <span>今日 <strong className="text-gray-600 dark:text-gray-300">{board.today_count}</strong></span>
            </div>
          </div>

          <div className="flex shrink-0 gap-2">
            <Link
              href={`/posts/new?board_id=${board.id}`}
              className="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl font-medium transition-colors"
            >
              <PlusIcon className="w-4 h-4" />
              发帖
            </Link>
            {user && (
              <Link
                href={`/boards/${slug}/apply`}
                className="flex items-center gap-1.5 px-4 py-2 border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 text-sm text-gray-600 dark:text-gray-300 rounded-xl transition-colors"
              >
                <ShieldCheckIcon className="w-4 h-4" />
                申请版主
              </Link>
            )}
          </div>
        </div>
      </div>

      {/* Posts */}
      {loadingPosts ? (
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="bg-white dark:bg-gray-800 rounded-xl p-5 h-24 border border-gray-100 dark:border-gray-700 animate-pulse" />
          ))}
        </div>
      ) : posts.length === 0 ? (
        <div className="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
          <ChatBubbleLeftRightIcon className="w-10 h-10 text-gray-300 mx-auto mb-3" />
          <p className="text-gray-500 dark:text-gray-400 mb-4">还没有帖子，来发第一帖吧</p>
          <Link
            href={`/posts/new?board_id=${board.id}`}
            className="inline-flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl transition-colors"
          >
            <PlusIcon className="w-4 h-4" />
            发布新帖
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map(post => <BoardPostCard key={post.id} post={post} />)}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2 mt-8">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            <ChevronLeftIcon className="w-4 h-4" />
          </button>
          <span className="text-sm text-gray-500 px-2">
            第 <strong className="text-gray-900 dark:text-white">{page}</strong> / {totalPages} 页
          </span>
          <button
            onClick={() => setPage(p => Math.min(totalPages, p + 1))}
            disabled={page >= totalPages}
            className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            <ChevronRightIcon className="w-4 h-4" />
          </button>
        </div>
      )}
    </div>
  );
}