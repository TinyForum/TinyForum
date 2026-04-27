// src/app/[locale]/boards/[slug]/page.tsx
"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useQuery, keepPreviousData } from "@tanstack/react-query"; // 导入 keepPreviousData
import { boardApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import {
  PlusIcon,
  ShieldCheckIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChatBubbleLeftRightIcon,
} from "@heroicons/react/24/outline";
import { BoardPostCard } from "@/components/boards/BoardPostCard";
import { CreateBoardInline } from "@/components/boards/CreateBoardInline";
import type { Board, PageData } from "@/lib/api/types";
import { BoardPostListItem } from "@/lib/api/modules/boards";

const PAGE_SIZE = 20;

export default function BoardDetailPage() {
  const params = useParams();
  const router = useRouter();
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
    retry: (failureCount, error: any) => {
      if (error.response?.status === 404) return false;
      return failureCount < 2;
    },
  });

  // 查询帖子列表 - 添加泛型类型
  const {
    data: postsData,
    isLoading: loadingPosts,
    error: postsError,
  } = useQuery<PageData<BoardPostListItem>>({
    queryKey: ["board-posts", slug, page],
    queryFn: () =>
      boardApi
        .getPostsBySlug(slug, { page, page_size: PAGE_SIZE })
        .then((res) => res.data.data), // 返回类型应为 PageData<BoardPostListItem>
    enabled: !!board,
    placeholderData: keepPreviousData, // 替代 keepPreviousData: true
  });

  const posts = postsData?.list ?? [];
  const total = postsData?.total ?? 0;
  const totalPages = Math.ceil(total / PAGE_SIZE);

  // 当板块 slug 变化时重置页码
  useEffect(() => {
    setPage(1);
  }, [slug]);

  const handleBoardCreated = (newBoard: Board) => {
    // 板块创建成功后的回调（可留空，或使用 router.refresh() 刷新）
    console.log("Board created", newBoard);
  };

  const isNotFound = boardError && (boardError as any).response?.status === 404;

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
            {posts.map((post) => (
              <BoardPostCard key={post.id} post={post} /> // 现在 BoardPostCard 应兼容 BoardPostListItem
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

// 其余子组件 (BoardHeader, EmptyState, Pagination, LoadingSkeleton, PostsSkeleton) 与之前相同，略...

// 子组件：板块头部
function BoardHeader({
  board,
  slug,
  user,
}: {
  board: Board;
  slug: string;
  user: any;
}) {
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
