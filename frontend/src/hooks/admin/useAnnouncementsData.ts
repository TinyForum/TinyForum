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

// ============ 配置选项 ============
interface UseAnnouncementsDataOptions {
  enabled?: boolean;
  defaultParams?: AnnouncementListParams;
  autoLoadPinned?: boolean;
}

// ============ Hook 返回值类型 ============
interface UseAnnouncementsDataReturn {
  // 数据状态
  announcements: Announcement[];
  pinnedAnnouncements: Announcement[];
  total: number;
  loading: boolean;
  submitting: boolean;
  refreshing: boolean;
  isLoading:boolean;

  // 分页
  page: number;
  pageSize: number;

  // 操作方法
  fetchAnnouncements: (params?: AnnouncementListParams) => Promise<void>;
  fetchPinnedAnnouncements: (boardId?: number) => Promise<void>;
  getAnnouncementById: (id: number) => Promise<Announcement | null>;
  createAnnouncement: (
    params: CreateAnnouncementPayload
  ) => Promise<Announcement | null>;
  updateAnnouncement: (
    id: number,
    params: UpdateAnnouncementPayload
  ) => Promise<Announcement | null>;
  deleteAnnouncement: (id: number) => Promise<boolean>;
  publishAnnouncement: (id: number) => Promise<boolean>;
  archiveAnnouncement: (id: number) => Promise<boolean>;
  pinAnnouncement: (id: number, pinned: boolean) => Promise<boolean>;

  // 状态设置
  setPage: (page: number) => void;
  setPageSize: (pageSize: number) => void;
  setFilters: (filters: AnnouncementListParams) => void;
  resetFilters: () => void;
}

