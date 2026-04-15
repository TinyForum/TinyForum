import { announcementApi } from "@/lib/api";
import { Announcement, AnnouncementListResponse, CreateAnnouncementPayload } from "@/lib/api/modules/announcements";
// import { Announcement } from "@/type/admin.types";
import { useState, useEffect } from "react";
// import { announcementApi } from "@/api/announcements";

export function useAnnouncements() {
  const [announcementsList, setAnnouncementsList] = useState<Announcement[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 获取公告列表
  const fetchAnnouncements = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await announcementApi.list();
      setAnnouncementsList(response.data.data.list || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "获取公告失败");
    } finally {
      setIsLoading(false);
    }
  };

  // 创建公告
  const createAnnouncement = async (data: CreateAnnouncementPayload) => {
    try {
      const response = await announcementApi.create(data);
      await fetchAnnouncements(); // 刷新列表
      return response.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "创建公告失败");
      return null;
    }
  };

  // 更新公告
  const updateAnnouncement = async (id: number, data: Partial<Announcement>) => {
    try {
      const response = await announcementApi.update(id, data);
      await fetchAnnouncements(); // 刷新列表
      return response.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "更新公告失败");
      return null;
    }
  };

  // 删除公告
  const deleteAnnouncement = async (id: number) => {
    try {
      await announcementApi.delete(id);
      await fetchAnnouncements(); // 刷新列表
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "删除公告失败");
      return false;
    }
  };

  // 初始加载
  useEffect(() => {
    fetchAnnouncements();
  }, []);

  return {
    announcementsList,
    isLoading,
    error,
    fetchAnnouncements,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
  };
}