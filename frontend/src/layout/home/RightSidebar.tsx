"use client";

import { LeaderboardItemResponse } from "@/shared/api/modules/user";
import { GuestCard } from "./right/GuestCard";
import { Leaderboard } from "./right/Leaderboard";
import { TimelineEvents } from "./right/TimelineEvents";
import { UserProfileCard } from "./right/UserProfileCard";
import { UserDO } from "@/shared/api/types/user.model";

interface TimelineEvent {
  id: number;
  action: string;
  target_type: string;
  target_id: number;
  created_at: string;
  actor?: {
    id: number;
    username: string;
    avatar: string;
  };
}

interface RightSidebarProps {
  isAuthenticated: boolean;
  userProfile: UserDO | null;
  leaderboard: LeaderboardItemResponse[];
  unreadCount: number;
  timelineEvents: TimelineEvent[];
}

export default function RightSidebar({
  isAuthenticated,
  userProfile,
  leaderboard,
  unreadCount,
  timelineEvents,
}: RightSidebarProps) {
  console.log("RightSidebar", {
    isAuthenticated,
    userProfile,
    leaderboard,
    unreadCount,
    timelineEvents,
  });
  return (
    <aside className="space-y-4">
      {/* 用户信息卡片 */}
      {isAuthenticated && userProfile ? (
        <UserProfileCard userProfile={userProfile} unreadCount={unreadCount} />
      ) : (
        <GuestCard />
      )}

      {/* 时间线事件（已登录时） */}
      {isAuthenticated && timelineEvents && timelineEvents.length > 0 && (
        <TimelineEvents timelineEvents={timelineEvents} />
      )}

      {/* 用户排行榜 */}
      {leaderboard && leaderboard.length > 0 && (
        <Leaderboard leaderboard={leaderboard} />
      )}
    </aside>
  );
}
