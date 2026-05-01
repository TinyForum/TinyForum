export enum AnnouncementType {
  Normal = 0, // 普通
  Important = 1, // 重要
  Emergency = 2, // 紧急
  Event = 3, // 活动
}
export enum AnnouncementStatus {
  /** 仅用于查询：所有状态 */
  All = -1,
  /** 草稿 */
  Draft = 0,
  /** 已发布 */
  Published = 1,
  /** 已归档 */
  Archived = 2,
}

export type CreateAnnouncementStatus = Exclude<
  AnnouncementStatus,
  AnnouncementStatus.All
>;

// 公告数据结构
export interface AnnouncementDO {
  id: number;
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: AnnouncementType;
  status: AnnouncementStatus;
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
    avatar?: string;
  } | null;
}

// ============ 请求参数类型 ============

// 创建公告请求
export interface CreateAnnouncementPayload {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  type?: AnnouncementType;
  is_pinned?: boolean;
  is_global?: boolean;
  status?: CreateAnnouncementStatus;
  board_id?: number | null;
  published_at?: string | null;
  expired_at?: string | null;
}

// 更新公告请求
export interface UpdateAnnouncementPayload {
  title?: string;
  content?: string;
  summary?: string;
  cover?: string;
  type?: AnnouncementType;
  is_pinned?: boolean;
  status?: AnnouncementStatus;
  is_global?: boolean;
  board_id?: number | null;
  published_at?: string | null;
  expired_at?: string | null;
}

// 公告列表查询参数
export interface AnnouncementListParams {
  page?: number;
  page_size?: number;
  board_id?: number;
  type?: AnnouncementType;
  status?: AnnouncementStatus;
  is_pinned?: boolean;
  is_global?: boolean;
  keyword?: string;
  start_time?: string;
  end_time?: string;
}

// 置顶/取消置顶请求
export interface PinAnnouncementPayload {
  pinned: boolean;
}

// ============ API 响应类型 ============

// 分页响应数据
export interface AnnouncementListResponse {
  list: AnnouncementDO[];
  total: number;
  page: number;
  page_size: number;
}
