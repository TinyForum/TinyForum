// src/app/users/[id]/UserProfileClient.tsx
'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userApi, postApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import PostCard from '@/components/post/PostCard';
import Image from 'next/image';
import { formatDate } from '@/lib/utils';
import { Users, FileText, Star, UserPlus, UserMinus, Settings } from 'lucide-react';
import toast from 'react-hot-toast';
import Link from 'next/link';
import Avatar from '@/components/user/Avatar';

export default function UserProfileClient({ userId }: { userId: number }) {
  const { user: currentUser, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<'posts' | 'articles'>('posts');

  const { data: profile, isLoading } = useQuery({
    queryKey: ['user', userId],
    queryFn: () => userApi.getProfile(userId).then((r) => r.data.data),
  });

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

  const followMutation = useMutation({
    mutationFn: () =>
      profile?.is_following ? userApi.unfollow(userId) : userApi.follow(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', userId] });
      toast.success(profile?.is_following ? '已取消关注' : '关注成功');
    },
    onError: () => toast.error('操作失败'),
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
    return <div className="text-center py-20 text-base-content/40">用户不存在</div>;
  }

  const isSelf = currentUser?.id === userId;
  const posts = postsData?.list ?? [];

  return (
    <div className="max-w-4xl mx-auto">
      {/* Profile header */}
      <div className="card bg-base-100 border border-base-300 shadow-sm mb-6">
        <div className="card-body p-6">
          <div className="flex flex-col sm:flex-row items-center sm:items-start gap-5">
            {/* Avatar */}
            <div className="avatar">
              <div className="w-24 h-24 rounded-2xl ring ring-primary ring-offset-base-100 ring-offset-4">
            

                <Avatar 
  username={profile.username} 
  avatarUrl={profile.avatar}  // 数据库中的头像
  size="md" 
/>

              </div>
            </div>

            {/* Info */}
            <div className="flex-1 text-center sm:text-left">
              <div className="flex items-center justify-center sm:justify-start gap-3 flex-wrap">
                <h1 className="text-2xl font-bold">{profile.username}</h1>
                {profile.role === 'admin' && (
                  <span className="badge badge-warning badge-sm">管理员</span>
                )}
              </div>
              {profile.bio && (
                <p className="text-base-content/60 mt-2 text-sm max-w-lg">{profile.bio}</p>
              )}
              <div className="text-xs text-base-content/40 mt-2">
                加入于 {formatDate(profile.created_at)}
              </div>

              {/* Stats */}
              <div className="flex items-center justify-center sm:justify-start gap-6 mt-4">
                <div className="text-center">
                  <div className="text-xl font-bold text-primary">{profile.score ?? 0}</div>
                  <div className="text-xs text-base-content/40 flex items-center gap-1"><Star className="w-3 h-3" /> 积分</div>
                </div>
                <div className="text-center">
                  <div className="text-xl font-bold">{profile.follower_count ?? 0}</div>
                  <div className="text-xs text-base-content/40 flex items-center gap-1"><Users className="w-3 h-3" /> 粉丝</div>
                </div>
                <div className="text-center">
                  <div className="text-xl font-bold">{profile.following_count ?? 0}</div>
                  <div className="text-xs text-base-content/40 flex items-center gap-1"><Users className="w-3 h-3" /> 关注</div>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-2">
              {isSelf ? (
                <Link href="/settings" className="btn btn-outline btn-sm gap-1">
                  <Settings className="w-4 h-4" /> 编辑资料
                </Link>
              ) : isAuthenticated ? (
                <button
                  className={`btn btn-sm gap-1 ${profile.is_following ? 'btn-ghost' : 'btn-primary'}`}
                  onClick={() => followMutation.mutate()}
                  disabled={followMutation.isPending}
                >
                  {profile.is_following ? (
                    <><UserMinus className="w-4 h-4" /> 取消关注</>
                  ) : (
                    <><UserPlus className="w-4 h-4" /> 关注</>
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
          <FileText className="w-4 h-4" /> 帖子
        </button>
        <button
          className={`tab gap-2 ${tab === 'articles' ? 'tab-active' : ''}`}
          onClick={() => setTab('articles')}
        >
          <FileText className="w-4 h-4" /> 文章
        </button>
      </div>

      {posts.length === 0 ? (
        <div className="text-center py-16 text-base-content/40">
          <FileText className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>还没有{tab === 'articles' ? '文章' : '帖子'}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}
    </div>
  );
}