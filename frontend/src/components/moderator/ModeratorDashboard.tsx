// components/moderator/ModeratorDashboard.tsx

import { ModeratorBoard } from "@/lib/api/modules/moderator";

interface ModeratorPermissions {
  canDeletePost: boolean;
  canPinPost: boolean;
  canEditAnyPost: boolean;
  canManageModerator: boolean;
  canBanUser: boolean;
}

interface DashboardStats {
  postCount: number;
  reportCount: number;
  bannedCount: number;
}

interface ModeratorDashboardProps {
  board: ModeratorBoard;
  permissions: ModeratorPermissions;
  stats: DashboardStats;
  t: (key: string) => string;
}

export function ModeratorDashboard({
  stats,
  t,
}: ModeratorDashboardProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div className="stat bg-base-100 rounded-lg border border-base-300 p-4 shadow-sm hover:shadow-md transition-shadow">
        <div className="stat-title text-base-content/60">{t("total_posts")}</div>
        <div className="stat-value text-primary text-3xl font-bold mt-2">
          {stats.postCount || 0}
        </div>
        <div className="stat-desc text-xs text-base-content/40 mt-1">
          {t("in_current_board")}
        </div>
      </div>
      
      <div className="stat bg-base-100 rounded-lg border border-base-300 p-4 shadow-sm hover:shadow-md transition-shadow">
        <div className="stat-title text-base-content/60">{t("pending_reports")}</div>
        <div className="stat-value text-warning text-3xl font-bold mt-2">
          {stats.reportCount || 0}
        </div>
        <div className="stat-desc text-xs text-base-content/40 mt-1">
          {t("awaiting_review")}
        </div>
      </div>
      
      <div className="stat bg-base-100 rounded-lg border border-base-300 p-4 shadow-sm hover:shadow-md transition-shadow">
        <div className="stat-title text-base-content/60">{t("banned_users")}</div>
        <div className="stat-value text-error text-3xl font-bold mt-2">
          {stats.bannedCount || 0}
        </div>
        <div className="stat-desc text-xs text-base-content/40 mt-1">
          {t("currently_banned")}
        </div>
      </div>
    </div>
  );
}