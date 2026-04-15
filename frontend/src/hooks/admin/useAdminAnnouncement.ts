import { useState, useEffect, useCallback } from "react";
// import { toast } from "antd";
import {
  announcementApi,
  type Announcement,
  type AnnouncementListParams,
  type AnnouncementListResponse,
  type CreateAnnouncementPayload,
  type UpdateAnnouncementPayload,
  type AnnouncementStatus,
} from "@/lib/api/modules/announcements";
import toast from "react-hot-toast";


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
  options?: UseAdminAnnouncementOptions
): UseAnnouncementReturn {
  const { autoLoad = true } = options || {};
  const [announcement, setAnnouncement] = useState<Announcement | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const fetch = useCallback(async (announcementId: number): Promise<Announcement | null> => {
    setLoading(true);
    try {
      const response = await announcementApi.getById(announcementId);
      if (response.data.code === 0) {
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
  }, []);

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