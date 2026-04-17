import { User } from "@/lib/api";
import Link from "next/link";
import { TrophyIcon, UserCircleIcon } from "@heroicons/react/24/outline";
import { LeaderboardItemResponse } from "@/lib/api/modules/users";

// 活跃用户卡片
interface ActiveUserCardProps {
  user: LeaderboardItemResponse;
  rank?: number; // 排名
}

export function ActiveUserCard({ user, rank }: ActiveUserCardProps) {
  // 根据排名获取奖牌颜色
  const getRankColor = (rank: number) => {
    switch (rank) {
      case 1:
        return "text-yellow-500 bg-yellow-100 dark:bg-yellow-900/20";
      case 2:
        return "text-gray-500 bg-gray-100 dark:bg-gray-800";
      case 3:
        return "text-orange-500 bg-orange-100 dark:bg-orange-900/20";
      default:
        return "text-base-content/40 bg-base-200";
    }
  };

  // 根据排名获取奖牌图标
  const getRankIcon = (rank: number) => {
    switch (rank) {
      case 1:
        return "🥇";
      case 2:
        return "🥈";
      case 3:
        return "🥉";
      default:
        return null;
    }
  };

  return (
    <Link href={`/users/${user.id}`} className="block">
      <div className="group card bg-base-100 shadow-sm border border-base-200 hover:shadow-md hover:border-primary/20 transition-all duration-200 cursor-pointer">
        <div className="card-body p-3">
          <div className="flex items-center gap-3">
            {/* 排名徽章 */}
            {rank && rank <= 3 && (
              <div className={`w-8 h-8 rounded-full flex items-center justify-center ${getRankColor(rank)}`}>
                <span className="text-lg">{getRankIcon(rank)}</span>
              </div>
            )}
            
            {/* 头像 */}
            {user.avatar ? (
              <div className="avatar">
                <div className="w-10 h-10 rounded-full ring-1 ring-base-200 group-hover:ring-primary transition-all">
                  <img 
                    src={user.avatar} 
                    alt={user.username} 
                    className="object-cover"
                  />
                </div>
              </div>
            ) : (
              <div className="avatar placeholder">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-red-100 to-red-200 dark:from-red-900/30 dark:to-red-800/20 group-hover:from-red-200 group-hover:to-red-300 transition-all">
                  <span className="text-red-600 dark:text-red-400 font-medium">
                    {user.username.charAt(0).toUpperCase()}
                  </span>
                </div>
              </div>
            )}
            
            {/* 用户信息 */}
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <span className="font-medium text-base-content truncate group-hover:text-primary transition-colors">
                  {user.username}
                </span>
                {/* {user.role === "admin" && (
                  <span className="badge badge-primary badge-xs">管理员</span>
                )} */}
              </div>
              {/* <div className="flex items-center gap-2 text-xs text-base-content/50">
                <span>积分: {user.score || 0}</span>
                {user.post_count !== undefined && (
                  <>
                    <span>•</span>
                    <span>帖子: {user.post_count}</span>
                  </>
                )}
              </div> */}
            </div>
            
            {/* 箭头指示 */}
            <div className="opacity-0 group-hover:opacity-100 transition-opacity">
              <svg 
                className="w-4 h-4 text-base-content/40 group-hover:text-primary transition-colors" 
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}