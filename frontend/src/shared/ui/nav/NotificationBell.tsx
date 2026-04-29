"use client";

import Link from "next/link";
import { Bell, X } from "lucide-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { notificationApi } from "@/shared/api";
import { useTranslations } from "next-intl";
import { ApiResponse } from "@/shared/api/types";
import { Transition, Dialog } from "@headlessui/react";
import { Fragment, useState, useRef, useEffect } from "react";
import { createPortal } from "react-dom";

interface Notification {
  id: number;
  content: string;
  is_read: boolean;
  created_at: string;
  user_id: number;
  sender_id?: number;
  type: string;
  target_id?: number;
  target_type: string;
}

interface NotificationListResponse {
  list: Notification[];
  total: number;
  page: number;
  page_size: number;
}

interface NotificationBellProps {
  unreadCount: number;
  compact?: boolean;
}

export default function NotificationBell({
  unreadCount,
  compact = false,
}: NotificationBellProps) {
  const t = useTranslations("Notifications");
  const queryClient = useQueryClient();
  const [selectedNotification, setSelectedNotification] =
    useState<Notification | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, right: 0 });
  const [isMounted, setIsMounted] = useState(false); // SSR 守卫
  const buttonRef = useRef<HTMLButtonElement>(null);
  const popoverRef = useRef<HTMLDivElement>(null);

  // SSR 守卫：客户端挂载后才允许 Portal 渲染
  useEffect(() => {
    setIsMounted(true);
  }, []);

  const { data: previewData, refetch } = useQuery({
    queryKey: ["notifications", "preview"],
    queryFn: () =>
      notificationApi
        .list({ page: 1, page_size: 5 })
        .then(
          (r: { data: ApiResponse<NotificationListResponse> }) => r.data.data,
        ),
    enabled: false,
  });

  const markAsReadMutation = useMutation({
    mutationFn: (notificationId: number) =>
      notificationApi.markRead(notificationId),
    onSuccess: () => {
      queryClient.setQueryData(["notifications", "preview"], (oldData: any) => {
        if (!oldData) return oldData;
        return {
          ...oldData,
          list: oldData.list.map((n: Notification) =>
            n.id === selectedNotification?.id ? { ...n, is_read: true } : n,
          ),
        };
      });
      queryClient.invalidateQueries({ queryKey: ["unreadCount"] });
    },
  });

  const notifications: Notification[] = previewData?.list || [];

  const updatePosition = () => {
    if (buttonRef.current) {
      const rect = buttonRef.current.getBoundingClientRect();
      setPosition({
        top: rect.bottom + window.scrollY + 8,
        right: window.innerWidth - rect.right + window.scrollX,
      });
    }
  };

  const handleOpen = () => {
    if (!isOpen) {
      refetch();
      updatePosition();
      setIsOpen(true);
    } else {
      setIsOpen(false);
    }
  };

  const handleClose = () => setIsOpen(false);

  // 用 mousedown + ref 判断点击是否在弹出层或按钮内，替代透明遮罩
  useEffect(() => {
    if (!isOpen) return;

    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as Node;
      const inPopover = popoverRef.current?.contains(target);
      const inButton = buttonRef.current?.contains(target);
      if (!inPopover && !inButton) {
        handleClose();
      }
    };

    // mousedown 比 click 更早触发，体验更好，也不会被子元素 stopPropagation 影响
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isOpen]);

  // 滚动 / resize 时同步更新位置
  useEffect(() => {
    if (!isOpen) return;
    window.addEventListener("scroll", updatePosition, true); // capture 模式，捕获所有滚动容器
    window.addEventListener("resize", updatePosition);
    return () => {
      window.removeEventListener("scroll", updatePosition, true);
      window.removeEventListener("resize", updatePosition);
    };
  }, [isOpen]);

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    if (diffMins < 1) return "刚刚";
    if (diffMins < 60) return `${diffMins}分钟前`;
    if (diffHours < 24) return `${diffHours}小时前`;
    if (diffDays < 7) return `${diffDays}天前`;
    return date.toLocaleDateString();
  };

  const truncateContent = (content: string, maxLength = 50) => {
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + "...";
  };

  const handleNotificationClick = (notification: Notification) => {
    setSelectedNotification(notification);
    setIsModalOpen(true);
    handleClose();
    if (!notification.is_read) {
      markAsReadMutation.mutate(notification.id);
    }
  };

  const popoverContent = (
    <div
      ref={popoverRef}
      className="absolute z-50 w-[420px] origin-top-right"
      style={{
        top: `${position.top}px`,
        right: `${position.right}px`,
      }}
    >
      <div className="overflow-hidden rounded-lg shadow-xl ring-1 ring-black ring-opacity-5 bg-white dark:bg-gray-800">
        {/* 头部 */}
        <div className="p-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
          <h3 className="font-semibold text-gray-900 dark:text-gray-100">
            {t("title")}
          </h3>
          {unreadCount > 0 && (
            <button
              onClick={handleClose}
              className="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 hover:underline"
            >
              {t("mark_all_read")}
            </button>
          )}
        </div>

        {/* 通知列表 */}
        <div className="max-h-96 overflow-y-auto">
          {notifications.length === 0 ? (
            <div className="p-8 text-center text-gray-500 dark:text-gray-400">
              <Bell className="w-8 h-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">{t("no_notifications")}</p>
            </div>
          ) : (
            notifications.map((notif: Notification) => (
              <div
                key={notif.id}
                onClick={() => handleNotificationClick(notif)}
                className={`
                  cursor-pointer p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50
                  border-b border-gray-100 dark:border-gray-700 last:border-0
                  transition-colors
                  ${!notif.is_read ? "bg-blue-50/50 dark:bg-blue-900/20" : ""}
                `}
              >
                <div className="flex items-start gap-3">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-gray-900 dark:text-gray-100 leading-relaxed break-words">
                      {truncateContent(notif.content, 50)}
                    </p>
                    <div className="flex items-center gap-2 mt-1.5">
                      <span className="text-xs text-gray-400 dark:text-gray-500">
                        {formatTime(notif.created_at)}
                      </span>
                      {!notif.is_read && (
                        <span className="w-1.5 h-1.5 bg-blue-500 rounded-full" />
                      )}
                    </div>
                  </div>
                  {!notif.is_read && (
                    <div className="w-2 h-2 bg-blue-500 rounded-full mt-1.5 flex-shrink-0" />
                  )}
                </div>
              </div>
            ))
          )}
        </div>

        {/* 底部 */}
        <div className="p-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <Link
            href="/notifications"
            onClick={handleClose}
            className="block text-center text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 hover:underline py-1"
          >
            {t("view_all")}
          </Link>
        </div>
      </div>
    </div>
  );

  return (
    <>
      {/* 触发器按钮 */}
      <button
        ref={buttonRef}
        onClick={handleOpen}
        className={`btn btn-ghost btn-sm relative focus:outline-none ${compact ? "btn-square" : "btn-circle"}`}
        aria-label="通知"
      >
        <Bell className="w-5 h-5" />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 badge badge-error badge-xs min-w-[16px] h-4 text-[10px] animate-pulse">
            {unreadCount > 99 ? "99+" : unreadCount}
          </span>
        )}
      </button>

      {/* 弹出层：isMounted 守卫替代 typeof document !== "undefined" */}
      {isMounted && isOpen && createPortal(popoverContent, document.body)}

      {/* 通知详情模态框：同样用 Portal 挂到 body，防止被父级 stacking context 遮挡 */}
      {isMounted &&
        createPortal(
          <Transition appear show={isModalOpen} as={Fragment}>
            <Dialog
              className="relative z-[99999]"
              onClose={() => setIsModalOpen(false)}
            >
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0"
                enterTo="opacity-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" />
              </Transition.Child>

              <div className="fixed inset-0 flex items-center justify-center p-4">
                <Transition.Child
                  as={Fragment}
                  enter="ease-out duration-300"
                  enterFrom="opacity-0 scale-95"
                  enterTo="opacity-100 scale-100"
                  leave="ease-in duration-200"
                  leaveFrom="opacity-100 scale-100"
                  leaveTo="opacity-0 scale-95"
                >
                  <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white dark:bg-gray-800 p-6 shadow-xl transition-all">
                    <div className="flex items-center justify-between mb-4">
                      <Dialog.Title className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                        通知详情
                      </Dialog.Title>
                      <button
                        onClick={() => setIsModalOpen(false)}
                        className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors"
                      >
                        <X className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                      </button>
                    </div>

                    {selectedNotification && (
                      <div className="space-y-4">
                        <div className="prose prose-sm max-w-none dark:prose-invert">
                          <p className="text-base leading-relaxed whitespace-pre-wrap break-words text-gray-700 dark:text-gray-300">
                            {selectedNotification.content}
                          </p>
                        </div>
                        <div className="flex items-center justify-between pt-3 border-t border-gray-200 dark:border-gray-700">
                          <div className="text-xs text-gray-400 dark:text-gray-500 space-y-1">
                            <div>
                              发送时间：{selectedNotification.created_at}
                            </div>
                            {selectedNotification.type && (
                              <div>类型：{selectedNotification.type}</div>
                            )}
                          </div>
                          {selectedNotification.target_id && (
                            <Link
                              href={`/${selectedNotification.target_type}/${selectedNotification.target_id}`}
                              onClick={() => setIsModalOpen(false)}
                              className="px-3 py-1.5 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                            >
                              查看详情
                            </Link>
                          )}
                        </div>
                      </div>
                    )}
                  </Dialog.Panel>
                </Transition.Child>
              </div>
            </Dialog>
          </Transition>,
          document.body,
        )}
    </>
  );
}
