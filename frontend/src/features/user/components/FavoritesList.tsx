// components/user/FavoritesList.tsx
"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";
import { useMeFavorites } from "@/features/user/hooks/useMeFavorites";

export function FavoritesList() {
  const t = useTranslations("User");
  const {
    list: favorites,
    isLoading,
    error,
    has_more: hasMore,
    loadMore,
    refresh,
    removeFavorite,
    isRemoving,
  } = useMeFavorites();

  if (isLoading && favorites.length === 0) {
    return (
      <div className="grid gap-4">
        {[...Array(3)].map((_, idx) => (
          <div
            key={idx}
            className="skeleton h-28 w-full rounded-2xl border border-base-300"
          ></div>
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 rounded-2xl border border-error/20 bg-error/5 p-8 text-center shadow-sm">
        <p className="text-base-content/80">{error}</p>
        <button className="btn btn-primary btn-sm" onClick={refresh}>
          {t("retry")}
        </button>
      </div>
    );
  }

  if (favorites.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_favorites")}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-5">
      <div className="grid gap-4">
        {favorites.map((post) => (
          <div
            key={post.id}
            className="card bg-base-100 border border-base-300 transition-all hover:shadow-md"
          >
            <div className="card-body">
              <h3 className="card-title">
                <Link
                  href={`/posts/${post.target_id}`}
                  className="hover:text-primary hover:underline"
                >
                  {post.title}
                </Link>
              </h3>
              <p className="text-sm text-base-content/60">
                {t("board")}: {post.board_name || "—"} | {t("likes")}:{" "}
                {post.like_count}
              </p>
              {post.summary && (
                <p className="text-sm mt-2 line-clamp-2">{post.summary}</p>
              )}
              <div className="card-actions justify-end">
                <button
                  className="btn btn-sm btn-outline btn-error"
                  onClick={() => removeFavorite(post.target_id)}
                  disabled={isRemoving}
                >
                  {t("remove")}
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {hasMore && (
        <div className="flex justify-center pt-2">
          <button
            className="btn btn-outline btn-primary rounded-full px-8"
            onClick={loadMore}
            disabled={isLoading}
          >
            {isLoading ? t("loading") : t("load_more")}
          </button>
        </div>
      )}
    </div>
  );
}
