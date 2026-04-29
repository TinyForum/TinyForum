// components/member/MemberFavorites.tsx
"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";

interface MemberFavoritesProps {
  onRemove: (id: number) => void;
  favorites: Favorite[];
}
interface Favorite {
  id: number;
  post_id: number;
  post_title: string;
  board_name: string;
  like_count: number;
  excerpt: string;
  created_at: string;
}

export function MemberFavorites({ onRemove, favorites }: MemberFavoritesProps) {
  // TODO: 从 props 获取数据

  const t = useTranslations("Member");
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
      {favorites.map((item: Favorite) => (
        <div key={item.id} className="card bg-base-100 border border-base-300">
          <div className="card-body">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <h3 className="card-title">
                  <Link
                    href={`/post/${item.post_id}`}
                    className="hover:link-hover"
                  >
                    {item.post_title}
                  </Link>
                </h3>
                <p className="text-sm text-base-content/60 mt-1">
                  {t("board")}: {item.board_name} | {t("likes")}:{" "}
                  {item.like_count}
                </p>
                <p className="text-sm mt-2 line-clamp-2">{item.excerpt}</p>
                <p className="text-xs text-base-content/40 mt-2">
                  {t("favorited_at")}:{" "}
                  {new Date(item.created_at).toLocaleDateString()}
                </p>
              </div>
              <button
                onClick={() => onRemove(item.id)}
                className="btn btn-sm btn-outline btn-error"
              >
                {t("remove")}
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
