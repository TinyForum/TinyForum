import { AnnouncementFormValues } from "@/shared/type/announcement.type";
import { fromDateTimeLocal } from "../format/fromDateTimeLocal";
import { CreateAnnouncementPayload } from "@/shared/api/types/announcement.model";

// 转换表单提交数据
export const formValuesToPayload = (
  values: AnnouncementFormValues,
): CreateAnnouncementPayload => {
  return {
    title: values.title,
    content: values.content,
    summary: values.summary || "",
    cover: values.cover || "",
    type: values.type,
    status: values.status,
    is_pinned: values.is_pinned,
    is_global: values.is_global,
    board_id: values.is_global ? null : values.board_id || null,
    published_at: fromDateTimeLocal(values.published_at),
    expired_at: fromDateTimeLocal(values.expired_at),
  };
};
