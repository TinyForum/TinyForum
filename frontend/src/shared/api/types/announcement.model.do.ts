// 公告数据结构
export interface AnnouncementDO {
  id: number;
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: AnnouncementType;
  status: CreateAnnouncementStatus;
  is_pinned: boolean;
  is_global: boolean;
  board_id: number | null;
  published_at: string | null;
  expired_at: string | null;
  view_count: number;
  created_by: number;
  updated_by: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;

  // 关联数据
  board?: {
    id: number;
    name: string;
    slug: string;
  } | null;
  creator?: {
    id: number;
    username: string;
    avatar_url?: string;
  } | null;
}

export enum AnnouncementType {
  Normal = "normal", // 普通
  Important = "important", // 重要
  Emergency = "emergency", // 紧急
  Event = "event", // 活动
}
export enum AnnouncementStatus {
  /** 仅用于查询：所有状态 */
  All = "all",
  /** 草稿 */
  Draft = "draft",
  /** 已发布 */
  Published = "published",
  /** 已归档 */
  Archived = "archived",
}

export type CreateAnnouncementStatus = Exclude<
  AnnouncementStatus,
  AnnouncementStatus.All
>;

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
