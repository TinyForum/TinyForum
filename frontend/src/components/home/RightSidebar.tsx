"use client";

import Link from "next/link";
import { 
  User, 
  Bell, 
  Activity, 
  AlertTriangle, 
  Award,
  TrendingUp,
  Clock,
  Heart,
  MessageSquare,
  UserPlus,
  CheckCircle,
  Zap,
  FileText
} from "lucide-react";
import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";
import { formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";

interface UserProfile {
  id: number;
  username: string;
  avatar: string;
  bio: string;
  score: number;
  created_at: string;
}

interface TimelineEvent {
  id: number;
  action: string;
  target_type: string;
  target_id: number;
  created_at: string;
  actor?: {
    id: number;
    username: string;
    avatar: string;
  };
}

interface RightSidebarProps {
  isAuthenticated: boolean;
  user: any;
  userProfile?: UserProfile;
  leaderboard: any[];
  unreadCount: number;
  timelineEvents: TimelineEvent[];
}

export default function RightSidebar({
  isAuthenticated,
  user,
  userProfile,
  leaderboard,
  unreadCount,
  timelineEvents,
}: RightSidebarProps) {
  const t = useTranslations("sidebar");

  const getEventIcon = (action: string) => {
    switch (action) {
      case "create_post":
        return <FileText className="w-3 h-3" />;
      case "create_comment":
        return <MessageSquare className="w-3 h-3" />;
      case "like_post":
      case "like_comment":
        return <Heart className="w-3 h-3" />;
      case "follow_user":
        return <UserPlus className="w-3 h-3" />;
      case "accept_answer":
        return <CheckCircle className="w-3 h-3" />;
      default:
        return <Activity className="w-3 h-3" />;
    }
  };

  const getEventText = (event: TimelineEvent) => {
    const actor = event.actor?.username || "用户";
    switch (event.action) {
      case "create_post":
        return `${actor} 发布了新帖子`;
      case "create_comment":
        return `${actor} 发表了评论`;
      case "like_post":
        return `${actor} 赞了一个帖子`;
      case "like_comment":
        return `${actor} 赞了一条评论`;
      case "follow_user":
        return `${actor} 关注了你`;
      case "accept_answer":
        return `${actor} 采纳了你的回答`;
      default:
        return `${actor} 有了新动态`;
    }
  };

  return (
    <aside className="space-y-4">
      {/* 用户信息卡片 */}
      {isAuthenticated && userProfile ? (
        <div className="rounded-lg border bg-card">
          <div className="p-4 text-center border-b">
            <div className="relative inline-block">
              <img
                src={userProfile.avatar || "/default-avatar.png"}
                alt={userProfile.username}
                className="w-16 h-16 rounded-full mx-auto mb-2 object-cover"
              />
              {unreadCount > 0 && (
                <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                  {unreadCount > 99 ? "99+" : unreadCount}
                </span>
              )}
            </div>
            <h3 className="font-semibold">{userProfile.username}</h3>
            <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
              {userProfile.bio || "这个人很懒，什么都没写"}
            </p>
          </div>
          <div className="p-3 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground flex items-center gap-1">
                <Award className="w-4 h-4" />
                积分
              </span>
              <span className="font-medium text-primary">{userProfile.score}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground flex items-center gap-1">
                <Clock className="w-4 h-4" />
                注册时间
              </span>
              <span className="text-xs">
                {new Date(userProfile.created_at).toLocaleDateString()}
              </span>
            </div>
          </div>
          <div className="p-3 border-t space-y-2">
            <Link
              href="/notifications"
              className="flex items-center justify-between p-2 rounded-lg hover:bg-muted transition-colors"
            >
              <span className="flex items-center gap-2 text-sm">
                <Bell className="w-4 h-4" />
                通知
              </span>
              {unreadCount > 0 && (
                <span className="bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">
                  {unreadCount}
                </span>
              )}
            </Link>
            <Link
              href={`/profile/${userProfile.id}`}
              className="flex items-center gap-2 p-2 rounded-lg hover:bg-muted transition-colors text-sm"
            >
              <User className="w-4 h-4" />
              个人主页
            </Link>
            <Link
              href="/settings"
              className="flex items-center gap-2 p-2 rounded-lg hover:bg-muted transition-colors text-sm"
            >
              <Activity className="w-4 h-4" />
              设置
            </Link>
          </div>
        </div>
      ) : (
        <div className="rounded-lg border bg-card p-4 text-center">
          <User className="w-12 h-12 mx-auto text-muted-foreground mb-2" />
          <p className="text-sm text-muted-foreground mb-3">登录后查看更多功能</p>
          <Link
            href="/login"
            className="inline-block w-full bg-primary text-primary-foreground rounded-lg px-4 py-2 text-sm hover:bg-primary/90 transition-colors"
          >
            立即登录
          </Link>
          <Link
            href="/register"
            className="inline-block w-full mt-2 text-sm text-muted-foreground hover:text-primary transition-colors"
          >
            注册账号
          </Link>
        </div>
      )}

      {/* 违规情况（已登录时） */}
      {isAuthenticated && (
        <div className="rounded-lg border bg-card">
          <div className="p-3 border-b">
            <h3 className="font-semibold flex items-center gap-2">
              <AlertTriangle className="w-4 h-4 text-yellow-500" />
              {t("violation_status")}
            </h3>
          </div>
          <div className="p-3 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">违规次数</span>
              <span className="font-medium text-green-600">0</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">禁言状态</span>
              <span className="font-medium text-green-600">正常</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">警告等级</span>
              <span className="font-medium text-green-600">无</span>
            </div>
            <Link
              href="/violations"
              className="block text-xs text-center text-muted-foreground hover:text-primary mt-2"
            >
              查看详情 →
            </Link>
          </div>
        </div>
      )}

      {/* 时间线事件（已登录时） */}
      {isAuthenticated && timelineEvents && timelineEvents.length > 0 && (
        <div className="rounded-lg border bg-card">
          <div className="p-3 border-b">
            <h3 className="font-semibold flex items-center gap-2">
              <Zap className="w-4 h-4 text-blue-500" />
              {t("recent_updates")}
            </h3>
          </div>
          <div className="p-2 space-y-1 max-h-[300px] overflow-y-auto">
            {timelineEvents.slice(0, 5).map((event) => (
              <div
                key={event.id}
                className="flex items-start gap-2 p-2 rounded-lg hover:bg-muted transition-colors text-sm"
              >
                <div className="flex-shrink-0 mt-0.5 text-muted-foreground">
                  {getEventIcon(event.action)}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm truncate">{getEventText(event)}</p>
                  <p className="text-xs text-muted-foreground">
                    {formatDistanceToNow(new Date(event.created_at), {
                      addSuffix: true,
                      locale: zhCN,
                    })}
                  </p>
                </div>
              </div>
            ))}
          </div>
          <div className="p-2 border-t">
            <Link
              href="/timeline"
              className="block text-xs text-center text-muted-foreground hover:text-primary"
            >
              查看更多动态 →
            </Link>
          </div>
        </div>
      )}

      {/* 用户排行榜 */}
      {leaderboard && leaderboard.length > 0 && (
        <div className="rounded-lg border bg-card">
          <div className="p-3 border-b">
            <h3 className="font-semibold flex items-center gap-2">
              <TrendingUp className="w-4 h-4" />
              {t("leaderboard")}
            </h3>
          </div>
          <div className="p-2 space-y-1">
            {leaderboard.slice(0, 5).map((user, index) => (
              <Link
                key={user.id}
                href={`/profile/${user.id}`}
                className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted transition-colors"
              >
                <div className="flex-shrink-0 w-6 text-center">
                  {index === 0 && <span className="text-yellow-500 font-bold">🥇</span>}
                  {index === 1 && <span className="text-gray-400 font-bold">🥈</span>}
                  {index === 2 && <span className="text-amber-600 font-bold">🥉</span>}
                  {index > 2 && <span className="text-muted-foreground">{index + 1}</span>}
                </div>
                <img
                  src={user.avatar || "/default-avatar.png"}
                  alt={user.username}
                  className="w-6 h-6 rounded-full object-cover"
                />
                <span className="flex-1 text-sm truncate">{user.username}</span>
                <span className="text-xs font-medium text-primary">{user.score}</span>
              </Link>
            ))}
          </div>
          <div className="p-2 border-t">
            <Link
              href="/leaderboard"
              className="block text-xs text-center text-muted-foreground hover:text-primary"
            >
              完整排行榜 →
            </Link>
          </div>
        </div>
      )}
    </aside>
  );
}