'use client';

import { useLeaderboard } from '@/hooks/useLeaderboard';
import Link from 'next/link';
import { Trophy, Star, Crown, Medal, TrendingUp, Users } from 'lucide-react';
import Avatar from '@/components/user/Avatar';
import { useTranslations } from 'next-intl';
import React from 'react';

// 加载骨架屏组件
function LoadingSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 10 }).map((_, i) => (
        <div key={i} className="skeleton h-20 w-full rounded-xl" />
      ))}
    </div>
  );
}

// 排名徽章组件
function RankBadge({ rank, isTopThree }: { rank: number; isTopThree: boolean }) {
  if (!isTopThree) {
    return (
      <div className="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm bg-base-200 text-base-content/50">
        {rank}
      </div>
    );
  }
  
  const config = {
    1: { bg: 'bg-gradient-to-br from-yellow-400 to-yellow-500', text: 'text-yellow-900', icon: Crown },
    2: { bg: 'bg-gradient-to-br from-gray-300 to-gray-400', text: 'text-gray-700', icon: Medal },
    3: { bg: 'bg-gradient-to-br from-amber-500 to-amber-600', text: 'text-white', icon: Medal },
  };
  
  const { bg, text, icon: Icon } = config[rank as keyof typeof config];
  
  return (
    <div className={`w-10 h-10 rounded-full flex items-center justify-center ${bg} ${text} shadow-md`}>
      <Icon className="w-5 h-5" />
    </div>
  );
}

// 用户卡片组件
function UserCard({ user, rank }: { user: any; rank: number }) {
  const t = useTranslations("Leaderboard")
  const isTopThree = rank < 3;
  const cardStyles = {
    0: 'border-yellow-400/50 bg-gradient-to-r from-yellow-50/50 to-transparent dark:from-yellow-900/10',
    1: 'border-gray-300/50 bg-gradient-to-r from-gray-50/50 to-transparent dark:from-gray-900/10',
    2: 'border-amber-500/50 bg-gradient-to-r from-amber-50/50 to-transparent dark:from-amber-900/10',
  }[rank] || 'border-base-200 bg-base-100';
  
  const scoreColor = isTopThree ? 'text-warning' : 'text-base-content/40';

  return (
    <Link href={`/users/${user.id}`}>
      <div className={`card border shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 ${cardStyles}`}>
        <div className="card-body p-4">
          <div className="flex items-center gap-4">
            {/* 排名 */}
            <RankBadge rank={rank + 1} isTopThree={isTopThree} />

            {/* 头像 */}
            <div className="avatar">
              <div className="w-11 h-11 rounded-full ring-2 ring-primary/20 ring-offset-2">
                <Avatar 
                  username={user.username} 
                  avatarUrl={user.avatar}
                  size="md" 
                />
              </div>
            </div>

            {/* 用户信息 */}
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-base-content truncate">
                  {user.username}
                </span>
                {rank === 0 && (
                  <span className="badge badge-warning badge-sm gap-1">
                    <Crown className="w-3 h-3" />
                 
                    {t("first")}
                  </span>
                )}
              </div>
              {user.bio && (
                <div className="text-xs text-base-content/40 truncate max-w-xs">
                  {user.bio}
                </div>
              )}
            </div>

            {/* 积分 */}
            <div className="flex items-center gap-1.5">
              <Star className={`w-5 h-5 ${isTopThree ? 'text-warning' : 'text-base-content/30'}`} />
              <span className={`font-bold text-lg ${scoreColor}`}>
                {user.score?.toLocaleString() || 0}
              </span>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}

// 空状态组件
function EmptyState() {
  const t = useTranslations("Leaderboard");
  return (
    <div className="bg-base-100 rounded-2xl shadow-sm p-12 text-center border border-base-200">
      <div className="text-6xl mb-4 opacity-50">📊</div>
      <h3 className="text-lg font-semibold text-base-content mb-2">
        {t("no_data")}
      </h3>
      <p className="text-base-content/60 text-sm">
        {t("no_data_description")}
      </p>
    </div>
  );
}

// 统计卡片组件
function StatsCards({ totalUsers, topScore }: { totalUsers: number; topScore: number }) {
    const t = useTranslations("Leaderboard");
  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
      <div className="bg-gradient-to-br from-primary/10 to-primary/5 rounded-2xl p-4 text-center border border-primary/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <Users className="w-5 h-5 text-primary" />
          <span className="text-sm font-medium text-primary">总用户数</span>
        </div>
        <div className="text-2xl font-bold text-base-content">{totalUsers}</div>
      </div>
      
      <div className="bg-gradient-to-br from-warning/10 to-warning/5 rounded-2xl p-4 text-center border border-warning/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <Trophy className="w-5 h-5 text-warning" />
          <span className="text-sm font-medium text-warning">最高积分</span>
        </div>
        <div className="text-2xl font-bold text-base-content">{topScore?.toLocaleString() || 0}</div>
      </div>
      
      <div className="bg-gradient-to-br from-secondary/10 to-secondary/5 rounded-2xl p-4 text-center border border-secondary/20">
        <div className="flex items-center justify-center gap-2 mb-2">
          <TrendingUp className="w-5 h-5 text-secondary" />
          <span className="text-sm font-medium text-secondary">活跃排名</span>
        </div>
        <div className="text-2xl font-bold text-base-content">实时更新</div>
      </div>
    </div>
  );
}

