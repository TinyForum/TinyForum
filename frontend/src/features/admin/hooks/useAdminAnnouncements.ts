// hooks/admin/useAdminAnnouncements.ts
import { announcementApi } from "@/shared/api";
import {
  Announcement,
  CreateAnnouncementPayload,
  UpdateAnnouncementPayload,
} from "@/shared/api/modules/announcements";
import { useState, useEffect, useCallback } from "react";
import toast from "react-hot-toast";

interface UseAdminAnnouncementsReturn {
  // 数据
  announcements: Announcement[];
  pinnedAnnouncements: Announcement[];
  total: number;

  // 状态
  isLoading: boolean;
  isSubmitting: boolean;

  // 分页
  page: number;
  pageSize: number;
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;

  // 操作方法
  refetch: () => Promise<void>;
  getAnnouncementById: (id: number) => Promise<Announcement | null>;
  createAnnouncement: (
    data: CreateAnnouncementPayload,
  ) => Promise<Announcement | null>;
  updateAnnouncement: (
    id: number,
    data: UpdateAnnouncementPayload,
  ) => Promise<Announcement | null>;
  deleteAnnouncement: (id: number) => Promise<boolean>;
  publishAnnouncement: (id: number) => Promise<boolean>;
  archiveAnnouncement: (id: number) => Promise<boolean>;
  pinAnnouncement: (id: number, isPinned: boolean) => Promise<boolean>;
}

export function useAdminAnnouncements(): UseAdminAnnouncementsReturn {
  const [announcements, setAnnouncements] = useState<Announcement[]>([]);
  const [pinnedAnnouncements, setPinnedAnnouncements] = useState<
    Announcement[]
  >([]);
  const [total, setTotal] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);

  // 获取公告列表
  const fetchAnnouncements = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await announcementApi.adminList({
        page,
        page_size: pageSize,
      });

      if (response.data.code === 0 || response.data.code === 0) {
        // 添加安全检查，确保 response.data.data 存在
        const data = response.data.data;
        const list = data?.list || [];
        setAnnouncements(list);
        setTotal(data?.total || 0);
        // 筛选置顶公告
        setPinnedAnnouncements(
          list.filter((ann: Announcement) => ann.is_pinned === true),
        );
      } else {
        toast.error(response.data.message || "获取公告列表失败");
      }
    } catch (error) {
      console.error("获取公告列表失败:", error);
      toast.error("获取公告列表失败");
    } finally {
      setIsLoading(false);
    }
  }, [page, pageSize]);

  // 根据 ID 获取公告详情
  const getAnnouncementById = useCallback(
    async (id: number): Promise<Announcement | null> => {
      try {
        const response = await announcementApi.getById(id);
        if (response.data.code === 0 || response.data.code === 0) {
          return response.data.data || null;
        }
        toast.error(response.data.message || "获取公告详情失败");
        return null;
      } catch (error) {
        console.error("获取公告详情失败:", error);
        toast.error("获取公告详情失败");
        return null;
      }
    },
    [],
  );

  // 创建公告
  const createAnnouncement = useCallback(
    async (data: CreateAnnouncementPayload): Promise<Announcement | null> => {
      setIsSubmitting(true);
      try {
        const response = await announcementApi.create(data);
        if (response.data.code === 0 || response.data.code === 0) {
          toast.success("创建公告成功");
          await fetchAnnouncements();
          return response.data.data || null;
        }
        toast.error(response.data.message || "创建公告失败");
        return null;
      } catch (error) {
        console.error("创建公告失败:", error);
        toast.error("创建公告失败");
        return null;
      } finally {
        setIsSubmitting(false);
      }
    },
    [fetchAnnouncements],
  );

  // 更新公告
  const updateAnnouncement = useCallback(
    async (
      id: number,
      data: UpdateAnnouncementPayload,
    ): Promise<Announcement | null> => {
      setIsSubmitting(true);
      try {
        const response = await announcementApi.update(id, data);
        if (response.data.code === 0 || response.data.code === 0) {
          toast.success("更新公告成功");
          await fetchAnnouncements();
          return response.data.data || null;
        }
        toast.error(response.data.message || "更新公告失败");
        return null;
      } catch (error) {
        console.error("更新公告失败:", error);
        toast.error("更新公告失败");
        return null;
      } finally {
        setIsSubmitting(false);
      }
    },
    [fetchAnnouncements],
  );

  // 删除公告
  const deleteAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      setIsSubmitting(true);
      try {
        const response = await announcementApi.delete(id);
        if (response.data.code === 0 || response.data.code === 0) {
          toast.success("删除公告成功");
          await fetchAnnouncements();
          return true;
        }
        toast.error(response.data.message || "删除公告失败");
        return false;
      } catch (error) {
        console.error("删除公告失败:", error);
        toast.error("删除公告失败");
        return false;
      } finally {
        setIsSubmitting(false);
      }
    },
    [fetchAnnouncements],
  );

  // 发布公告：设置 status 为 published，并设置 published_at 为当前时间
  const publishAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      setIsSubmitting(true);
      try {
        const response = await announcementApi.update(id, {
          status: "published",
          published_at: new Date().toISOString(),
        });
        if (response.data.code === 0 || response.data.code === 0) {
          toast.success("发布公告成功");
          await fetchAnnouncements();
          return true;
        }
        toast.error(response.data.message || "发布公告失败");
        return false;
      } catch (error) {
        console.error("发布公告失败:", error);
        toast.error("发布公告失败");
        return false;
      } finally {
        setIsSubmitting(false);
      }
    },
    [fetchAnnouncements],
  );

  // 归档公告：设置 expired_at 为当前时间之前（标记为过期）
  const archiveAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      setIsSubmitting(true);
      try {
        const response = await announcementApi.update(id, {
          expired_at: new Date().toISOString(),
        });
        if (response.data.code === 0 || response.data.code === 0) {
          toast.success("归档公告成功");
          await fetchAnnouncements();
          return true;
        }
        toast.error(response.data.message || "归档公告失败");
        return false;
      } catch (error) {
        console.error("归档公告失败:", error);
        toast.error("归档公告失败");
        return false;
      } finally {
        setIsSubmitting(false);
      }
    },
    [fetchAnnouncements],
  );

  // 置顶/取消置顶
  const pinAnnouncement = useCallback(
    async (id: number, isPinned: boolean): Promise<boolean> => {
      setIsSubmitting(true);
      try {
        const response = await announcementApi.update(id, {
          is_pinned: isPinned,
        });
        if (response.data.code === 0 || response.data.code === 0) {
          await fetchAnnouncements();
          return true;
        }
        toast.error(response.data.message || "操作失败");
        return false;
      } catch (error) {
        console.error("置顶操作失败:", error);
        toast.error("操作失败");
        return false;
      } finally {
        setIsSubmitting(false);
      }
    },
    [fetchAnnouncements],
  );

  // 初始加载
  useEffect(() => {
    fetchAnnouncements();
  }, [fetchAnnouncements]);

  return {
    // 数据
    announcements,
    pinnedAnnouncements,
    total,

    // 状态
    isLoading,
    isSubmitting,

    // 分页
    page,
    pageSize,
    setPage,
    setPageSize,

    // 操作方法
    refetch: fetchAnnouncements,
    getAnnouncementById,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
    publishAnnouncement,
    archiveAnnouncement,
    pinAnnouncement,
  };
}
