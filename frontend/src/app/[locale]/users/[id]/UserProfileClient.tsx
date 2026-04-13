// src/app/users/[id]/UserProfileClient.tsx
'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userApi, postApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import PostCard from '@/components/post/PostCard';
import { formatDate } from '@/lib/utils';
import { Users, FileText, Star, UserPlus, UserMinus, Settings, X } from 'lucide-react';
import toast from 'react-hot-toast';
import Link from 'next/link';
import Avatar from '@/components/user/Avatar';
import { useTranslations } from 'next-intl';
import type { User } from '@/lib/api/types';

// 用户列表弹窗组件
function UserListModal({ 
  title, 
  users, 
  total, 
  onClose,
  onUserClick,
  isLoading,
}: { 
  title: string; 
  users: User[]; 
  total: number;
  onClose: () => void;
  onUserClick: (userId: number) => void;
  isLoading?: boolean;
}) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-base-100 rounded-2xl shadow-xl max-w-md w-full max-h-[80vh] flex flex-col" onClick={(e) => e.stopPropagation()}>
        <div className="flex items-center justify-between p-4 border-b border-base-200">
          <h3 className="text-lg font-semibold">
            {title} ({total})
          </h3>
          <button onClick={onClose} className="btn btn-sm btn-ghost btn-circle">
            <X className="w-4 h-4" />
          </button>
        </div>
        
        <div className="flex-1 overflow-y-auto p-2">
          {isLoading ? (
            <div className="space-y-2 p-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="flex items-center gap-3">
                  <div className="skeleton w-10 h-10 rounded-full" />
                  <div className="flex-1">
                    <div className="skeleton h-4 w-24 mb-1" />
                    <div className="skeleton h-3 w-32" />
                  </div>
                </div>
              ))}
            </div>
          ) : users.length === 0 ? (
            <div className="text-center py-8 text-base-content/40">暂无数据</div>
          ) : (
            <div className="space-y-1">
              {users.map((user) => (
                <button
                  key={user.id}
                  onClick={() => {
                    onUserClick(user.id);
                    onClose();
                  }}
                  className="w-full flex items-center gap-3 p-3 rounded-xl hover:bg-base-200 transition-colors text-left"
                >
                  <div className="avatar">
                    <div className="w-10 h-10 rounded-full">
                      <Avatar username={user.username} avatarUrl={user.avatar} size="md" />
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm truncate">{user.username}</p>
                    {user.bio && <p className="text-xs text-base-content/40 truncate">{user.bio}</p>}
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default function UserProfileClient({ userId }: { userId: number }) {
  const { user: currentUser, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<'posts' | 'articles'>('posts');
  const [showFollowers, setShowFollowers] = useState(false);
  const [showFollowing, setShowFollowing] = useState(false);

  const t = useTranslations("Profile");
  
  // 获取用户资料
  const { data: profile, isLoading } = useQuery({
    queryKey: ['user', userId],
    queryFn: () => userApi.getProfile(userId).then((r) => r.data.data),
  });

  // 获取粉丝列表（用于获取粉丝数）
  const { data: followersData } = useQuery({
    queryKey: ['user-followers-count', userId],
    queryFn: () => userApi.follwowers(userId, { page: 1, page_size: 1 }).then((r) => r.data.data),
  });

  // 获取关注列表（用于获取关注数）
  const { data: followingData } = useQuery({
    queryKey: ['user-following-count', userId],
    queryFn: () => userApi.following(userId, { page: 1, page_size: 1 }).then((r) => r.data.data),
  });

  // 获取用户的帖子/文章
  const { data: postsData } = useQuery({
    queryKey: ['user-posts', userId, tab],
    queryFn: () =>
      postApi.list({
        author_id: userId,
        type: tab === 'articles' ? 'article' : 'post',
        page: 1,
        page_size: 20,
      }).then((r) => r.data.data),
  });

  // 获取粉丝详细列表（弹窗用）
  const { data: followersDetail, isLoading: followersLoading } = useQuery({
    queryKey: ['user-followers-detail', userId],
    queryFn: () => userApi.follwowers(userId, { page: 1, page_size: 100 }).then((r) => r.data.data),
    enabled: showFollowers,
  });

  // 获取关注详细列表（弹窗用）
  const { data: followingDetail, isLoading: followingLoading } = useQuery({
    queryKey: ['user-following-detail', userId],
    queryFn: () => userApi.following(userId, { page: 1, page_size: 100 }).then((r) => r.data.data),
    enabled: showFollowing,
  });

  // 检查当前用户是否关注了该用户
  const { data: isFollowing } = useQuery({
    queryKey: ['check-following', userId, currentUser?.id],
    queryFn: async () => {
      if (!currentUser || currentUser.id === userId) return false;
      const res = await userApi.following(currentUser.id, { page: 1, page_size: 100 });
      return res.data.data.list?.some(u => u.id === userId) ?? false;
    },
    enabled: !!currentUser && currentUser.id !== userId,
  });

  // 关注/取消关注 mutation
  const followMutation = useMutation({
    mutationFn: () => userApi.follow(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-followers-count', userId] });
      queryClient.invalidateQueries({ queryKey: ['user-following-count', userId] });
      queryClient.invalidateQueries({ queryKey: ['check-following', userId] });
      queryClient.invalidateQueries({ queryKey: ['user-followers-detail', userId] });
      queryClient.invalidateQueries({ queryKey: ['user-following-detail', userId] });
      toast.success(t("followed"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  const unfollowMutation = useMutation({
    mutationFn: () => userApi.unfollow(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-followers-count', userId] });
      queryClient.invalidateQueries({ queryKey: ['user-following-count', userId] });
      queryClient.invalidateQueries({ queryKey: ['check-following', userId] });
      queryClient.invalidateQueries({ queryKey: ['user-followers-detail', userId] });
      queryClient.invalidateQueries({ queryKey: ['user-following-detail', userId] });
      toast.success(t("unfollowed"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto space-y-4">
        <div className="skeleton h-40 w-full rounded-xl" />
        <div className="skeleton h-20 w-full rounded-xl" />
      </div>
    );
  }

  if (!profile) {
    return <div className="text-center py-20 text-base-content/40">{t("user_not_found")}</div>;
  }

  const isSelf = currentUser?.id === userId;
  const posts = postsData?.list ?? [];
  const followerCount = followersData?.total ?? 0;
  const followingCount = followingData?.total ?? 0;

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

  return (
    <div className="max-w-4xl mx-auto">
      {/* Profile header */}
      <div className="card bg-base-100 border border-base-300 shadow-sm mb-6">
        <div className="card-body p-6">
          <div className="flex flex-col sm:flex-row items-center sm:items-start gap-5">
            {/* Avatar */}
            <div className="avatar">
              <div className="w-24 h-24 rounded-2xl ring ring-primary ring-offset-base-100 ring-offset-4">
                <Avatar username={profile.username} avatarUrl={profile.avatar} shape="square" />
              </div>
            </div>

            {/* Info */}
            <div className="flex-1 text-center sm:text-left">
              <div className="flex items-center justify-center sm:justify-start gap-3 flex-wrap">
                <h1 className="text-2xl font-bold">{profile.username}</h1>
                {profile.role === 'admin' && (
                  <span className="badge badge-warning badge-sm">{t("administrator")}</span>
                )}
              </div>
              {profile.bio && (
                <p className="text-base-content/60 mt-2 text-sm max-w-lg">{profile.bio}</p>
              )}
              <div className="text-xs text-base-content/40 mt-2">
                {t("joined")} {formatDate(profile.created_at)}
              </div>

              {/* Stats */}
              <div className="flex items-center justify-center sm:justify-start gap-6 mt-4">
                <div className="text-center">
                  <div className="text-xl font-bold text-primary">{profile.score ?? 0}</div>
                  <div className="text-xs text-base-content/40 flex items-center gap-1">
                    <Star className="w-3 h-3" /> {t("score")}
                  </div>
                </div>
                
                <button 
                  onClick={() => setShowFollowers(true)}
                  className="text-center hover:bg-base-200 rounded-lg px-3 py-1 transition-colors"
                >
                  <div className="text-xl font-bold">{followerCount}</div>
                  <div className="text-xs text-base-content/40 flex items-center gap-1">
                    <Users className="w-3 h-3" /> {t("followers")}
                  </div>
                </button>
                
                <button 
                  onClick={() => setShowFollowing(true)}
                  className="text-center hover:bg-base-200 rounded-lg px-3 py-1 transition-colors"
                >
                  <div className="text-xl font-bold">{followingCount}</div>
                  <div className="text-xs text-base-content/40 flex items-center gap-1">
                    <Users className="w-3 h-3" /> {t("following")}
                  </div>
                </button>
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-2">
              {isSelf ? (
                <Link href="/settings" className="btn btn-outline btn-sm gap-1">
                  <Settings className="w-4 h-4" /> {t("edit_profile")}
                </Link>
              ) : isAuthenticated ? (
                <button
                  className={`btn btn-sm gap-1 ${isFollowing ? 'btn-ghost' : 'btn-primary'}`}
                  onClick={handleFollowAction}
                  disabled={followMutation.isPending || unfollowMutation.isPending}
                >
                  {isFollowing ? (
                    <><UserMinus className="w-4 h-4" /> {t("unfollow")}</>
                  ) : (
                    <><UserPlus className="w-4 h-4" /> {t("follow")}</>
                  )}
                </button>
              ) : null}
            </div>
          </div>
        </div>
      </div>

      {/* Posts tabs */}
      <div className="tabs tabs-boxed bg-base-100 border border-base-300 mb-4 p-1">
        <button
          className={`tab gap-2 ${tab === 'posts' ? 'tab-active' : ''}`}
          onClick={() => setTab('posts')}
        >
          <FileText className="w-4 h-4" /> {t("the_posts")}
        </button>
        <button
          className={`tab gap-2 ${tab === 'articles' ? 'tab-active' : ''}`}
          onClick={() => setTab('articles')}
        >
          <FileText className="w-4 h-4" /> {t("the_articles")}
        </button>
      </div>

      {posts.length === 0 ? (
        <div className="text-center py-16 text-base-content/40">
          <FileText className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>{tab === 'articles' ? t("no_articles") : t("no_posts")}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}

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