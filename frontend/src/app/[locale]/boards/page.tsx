"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { Board, boardApi } from "@/shared/api";
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

export default function BoardsPage() {
  const { user } = useAuthStore();
  const [boards, setBoards] = useState<Board[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [page] = useState(1);
  const [pageSize] = useState(20);
  const [total, setTotal] = useState(0);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await boardApi.list({ page, page_size: pageSize });
      const pageData = res.data.data;

      // 添加安全检查
      if (pageData) {
        setBoards(pageData.list || []);
        setTotal(pageData.total || 0);
      } else {
        setBoards([]);
        setTotal(0);
      }
    } catch (err) {
      console.error("Failed to load boards:", err);
      setBoards([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  useEffect(() => {
    load();
  }, [load]);

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
            探索板块
          </h1>
          <p className="text-base-content/60 mt-2 text-sm md:text-base">
            发现你感兴趣的讨论区，参与精彩对话
          </p>
        </div>

        {/* Stats Cards */}
        {!loading && boards.length > 0 && (
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
            <div className="card bg-base-100 shadow-md border border-base-200 hover:shadow-lg transition-shadow">
              <div className="card-body p-4 flex-row items-center justify-between">
                <div>
                  <p className="text-base-content/60 text-sm">总板块</p>
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
                  <p className="text-base-content/60 text-sm">总帖子</p>
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
                  <p className="text-base-content/60 text-sm">今日活跃</p>
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
              placeholder="搜索板块名称或描述..."
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
              href="/boards/new"
              className="inline-flex items-center justify-center gap-2 px-5 py-3 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white text-sm font-medium rounded-xl shadow-md hover:shadow-lg transition-all"
            >
              <PlusIcon className="w-4 h-4" />
              创建新板块
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
                    未找到相关板块
                  </h3>
                  <p className="text-base-content/60 text-sm mb-4">
                    {`没有找到匹配 ${search} 的板块`}
                  </p>
                  <button
                    onClick={() => setSearch("")}
                    className="btn btn-primary btn-sm"
                  >
                    清除搜索
                  </button>
                </>
              ) : (
                <>
                  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-50 dark:bg-red-900/20 flex items-center justify-center">
                    <FolderPlusIcon className="w-10 h-10 text-red-400" />
                  </div>
                  <h3 className="text-lg font-semibold text-base-content mb-2">
                    暂无板块
                  </h3>
                  <p className="text-base-content/60 text-sm mb-4">
                    还没有任何板块，快来创建第一个吧
                  </p>
                  {user ? (
                    <Link href="/boards/new" className="btn btn-primary">
                      <PlusIcon className="w-4 h-4" />
                      创建第一个板块
                    </Link>
                  ) : (
                    <Link href="/auth/login" className="btn btn-primary">
                      登录后创建
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
