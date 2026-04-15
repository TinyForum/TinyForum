"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { TrendingUp } from "lucide-react";

interface LeaderboardUser {
  id: number;
  username: string;
  avatar: string;
  score: number;
}

interface LeaderboardProps {
  leaderboard: LeaderboardUser[];
}

export function Leaderboard({ leaderboard }: LeaderboardProps) {
  const t = useTranslations("Sidebar");

  const getRankIcon = (index: number) => {
    switch (index) {
      case 0:
        return <span className="text-yellow-500 font-bold">🥇</span>;
      case 1:
        return <span className="text-gray-400 font-bold">🥈</span>;
      case 2:
        return <span className="text-amber-600 font-bold">🥉</span>;
      default:
        return <span className="text-muted-foreground">{index + 1}</span>;
    }
  };

  return (
    <div className="rounded-lg border bg-card">
      <div className="p-3 border-b">
        <h3 className="font-semibold flex items-center gap-2">
          <TrendingUp className="w-4 h-4" />
          {t("leaderboard")}
        </h3>
      </div>
      <div className="p-2 space-y-1">
        {leaderboard.slice(0, 5).map((user, index) => (
          <Link
            key={user.id}
            href={`/profile/${user.id}`}
            className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted transition-colors"
          >
            <div className="flex-shrink-0 w-6 text-center">
              {getRankIcon(index)}
            </div>
            <img
              src={user.avatar || "/default-avatar.png"}
              alt={user.username}
              className="w-6 h-6 rounded-full object-cover"
            />
            <span className="flex-1 text-sm truncate">{user.username}</span>
            <span className="text-xs font-medium text-primary">{user.score}</span>
          </Link>
        ))}
      </div>
      <div className="p-2 border-t">
        <Link
          href="/leaderboard"
          className="block text-xs text-center text-muted-foreground hover:text-primary"
        >
          {t("view_the_full_rankings")} →
        </Link>
      </div>
    </div>
  );
}