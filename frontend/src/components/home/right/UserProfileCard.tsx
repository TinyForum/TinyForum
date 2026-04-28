"use client";

import Link from "next/link";
import Image from "next/image";
import { useTranslations } from "next-intl";
import { User, Bell, Activity, Award, Clock } from "lucide-react";

interface UserProfile {
  id: number;
  username: string;
  avatar: string;
  bio: string;
  score: number;
  created_at: string;
}

interface UserProfileCardProps {
  userProfile: UserProfile;
  unreadCount: number;
}

export function UserProfileCard({
  userProfile,
  unreadCount,
}: UserProfileCardProps) {
  const t = useTranslations("Sidebar");

  return (
    <div className="rounded-lg border bg-card shadow-sm hover:shadow-md transition-shadow duration-200">
      <div className="p-4 text-center border-b">
        <div className="relative inline-block">
          <div className="relative w-16 h-16 mx-auto mb-2">
            <Image
              src={userProfile.avatar || "/default-avatar.png"}
              alt={userProfile.username}
              fill
              className="rounded-full object-cover"
              sizes="64px"
              priority={false}
            />
          </div>
          {unreadCount > 0 && (
            <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center font-medium shadow-sm">
              {unreadCount > 99 ? "99+" : unreadCount}
            </span>
          )}
        </div>
        <h3 className="font-semibold text-base">{userProfile.username}</h3>
        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
          {userProfile.bio || t("no_bio")}
        </p>
      </div>

      <div className="p-3 space-y-2">
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground flex items-center gap-1">
            <Award className="w-4 h-4" />
            {t("score")}
          </span>
          <span className="font-medium text-primary">{userProfile.score}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground flex items-center gap-1">
            <Clock className="w-4 h-4" />
            {t("registered_at")}
          </span>
          <span className="text-xs">
            {new Date(userProfile.created_at).toLocaleDateString()}
          </span>
        </div>
      </div>

      <div className="p-3 border-t space-y-1">
        <Link
          href="/notifications"
          className="flex items-center justify-between p-2 rounded-lg hover:bg-muted transition-colors"
        >
          <span className="flex items-center gap-2 text-sm">
            <Bell className="w-4 h-4" />
            {t("notifications")}
          </span>
          {unreadCount > 0 && (
            <span className="bg-red-500 text-white text-xs px-2 py-0.5 rounded-full font-medium">
              {unreadCount > 99 ? "99+" : unreadCount}
            </span>
          )}
        </Link>
        <Link
          href={`/users/${userProfile.id}`}
          className="flex items-center gap-2 p-2 rounded-lg hover:bg-muted transition-colors text-sm"
        >
          <User className="w-4 h-4" />
          {t("profile")}
        </Link>
        <Link
          href="/settings"
          className="flex items-center gap-2 p-2 rounded-lg hover:bg-muted transition-colors text-sm"
        >
          <Activity className="w-4 h-4" />
          {t("settings")}
        </Link>
      </div>
    </div>
  );
}
