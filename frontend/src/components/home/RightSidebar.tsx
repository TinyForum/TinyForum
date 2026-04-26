"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { GuestCard } from "./right/GuestCard";
import { Leaderboard } from "./right/Leaderboard";
import { TimelineEvents } from "./right/TimelineEvents";
import { UserProfileCard } from "./right/UserProfileCard";
import { LeaderboardItemResponse } from "@/lib/api/modules/users";


interface UserProfile {
  id: number;
  username: string;
  avatar: string;
  bio: string;
  score: number;
  created_at: string;
}

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
  user: any;
  userProfile?: UserProfile;
  leaderboard: LeaderboardItemResponse[];
  unreadCount: number;
  timelineEvents: TimelineEvent[];
}

export default function RightSidebar({
  isAuthenticated,
  user,
  userProfile,
  leaderboard,
  unreadCount,
  timelineEvents,
}: RightSidebarProps) {
  const t = useTranslations("Sidebar");

  return (
    <aside className="space-y-4">
      {/* 用户信息卡片 */}
      {isAuthenticated && userProfile ? (
        <UserProfileCard 
          userProfile={userProfile} 
          unreadCount={unreadCount} 
        />
      ) : (
        <GuestCard />
      )}

      {/* 违规情况（已登录时） */}
   

      {/* 时间线事件（已登录时） */}
      {isAuthenticated && timelineEvents && timelineEvents.length > 0 && (
        <TimelineEvents timelineEvents={timelineEvents} />
      )}

      {/* 用户排行榜 */}
      {leaderboard && (
        <Leaderboard leaderboard={leaderboard} />
      )}
    </aside>
  );
}