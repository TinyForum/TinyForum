import { UserDO } from "./user.model.do";

export interface CreateCommentPayload {
  post_id: number;
  content: string;
  parent_id?: number;
}

/**
 * Comment
 */
export interface CommentResponse {
  id: number;
  created_at: string;
  deleted_at: null;
  dislike_count: number;
  is_anonymous: boolean;
  is_pinned: boolean;
  reply: Reply;
  reply_id: number;
  report_count: number;
  sort_weight: number;
  updated_at: string;
  works_id: number;
}

/**
 * reply
 */
export interface Reply {
  id: number;
  author: Author;
  author_id: number;
  content: string;
  created_at: string;
  deleted_at: null;
  like_count: number;
  parent_id: null;
  status: string;
  target_id: number;
  target_type: string;
  updated_at: string;
  replies: Reply[];
}

/**
 * author
 */
export interface Author {
  id: number;
  avatar_url: string;
  bio: string;
  created_at: string;
  deleted_at: null;
  email: string;
  invited_by_id: null;
  is_active: boolean;
  is_blocked: boolean;
  last_login: null;
  role: string;
  score: number;
  updated_at: string;
  username: string;
}
