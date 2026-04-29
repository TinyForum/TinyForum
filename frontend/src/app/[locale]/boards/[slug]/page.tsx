// src/app/[locale]/boards/[slug]/page.tsx
"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { boardApi } from "@/shared/api";
import { useAuthStore } from "@/store/auth";
import {
  PlusIcon,
  ShieldCheckIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChatBubbleLeftRightIcon,
} from "@heroicons/react/24/outline";

// import type { Board, PageData } from "@/shared/api/types";
import type { BoardPostListItem } from "@/shared/api/modules/boards";
import { Board, PageData } from "@/shared/api/types";
import { CreateBoardInline } from "@/features/boards/components/CreateBoardInline";
import { BoardPostCard } from "@/features/boards/components/BoardPostCard";

const PAGE_SIZE = 20;

// 定义错误响应类型
interface ErrorResponse {
  response?: {
    status: number;
  };
}

export default function BoardDetailPage() {
  const params = useParams();
  const { user } = useAuthStore();
  const slug = params.slug as string;

  const [page, setPage] = useState(1);

  // 查询板块信息
  const {
    data: board,
    isLoading: loadingBoard,
    error: boardError,
  } = useQuery({
    queryKey: ["board", slug],
    queryFn: () => boardApi.getBySlug(slug).then((res) => res.data.data),
    retry: (failureCount, error: ErrorResponse) => {
      if (error.response?.status === 404) return false;
      return failureCount < 2;
    },
  });

  // 查询帖子列表 - 修复类型问题
  const { data: postsData, isLoading: loadingPosts } = useQuery({
    queryKey: ["board-posts", slug, page],
    queryFn: async () => {
      const res = await boardApi.getPostsBySlug(slug, {
        page,
        page_size: PAGE_SIZE,
      });
      // 确保返回的数据不为 undefined
      return res.data.data as PageData<BoardPostListItem>;
    },
    enabled: !!board,
    placeholderData: keepPreviousData,
  });

  // 安全地访问数据
  const posts = postsData?.list ?? [];
  const total = postsData?.total ?? 0;
  const totalPages = Math.ceil(total / PAGE_SIZE);

  // 当板块 slug 变化时重置页码
  useEffect(() => {
    setPage(1);
  }, [slug]);

  const handleBoardCreated = (newBoard: Board) => {
    console.log("Board created", newBoard);
  };

  // 检查是否为 404 错误
  const isNotFound =
    boardError && (boardError as ErrorResponse).response?.status === 404;

  if (loadingBoard) {
    return <LoadingSkeleton />;
  }

  if (isNotFound) {
    return <CreateBoardInline slug={slug} onCreated={handleBoardCreated} />;
  }

  if (!board) {
    return <div className="text-center py-10">板块加载失败</div>;
  }

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <nav className="mb-6 flex items-center gap-1.5 text-sm text-gray-400">
        <Link href="/" className="hover:text-blue-500 transition-colors">
          首页
        </Link>
        <span>/</span>
        <Link href="/boards" className="hover:text-blue-500 transition-colors">
          板块
        </Link>
        <span>/</span>
        <span className="text-gray-700 dark:text-gray-200 font-medium">
          {board.name}
        </span>
      </nav>

      <BoardHeader board={board} slug={slug} user={user} />

      {loadingPosts ? (
        <PostsSkeleton />
      ) : posts.length === 0 ? (
        <EmptyState boardId={board.id} />
      ) : (
        <>
          <div className="space-y-3">
            {posts.map((post: BoardPostListItem) => (
              <BoardPostCard key={post.id} post={post} />
            ))}
          </div>
          {totalPages > 1 && (
            <Pagination
              page={page}
              totalPages={totalPages}
              onPageChange={setPage}
            />
          )}
        </>
      )}
    </div>
  );
}

// 子组件：板块头部
interface BoardHeaderProps {
  board: Board;
  slug: string;
  user: {
    id: number;
    username: string;
    role?: string;
  } | null;
}

function BoardHeader({ board, slug, user }: BoardHeaderProps) {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm p-6 mb-6">
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <h1 className="text-xl font-bold text-gray-900 dark:text-white truncate">
            {board.name}
          </h1>
          {board.description && (
            <p className="mt-1.5 text-sm text-gray-500 dark:text-gray-400 leading-relaxed">
              {board.description}
            </p>
          )}
          <div className="flex items-center gap-4 mt-3 text-xs text-gray-400">
            <span>
              帖子{" "}
              <strong className="text-gray-600 dark:text-gray-300">
                {board.post_count}
              </strong>
            </span>
            <span>
              主题{" "}
              <strong className="text-gray-600 dark:text-gray-300">
                {board.thread_count}
              </strong>
            </span>
            <span>
              今日{" "}
              <strong className="text-gray-600 dark:text-gray-300">
                {board.today_count}
              </strong>
            </span>
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
  );
}

// 空状态组件
function EmptyState({ boardId }: { boardId: number }) {
  return (
    <div className="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
      <ChatBubbleLeftRightIcon className="w-10 h-10 text-gray-300 mx-auto mb-3" />
      <p className="text-gray-500 dark:text-gray-400 mb-4">
        还没有帖子，来发第一帖吧
      </p>
      <Link
        href={`/posts/new?board_id=${boardId}`}
        className="inline-flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl transition-colors"
      >
        <PlusIcon className="w-4 h-4" />
        发布新帖
      </Link>
    </div>
  );
}

// 分页组件
function Pagination({
  page,
  totalPages,
  onPageChange,
}: {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page === 1}
        className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        aria-label="上一页"
      >
        <ChevronLeftIcon className="w-4 h-4" />
      </button>
      <span className="text-sm text-gray-500 px-2">
        第 <strong className="text-gray-900 dark:text-white">{page}</strong> /{" "}
        {totalPages} 页
      </span>
      <button
        onClick={() => onPageChange(page + 1)}
        disabled={page >= totalPages}
        className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        aria-label="下一页"
      >
        <ChevronRightIcon className="w-4 h-4" />
      </button>
    </div>
  );
}

// 加载骨架屏
function LoadingSkeleton() {
  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-pulse">
      <div className="h-4 w-32 bg-gray-200 dark:bg-gray-700 rounded mb-6" />
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 mb-6">
        <div className="h-7 w-48 bg-gray-200 dark:bg-gray-700 rounded mb-3" />
        <div className="h-4 w-72 bg-gray-100 dark:bg-gray-700 rounded" />
      </div>
      <div className="space-y-3">
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            className="bg-white dark:bg-gray-800 rounded-xl p-5 h-24 border border-gray-100 dark:border-gray-700"
          />
        ))}
      </div>
    </div>
  );
}

// 帖子列表骨架屏
function PostsSkeleton() {
  return (
    <div className="space-y-3">
      {[...Array(5)].map((_, i) => (
        <div
          key={i}
          className="bg-white dark:bg-gray-800 rounded-xl p-5 h-24 border border-gray-100 dark:border-gray-700 animate-pulse"
        />
      ))}
    </div>
  );
}
