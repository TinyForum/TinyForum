// components/moderator/ModeratorDashboard.tsx
interface ModeratorDashboardProps {
  board: any;
  permissions: any;
  stats: {
    postCount: number;
    reportCount: number;
    bannedCount: number;
  };
  t: (key: string) => string;
}

export function ModeratorDashboard({
  board,
  permissions,
  stats,
  t,
}: ModeratorDashboardProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div className="stat bg-base-100 rounded-lg border border-base-300 p-4">
        <div className="stat-title">{t("total_posts")}</div>
        <div className="stat-value text-primary">{stats.postCount || 0}</div>
      </div>
      <div className="stat bg-base-100 rounded-lg border border-base-300 p-4">
        <div className="stat-title">{t("pending_reports")}</div>
        <div className="stat-value text-warning">{stats.reportCount || 0}</div>
      </div>
      <div className="stat bg-base-100 rounded-lg border border-base-300 p-4">
        <div className="stat-title">{t("banned_users")}</div>
        <div className="stat-value text-error">{stats.bannedCount || 0}</div>
      </div>
    </div>
  );
}
