// hooks/useNotifications.ts
import { useState, useEffect, useCallback } from "react";
import { toast } from "react-hot-toast";
import { notificationApi } from "@/shared/api";
import { Notification } from "@/shared/api/types";

interface UseNotificationsOptions {
  pageSize?: number;
  autoLoad?: boolean;
}

export function useNotifications(options: UseNotificationsOptions = {}) {
  const { pageSize = 10, autoLoad = true } = options;
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [total, setTotal] = useState<number>(0);

  const loadNotifications = useCallback(
    async (page: number = 1) => {
      setLoading(true);
      try {
        const response = await notificationApi.list({
          page,
          page_size: pageSize,
        });
        const { code, data, message } = response.data;
        if (code === 0 && data) {
          const {
            list,
            page: resPage,
            page_size: resPageSize,
            total: resTotal,
          } = data;
          setNotifications(list as Notification[]);
          setTotal(resTotal);
          setCurrentPage(resPage);
          // 判断是否还有更多：当前页 * 每页条数 < 总数
          setHasMore(resPage * resPageSize < resTotal);
        } else {
          toast.error(message || "加载通知失败");
        }
      } catch (error) {
        console.error("Failed to load notifications:", error);
        toast.error("加载通知失败");
      } finally {
        setLoading(false);
      }
    },
    [pageSize],
  );

  const refresh = useCallback(() => {
    setCurrentPage(1);
    loadNotifications(1);
  }, [loadNotifications]);

  const loadMore = useCallback(() => {
    if (!loading && hasMore) {
      loadNotifications(currentPage + 1);
    }
  }, [loading, hasMore, currentPage, loadNotifications]);

  const markAsRead = useCallback(async (notificationId: number) => {
    try {
      const response = await notificationApi.markRead(notificationId);
      if (response.data.code === 0) {
        setNotifications((prev) =>
          prev.map((n) =>
            n.id === notificationId ? { ...n, is_read: true } : n,
          ),
        );
      } else {
        toast.error(response.data.message || "标记失败");
      }
    } catch (error) {
      console.error("Failed to mark as read:", error);
      toast.error("操作失败");
    }
  }, []);

  const markAllAsRead = useCallback(async () => {
    try {
      const response = await notificationApi.markAllRead();
      if (response.data.code === 0) {
        setNotifications((prev) => prev.map((n) => ({ ...n, is_read: true })));
      } else {
        toast.error(response.data.message || "操作失败");
      }
    } catch (error) {
      console.error("Failed to mark all as read:", error);
      toast.error("操作失败");
    }
  }, []);

  useEffect(() => {
    if (autoLoad) {
      refresh();
    }
  }, [autoLoad, refresh]);

  return {
    notifications,
    loading,
    hasMore,
    total,
    refresh,
    loadMore,
    markAsRead,
    markAllAsRead,
  };
}