// ============ Hook 实现 ============
export function useAnnouncementsData(
  enabled: boolean = true,
  options?: UseAnnouncementsDataOptions
): UseAnnouncementsDataReturn {
  const { defaultParams = {}, autoLoadPinned = true } = options || {};

  // 数据状态
  const [announcements, setAnnouncements] = useState<Announcement[]>([]);
  const [pinnedAnnouncements, setPinnedAnnouncements] = useState<Announcement[]>(
    []
  );
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [refreshing, setRefreshing] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // 分页状态
  const [page, setPage] = useState<number>(defaultParams.page || 1);
  const [pageSize, setPageSize] = useState<number>(defaultParams.page_size || 20);
  const [filters, setFilters] = useState<AnnouncementListParams>(() => ({
    page,
    page_size: pageSize,
    ...defaultParams,
  }));

  // 获取公告列表
  const fetchAnnouncements = useCallback(
    async (params?: AnnouncementListParams) => {
      if (!enabled) return;

      const isRefresh = params?.page === 1 || params?.page === undefined;
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }

      try {
        const queryParams = { ...filters, ...params, page, page_size: pageSize };
        // 移除 undefined 值
        Object.keys(queryParams).forEach((key) => {
          if (queryParams[key as keyof AnnouncementListParams] === undefined) {
            delete queryParams[key as keyof AnnouncementListParams];
          }
        });

        const response = await announcementApi.list(queryParams);

        if (response.data.code === 0 && response.data.data) {
          setAnnouncements(response.data.data.list || []);
          setTotal(response.data.data.total || 0);
        } else {
          toast.error(response.data.message || "获取公告列表失败");
        }
        setIsLoading(false);
      } catch (error) {
        console.error("获取公告列表失败:", error);
        toast.error("获取公告列表失败，请稍后重试");
      } finally {
        setLoading(false);
        setRefreshing(false);
      }
    },
    [enabled, filters, page, pageSize]
  );

  // 获取置顶公告
  const fetchPinnedAnnouncements = useCallback(
    async (boardId?: number) => {
      if (!enabled) return;

      try {
        const response = await announcementApi.getPinned(boardId);

        if (response.data.code === 0) {
          setPinnedAnnouncements(response.data.data || []);
        }
      } catch (error) {
        console.error("获取置顶公告失败:", error);
      }
    },
    [enabled]
  );

  // 根据 ID 获取公告详情
  const getAnnouncementById = useCallback(
    async (id: number): Promise<Announcement | null> => {
      try {
        const response = await announcementApi.getById(id);

        if (response.data.code === 0) {
          return response.data.data;
        } else {
          toast.error(response.data.message || "获取公告详情失败");
          return null;
        }
      } catch (error) {
        console.error("获取公告详情失败:", error);
        toast.error("获取公告详情失败，请稍后重试");
        return null;
      }
    },
    []
  );

  // 创建公告
  const createAnnouncement = useCallback(
    async (params: CreateAnnouncementPayload): Promise<Announcement | null> => {
      setSubmitting(true);
      try {
        const response = await announcementApi.create(params);

        console.log("创建公告: ",response);
        if (response.data.code === 0) {
          toast.success("创建公告成功");
          await fetchAnnouncements();
          return response.data.data;
        } else {
          toast.error(response.data.message || "创建公告失败");
          return null;
        }
      } catch (error) {
        console.error("创建公告失败:", error);
        toast.error("创建公告失败，请稍后重试");
        return null;
      } finally {
        setSubmitting(false);
      }
    },
    [fetchAnnouncements]
  );

  // 更新公告
  const updateAnnouncement = useCallback(
    async (id: number, params: UpdateAnnouncementPayload): Promise<Announcement | null> => {
      setSubmitting(true);
      try {
        const response = await announcementApi.update(id, params);

        if (response.data.code === 0) {
          toast.success("更新公告成功");
          await fetchAnnouncements();
          // 如果影响置顶状态，刷新置顶列表
          if (params.is_pinned !== undefined) {
            await fetchPinnedAnnouncements();
          }
          return response.data.data;
        } else {
          toast.error(response.data.message || "更新公告失败");
          return null;
        }
      } catch (error) {
        console.error("更新公告失败:", error);
        toast.error("更新公告失败，请稍后重试");
        return null;
      } finally {
        setSubmitting(false);
      }
    },
    [fetchAnnouncements, fetchPinnedAnnouncements]
  );

  // 删除公告
  const deleteAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      try {
        const response = await announcementApi.delete(id);

        if (response.data.code === 0) {
          toast.success("删除公告成功");
          await fetchAnnouncements();
          await fetchPinnedAnnouncements();
          return true;
        } else {
          toast.error(response.data.message || "删除公告失败");
          return false;
        }
      } catch (error) {
        console.error("删除公告失败:", error);
        toast.error("删除公告失败，请稍后重试");
        return false;
      }
    },
    [fetchAnnouncements, fetchPinnedAnnouncements]
  );

  // 发布公告
  const publishAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      setSubmitting(true);
      try {
        const response = await announcementApi.publish(id);

        if (response.data.code === 0) {
          toast.success("发布公告成功");
          await fetchAnnouncements();
          await fetchPinnedAnnouncements();
          return true;
        } else {
          toast.error(response.data.message || "发布公告失败");
          return false;
        }
      } catch (error) {
        console.error("发布公告失败:", error);
        toast.error("发布公告失败，请稍后重试");
        return false;
      } finally {
        setSubmitting(false);
      }
    },
    [fetchAnnouncements, fetchPinnedAnnouncements]
  );

  // 归档公告
  const archiveAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      setSubmitting(true);
      try {
        const response = await announcementApi.archive(id);

        if (response.data.code === 0) {
          toast.success("归档公告成功");
          await fetchAnnouncements();
          return true;
        } else {
          toast.error(response.data.message || "归档公告失败");
          return false;
        }
      } catch (error) {
        console.error("归档公告失败:", error);
        toast.error("归档公告失败，请稍后重试");
        return false;
      } finally {
        setSubmitting(false);
      }
    },
    [fetchAnnouncements]
  );

  // 置顶/取消置顶
  const pinAnnouncement = useCallback(
    async (id: number, pinned: boolean): Promise<boolean> => {
      setSubmitting(true);
      try {
        const response = await announcementApi.pin(id, pinned);

        if (response.data.code === 0) {
          toast.success(pinned ? "置顶成功" : "取消置顶成功");
          await fetchAnnouncements();
          await fetchPinnedAnnouncements();
          return true;
        } else {
          toast.error(response.data.message || "操作失败");
          return false;
        }
      } catch (error) {
        console.error("置顶操作失败:", error);
        toast.error("操作失败，请稍后重试");
        return false;
      } finally {
        setSubmitting(false);
      }
    },
    [fetchAnnouncements, fetchPinnedAnnouncements]
  );

  // 重置筛选条件
  const resetFilters = useCallback(() => {
    const newFilters = {
      page: 1,
      page_size: pageSize,
    };
    setFilters(newFilters);
    setPage(1);
  }, [pageSize]);

  // 监听分页和筛选变化，重新加载数据
  useEffect(() => {
    if (enabled) {
      fetchAnnouncements();
    }
  }, [enabled, page, pageSize, filters, fetchAnnouncements]);

  // 初始加载置顶公告
  useEffect(() => {
    if (enabled && autoLoadPinned) {
      fetchPinnedAnnouncements();
    }
  }, [enabled, autoLoadPinned, fetchPinnedAnnouncements]);

  return {
    // 数据状态
    announcements,
    pinnedAnnouncements,
    total,
    loading,
    submitting,
    refreshing,

    // 分页
    page,
    pageSize,


    // 操作方法
    fetchAnnouncements,
    fetchPinnedAnnouncements,
    getAnnouncementById,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
    publishAnnouncement,
    archiveAnnouncement,
    pinAnnouncement,

    // 状态设置
    setPage,
    setPageSize,
    setFilters,
    resetFilters,
    isLoading

  };
}

