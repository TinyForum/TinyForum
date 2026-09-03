// boards/page.tsx 客户端组件
"use client";

import { useState } from "react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useAuthStore } from "@/store/auth";
import {
  MagnifyingGlassIcon,
  FolderPlusIcon,
  PlusIcon,
  HashtagIcon,
  ChatBubbleLeftRightIcon,
  FireIcon,
  RectangleGroupIcon,
} from "@heroicons/react/24/outline";
import { BoardCard } from "@/features/boards/components/BoardCard";
import { useBoard } from "@/features/boards/hooks/useBoard";

export default function BoardsClient() {
  const { user } = useAuthStore();
  const t = useTranslations("Board");
  const [search, setSearch] = useState("");
  const [page] = useState(1);
  const [pageSize] = useState(20);

  // 查询：分页获取板块列表
  const { boards, total, loading } = useBoard({ page, pageSize });

  // 搜索过滤
  const filtered = search.trim()
    ? boards.filter((board) => {
        const q = search.toLowerCase();
        return (
          board.name.toLowerCase().includes(q) ||
          (board.description && board.description.toLowerCase().includes(q))
        );
      })
    : boards;

  // 统计数据
  const totalPosts = boards.reduce(
    (sum, board) => sum + (board.post_count || 0),
    0,
  );
  const totalToday = boards.reduce(
    (sum, board) => sum + (board.today_count || 0),
    0,
  );

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-200 via-base-100 to-base-200">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12">
        {/* Hero Section */}
        <div className="text-center mb-10 md:mb-12">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-red-500 to-red-600 shadow-lg mb-4">
            <HashtagIcon className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-red-600 to-red-500 bg-clip-text text-transparent">
             {t("explore_boards")}
          </h1>
          <p className="text-base-content/60 mt-2 text-sm md:text-base">
             {t("explore_boards_subtitle")}
          </p>
        </div>

        {/* Stats Cards */}
        {!loading && boards.length > 0 && (
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
            <div className="card bg-base-100 shadow-md border border-base-200 hover:shadow-lg transition-shadow">
              <div className="card-body p-4 flex-row items-center justify-between">
                <div>
                   <p className="text-base-content/60 text-sm">{t("total_boards")}</p>
                  <p className="text-2xl font-bold text-primary">{total}</p>
                </div>
                <div className="w-10 h-10 rounded-lg bg-red-50 dark:bg-red-900/20 flex items-center justify-center">
                  <RectangleGroupIcon className="w-5 h-5 text-red-500" />
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-md border border-base-200 hover:shadow-lg transition-shadow">
              <div className="card-body p-4 flex-row items-center justify-between">
                <div>
                   <p className="text-base-content/60 text-sm">{t("total_posts")}</p>
                  <p className="text-2xl font-bold text-primary">
                    {totalPosts}
                  </p>
                </div>
                <div className="w-10 h-10 rounded-lg bg-red-50 dark:bg-red-900/20 flex items-center justify-center">
                  <ChatBubbleLeftRightIcon className="w-5 h-5 text-red-500" />
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-md border border-base-200 hover:shadow-lg transition-shadow">
              <div className="card-body p-4 flex-row items-center justify-between">
                <div>
                   <p className="text-base-content/60 text-sm">{t("active_today")}</p>
                  <p className="text-2xl font-bold text-orange-500">
                    {totalToday}
                  </p>
                </div>
                <div className="w-10 h-10 rounded-lg bg-orange-50 dark:bg-orange-900/20 flex items-center justify-center">
                  <FireIcon className="w-5 h-5 text-orange-500" />
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Search and Create Bar */}
        <div className="flex flex-col sm:flex-row gap-4 mb-8">
          <div className="flex-1 relative">
            <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
               placeholder={t("search_placeholder")}
              className="w-full pl-11 pr-4 py-3 rounded-xl border border-base-200 bg-base-100 text-base-content placeholder-base-content/40 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
            />
            {search && (
              <button
                onClick={() => setSearch("")}
                className="absolute right-4 top-1/2 -translate-y-1/2 text-base-content/40 hover:text-base-content/60 transition-colors"
              >
                <span className="text-xl leading-none">×</span>
              </button>
            )}
          </div>

          {user && (
            <Link
              href="/posts/new"
              className="inline-flex items-center justify-center gap-2 px-5 py-3 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white text-sm font-medium rounded-xl shadow-md hover:shadow-lg transition-all"
            >
              <PlusIcon className="w-4 h-4" />
              {t("publish_new_post")}
            </Link>
          )}
        </div>

        {/* Board List */}
        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {[...Array(6)].map((_, i) => (
              <div
                key={i}
                className="card bg-base-100 shadow-md border border-base-200"
              >
                <div className="card-body p-5">
                  <div className="h-6 bg-base-200 rounded-lg w-1/3 mb-3 animate-pulse" />
                  <div className="h-4 bg-base-200 rounded-lg w-full mb-2 animate-pulse" />
                  <div className="h-4 bg-base-200 rounded-lg w-2/3 animate-pulse" />
                  <div className="flex gap-4 mt-3">
                    <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                    <div className="h-3 bg-base-200 rounded w-16 animate-pulse" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : filtered.length === 0 ? (
          <div className="card bg-base-100 shadow-md border border-base-200">
            <div className="card-body py-16 text-center">
              {search ? (
                <>
                  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-50 dark:bg-red-900/20 flex items-center justify-center">
                    <MagnifyingGlassIcon className="w-10 h-10 text-red-400" />
                  </div>
                   <h3 className="text-lg font-semibold text-base-content mb-2">
                      {t("no_boards_found")}
                   </h3>
                   <p className="text-base-content/60 text-sm mb-4">
                     {t("no_board_match", { search })}
                  </p>
                  <button
                    onClick={() => setSearch("")}
                    className="btn btn-primary btn-sm"
                  >
                     {t("clear_search")}
                  </button>
                </>
              ) : (
                <>
                  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-50 dark:bg-red-900/20 flex items-center justify-center">
                    <FolderPlusIcon className="w-10 h-10 text-red-400" />
                  </div>
                  <h3 className="text-lg font-semibold text-base-content mb-2">
                     {t("no_boards_yet")}
                  </h3>
                  <p className="text-base-content/60 text-sm mb-4">
                     {t("no_boards_desc")}
                  </p>
                  {user ? (
                    <Link href="/boards/new" className="btn btn-primary">
                      <PlusIcon className="w-4 h-4" />
                       {t("create_first_board")}
                    </Link>
                  ) : (
                    <Link href="/auth/login" className="btn btn-primary">
                       {t("login_to_create")}
                    </Link>
                  )}
                </>
              )}
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {filtered.map((board) => (
              <BoardCard key={board.id} board={board} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
