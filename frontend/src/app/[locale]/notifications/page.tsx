"use client";

import { useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { notificationApi } from "@/lib/api";
import type { Notification } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";
import { timeAgo } from "@/lib/utils";
import {
  Bell,
  CheckCheck,
  Heart,
  MessageSquare,
  UserPlus,
  Info,
} from "lucide-react";
import toast from "react-hot-toast";
import Avatar from "@/components/user/Avatar";
import { useTranslations } from "next-intl";

// 通知图标组件
const NotifIcon = ({ type }: { type: Notification["type"] }) => {
  const cls = "w-4 h-4";
  switch (type) {
    case "like":
      return <Heart className={`${cls} text-error`} />;
    case "comment":
      return <MessageSquare className={`${cls} text-info`} />;
    case "reply":
      return <MessageSquare className={`${cls} text-success`} />;
    case "follow":
      return <UserPlus className={`${cls} text-primary`} />;
    default:
      return <Info className={`${cls} text-warning`} />;
  }
};

export default function NotificationsPage() {
  const { isAuthenticated } = useAuthStore();
  const router = useRouter();
  const queryClient = useQueryClient();
  const t = useTranslations("notifications");

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/auth/login");
    }
  }, [isAuthenticated, router]);

  // 获取通知列表
  const { data, isLoading } = useQuery({
    queryKey: ["notifications"],
    queryFn: async () => {
      const response = await notificationApi.list({ page: 1, page_size: 50 });
      return response.data.data;
    },
    enabled: isAuthenticated,
  });

  // 全部标记已读
  const markAllMutation = useMutation({
    mutationFn: () => notificationApi.markAllRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "unread"] });
      toast.success(t("mark_all_success") || "已全部标记为已读");
    },
    onError: () => {
      toast.error(t("mark_all_error") || "操作失败，请重试");
    },
  });

  const notifications: Notification[] = data?.list ?? [];
  const unreadCount = notifications.filter((n) => !n.is_read).length;

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-6">
      {/* 头部 */}
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold flex items-center gap-2">
          <Bell className="w-6 h-6 text-primary" />
          {t("title")}
          {unreadCount > 0 && (
            <span className="badge badge-error badge-sm">{unreadCount}</span>
          )}
        </h1>
        {unreadCount > 0 && (
          <button
            className="btn btn-ghost btn-sm gap-1"
            onClick={() => markAllMutation.mutate()}
            disabled={markAllMutation.isPending}
          >
            <CheckCheck className="w-4 h-4" />
            {markAllMutation.isPending ? t("marking") : t("mark_all_read")}
          </button>
        )}
      </div>

      {/* 加载状态 */}
      {isLoading ? (
        <div className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="skeleton h-20 w-full rounded-xl" />
          ))}
        </div>
      ) : notifications.length === 0 ? (
        /* 空状态 */
        <div className="text-center py-20 text-base-content/40">
          <Bell className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>{t("no_notifications")}</p>
        </div>
      ) : (
        /* 通知列表 */
        <div className="space-y-2">
          {notifications.map((notification) => (
            <NotificationCard
              key={notification.id}
              notification={notification}
              t={t}
            />
          ))}
        </div>
      )}
    </div>
  );
}

// 抽取通知卡片组件
function NotificationCard({
  notification,
  t,
}: {
  notification: Notification;
  t: (key: string) => string;
}) {
  const queryClient = useQueryClient();

  // 标记单条已读
  const markReadMutation = useMutation({
    mutationFn: () => notificationApi.markRead(notification.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "unread"] });
    },
  });

  const handleClick = () => {
    if (!notification.is_read) {
      markReadMutation.mutate();
    }
    // 跳转到相关页面
    if (notification.target_id && notification.target_type) {
      // 根据类型跳转
      // router.push(`/posts/${notification.target_id}`);
    }
  };

  return (
    <div
      className={`card border transition-colors cursor-pointer hover:shadow-md ${
        notification.is_read
          ? "bg-base-100 border-base-300"
          : "bg-primary/5 border-primary/20"
      }`}
      onClick={handleClick}
    >
      <div className="card-body p-4">
        <div className="flex items-start gap-3">
          {/* 发送者头像 */}
          {notification.sender ? (
            <div className="avatar flex-none">
              <div className="w-9 h-9 rounded-full">
                <Avatar
                  username={notification.sender.username}
                  avatarUrl={notification.sender.avatar}
                  size="md"
                />
              </div>
            </div>
          ) : (
            <div className="w-9 h-9 rounded-full bg-base-200 flex items-center justify-center flex-none">
              <Bell className="w-4 h-4 text-base-content/40" />
            </div>
          )}

          {/* 通知内容 */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <NotifIcon type={notification.type} />
              <p className="text-sm text-base-content/80 flex-1">
                {notification.content}
              </p>
              {!notification.is_read && (
                <span className="w-2 h-2 bg-primary rounded-full flex-none" />
              )}
            </div>
            <p className="text-xs text-base-content/40 mt-1">
              {timeAgo(notification.created_at)}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
