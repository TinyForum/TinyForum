import { useState, useEffect, useCallback } from "react";
import {
  announcementApi,
  type Announcement,
} from "@/shared/api/modules/announcements";
import toast from "react-hot-toast";
import { ApiResponse } from "@/shared/api/types";

// ============ 用于单个公告的 Hook ============
interface UseAdminAnnouncementOptions {
  autoLoad?: boolean;
}

interface UseAnnouncementReturn {
  announcement: Announcement | null;
  loading: boolean;
  fetch: (id: number) => Promise<Announcement | null>;
  clear: () => void;
}

export function useAdminAnnouncement(
  id?: number,
  options?: UseAdminAnnouncementOptions,
): UseAnnouncementReturn {
  const { autoLoad = true } = options || {};
  const [announcement, setAnnouncement] = useState<Announcement | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const fetch = useCallback(
    async (announcementId: number): Promise<Announcement | null> => {
      setLoading(true);
      try {
        const response: { data: ApiResponse<Announcement> } =
          await announcementApi.getById(announcementId);

        if (response.data.code === 0 && response.data.data) {
          setAnnouncement(response.data.data);
          return response.data.data;
        } else {
          toast.error(response.data.message || "获取公告详情失败");
          return null;
        }
      } catch (error) {
        console.error("获取公告详情失败:", error);
        toast.error("获取公告详情失败，请稍后重试");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const clear = useCallback(() => {
    setAnnouncement(null);
  }, []);

  useEffect(() => {
    if (autoLoad && id) {
      fetch(id);
    }
  }, [autoLoad, id, fetch]);

  return {
    announcement,
    loading,
    fetch,
    clear,
  };
}
