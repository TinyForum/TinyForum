// hooks/useUnreadCount.ts
import { useState, useEffect, useCallback } from "react";
import { toast } from "react-hot-toast";
import { notificationApi } from "@/shared/api";

export function useUnreadCount() {
  const [unreadCount, setUnreadCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchUnreadCount = useCallback(async () => {
    setLoading(true);
    try {
      const response = await notificationApi.unreadCount();
      if (response.data.code === 0) {
        setUnreadCount(response.data.data?.count ?? 0);
      } else {
        toast.error(response.data.message || "获取未读数量失败");
      }
    } catch (error) {
      console.error("Failed to fetch unread count:", error);
      toast.error("获取未读数量失败");
    } finally {
      setLoading(false);
    }
  }, []);

  // 手动增加/减少计数（用于本地乐观更新）
  const incrementUnread = useCallback(() => {
    setUnreadCount((prev) => prev + 1);
  }, []);

  const decrementUnread = useCallback(() => {
    setUnreadCount((prev) => Math.max(0, prev - 1));
  }, []);

  const resetUnread = useCallback(() => {
    setUnreadCount(0);
  }, []);

  useEffect(() => {
    fetchUnreadCount();
    // 可选：轮询或通过事件监听更新
  }, [fetchUnreadCount]);

  return {
    unreadCount,
    loading,
    refetch: fetchUnreadCount,
    incrementUnread,
    decrementUnread,
    resetUnread,
  };
}
