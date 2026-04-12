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

// ─── CreateBoardInline ─────────────────────────────────────────────────────────

function CreateBoardInline({ slug, onCreated }: { slug: string; onCreated: (board: Board) => void }) {
  const { user } = useAuthStore();
  const [form, setForm] = useState({
    name: '',
    slug: slug,
    description: '',
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
    setError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) { setError('请填写板块名称'); return; }
    if (!form.slug.trim()) { setError('请填写板块标识'); return; }
    if (!/^[a-z0-9-]+$/.test(form.slug)) { setError('标识只能包含小写字母、数字和连字符'); return; }

    setSubmitting(true);
    try {
      const res = await boardApi.create({
        name: form.name.trim(),
        slug: form.slug.trim(),
        description: form.description.trim(),
      });
      onCreated(res.data.data);
    } catch (err: any) {
      setError(err.response?.data?.message || '创建失败，请稍后重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="w-full max-w-lg">
        {/* Icon + heading */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-blue-50 dark:bg-blue-900/30 mb-4">
            <ExclamationCircleIcon className="w-8 h-8 text-blue-500" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">板块不存在</h1>
          <p className="mt-2 text-gray-500 dark:text-gray-400">
            <span className="font-mono bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded text-sm text-blue-600 dark:text-blue-400">/{slug}</span>
            {' '}尚未创建，你可以现在创建它。
          </p>
        </div>

        {!user ? (
          <div className="text-center space-y-3">
            <p className="text-gray-500">需要登录才能创建板块</p>
            <Link
              href={`/auth/login?redirect=/boards/${slug}`}
              className="inline-block px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-medium transition-colors"
            >
              去登录
            </Link>
            <div className="pt-2">
              <Link href="/boards" className="text-sm text-gray-400 hover:text-blue-500 transition-colors">
                ← 返回板块列表
              </Link>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 space-y-5">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                板块名称 <span className="text-red-500">*</span>
              </label>
              <input
                name="name"
                value={form.name}
                onChange={handleChange}
                placeholder="例如：前端技术交流"
                className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                板块标识 <span className="text-red-500">*</span>
              </label>
              <div className="flex items-center rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 overflow-hidden focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-transparent transition">
                <span className="pl-4 pr-1 text-gray-400 text-sm select-none">/boards/</span>
                <input
                  name="slug"
                  value={form.slug}
                  onChange={handleChange}
                  placeholder="frontend"
                  className="flex-1 py-2.5 pr-4 bg-transparent text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none text-sm"
                />
              </div>
              <p className="mt-1 text-xs text-gray-400">仅小写字母、数字、连字符</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                板块描述
              </label>
              <textarea
                name="description"
                value={form.description}
                onChange={handleChange}
                placeholder="简单介绍这个板块的用途..."
                rows={3}
                className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition resize-none"
              />
            </div>

            {error && (
              <p className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 px-4 py-2 rounded-lg">
                {error}
              </p>
            )}

            <div className="flex items-center justify-between pt-1">
              <Link href="/boards" className="text-sm text-gray-400 hover:text-blue-500 transition-colors flex items-center gap-1">
                <ArrowLeftIcon className="w-3.5 h-3.5" />
                返回列表
              </Link>
              <button
                type="submit"
                disabled={submitting}
                className="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed text-white rounded-xl font-medium transition-colors"
              >
                <FolderPlusIcon className="w-4 h-4" />
                {submitting ? '创建中...' : '创建板块'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}

// ─── PostCard ──────────────────────────────────────────────────────────────────

function PostCard({ post }: { post: Post }) {
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(post.like_count);

  const handleLike = async (e: React.MouseEvent) => {
    e.preventDefault();
    try {
      if (liked) {
        await postApi.unlike(post.id);
        setLikeCount(prev => prev - 1);
      } else {
        await postApi.like(post.id);
        setLikeCount(prev => prev + 1);
      }
      setLiked(!liked);
    } catch (error) {
      console.error('Failed to like:', error);
    }
  };

  const excerpt = post.summary || post.content.replace(/<[^>]*>/g, '').slice(0, 180);

  return (
    <div className="group bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 hover:border-blue-200 dark:hover:border-blue-700 hover:shadow-md transition-all duration-200 p-5">
      <Link href={`/posts/${post.id}`} className="block mb-3">
        <h3 className="text-base font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 line-clamp-1 transition-colors">
          {post.title}
        </h3>
        {excerpt && (
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400 line-clamp-2 leading-relaxed">
            {excerpt}
          </p>
        )}
      </Link>

      <div className="flex items-center justify-between text-xs text-gray-400">
        <div className="flex items-center gap-3">
          <Link href={`/users/${post.author_id}`} className="flex items-center gap-1.5 hover:text-blue-500 transition-colors">
            <img
              src={post.author?.avatar || '/default-avatar.png'}
              alt={post.author?.username}
              className="w-5 h-5 rounded-full object-cover"
            />
            <span>{post.author?.username}</span>
          </Link>
          <span>{new Date(post.created_at).toLocaleDateString('zh-CN')}</span>
        </div>

        <div className="flex items-center gap-3">
          <button onClick={handleLike} className="flex items-center gap-1 hover:text-red-500 transition-colors">
            {liked
              ? <HeartSolidIcon className="w-4 h-4 text-red-500" />
              : <HeartIcon className="w-4 h-4" />}
            <span>{likeCount}</span>
          </button>
          <div className="flex items-center gap-1">
            <ChatBubbleLeftRightIcon className="w-4 h-4" />
            <span>{post.like_count}</span>
          </div>
          <div className="flex items-center gap-1">
            <EyeIcon className="w-4 h-4" />
            <span>{post.view_count}</span>
          </div>
        </div>
      </div>
    </div>
  );
}

// ─── BoardDetailPage ───────────────────────────────────────────────────────────

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
      console.log(res.data.data.slug)
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
      const res = await boardApi.getPosts(slug, { page, page_size: PAGE_SIZE });
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
          {posts.map(post => <PostCard key={post.id} post={post} />)}
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