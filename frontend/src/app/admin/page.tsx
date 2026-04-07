'use client';

import { useEffect, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import { formatDate } from '@/lib/utils';
import {
  LayoutDashboard, Users, FileText, Pin, PinOff,
  ShieldCheck, ShieldOff, Search,
} from 'lucide-react';
import toast from 'react-hot-toast';

export default function AdminPage() {
  const { user, isAuthenticated } = useAuthStore();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<'users' | 'posts'>('users');
  const [keyword, setKeyword] = useState('');
  const [page, setPage] = useState(1);

  useEffect(() => {
    if (!isAuthenticated || user?.role !== 'admin') {
      router.push('/');
    }
  }, [isAuthenticated, user, router]);

  const { data: usersData, isLoading: usersLoading } = useQuery({
    queryKey: ['admin-users', page, keyword],
    queryFn: () => adminApi.listUsers({ page, page_size: 20, keyword }).then((r) => r.data.data),
    enabled: tab === 'users' && user?.role === 'admin',
  });

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ['admin-posts', page, keyword],
    queryFn: () => adminApi.listPosts({ page, page_size: 20, keyword }).then((r) => r.data.data),
    enabled: tab === 'posts' && user?.role === 'admin',
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, active }: { id: number; active: boolean }) =>
      adminApi.setUserActive(id, active),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] });
      toast.success('操作成功');
    },
    onError: () => toast.error('操作失败'),
  });

  const togglePinMutation = useMutation({
    mutationFn: (id: number) => adminApi.togglePin(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-posts'] });
      toast.success('操作成功');
    },
    onError: () => toast.error('操作失败'),
  });

  if (!isAuthenticated || user?.role !== 'admin') return null;

  const handleTabChange = (t: 'users' | 'posts') => {
    setTab(t);
    setPage(1);
    setKeyword('');
  };

  return (
    <div className="max-w-6xl mx-auto">
      {/* Header */}
      <div className="flex items-center gap-3 mb-6">
        <LayoutDashboard className="w-6 h-6 text-primary" />
        <h1 className="text-2xl font-bold">管理后台</h1>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-boxed bg-base-100 border border-base-300 mb-4 p-1 w-fit">
        <button
          className={`tab gap-2 ${tab === 'users' ? 'tab-active' : ''}`}
          onClick={() => handleTabChange('users')}
        >
          <Users className="w-4 h-4" /> 用户管理
        </button>
        <button
          className={`tab gap-2 ${tab === 'posts' ? 'tab-active' : ''}`}
          onClick={() => handleTabChange('posts')}
        >
          <FileText className="w-4 h-4" /> 帖子管理
        </button>
      </div>

      {/* Search bar */}
      <div className="flex gap-3 mb-4">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
          <input
            type="text"
            placeholder={tab === 'users' ? '搜索用户名/邮箱...' : '搜索帖子标题...'}
            value={keyword}
            onChange={(e) => { setKeyword(e.target.value); setPage(1); }}
            className="input input-bordered input-sm w-full pl-9 focus:outline-none focus:border-primary"
          />
        </div>
      </div>

      {/* Users table */}
      {tab === 'users' && (
        <div className="card bg-base-100 border border-base-300 shadow-sm overflow-hidden">
          {usersLoading ? (
            <div className="flex justify-center py-12">
              <span className="loading loading-spinner loading-lg text-primary" />
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>用户</th>
                    <th>邮箱</th>
                    <th>角色</th>
                    <th>积分</th>
                    <th>注册时间</th>
                    <th>状态</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {(usersData?.list ?? []).map((u) => (
                    <tr key={u.id}>
                      <td>
                        <div className="flex items-center gap-2">
                          <div className="avatar">
                            <div className="w-8 h-8 rounded-full">
                              <Image
                                src={u.avatar || `https://api.dicebear.com/8.x/initials/svg?seed=${u.username}`}
                                alt={u.username}
                                width={32}
                                height={32}
                                className="rounded-full"
                              />
                            </div>
                          </div>
                          <span className="font-medium text-sm">{u.username}</span>
                        </div>
                      </td>
                      <td className="text-sm text-base-content/60">{u.email}</td>
                      <td>
                        <span className={`badge badge-sm ${u.role === 'admin' ? 'badge-warning' : 'badge-ghost'}`}>
                          {u.role === 'admin' ? '管理员' : '用户'}
                        </span>
                      </td>
                      <td className="text-sm font-medium text-warning">{u.score}</td>
                      <td className="text-xs text-base-content/50">{formatDate(u.created_at)}</td>
                      <td>
                        <span className={`badge badge-sm ${u.is_active ? 'badge-success' : 'badge-error'}`}>
                          {u.is_active ? '正常' : '已封禁'}
                        </span>
                      </td>
                      <td>
                        {u.id !== user?.id && (
                          <button
                            className={`btn btn-xs gap-1 ${u.is_active ? 'btn-error btn-outline' : 'btn-success btn-outline'}`}
                            onClick={() => toggleActiveMutation.mutate({ id: u.id, active: !u.is_active })}
                            disabled={toggleActiveMutation.isPending}
                          >
                            {u.is_active ? (
                              <><ShieldOff className="w-3 h-3" /> 封禁</>
                            ) : (
                              <><ShieldCheck className="w-3 h-3" /> 解封</>
                            )}
                          </button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {/* Pagination */}
          {(usersData?.total ?? 0) > 20 && (
            <div className="flex justify-center p-4">
              <div className="join">
                <button className="join-item btn btn-sm" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>«</button>
                <button className="join-item btn btn-sm btn-active">{page}</button>
                <button
                  className="join-item btn btn-sm"
                  disabled={page * 20 >= (usersData?.total ?? 0)}
                  onClick={() => setPage((p) => p + 1)}
                >»</button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Posts table */}
      {tab === 'posts' && (
        <div className="card bg-base-100 border border-base-300 shadow-sm overflow-hidden">
          {postsLoading ? (
            <div className="flex justify-center py-12">
              <span className="loading loading-spinner loading-lg text-primary" />
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>标题</th>
                    <th>作者</th>
                    <th>类型</th>
                    <th>状态</th>
                    <th>浏览/点赞</th>
                    <th>发布时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {(postsData?.list ?? []).map((post) => (
                    <tr key={post.id}>
                      <td className="max-w-xs">
                        <a
                          href={`/posts/${post.id}`}
                          target="_blank"
                          rel="noreferrer"
                          className="text-sm hover:text-primary transition-colors line-clamp-1"
                        >
                          {post.pin_top && <Pin className="w-3 h-3 inline mr-1 text-primary" />}
                          {post.title}
                        </a>
                      </td>
                      <td className="text-sm text-base-content/60">{post.author?.username}</td>
                      <td>
                        <span className={`badge badge-sm ${
                          post.type === 'article' ? 'badge-secondary' :
                          post.type === 'topic' ? 'badge-accent' : 'badge-ghost'
                        }`}>
                          {post.type === 'article' ? '文章' : post.type === 'topic' ? '话题' : '帖子'}
                        </span>
                      </td>
                      <td>
                        <span className={`badge badge-sm ${post.status === 'published' ? 'badge-success' : 'badge-warning'}`}>
                          {post.status === 'published' ? '已发布' : post.status === 'draft' ? '草稿' : '隐藏'}
                        </span>
                      </td>
                      <td className="text-xs text-base-content/50">
                        {post.view_count} / {post.like_count}
                      </td>
                      <td className="text-xs text-base-content/50">{formatDate(post.created_at)}</td>
                      <td>
                        <button
                          className={`btn btn-xs gap-1 ${post.pin_top ? 'btn-ghost' : 'btn-primary btn-outline'}`}
                          onClick={() => togglePinMutation.mutate(post.id)}
                          disabled={togglePinMutation.isPending}
                        >
                          {post.pin_top ? (
                            <><PinOff className="w-3 h-3" /> 取消置顶</>
                          ) : (
                            <><Pin className="w-3 h-3" /> 置顶</>
                          )}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {(postsData?.total ?? 0) > 20 && (
            <div className="flex justify-center p-4">
              <div className="join">
                <button className="join-item btn btn-sm" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>«</button>
                <button className="join-item btn btn-sm btn-active">{page}</button>
                <button
                  className="join-item btn btn-sm"
                  disabled={page * 20 >= (postsData?.total ?? 0)}
                  onClick={() => setPage((p) => p + 1)}
                >»</button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
