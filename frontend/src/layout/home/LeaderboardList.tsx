// src/components/home/LeaderboardList.tsx
"use client";

import Link from "next/link";
import { Trophy, ChevronRight } from "lucide-react";
import { useTranslations } from "next-intl";
import Avatar from "@/features/user/components/Avatar";

interface User {
  id: string | number;
  username: string;
  avatar?: string;
  score: number;
}

interface LeaderboardListProps {
  leaderboard: User[];
}

export default function LeaderboardList({ leaderboard }: LeaderboardListProps) {
  const t = useTranslations("post");

  // 安全检查
  if (!leaderboard || !Array.isArray(leaderboard) || leaderboard.length === 0) {
    return null;
  }

  // 过滤有效用户并限制数量
  const validUsers = leaderboard
    .filter((user) => user && user.id && user.username !== undefined)
    .slice(0, 8);

  if (validUsers.length === 0) {
    return null;
  }

  return (
    <div className="card bg-base-100 border border-base-300 shadow-sm">
      <div className="card-body p-4">
        <h3 className="font-bold flex items-center gap-2 mb-3">
          <Trophy className="w-4 h-4 text-warning" /> {t("leaderboard")}
        </h3>
        <div className="space-y-2">
          {validUsers.map((user, index) => (
            <LeaderboardItem
              key={`${user.id}-${index}`} // 组合 key 确保唯一性
              user={user}
              rank={index + 1}
            />
          ))}
        </div>
        <Link href="/leaderboard" className="btn btn-ghost btn-xs mt-2 gap-1">
          {t("view_the_full_rankings")} <ChevronRight className="w-3 h-3" />
        </Link>
      </div>
    </div>
  );
}

function LeaderboardItem({ user, rank }: { user: User; rank: number }) {
  const getRankStyles = () => {
    switch (rank) {
      case 1:
        return "bg-yellow-400 text-yellow-900";
      case 2:
        return "bg-gray-300 text-gray-700";
      case 3:
        return "bg-amber-600 text-white";
      default:
        return "text-base-content/40";
    }
  };

  return (
    <Link
      href={`/users/${user.id}`}
      className="flex items-center gap-2 hover:bg-base-200 rounded-lg p-1.5 transition-colors"
    >
      <span
        className={`w-5 h-5 text-xs font-bold flex items-center justify-center rounded-full ${getRankStyles()}`}
      >
        {rank}
      </span>

      <Avatar username={user.username} avatarUrl={user.avatar} size="md" />

      <span className="flex-1 text-sm truncate">{user.username}</span>
      <span className="text-xs text-warning font-medium">{user.score}</span>
    </Link>
  );
}
