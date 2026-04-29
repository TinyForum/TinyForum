"use client";

import { useTranslations } from "next-intl";

interface ReviewerStatsProps {
  stats: {
    pending: number;
    reported: number;
    reviewedToday: number;
  };
}

export function ReviewerStats({ stats }: ReviewerStatsProps) {
  const t = useTranslations("Review");
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-base-content/60">
                {t("pending_review")}
              </p>
              <p className="text-2xl font-bold text-warning">{stats.pending}</p>
            </div>
            <div className="text-3xl">📋</div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-base-content/60">
                {t("reported_content")}
              </p>
              <p className="text-2xl font-bold text-error">{stats.reported}</p>
            </div>
            <div className="text-3xl">🚫</div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-base-content/60">
                {t("reviewed_today")}
              </p>
              <p className="text-2xl font-bold text-success">
                {stats.reviewedToday}
              </p>
            </div>
            <div className="text-3xl">✅</div>
          </div>
        </div>
      </div>
    </div>
  );
}
