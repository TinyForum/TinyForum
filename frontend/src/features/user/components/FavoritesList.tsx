// components/user/FavoritesList.tsx
"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";

interface FavoritesListProps {
  favorites: favorite[];
}
interface favorite {
  id: number;
  title: string;
  board_name: string;
  like_count: number;
  excerpt: string;
}

export function FavoritesList({ favorites }: FavoritesListProps) {
  const t = useTranslations("User");
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
    <div className="grid gap-4">
      {favorites.map((post: favorite) => (
        <div key={post.id} className="card bg-base-100 border border-base-300">
          <div className="card-body">
            <h3 className="card-title">
              <Link href={`/post/${post.id}`} className="hover:link-hover">
                {post.title}
              </Link>
            </h3>
            <p className="text-sm text-base-content/60">
              {t("board")}: {post.board_name} | {t("likes")}: {post.like_count}
            </p>
            <p className="text-sm mt-2">{post.excerpt}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-sm btn-outline btn-error">
                {t("remove")}
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
