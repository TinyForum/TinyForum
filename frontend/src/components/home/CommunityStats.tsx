// components/home/CommunityStats.tsx
import { useStatistics } from "@/hooks/useStatistics";
import { Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";
// import { useTranslation } from "react-i18next";

interface CommunityStatsProps {
  className?: string;
}

export const CommunityStats = ({ className = "" }: CommunityStatsProps) => {
  const t  = useTranslations("Sidebar");
  const { dayStats, totalStats, isLoading } = useStatistics({
    autoFetch: true,
  });

  // 骨架屏组件
  const SkeletonItem = () => (
    <div className="flex justify-between">
      <span className="text-muted-foreground">——</span>
      <span className="font-medium text-muted-foreground/50">——</span>
    </div>
  );

  // 统计项组件
  const StatItem = ({ label, value }: { label: string; value: number }) => (
    <div className="flex justify-between">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium">{value.toLocaleString()}</span>
    </div>
  );

  // 根据实际类型定义获取数据
  // totalStats 是 StatsInfoResp 类型
  const totalUsers = totalStats?.base_info?.total_user || 0;
  const totalPosts = totalStats?.base_info?.total_article || 0;
  
  // dayStats 是 StatsTodayInfo 类型
  const todayPosts = dayStats?.new_article || 0;
  const yesterdayActive = dayStats?.active_user || 0;

  return (
    <div className={`rounded-lg border bg-card ${className}`}>
      <div className="p-3 border-b">
        <h3 className="font-semibold flex items-center gap-2">
          <Sparkles className="w-4 h-4" />
          {t("community_stats")}
        </h3>
      </div>
      <div className="p-3 space-y-2 text-sm">
        {isLoading ? (
          <>
            <SkeletonItem />
            <SkeletonItem />
            <SkeletonItem />
            <SkeletonItem />
          </>
        ) : (
          <>
            <StatItem 
              label={t("today_posts")} 
              value={todayPosts} 
            />
            <StatItem 
              label={t("yesterday_active")} 
              value={yesterdayActive} 
            />
            <StatItem 
              label={t("total_users")} 
              value={totalUsers} 
            />
            <StatItem 
              label={t("total_posts")} 
              value={totalPosts} 
            />
          </>
        )}
      </div>
    </div>
  );
};