// components/member/MemberNotifications.tsx
"use client";

import { useTranslations } from "next-intl";

interface MemberNotificationsProps {
  onMarkRead: (id: number) => void;
  onMarkAllRead: () => void;
  notifications: Notification[];
}

interface Notification {
  id: number;
  title: string;
  content: string;
  is_read: boolean;
  created_at: string;
}

export function MemberNotifications({
  onMarkRead,
  onMarkAllRead,
  notifications,
}: MemberNotificationsProps) {
  const t = useTranslations("Member");
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
      <div className="flex justify-end mb-4">
        <button onClick={onMarkAllRead} className="btn btn-ghost btn-sm">
          {t("mark_all_read")}
        </button>
      </div>

      {notifications.map((notif: Notification) => (
        <div
          key={notif.id}
          className={`card bg-base-100 border p-4 ${
            !notif.is_read ? "border-primary bg-primary/5" : "border-base-300"
          }`}
        >
          <div className="flex justify-between items-start">
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <p className="font-medium">{notif.title}</p>
                {!notif.is_read && (
                  <span className="badge badge-sm badge-primary">
                    {t("new")}
                  </span>
                )}
              </div>
              <p className="text-sm text-base-content/60 mt-1">
                {notif.content}
              </p>
              <p className="text-xs text-base-content/40 mt-2">
                {new Date(notif.created_at).toLocaleString()}
              </p>
            </div>
            {!notif.is_read && (
              <button
                onClick={() => onMarkRead(notif.id)}
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
