"use client";

import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { ProfileSidebar } from "./ProfileSidebar";
import { ProfileContent } from "./ProfileContent";
import { UserListModal } from "./UserListModal";
import { useUserProfile } from "../hooks/useUserProfile";
import { userApi } from "@/shared/api/modules/user";
import { UserDO } from "@/shared/api/types/user.model.do";
import { useFollowAction } from "../hooks/useFollowAction";
import { useFollowList } from "../hooks/useFollowList";

export default function UserProfileClient({ userId }: { userId: number }) {
  const { user: currentUser, isAuthenticated } = useAuthStore();
  const t = useTranslations("Profile");

  // 弹窗状态
  const [showFollowers, setShowFollowers] = useState(false);
  const [showFollowing, setShowFollowing] = useState(false);

  // 1. 用户资料
  const {
    profile: profile,
    loading: profileLoading,
    loadProfile,
  } = useUserProfile();

  // 2. 粉丝列表（独立实例）
  const {
    users: followersList,
    total: followersTotal,
    loading: followersLoading,
    loadFollowers,
  } = useFollowList();

  // 3. 关注列表（独立实例）
  const {
    users: followingList,
    total: followingTotal,
    loading: followingLoading,
    loadFollowing,
  } = useFollowList();

  // 4. 关注操作
  const { follow, unfollow, loading: followActionLoading } = useFollowAction();

  // 5. 当前用户是否已关注该用户（通过查询当前用户的关注列表判断）
  const { data: isFollowing, refetch: refetchFollowStatus } = useQuery({
    queryKey: ["check-following", userId, currentUser?.id],
    queryFn: async () => {
      if (!currentUser || currentUser.id === userId) return false;
      const res = await userApi.getFollowing(currentUser.id, {
        page: 1,
        page_size: 100,
      });
      const data = res.data.data;
      return data?.list?.some((u: UserDO) => u.id === userId) ?? false;
    },
    enabled: !!currentUser && currentUser.id !== userId,
  });

  // 加载用户资料
  useEffect(() => {
    if (userId) {
      loadProfile(userId);
    }
  }, [userId, loadProfile]);

  // 弹窗打开时加载对应列表
  useEffect(() => {
    if (showFollowers) {
      loadFollowers(userId, 1, 100);
    }
  }, [showFollowers, userId, loadFollowers]);

  useEffect(() => {
    if (showFollowing) {
      loadFollowing(userId, 1, 100);
    }
  }, [showFollowing, userId, loadFollowing]);

  // 关注/取消关注后的回调：刷新关注状态和粉丝/关注数量
  const handleFollowAction = async () => {
    if (isFollowing) {
      const success = await unfollow(userId);
      if (success) {
        refetchFollowStatus();
        // 刷新粉丝/关注列表（如果有弹窗打开）
        if (showFollowers) loadFollowers(userId, 1, 100);
        if (showFollowing) loadFollowing(userId, 1, 100);
      }
    } else {
      const success = await follow(userId);
      if (success) {
        refetchFollowStatus();
        if (showFollowers) loadFollowers(userId, 1, 100);
        if (showFollowing) loadFollowing(userId, 1, 100);
      }
    }
  };

  const handleUserClick = (clickedUserId: number) => {
    window.location.href = `/users/${clickedUserId}`;
  };

  if (profileLoading) {
    return (
      <div className="flex gap-6 max-w-6xl mx-auto">
        <div className="w-80 flex-shrink-0">
          <div className="skeleton h-96 w-full rounded-xl" />
        </div>
        <div className="flex-1">
          <div className="skeleton h-40 w-full rounded-xl mb-4" />
          <div className="skeleton h-20 w-full rounded-xl" />
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="text-center py-20 text-base-content/40">
        {t("user_not_found")}
      </div>
    );
  }

  const isSelf = currentUser?.id === userId;
  // 粉丝/关注总数可从 profile 中获取（若有），否则从列表中取
  const followerCount = profile.follower_count ?? followersTotal;
  const followingCount = profile.following_count ?? followingTotal;

  return (
    <div className="flex gap-6 max-w-6xl mx-auto">
      {/* 左侧边栏 */}
      <div className="w-80 flex-shrink-0 sticky top-20 h-fit">
        <ProfileSidebar
          profile={profile}
          isSelf={isSelf}
          isAuthenticated={isAuthenticated}
          isFollowing={!!isFollowing}
          followerCount={followerCount}
          followingCount={followingCount}
          onFollow={handleFollowAction}
          onShowFollowers={() => setShowFollowers(true)}
          onShowFollowing={() => setShowFollowing(true)}
          isFollowPending={followActionLoading}
        />
      </div>

      {/* 右侧内容区 */}
      <div className="flex-1 min-w-0">
        <ProfileContent userId={userId} isAuthenticated={isAuthenticated} />
      </div>

      {/* 粉丝列表弹窗 */}
      {showFollowers && (
        <UserListModal
          title={t("followers")}
          users={followersList}
          total={followersTotal}
          isLoading={followersLoading}
          onClose={() => setShowFollowers(false)}
          onUserClick={handleUserClick}
        />
      )}

      {/* 关注列表弹窗 */}
      {showFollowing && (
        <UserListModal
          title={t("following")}
          users={followingList}
          total={followingTotal}
          isLoading={followingLoading}
          onClose={() => setShowFollowing(false)}
          onUserClick={handleUserClick}
        />
      )}
    </div>
  );
}
