"use client";

import Link from "next/link";
import { Bell } from "lucide-react";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { notificationApi } from "@/lib/api";

interface NotificationBellProps {
  unreadCount: number;
}

export default function NotificationBell({ unreadCount }: NotificationBellProps) {
  const [isOpen, setIsOpen] = useState(false);

  // 获取最新通知预览
  const { data: previewData } = useQuery({
    queryKey: ["notifications", "preview"],
    queryFn: () => notificationApi.list({ page: 1, page_size: 5 }).then((r) => r.data.data),
    enabled: isOpen,
  });

  const notifications = previewData?.list || [];

  return (
    <div className="dropdown dropdown-end">
      <div
        tabIndex={0}
        role="button"
        className="btn btn-ghost btn-sm btn-circle relative"
        onClick={() => setIsOpen(!isOpen)}
      >
        <Bell className="w-5 h-5" />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 badge badge-error badge-xs min-w-[16px] h-4 text-[10px] animate-bounce">
            {unreadCount > 99 ? "99+" : unreadCount}
          </span>
        )}
      </div>
      
      <div className="dropdown-content mt-2 w-80 bg-base-100 rounded-lg shadow-xl border border-base-200 z-50">
        <div className="p-3 border-b border-base-200 flex justify-between items-center">
          <h3 className="font-semibold">通知中心</h3>
          {unreadCount > 0 && (
            <Link href="/notifications" className="text-xs text-primary hover:underline">
              标记全部已读
            </Link>
          )}
        </div>
        
        <div className="max-h-96 overflow-y-auto">
          {notifications.length === 0 ? (
            <div className="p-8 text-center text-base-content/50">
              <Bell className="w-8 h-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">暂无通知</p>
            </div>
          ) : (
            notifications.map((notif: any) => (
              <Link
                key={notif.id}
                href={`/notifications`}
                className={`block p-3 hover:bg-base-200 transition-colors border-b border-base-200 last:border-0 ${
                  !notif.is_read ? "bg-primary/5" : ""
                }`}
              >
                <p className="text-sm line-clamp-2">{notif.content}</p>
                <span className="text-xs text-base-content/40 mt-1 block">
                  {new Date(notif.created_at).toLocaleDateString()}
                </span>
              </Link>
            ))
          )}
        </div>
        
        <div className="p-2 border-t border-base-200">
          <Link
            href="/notifications"
            className="block text-center text-sm text-primary hover:underline py-1"
          >
            查看所有通知
          </Link>
        </div>
      </div>
    </div>
  );
}