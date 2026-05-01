import { AnnouncementDO } from "@/shared/api/types/announcement.model";
import { AnnouncementFormValues } from "@/shared/type/announcement.type";
import { toDateTimeLocal } from "./toDateTimeLocal";

// 转换公告数据为表单数据
export const announcementToFormValues = (
  announcement: AnnouncementDO,
): AnnouncementFormValues => {
  return {
    title: announcement.title,
    content: announcement.content,
    summary: announcement.summary || "",
    cover: announcement.cover || "",
    type: announcement.type,
    is_pinned: announcement.is_pinned,
    status: announcement.status,
    is_global: announcement.is_global,
    board_id: announcement.board_id || null,
    published_at: toDateTimeLocal(announcement.published_at),
    expired_at: toDateTimeLocal(announcement.expired_at),
  };
};
