// components/member/MemberStats.tsx
"use client";

import { useTranslations } from "next-intl";

interface MemberStatsProps {
  stats: {
    posts: number;
    comments: number;
    favorites: number;
    unreadNotif: number;
  };
}

export function MemberStats({ stats }: MemberStatsProps) {
  const t = useTranslations("Member");
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div className="card bg-base-100 border border-base-300 hover:shadow-lg transition-shadow">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-base-content/60 text-sm">{t("total_posts")}</p>
              <p className="text-3xl font-bold mt-1">{stats.posts}</p>
            </div>
            <div className="text-3xl text-primary">📝</div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300 hover:shadow-lg transition-shadow">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-base-content/60 text-sm">
                {t("total_comments")}
              </p>
              <p className="text-3xl font-bold mt-1">{stats.comments}</p>
            </div>
            <div className="text-3xl text-secondary">💬</div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300 hover:shadow-lg transition-shadow">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-base-content/60 text-sm">
                {t("total_favorites")}
              </p>
              <p className="text-3xl font-bold mt-1">{stats.favorites}</p>
            </div>
            <div className="text-3xl text-error">❤️</div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300 hover:shadow-lg transition-shadow">
        <div className="card-body">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-base-content/60 text-sm">
                {t("unread_notifications")}
              </p>
              <p className="text-3xl font-bold mt-1">{stats.unreadNotif}</p>
            </div>
            <div className="text-3xl text-warning">🔔</div>
          </div>
        </div>
      </div>
    </div>
  );
}
