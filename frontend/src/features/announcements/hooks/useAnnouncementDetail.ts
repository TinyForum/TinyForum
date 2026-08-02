// hooks/useAnnouncementDetail.ts
import { useState, useCallback, useEffect } from "react";
import { announcementApi } from "@/shared/api/modules/announcements";
import { AnnouncementDO } from "@/shared/api/types/announcement.model.do";

// 错误响应类型
interface ErrorResponse {
  response?: {
    status?: number;
    data?: {
      message?: string;
    };
  };
  message?: string;
}

// 公告详情 Hook（id 为 null 时视为无效 ID，不做请求）
export function useAnnouncementDetail(id: number | null) {
  const [announcement, setAnnouncement] = useState<AnnouncementDO | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    if (id === null) {
      setError("无效的公告ID");
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await announcementApi.getById(id);
      if (response.data.code === 0) {
        if (response.data.data) {
          setAnnouncement(response.data.data);
        }
      } else {
        setError(response.data.message || "公告不存在");
      }
    } catch (err: unknown) {
      console.error("Failed to load announcement:", err);
      const e = err as ErrorResponse;
      if (e.response?.status === 404) {
        setError("公告不存在");
      } else {
        setError("加载失败，请稍后重试");
      }
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  return { announcement, loading, error };
}
