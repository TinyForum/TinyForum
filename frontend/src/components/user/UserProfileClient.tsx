"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userApi, postApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { formatDate } from "@/lib/utils";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import type { PostType, User } from "@/lib/api/types";
import { ProfileSidebar } from "./ProfileSidebar";
import { ProfileContent } from "./ProfileContent";
import { UserListModal } from "./UserListModal";

export default function UserProfileClient({ userId }: { userId: number }) {
  const { user: currentUser, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
  const [showFollowers, setShowFollowers] = useState(false);
  const [showFollowing, setShowFollowing] = useState(false);
  const t = useTranslations("Profile");

  // 获取用户资料
  const { data: profile, isLoading } = useQuery({
    queryKey: ["user", userId],
    queryFn: () => userApi.getProfile(userId).then((r) => r.data.data),
  });

  // 获取粉丝列表（用于获取粉丝数）
  const { data: followersData } = useQuery({
    queryKey: ["user-followers-count", userId],
    queryFn: () =>
      userApi
        .follwowers(userId, { page: 1, page_size: 1 })
        .then((r) => r.data.data),
  });

  // 获取关注列表（用于获取关注数）
  const { data: followingData } = useQuery({
    queryKey: ["user-following-count", userId],
    queryFn: () =>
      userApi
        .following(userId, { page: 1, page_size: 1 })
        .then((r) => r.data.data),
  });

  // 获取粉丝详细列表（弹窗用）
  const { data: followersDetail, isLoading: followersLoading } = useQuery({
    queryKey: ["user-followers-detail", userId],
    queryFn: () =>
      userApi
        .follwowers(userId, { page: 1, page_size: 100 })
        .then((r) => r.data.data),
    enabled: showFollowers,
  });

  // 获取关注详细列表（弹窗用）
  const { data: followingDetail, isLoading: followingLoading } = useQuery({
    queryKey: ["user-following-detail", userId],
    queryFn: () =>
      userApi
        .following(userId, { page: 1, page_size: 100 })
        .then((r) => r.data.data),
    enabled: showFollowing,
  });

  // 检查当前用户是否关注了该用户
  const { data: isFollowing } = useQuery({
    queryKey: ["check-following", userId, currentUser?.id],
    queryFn: async () => {
      if (!currentUser || currentUser.id === userId) return false;
      const res = await userApi.following(currentUser.id, {
        page: 1,
        page_size: 100,
      });
      return res.data.data.list?.some((u) => u.id === userId) ?? false;
    },
    enabled: !!currentUser && currentUser.id !== userId,
  });

  // 关注/取消关注 mutation
  const followMutation = useMutation({
    mutationFn: () => userApi.follow(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["user-followers-count", userId],
      });
      queryClient.invalidateQueries({
        queryKey: ["user-following-count", userId],
      });
      queryClient.invalidateQueries({ queryKey: ["check-following", userId] });
      queryClient.invalidateQueries({
        queryKey: ["user-followers-detail", userId],
      });
      queryClient.invalidateQueries({
        queryKey: ["user-following-detail", userId],
      });
      toast.success(t("followed"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  const unfollowMutation = useMutation({
    mutationFn: () => userApi.unfollow(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["user-followers-count", userId],
      });
      queryClient.invalidateQueries({
        queryKey: ["user-following-count", userId],
      });
      queryClient.invalidateQueries({ queryKey: ["check-following", userId] });
      queryClient.invalidateQueries({
        queryKey: ["user-followers-detail", userId],
      });
      queryClient.invalidateQueries({
        queryKey: ["user-following-detail", userId],
      });
      toast.success(t("unfollowed"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  const handleFollowAction = () => {
    if (isFollowing) {
      unfollowMutation.mutate();
    } else {
      followMutation.mutate();
    }
  };

  const handleUserClick = (clickedUserId: number) => {
    window.location.href = `/users/${clickedUserId}`;
  };

  if (isLoading) {
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
  const followerCount = followersData?.total ?? 0;
  const followingCount = followingData?.total ?? 0;

  return (
    <div className="flex gap-6 max-w-6xl mx-auto">
      {/* 左侧边栏 */}
      <div className="w-80 flex-shrink-0 sticky top-20 h-fit">
        <ProfileSidebar
          profile={profile}
          isSelf={isSelf}
          isAuthenticated={isAuthenticated}
          isFollowing={isFollowing}
          followerCount={followerCount}
          followingCount={followingCount}
          onFollow={handleFollowAction}
          onShowFollowers={() => setShowFollowers(true)}
          onShowFollowing={() => setShowFollowing(true)}
          isFollowPending={
            followMutation.isPending || unfollowMutation.isPending
          }
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
          users={followersDetail?.list || []}
          total={followersDetail?.total || 0}
          isLoading={followersLoading}
          onClose={() => setShowFollowers(false)}
          onUserClick={handleUserClick}
        />
      )}

      {/* 关注列表弹窗 */}
      {showFollowing && (
        <UserListModal
          title={t("following")}
          users={followingDetail?.list || []}
          total={followingDetail?.total || 0}
          isLoading={followingLoading}
          onClose={() => setShowFollowing(false)}
          onUserClick={handleUserClick}
        />
      )}
    </div>
  );
}
