// components/member/MemberNotifications.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Bell, CheckCheck } from "lucide-react";
import { timeAgo } from "@/shared/lib/utils";
import { useNotifications } from "@/features/notification/hooks/useNotification";

// 会员中心通知列表：标记已读 / 全部已读
export function MemberNotifications() {
  const t = useTranslations("Member");
  const { notifications, loading, markAsRead, markAllAsRead } =
    useNotifications({ pageSize: 50 });
  const [isMarkingAll, setIsMarkingAll] = useState(false);

  const unreadCount = notifications.filter((n) => !n.is_read).length;

  const handleMarkAllRead = async () => {
    if (isMarkingAll) return;
    setIsMarkingAll(true);
    await markAllAsRead();
    setIsMarkingAll(false);
  };

  if (loading && notifications.length === 0) {
    return (
      <div className="space-y-2">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="skeleton h-20 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  if (notifications.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_notifications")}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {unreadCount > 0 && (
        <div className="flex justify-end">
          <button
            onClick={handleMarkAllRead}
            disabled={isMarkingAll}
            className="btn btn-ghost btn-sm gap-1"
          >
            <CheckCheck className="w-4 h-4" />
            {t("mark_all_read")}
          </button>
        </div>
      )}

      {notifications.map((notif) => (
        <div
          key={notif.id}
          className={`card bg-base-100 border p-4 ${
            !notif.is_read ? "border-primary bg-primary/5" : "border-base-300"
          }`}
        >
          <div className="flex justify-between items-start gap-3">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <Bell className="w-4 h-4 text-base-content/40 flex-none" />
                <p className="text-sm text-base-content/80 flex-1">
                  {notif.content}
                </p>
                {!notif.is_read && (
                  <span className="badge badge-sm badge-primary">
                    {t("new")}
                  </span>
                )}
              </div>
              <p className="text-xs text-base-content/40 mt-1 ml-6">
                {timeAgo(notif.created_at)}
              </p>
            </div>
            {!notif.is_read && (
              <button
                onClick={() => markAsRead(notif.id)}
                className="btn btn-xs btn-primary"
              >
                {t("mark_as_read")}
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