export default function LeaderboardPage() {
  // 使用自定义 hook，指定需要返回的字段（包括 bio 用于显示简介）
  const { data, isLoading, error } = useLeaderboard({
    limit: 50,
    fields: 'id,username,avatar,score,bio' // 确保返回 bio 字段
  });
console.log(data);

  const t = useTranslations("Leaderboard");
  const users = data ?? [];
  const totalUsers = users.length;
  const topScore = users[0]?.score || 0;
  console.log(users,totalUsers,topScore);

  if (error) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="bg-error/10 rounded-2xl p-8 text-center border border-error/20">
          <div className="text-5xl mb-4">⚠️</div>
          <h3 className="text-lg font-semibold text-error mb-2">加载失败</h3>
          <p className="text-base-content/60 text-sm">请稍后重试</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-200 to-base-100">
      <div className="max-w-3xl mx-auto px-4 py-8 md:py-12">
        {/* 头部区域 */}
        <div className="text-center mb-10">
          <div className="relative inline-block mb-4">
            <div className="absolute inset-0 bg-gradient-to-r from-primary/20 to-warning/20 rounded-full blur-2xl" />
            <div className="relative bg-gradient-to-br from-primary to-warning p-4 rounded-full shadow-lg">
              <Trophy className="w-10 h-10 text-white" />
            </div>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-warning bg-clip-text text-transparent mb-2">
            {t("title")}
          </h1>
          <p className="text-base-content/60 text-sm">
            {t('description')}
          </p>
        </div>

        {/* 统计卡片 */}
        {!isLoading && users.length > 0 && (
          <StatsCards totalUsers={totalUsers} topScore={topScore} />
        )}

        {/* 排行榜列表 */}
        {isLoading ? (
          <LoadingSkeleton />
        ) : totalUsers === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-3">
            {/* 前三名特别展示 */}
            {users.slice(0, 3).map((user, index) => (
              <UserCard key={user.id} user={user} rank={index} />
            ))}
            
            {/* 分隔线 */}
            {users.length > 3 && (
              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-base-200"></div>
                </div>
                <div className="relative flex justify-center text-xs">
                  <span className="px-3 bg-base-200/50 text-base-content/40 rounded-full py-1">
                    更多优秀用户
                  </span>
                </div>
              </div>
            )}
            
            {/* 其他用户 */}
            {users.slice(3).map((user, index) => (
              <UserCard key={user.id} user={user} rank={index + 3} />
            ))}
            
            {/* 底部提示 */}
            <div className="text-center text-xs text-base-content/40 mt-6 pt-4 border-t border-base-200">
              仅显示前 {users.length} 名用户
            </div>
          </div>
        )}
      </div>
    </div>
  );
}