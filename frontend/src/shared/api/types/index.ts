import { UserRoleType } from "@/shared/type/roles.types";
import { Board } from "./board.model";

export * from "./basic.model";
/**
 * api/types/index.ts
 * 所有 API 共用的请求 / 响应类型
 * 从这里统一 re-export，业务层只需 import from '@/api/types'
 */

// ─── 通用包装 ─────────────────────────────────────────────────────────────────

// api/types/index.ts

// ─── 枚举 ─────────────────────────────────────────────────────────────────────

export type PostType = "post" | "article" | "topic" | "question";
export type PostStatus = "draft" | "published" | "hidden";
// export type UserRole =
//   | "guest"
//   | "user"
//   | "member"
//   | "moderator"
//   | "reviewer"
//   | "admin"
//   | "super_admin"
//   | "bot";
export type VoteType = "up" | "down";
export type NotificationType =
  | "comment"
  | "like"
  | "follow"
  | "reply"
  | "system";

// ─── 实体类型 ─────────────────────────────────────────────────────────────────

export interface User {
  id: number;
  username: string;
  email: string;
  avatar: string;
  bio: string;
  role: UserRoleType;
  score: number;
  is_active: boolean;
  is_blocked: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}

export interface Follow extends BaseModel {
  follower_id: number;
  following_id: number;
  follower?: User;
  following?: User;
}

export interface Tag {
  id: number;
  name: string;
  description: string;
  color: string;
  post_count: number;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: PostType;
  status: PostStatus;
  author_id: number;
  author?: User;
  view_count: number;
  like_count: number;
  pin_top: boolean;
  board_id: number;
  board?: Board;
  pin_in_board: boolean;
  question?: Question;
  tags?: Tag[];
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  content: string;
  post_id: number;
  author_id: number;
  author?: User;
  parent_id?: number;
  parent?: Comment;
  replies?: Comment[];
  like_count: number;
  is_answer: boolean;
  is_accepted: boolean;
  vote_count: number;
  created_at: string;
  updated_at: string;
}

export interface Question {
  id: number;
  post_id: number;
  accepted_answer_id?: number;
  accepted_answer?: Comment;
  reward_score: number;
  answer_count: number;
}

export interface QuestionResponse {
  id: number;
  post_id: number;
  accepted_answer_id?: number;
  accepted_answer?: Comment;
  reward_score: number;
  answer_count: number;
  post: Post;
  answers: Comment[];
  total: number;
}

export interface AnswerVoteResult {
  vote_count: number;
  user_vote?: VoteType;
}

export interface Notification {
  id: number;
  user_id: number;
  sender_id?: number;
  sender?: User;
  type: NotificationType;
  content: string;
  target_id?: number;
  target_type: string;
  is_read: boolean;
  created_at: string;
}

export interface ModeratorApplication {
  id: number;
  user_id: number;
  board_id: number;
  reason: string;
  status: "pending" | "approved" | "rejected";
  created_at: string;
}

export interface TimelineEvent {
  id: number;
  user_id: number;
  actor_id: number;
  actor?: User;
  action: string;
  target_id: number;
  target_type: string;
  payload: string;
  score: number;
  created_at: string;
}

export interface Subscription {
  id: number;
  subscriber_id: number;
  target_user_id: number;
  target_type: string;
  is_active: boolean;
}

export interface Topic {
  id: number;
  title: string;
  description: string;
  cover: string;
  creator_id: number;
  creator?: User;
  is_public: boolean;
  post_count: number;
  follower_count: number;
  created_at: string;
}

export interface TopicPost {
  id: number;
  topic_id: number;
  post_id: number;
  post?: Post;
  sort_order: number;
  added_by: number;
}

export interface TopicFollow {
  id: number;
  user_id: number;
  topic_id: number;
}

export interface AuthResult {
  user: User;
}

// lib/api/types/index.ts 中添加
// export interface Announcement {
//   id: number;
//   title: string;
//   content: string;
//   type: "system" | "feature" | "maintenance" | "policy";
//   is_pinned: boolean;
//   view_count: number;
//   author_id: number;
//   author?: User;
//   created_at: string;
//   updated_at: string;
// }

export interface QuestionSimple extends BaseModel {
  title: string;
  summary: string;
  view_count: number;
  answer_count: number;
  reward_score: number;
  accepted_answer_id: number | null;
  author: {
    id: number;
    username: string;
    avatar?: string;
  };
  tags: Array<{
    id: number;
    name: string;
  }>;
}

export interface BaseModel {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string;
}
