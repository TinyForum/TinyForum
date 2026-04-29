"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import {
  Settings,
  UserPlus,
  UserMinus,
  Star,
  Users,
  Calendar,
  Shield,
  Mail,
  MapPin,
  Link as LinkIcon,
} from "lucide-react";
import Avatar from "./Avatar";

interface ProfileSidebarProps {
  profile: {
    id: number;
    username: string;
    avatar: string;
    bio: string;
    score: number;
    created_at: string;
    role?: string;
    email?: string;
    location?: string;
    website?: string;
  };
  isSelf: boolean;
  isAuthenticated: boolean;
  isFollowing?: boolean;
  followerCount: number;
  followingCount: number;
  onFollow: () => void;
  onShowFollowers: () => void;
  onShowFollowing: () => void;
  isFollowPending: boolean;
}

export function ProfileSidebar({
  profile,
  isSelf,
  isAuthenticated,
  isFollowing,
  followerCount,
  followingCount,
  onFollow,
  onShowFollowers,
  onShowFollowing,
  isFollowPending,
}: ProfileSidebarProps) {
  const t = useTranslations("Profile");

  return (
    <div className="card bg-base-100 border border-base-300 shadow-sm">
      {/* 头像区域 */}
      <div className="relative">
        <div className="flex justify-center -mt-12 mb-4">
          <div className="avatar">
            <div className="w-28 h-28 rounded-2xl ring-4 ring-base-100 shadow-xl">
              <Avatar
                username={profile.username}
                avatarUrl={profile.avatar}
                shape="square"
              />
            </div>
          </div>
        </div>
      </div>

      {/* 用户信息 */}
      <div className="card-body pt-0">
        <div className="text-center">
          <h2 className="text-xl font-bold">{profile.username}</h2>
          {profile.role === "admin" && (
            <div className="badge badge-warning badge-sm mt-1">
              <Shield className="w-3 h-3 mr-1" />
              {t("administrator")}
            </div>
          )}
          {profile.bio && (
            <p className="text-sm text-base-content/60 mt-3">{profile.bio}</p>
          )}
        </div>

        {/* 统计数据 */}
        <div className="grid grid-cols-3 gap-2 py-4 border-y border-base-200 my-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-primary">
              {profile.score ?? 0}
            </div>
            <div className="text-xs text-base-content/40 flex items-center justify-center gap-1">
              <Star className="w-3 h-3" /> {t("score")}
            </div>
          </div>

          <button
            onClick={onShowFollowers}
            className="text-center hover:bg-base-200 rounded-lg py-2 transition-colors"
          >
            <div className="text-2xl font-bold">{followerCount}</div>
            <div className="text-xs text-base-content/40 flex items-center justify-center gap-1">
              <Users className="w-3 h-3" /> {t("followers")}
            </div>
          </button>

          <button
            onClick={onShowFollowing}
            className="text-center hover:bg-base-200 rounded-lg py-2 transition-colors"
          >
            <div className="text-2xl font-bold">{followingCount}</div>
            <div className="text-xs text-base-content/40 flex items-center justify-center gap-1">
              <Users className="w-3 h-3" /> {t("following")}
            </div>
          </button>
        </div>

        {/* 详细信息 */}
        <div className="space-y-2 text-sm">
          <div className="flex items-center gap-2 text-base-content/60">
            <Calendar className="w-4 h-4 flex-shrink-0" />
            <span>
              {t("joined")} {new Date(profile.created_at).toLocaleDateString()}
            </span>
          </div>

          {profile.location && (
            <div className="flex items-center gap-2 text-base-content/60">
              <MapPin className="w-4 h-4 flex-shrink-0" />
              <span>{profile.location}</span>
            </div>
          )}

          {profile.website && (
            <div className="flex items-center gap-2 text-base-content/60">
              <LinkIcon className="w-4 h-4 flex-shrink-0" />
              <a
                href={profile.website}
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-primary truncate"
              >
                {profile.website}
              </a>
            </div>
          )}

          {profile.email && isSelf && (
            <div className="flex items-center gap-2 text-base-content/60">
              <Mail className="w-4 h-4 flex-shrink-0" />
              <span className="truncate">{profile.email}</span>
            </div>
          )}
        </div>

        {/* 操作按钮 */}
        <div className="mt-4">
          {isSelf ? (
            <Link
              href="/settings"
              className="btn btn-outline btn-sm w-full gap-2"
            >
              <Settings className="w-4 h-4" /> {t("edit_profile")}
            </Link>
          ) : isAuthenticated ? (
            <button
              className={`btn btn-sm w-full gap-2 ${isFollowing ? "btn-ghost" : "btn-primary"}`}
              onClick={onFollow}
              disabled={isFollowPending}
            >
              {isFollowing ? (
                <>
                  <UserMinus className="w-4 h-4" /> {t("unfollow")}
                </>
              ) : (
                <>
                  <UserPlus className="w-4 h-4" /> {t("follow")}
                </>
              )}
            </button>
          ) : (
            <Link
              href="/auth/login"
              className="btn btn-outline btn-sm w-full gap-2"
            >
              {t("login_to_follow")}
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
