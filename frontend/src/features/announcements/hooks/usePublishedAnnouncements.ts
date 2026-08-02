// hooks/usePublishedAnnouncements.ts
import { useQuery } from "@tanstack/react-query";
import { announcementApi } from "@/shared/api/modules/announcements";
import { announcementKeys } from "./useAnnouncementKeys";
import { AnnouncementStatus } from "@/shared/api/types/announcement.model.do";

// 已发布公告列表查询 Hook（分页）
export function usePublishedAnnouncements(page: number, pageSize: number) {
  const params = {
    page,
    page_size: pageSize,
    status: AnnouncementStatus.Published,
  };

  return useQuery({
    queryKey: announcementKeys.list(params),
    queryFn: async () => {
      const res = await announcementApi.list(params);
      return res.data.data;
    },
  });
}
