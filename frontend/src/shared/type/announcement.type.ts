import {
  AnnouncementType,
  CreateAnnouncementStatus,
} from "../api/types/announcement.model";

// 类型定义 - 让可选字段真正可选
export interface AnnouncementFormValues {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  type: AnnouncementType;
  is_pinned: boolean;
  status: CreateAnnouncementStatus;
  is_global: boolean;
  board_id?: number | null;
  published_at?: string | null;
  expired_at?: string | null;
}
