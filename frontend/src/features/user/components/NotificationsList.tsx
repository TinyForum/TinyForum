// components/user/ViolationsList.tsx
"use client";

import { Notification } from "@/shared/api";
import { useTranslations } from "next-intl";

interface NotificationsListProps {
  notifications: Notification[];
}

export function ViolationsList({ notifications }: NotificationsListProps) {
  const t = useTranslations("Common");
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
        <button className="btn btn-ghost btn-sm">{t("mark_all_read")}</button>
      </div>
      {notifications.map((notif: Notification) => (
        <div
          key={notif.id}
          className={`card bg-base-100 border border-base-300 p-4 ${
            !notif.is_read ? "border-primary" : ""
          }`}
        >
          <div className="flex justify-between items-start">
            <div>
              <p className="font-medium">{notif.content}</p>
              <p className="text-sm text-base-content/60 mt-1">
                {notif.content}
              </p>
              <p className="text-xs text-base-content/40 mt-2">
                {new Date(notif.created_at).toLocaleString()}
              </p>
            </div>
            {!notif.is_read && (
              <button className="btn btn-xs btn-primary">
                {t("mark_as_read")}
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
