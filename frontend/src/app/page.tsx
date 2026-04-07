'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { postApi, tagApi, userApi } from '@/lib/api';
import PostCard from '@/components/post/PostCard';
import Link from 'next/link';
import Image from 'next/image';
import { Flame, Clock, Tag as TagIcon, Trophy, ChevronRight, PenSquare } from 'lucide-react';
import { useAuthStore } from '@/store/auth';
import { formatDate } from '@/lib/utils';

export default function HomePage() {
  const { isAuthenticated } = useAuthStore();
  const [sortBy, setSortBy] = useState<'' | 'hot'>('');
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [page, setPage] = useState(1);

  const { data: postsData, isLoading } = useQuery({
    queryKey: ['posts', sortBy, selectedTag, page],
    queryFn: () =>
      postApi.list({
        page,
        page_size: 15,
        sort_by: sortBy,
        tag_id: selectedTag ?? undefined,
      }).then((r) => r.data.data),
  });

  const { data: tags } = useQuery({
    queryKey: ['tags'],
    queryFn: () => tagApi.list().then((r) => r.data.data),
  });

  const { data: leaderboard } = useQuery({
    queryKey: ['leaderboard'],
    queryFn: () => userApi.leaderboard(10).then((r) => r.data.data),
  });

  const posts = postsData?.list ?? [];
  const total = postsData?.total ?? 0;
  const totalPages = Math.ceil(total / 15);

  return (
    <div className="flex flex-col lg:flex-row gap-6">
      {/* Main content */}
      <div className="flex-1 min-w-0">
        {/* Filter bar */}
        <div className="flex items-center justify-between mb-4 bg-base-100 rounded-xl p-3 border border-base-300">
          <div className="flex items-center gap-2">
            <button
              className={`btn btn-sm gap-1 ${sortBy === '' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => { setSortBy(''); setPage(1); }}
            >
              <Clock className="w-4 h-4" /> 最新
            </button>
            <button
              className={`btn btn-sm gap-1 ${sortBy === 'hot' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => { setSortBy('hot'); setPage(1); }}
            >
              <Flame className="w-4 h-4" /> 热门
            </button>
          </div>

          {isAuthenticated && (
            <Link href="/posts/new" className="btn btn-primary btn-sm gap-1">
              <PenSquare className="w-4 h-4" /> 发帖
            </Link>
          )}
        </div>

        {/* Posts */}
        {isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="skeleton h-28 w-full rounded-xl" />
            ))}
          </div>
        ) : posts.length === 0 ? (
          <div className="text-center py-20 text-base-content/40">
            <p className="text-lg">暂无帖子</p>
            {isAuthenticated && (
              <Link href="/posts/new" className="btn btn-primary mt-4">
                发布第一篇帖子
              </Link>
            )}
          </div>
        ) : (
          <div className="space-y-3">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center mt-6">
            <div className="join">
              <button
                className="join-item btn btn-sm"
                disabled={page === 1}
                onClick={() => setPage((p) => p - 1)}
              >
                «
              </button>
              {Array.from({ length: Math.min(totalPages, 7) }, (_, i) => i + 1).map((p) => (
                <button
                  key={p}
                  className={`join-item btn btn-sm ${page === p ? 'btn-active btn-primary' : ''}`}
                  onClick={() => setPage(p)}
                >
                  {p}
                </button>
              ))}
              <button
                className="join-item btn btn-sm"
                disabled={page === totalPages}
                onClick={() => setPage((p) => p + 1)}
              >
                »
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Sidebar */}
      <aside className="w-full lg:w-64 xl:w-72 flex-none space-y-4">
        {/* Tags */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h3 className="font-bold flex items-center gap-2 mb-3">
              <TagIcon className="w-4 h-4 text-primary" /> 热门标签
            </h3>
            <div className="flex flex-wrap gap-2">
              <button
                onClick={() => { setSelectedTag(null); setPage(1); }}
                className={`badge badge-lg cursor-pointer ${!selectedTag ? 'badge-primary' : 'badge-ghost hover:badge-primary'}`}
              >
                全部
              </button>
              {(tags ?? []).slice(0, 12).map((tag) => (
                <button
                  key={tag.id}
                  onClick={() => { setSelectedTag(selectedTag === tag.id ? null : tag.id); setPage(1); }}
                  className="badge badge-lg cursor-pointer hover:opacity-80 transition-opacity"
                  style={{
                    backgroundColor: selectedTag === tag.id ? tag.color : tag.color + '20',
                    color: tag.color,
                    borderColor: tag.color + '40',
                  }}
                >
                  {tag.name}
                  <span className="ml-1 opacity-60 text-xs">({tag.post_count})</span>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Leaderboard */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h3 className="font-bold flex items-center gap-2 mb-3">
              <Trophy className="w-4 h-4 text-warning" /> 积分排行榜
            </h3>
            <div className="space-y-2">
              {(leaderboard ?? []).slice(0, 8).map((u, i) => (
                <Link
                  key={u.id}
                  href={`/users/${u.id}`}
                  className="flex items-center gap-2 hover:bg-base-200 rounded-lg p-1.5 transition-colors"
                >
                  <span className={`w-5 h-5 text-xs font-bold flex items-center justify-center rounded-full ${
                    i === 0 ? 'bg-yellow-400 text-yellow-900' :
                    i === 1 ? 'bg-gray-300 text-gray-700' :
                    i === 2 ? 'bg-amber-600 text-white' :
                    'text-base-content/40'
                  }`}>
                    {i + 1}
                  </span>
                  <Image
                    src={u.avatar || `https://api.dicebear.com/8.x/initials/svg?seed=${u.username}`}
                    alt={u.username}
                    width={24}
                    height={24}
                    className="rounded-full"
                  />
                  <span className="flex-1 text-sm truncate">{u.username}</span>
                  <span className="text-xs text-warning font-medium">{u.score}</span>
                </Link>
              ))}
            </div>
            <Link href="/leaderboard" className="btn btn-ghost btn-xs mt-2 gap-1">
              查看完整排行 <ChevronRight className="w-3 h-3" />
            </Link>
          </div>
        </div>

        {/* Site info */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4 text-xs text-base-content/50 space-y-1">
            <p className="font-medium text-base-content/70">关于 BBS Forum</p>
            <p>一个现代化的技术交流社区，欢迎分享你的知识与想法。</p>
            <p className="pt-1">© {new Date().getFullYear()} BBS Forum</p>
          </div>
        </div>
      </aside>
    </div>
  );
}
