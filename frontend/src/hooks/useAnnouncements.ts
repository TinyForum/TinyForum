// hooks/useAnnouncements.ts
import { announcementApi } from "@/lib/api";
import {
  Announcement,
  AnnouncementListParams,
  AnnouncementListResponse,
} from "@/lib/api/modules/announcements";
import { ApiResponse } from "@/lib/api/types";
import { useState, useEffect, useCallback } from "react";

/**
 * 公告 Hook - 用于前台展示（左侧边栏、公告页面等）
 *
 * 说明：
 * - 所有用户（包括管理员）在前台看到的都是已发布且未过期的公告
 * - 管理员查看所有公告应该在后台管理面板使用 useAdminAnnouncements
 *
 * @param boardId - 可选，板块ID（版主用于获取板块公告）
 */
export function useAnnouncements(boardId?: number) {
  const [announcementsList, setAnnouncementsList] = useState<Announcement[]>(
    [],
  );
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAnnouncements = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    setError(null);
    try {
      // 使用模块中已定义的类型
      const params: AnnouncementListParams = {
        page: 1,
        page_size: 20,
        status: "published" as const, // 使用 const assertion 匹配 AnnouncementStatus 类型
      };

      // 如果传入了板块ID，获取该板块的公告
      if (boardId) {
        params.board_id = boardId;
        params.is_global = false;
      } else {
        // 默认获取全局公告
        params.is_global = true;
      }

      const response: { data: ApiResponse<AnnouncementListResponse> } =
        await announcementApi.list(params);

      console.log("前台公告列表:", response.data.data?.list);

      // 修复：response.data.data 可能为 undefined
      setAnnouncementsList(response.data.data?.list || []);
    } catch (err: unknown) {
      console.error("获取公告失败:", err);
      setError(err instanceof Error ? err.message : "获取公告失败");
    } finally {
      setIsLoading(false);
    }
  }, [boardId]);

  useEffect((): void => {
    fetchAnnouncements();
  }, [fetchAnnouncements]);

  return {
    announcementsList,
    isLoading,
    error,
    refetch: fetchAnnouncements,
  };
}
