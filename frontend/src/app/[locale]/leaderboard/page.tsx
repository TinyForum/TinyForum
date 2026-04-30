"use client";

import { useTranslations } from "next-intl";
import type { LeaderboardItemResponse } from "@/shared/api/modules/users";
import { useLeaderboard } from "@/features/leader/hooks/useLeaderboard";
import { UserRankCard } from "@/shared/ui/rank/UserRankCard";
import { LoadingSkeleton } from "@/shared/ui/common/LoadingSkeleton";
import { EmptyState } from "@/shared/ui/common/EmptyState";
import { StatsCards } from "@/shared/ui/rank/StatsCards";
import { Trophy } from "lucide-react";

export default function LeaderboardPage() {
  const { data, isLoading, error } = useLeaderboard({
    limit: 50,
    //  fields: "id,username,avatar,score,bio", // FIXME: 用户、成员参与排名，管理员、版主、审核员不参与
  });
  console.log(data);

  const t = useTranslations("Leaderboard");
  const users: LeaderboardItemResponse[] = data ?? [];
  const totalUsers = users.length;
  const topScore = users[0]?.score || 0;

  if (error) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="bg-error/10 rounded-2xl p-8 text-center border border-error/20">
          <div className="text-5xl mb-4">⚠️</div>
          <h3 className="text-lg font-semibold text-error mb-2">
            {t("load_failed")}
          </h3>
          <p className="text-base-content/60 text-sm">{t("retry_later")}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-3xl mx-auto px-4 py-8 md:py-12">
        {/* 头部区域 */}
        <div className="text-center mb-10">
          <div className="relative inline-block mb-4">
            <div className="absolute inset-0 bg-gradient-to-r from-primary/20 to-warning/20 rounded-full blur-2xl" />
            <div className="relative bg-gradient-to-br from-primary to-warning p-4 rounded-full shadow-lg">
              <Trophy className="w-10 h-10 text-white" />
            </div>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-warning bg-clip-text text-transparent mb-2">
            {t("title")}
          </h1>
          <p className="text-base-content/60 text-sm">{t("description")}</p>
        </div>

        {/* 统计卡片 */}
        {!isLoading && users.length > 0 && (
          <StatsCards totalUsers={totalUsers} topScore={topScore} />
        )}

        {/* 排行榜列表 */}
        {isLoading ? (
          <LoadingSkeleton />
        ) : totalUsers === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-3">
            {/* 前三名特别展示 */}
            {users.slice(0, 3).map((user, index) => (
              <>
                <div className="m-2">
                  <UserRankCard key={user.id} user={user} rank={index} />
                </div>
              </>
            ))}

            {/* 分隔线 */}
            {users.length > 3 && (
              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-base-200"></div>
                </div>
                <div className="relative flex justify-center text-xs">
                  <span className="px-3 bg-base-200/50 text-base-content/40 rounded-full py-1">
                    {t("more_users")}
                  </span>
                </div>
              </div>
            )}

            {/* 其他用户 */}
            {users.slice(3).map((user, index) => (
              <UserRankCard key={user.id} user={user} rank={index + 3} />
            ))}

            {/* 底部提示 */}
            <div className="text-center text-xs text-base-content/40 mt-6 pt-4 border-t border-base-200">
              {t("showing_top_users", { count: users.length })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
