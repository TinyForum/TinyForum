import { UserRoleType } from "@/shared/api/types/roles.model";

export interface Board {
  id: number;
  name: string;
  slug: string;
  description: string;
  icon: string;
  cover: string;
  parent_id?: number;
  parent?: Board;
  children?: Board[];
  sort_order: number;
  view_role: UserRoleType;
  post_role: UserRoleType;
  reply_role: UserRoleType;
  post_count: number;
  thread_count: number;
  today_count: number;
}

export interface BoardPostListItem {
  id: number;
  title: string;
  summary: string;
  cover: string;
  type: string;
  author_id: number;
  author_name: string;
  created_at: string; // ISO date string
}
// ─── Payloads ─────────────────────────────────────────────────────────────────

export interface CreateBoardPayload {
  name: string;
  slug: string;
  description?: string;
  icon?: string;
  cover?: string;
  parent_id?: number;
  sort_order?: number;
  view_role?: UserRoleType;
  post_role?: UserRoleType;
  reply_role?: UserRoleType;
}

export type UpdateBoardPayload = Partial<CreateBoardPayload>;

export type BoardCreationApplicationStatus =
  | "pending"
  | "approved"
  | "rejected"
  | "canceled";

export interface BoardCreationApplication {
  id: number;
  user_id: number;
  name: string;
  slug: string;
  description: string;
  reason: string;
  status: BoardCreationApplicationStatus;
  reviewer_id?: number;
  review_note: string;
  reviewed_at?: string;
  created_board_id?: number;
  created_at: string;
  updated_at: string;
  user?: { id: number; username: string };
  created_board?: Board;
}

export interface ApplyBoardCreationPayload {
  name: string;
  slug: string;
  description?: string;
  reason: string;
}

export interface AddModeratorPayload {
  user_id: number;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface BanUserPayload {
  user_id: number;
  reason: string;
  expires_at?: string;
}

export interface PinPostPayload {
  pin_in_board: boolean;
}
