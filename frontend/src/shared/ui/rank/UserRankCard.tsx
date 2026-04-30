import Avatar from "@/features/user/components/Avatar";
import { LeaderboardItemResponse } from "@/shared/api/modules/users";
import { Crown, Star } from "lucide-react";
import { useTranslations } from "next-intl";
import { RankBadge } from "./RankBadge";
import Link from "next/link";

// 用户卡片组件
export function UserRankCard({
  user,
  rank,
}: {
  user: LeaderboardItemResponse;
  rank: number;
}) {
  const t = useTranslations("Leaderboard");
  const isTopThree = rank < 3;
  const cardStyles =
    [
      "border-yellow-400/50 bg-gradient-to-r from-yellow-50/50 to-transparent dark:from-yellow-900/10",
      "border-gray-300/50 bg-gradient-to-r from-gray-50/50 to-transparent dark:from-gray-900/10",
      "border-amber-500/50 bg-gradient-to-r from-amber-50/50 to-transparent dark:from-amber-900/10",
    ][rank] || "border-base-200 bg-base-100";

  const scoreColor = isTopThree ? "text-warning" : "text-base-content/40";

  console.log(user);
  return (
    <Link href={`/users/${user.id}`}>
      <div
        className={`card border shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 ${cardStyles}`}
      >
        <div className="card-body p-4">
          <div className="flex items-center gap-4">
            {/* 排名 */}
            <RankBadge rank={rank + 1} isTopThree={isTopThree} />

            {/* 头像 */}
            <div className="avatar">
              <div className="w-11 h-11 rounded-full ring-2 ring-primary/20 ring-offset-2">
                <Avatar
                  // username={user.username}
                  avatarUrl={user.avatar}
                  size="md"
                />
              </div>
            </div>

            {/* 用户信息 */}
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-base-content truncate">
                  {user.username}
                </span>
                {rank === 0 && (
                  <span className="badge badge-warning badge-sm gap-1">
                    <Crown className="w-3 h-3" />
                    {t("champion")}
                  </span>
                )}
              </div>
              {user.bio && (
                <div className="text-xs text-base-content/40 truncate max-w-xs">
                  {user.bio}
                </div>
              )}
            </div>

            {/* 积分 */}
            <div className="flex items-center gap-1.5">
              <Star
                className={`w-5 h-5 ${isTopThree ? "text-warning" : "text-base-content/30"}`}
              />
              <span className={`font-bold text-lg ${scoreColor}`}>
                {user.score?.toLocaleString() || 0}
              </span>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}
