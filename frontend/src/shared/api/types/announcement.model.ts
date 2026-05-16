// ============ 请求参数类型 ============

import {
  AnnouncementType,
  CreateAnnouncementStatus,
  AnnouncementStatus,
  AnnouncementDO,
} from "./announcement.model.do";

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
