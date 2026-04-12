import { User, UserRole } from "./user.type";

export interface Board extends BaseModel {
  name: string;
  slug: string;
  description: string;
  icon: string;
  cover: string;
  parent_id: number | null;
  sort_order: number;
  view_role: UserRole;
  post_role: UserRole;
  reply_role: UserRole;
  post_count: number;
  thread_count: number;
  today_count: number;

  parent?: Board;
  children?: Board[];
  moderators?: Moderator[];
}

interface BoardTree {
  id: number;
  name: string;
  slug: string;
  description: string;
  icon: string;
  cover: string;
  parent_id: number | null;
  sort_order: number;
  view_role: UserRole;
  post_role: UserRole;
  reply_role: UserRole;
  post_count: number;
  thread_count: number;
  today_count: number;
  children?: BoardTree[];
}

interface Moderator extends BaseModel {
  user_id: number;
  board_id: number;
  permissions: string;
  can_delete_post: boolean;
  can_pin_post: boolean;
  can_edit_any_post: boolean;
  can_manage_moderator: boolean;
  can_ban_user: boolean;

  user?: User;
  board?: Board;
}

interface BoardBan extends BaseModel {
  user_id: number;
  board_id: number;
  banned_by: number;
  reason: string;
  expires_at: Date | null;

  user?: User;
  board?: Board;
  banner?: User;
}

interface ModeratorLog extends BaseModel {
  moderator_id: number;
  board_id: number;
  action: string;
  target_type: string;
  target_id: number;
  reason: string;
  old_value: string;
  new_value: string;

  moderator?: User;
  board?: Board;
}