"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";
import { useMePosts } from "@/features/user/hooks/useMePosts";
import { useEffect } from "react";

export function MyPostsTable() {
  const t = useTranslations("User");
  const {
    list: posts,
    isLoading,
    error,
    has_more: hasMore,
    loadMore,
    refresh,
  } = useMePosts();

  useEffect(() => {
    refresh();
  }, [refresh]);

  // 加载骨架屏
  if (isLoading && posts.length === 0) {
    return (
      <div className="space-y-4">
        <div className="flex justify-end">
          <div className="skeleton h-8 w-20 rounded-lg"></div>
        </div>
        <div className="overflow-x-auto rounded-xl border border-base-200 bg-base-100 shadow-sm">
          <table className="table">
            <thead>
              <tr className="bg-base-200/60">
                <th>{t("title")}</th>
                <th>{t("board")}</th>
                <th>{t("likes")}</th>
                <th>{t("comments")}</th>
                <th>{t("created_at")}</th>
                <th>{t("actions")}</th>
              </tr>
            </thead>
            <tbody>
              {[...Array(3)].map((_, idx) => (
                <tr key={idx} className="animate-pulse">
                  <td colSpan={6}>
                    <div className="skeleton h-12 w-full rounded-lg"></div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 rounded-2xl border border-error/20 bg-error/5 p-8 text-center shadow-sm">
        <div className="rounded-full bg-error/10 p-3">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-8 w-8 text-error"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
        </div>
        <p className="text-base-content/80">{error}</p>
        <button className="btn btn-primary btn-sm" onClick={refresh}>
          重试
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-5">
      {/* 顶部操作栏 */}
      <div className="flex items-center justify-between">
        <div className="text-sm text-base-content/60">
          {t("total_posts", { count: posts.length })}
        </div>
        <button
          className="btn btn-ghost btn-sm gap-1 text-base-content/70 transition-all hover:bg-base-200"
          onClick={refresh}
          disabled={isLoading}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`}
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
          {t("refresh")}
        </button>
      </div>

      {/* 表格容器 */}
      <div className="overflow-x-auto rounded-2xl border border-base-200 bg-base-100 shadow-md transition-all hover:shadow-lg">
        <table className="table table-zebra table-md">
          <thead className="bg-gradient-to-r from-red-50 to-base-100 dark:from-red-950/20 dark:to-base-200">
            <tr className="text-sm uppercase tracking-wide text-base-content/80">
              <th className="rounded-tl-2xl">{t("title")}</th>
              <th>{t("board")}</th>
              <th className="text-center">{t("likes")}</th>
              <th className="text-center">{t("comments")}</th>
              <th>{t("created_at")}</th>
              <th className="rounded-tr-2xl text-center">{t("actions")}</th>
            </tr>
          </thead>
          <tbody>
            {posts.length === 0 ? (
              <tr>
                <td colSpan={6} className="py-12 text-center">
                  <div className="flex flex-col items-center gap-2 text-base-content/50">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-12 w-12"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={1.5}
                        d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                      />
                    </svg>
                    <span>{t("no_posts")}</span>
                  </div>
                </td>
              </tr>
            ) : (
              posts.map((post) => (
                <tr
                  key={post.id}
                  className="group transition-colors hover:bg-red-50/30 dark:hover:bg-red-950/10"
                >
                  <td className="font-medium">
                    <Link
                      href={`/posts/${post.id}`}
                      className="line-clamp-2 hover:text-primary hover:underline"
                    >
                      {post.title}
                    </Link>
                  </td>
                  <td>
                    <span className="badge badge-ghost badge-sm rounded-full bg-base-200/70 px-3 font-normal">
                      {post.board_name || "—"}
                    </span>
                  </td>
                  <td className="text-center">
                    <div className="flex items-center justify-center gap-1">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-4 w-4 text-red-400"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          fillRule="evenodd"
                          d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="font-mono text-sm">
                        {post.likes_count}
                      </span>
                    </div>
                  </td>
                  <td className="text-center">
                    <div className="flex items-center justify-center gap-1">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-4 w-4 text-blue-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                        />
                      </svg>
                      <span className="font-mono text-sm">
                        {post.comment_count}
                      </span>
                    </div>
                  </td>
                  <td className="text-sm text-base-content/70">
                    {new Date(post.created_at).toLocaleDateString(undefined, {
                      year: "numeric",
                      month: "short",
                      day: "numeric",
                    })}
                  </td>
                  <td className="text-center">
                    <button
                      className="btn btn-error btn-xs min-h-0 h-8 gap-1 rounded-full px-4 shadow-sm transition-all hover:shadow-md"
                      onClick={() => {
                        // TODO: 实现删除逻辑
                        console.log("delete post", post.id);
                      }}
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-3.5 w-3.5"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          fillRule="evenodd"
                          d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                          clipRule="evenodd"
                        />
                      </svg>
                      {t("delete")}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* 加载更多区域 */}
      {hasMore && (
        <div className="flex justify-center pt-4">
          <button
            className="btn btn-outline btn-primary rounded-full px-8 shadow-sm transition-all hover:shadow-md"
            onClick={loadMore}
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <span className="loading loading-spinner loading-xs"></span>
                加载中...
              </>
            ) : (
              <>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="mr-2 h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 13l-7 7-7-7m14-8l-7 7-7-7"
                  />
                </svg>
                {t("load_more")}
              </>
            )}
          </button>
        </div>
      )}
    </div>
  );
}
