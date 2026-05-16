// export interface Moderator {
//   id: number;
//   user_id: number;
//   board_id: number;
//   user?: User;
//   board?: Board;
//   can_delete_post: boolean;
//   can_pin_post: boolean;
//   can_edit_any_post: boolean;
//   can_manage_moderator: boolean;
//   can_ban_user: boolean;
// }

export interface ApplyModeratorForm {
  reason: string;
  req_delete_post?: boolean;
  req_pin_post?: boolean;
  req_edit_any_post?: boolean;
  req_manage_moderator?: boolean;
  req_ban_user?: boolean;
}

export interface ReviewApplicationRequest {
  approve: boolean;
  review_note?: string;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface AddModeratorRequest {
  user_id: number;
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface UpdatePermissionsRequest {
  can_delete_post?: boolean;
  can_pin_post?: boolean;
  can_edit_any_post?: boolean;
  can_manage_moderator?: boolean;
  can_ban_user?: boolean;
}

export interface BanUserRequest {
  user_id: number;
  reason: string;
  expires_at?: string;
}

export interface ModeratorApplication {
  id: number;
  user_id: number;
  username: string;
  board_id: number;
  board_name: string;
  reason: string;
  status: "pending" | "approved" | "rejected" | "canceled";
  review_note: string;
  reviewed_by?: number;
  reviewed_at?: string;
  created_at: string;
  updated_at: string;
  // 添加缺失的字段
  board?: {
    id: number;
    name: string;
    slug: string;
  };
  req_delete_post: boolean;
  req_pin_post: boolean;
  req_edit_any_post: boolean;
  req_manage_moderator: boolean;
  req_ban_user: boolean;
}

export interface Moderator {
  id: number;
  user_id: number;
  created_at: string;
  board_id: number;
  permissions: {
    can_delete_post: boolean;
    can_pin_post: boolean;
    can_edit_any_post: boolean;
    can_manage_moderator: boolean;
    can_ban_user: boolean;
  };
  user?: {
    id: number;
    username: string;
    avatar?: string;
  };
  board?: {
    id: number;
    name: string;
  };
}

export interface BanRecord {
  id: number;
  user_id: number;
  board_id: number;
  moderator_id: number;
  reason: string;
  expires_at?: string;
  is_active: boolean;
  created_at: string;
}

// 申请状态
// type ApplicationStatus = "pending" | "approved" | "rejected" | "canceled";

// 申请状态详情
// interface ApplicationStatusDetailResponse {
//   has_application: boolean;
//   application_id?: number;
//   status?: ApplicationStatus;
//   reason?: string;
//   created_at?: string;
//   review_note?: string;
//   reviewer_id?: number;
//   reviewed_at?: string | null;
//   can_cancel: boolean;
//   can_resubmit: boolean;
//   requested_perms?: {
//     delete_post: boolean;
//     pin_post: boolean;
//     edit_any_post: boolean;
//     manage_moderator: boolean;
//     ban_user: boolean;
//   };
//   can_apply: boolean;
// }

// types/moderator.ts
export interface ModeratorBoard {
  id: number;
  name: string;
  slug: string;
  description?: string;
  icon?: string;
  cover?: string;
  parent_id?: number | null;
  sort_order?: number;
  post_count?: number;
  thread_count?: number;
  today_count?: number;
  created_at: string;
  updated_at: string;
  permissions: {
    can_delete_post: boolean;
    can_pin_post: boolean;
    can_edit_any_post: boolean;
    can_manage_moderator: boolean;
    can_ban_user: boolean;
  };
}

export interface ModeratorPost {
  id: number;
  title: string;
  content: string;
  author_id: number;
  author_name: string;
  author_avatar?: string;
  board_id: number;
  is_pinned: boolean;
  view_count: number;
  like_count: number;
  comment_count: number;
  status: "pending" | "published" | "deleted";
  created_at: string;
}

export interface ModeratorReport {
  id: number;
  target_type: "post" | "comment";
  target_id: number;
  target_title?: string;
  target_content?: string;
  reporter_id: number;
  reporter_name: string;
  reason: string;
  status: "pending" | "resolved" | "dismissed";
  created_at: string;
}

export interface BannedUser {
  id: number;
  user_id: number;
  username: string;
  avatar?: string;
  board_id: number;
  reason: string;
  expires_at?: string;
  is_active: boolean;
  created_at: string;
}

// export interface ModeratorApplication {
//   id: number;
//   user_id: number;
//   board_id: number;
//   reason: string;
//   status: "pending" | "approved" | "rejected";
//   created_at: string;
// }
