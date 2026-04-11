'use client';

import { useQuery } from '@tanstack/react-query';
import { userApi } from '@/lib/api';
import Image from 'next/image';
import Link from 'next/link';
import { Trophy, Star, Crown } from 'lucide-react';
import Avatar from '@/components/user/Avatar';

export default function LeaderboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['leaderboard', 50],
    queryFn: () => userApi.leaderboard(50).then((r) => r.data.data),
  });

  const users = data ?? [];

  return (
    <div className="max-w-2xl mx-auto">
      <div className="text-center mb-8">
        <Trophy className="w-12 h-12 text-warning mx-auto mb-3" />
        <h1 className="text-3xl font-bold">积分排行榜</h1>
        <p className="text-base-content/50 text-sm mt-2">参与讨论、发帖、评论都能获得积分</p>
      </div>

      {isLoading ? (
        <div className="space-y-2">
          {Array.from({ length: 10 }).map((_, i) => (
            <div key={i} className="skeleton h-16 w-full rounded-xl" />
          ))}
        </div>
      ) : (
        <div className="space-y-2">
          {users.map((u, i) => (
            <Link
              key={u.id}
              href={`/users/${u.id}`}
              className={`card border shadow-sm hover:shadow-md transition-all duration-200 hover:-translate-y-0.5 block ${
                i === 0 ? 'border-yellow-400 bg-yellow-50 dark:bg-yellow-900/10' :
                i === 1 ? 'border-gray-300 bg-gray-50 dark:bg-gray-900/10' :
                i === 2 ? 'border-amber-600 bg-amber-50 dark:bg-amber-900/10' :
                'border-base-300 bg-base-100'
              }`}
            >
              <div className="card-body p-4">
                <div className="flex items-center gap-4">
                  {/* Rank */}
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center font-black text-sm flex-none ${
                    i === 0 ? 'bg-yellow-400 text-yellow-900' :
                    i === 1 ? 'bg-gray-300 text-gray-700' :
                    i === 2 ? 'bg-amber-600 text-white' :
                    'bg-base-200 text-base-content/50'
                  }`}>
                    {i < 3 ? <Crown className="w-5 h-5" /> : i + 1}
                  </div>

                  {/* Avatar */}
                  <div className="avatar">
                    <div className="w-10 h-10 rounded-full">
                     
                       <Avatar 
  username={u.username} 
  avatarUrl={u.avatar}  // 数据库中的头像
  size="md" 
/>
                    </div>
                  </div>

                  {/* Name */}
                  <div className="flex-1">
                    <div className="font-semibold">{u.username}</div>
                    {u.bio && <div className="text-xs text-base-content/40 truncate max-w-xs">{u.bio}</div>}
                  </div>

                  {/* Score */}
                  <div className="flex items-center gap-1.5 font-bold text-lg">
                    <Star className={`w-5 h-5 ${i < 3 ? 'text-warning' : 'text-base-content/30'}`} />
                    <span className={i < 3 ? 'text-warning' : ''}>{u.score}</span>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
