import { Users, Trophy, TrendingUp } from "lucide-react";
import { useTranslations } from "next-intl";

// 统计卡片组件
export function StatsCards({
  totalUsers,
  topScore,
}: {
  totalUsers: number;
  topScore: number;
}) {
  const t = useTranslations("Leaderboard");
  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
      <div className="bg-gradient-to-br from-primary/10 to-primary/5 rounded-2xl p-4 text-center border border-primary/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <Users className="w-5 h-5 text-primary" />
          <span className="text-sm font-medium text-primary">
            {t("total_users")}
          </span>
        </div>
        <div className="text-2xl font-bold text-base-content">{totalUsers}</div>
      </div>

      <div className="bg-gradient-to-br from-warning/10 to-warning/5 rounded-2xl p-4 text-center border border-warning/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <Trophy className="w-5 h-5 text-warning" />
          <span className="text-sm font-medium text-warning">
            {t("highest_score")}
          </span>
        </div>
        <div className="text-2xl font-bold text-base-content">
          {topScore?.toLocaleString() || 0}
        </div>
      </div>

      <div className="bg-gradient-to-br from-secondary/10 to-secondary/5 rounded-2xl p-4 text-center border border-secondary/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <TrendingUp className="w-5 h-5 text-secondary" />
          <span className="text-sm font-medium text-secondary">
            {t("active_score")}
          </span>
        </div>
        <div className="text-2xl font-bold text-base-content">TODO</div>
      </div>
    </div>
  );
}
