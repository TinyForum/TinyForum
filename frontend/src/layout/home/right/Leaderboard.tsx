"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { TrendingUp } from "lucide-react";
import { LeaderboardItemResponse } from "@/shared/api/modules/users";
import Avatar from "@/features/user/components/Avatar";

interface LeaderboardProps {
  leaderboard: LeaderboardItemResponse[];
}

export function Leaderboard({ leaderboard }: LeaderboardProps) {
  const t = useTranslations("Sidebar");

  const getRankIcon = (index: number): React.ReactNode => {
    switch (index) {
      case 0:
        return <span className="text-yellow-500 font-bold text-lg">🥇</span>;
      case 1:
        return <span className="text-gray-400 font-bold text-lg">🥈</span>;
      case 2:
        return <span className="text-amber-600 font-bold text-lg">🥉</span>;
      default:
        return (
          <span className="text-muted-foreground font-medium text-sm">
            {index + 1}
          </span>
        );
    }
  };

  if (!leaderboard || leaderboard.length === 0) {
    return null;
  }

  return (
    <div className="rounded-lg border bg-card shadow-sm hover:shadow-md transition-shadow duration-200">
      <div className="p-3 border-b bg-gradient-to-r from-primary/5 to-transparent">
        <h3 className="font-semibold flex items-center gap-2 text-base">
          <TrendingUp className="w-4 h-4 text-primary" />
          {t("leaderboard")}
        </h3>
      </div>
      <div className="p-2 space-y-1">
        {leaderboard.slice(0, 5).map((user, index) => (
          <Link
            key={user.id}
            href={`/users/${user.id}`}
            className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted transition-all duration-200 hover:translate-x-0.5"
          >
            <div className="flex-shrink-0 w-7 text-center">
              {getRankIcon(index)}
            </div>
            <div className="relative w-7 h-7 flex-shrink-0">
              <Avatar username={user.username} avatarUrl={user.avatar} />
            </div>
            <span className="flex-1 text-sm truncate font-medium">
              {user.username}
            </span>
            <span className="text-xs font-bold text-primary bg-primary/10 px-2 py-0.5 rounded-full">
              {user.score}
            </span>
          </Link>
        ))}
      </div>
      <div className="p-2 border-t bg-muted/30">
        <Link
          href="/leaderboard"
          className="block text-xs text-center text-muted-foreground hover:text-primary transition-colors"
        >
          {t("view_the_full_rankings")} →
        </Link>
      </div>
    </div>
  );
}
