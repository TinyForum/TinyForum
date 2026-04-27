// import { User } from './user.type';

import { Board } from "./board.type";
import { Topic } from "./topic.type";
import { User } from "./user.type";

export type SortBy = "latest" | "hot" | "like" | "random";
// export interface Tag {
//   id: number;
//   name: string;
//   description: string;
//   color: string;
//   post_count: number;
// }

// export type PostType = 'post' | 'article' | 'topic' | "question" | "all";
// export type PostStatus = 'draft' | 'published' | 'hidden';

// export interface ApiResponse<T> {
//   code: number;
//   message: string;
//   data: T;
// }

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// export interface AuthResult {
//   token: string;
//   user: User;
// }

// // 添加以下类型定义

// // Board 板块
// export interface Board {
//   id: number;
//   name: string;
//   slug: string;
//   description: string;
//   icon: string;
//   cover: string;
//   parent_id: number | null;
//   sort_order: number;
//   view_role: 'user' | 'admin';
//   post_role: 'user' | 'admin';
//   reply_role: 'user' | 'admin';
//   post_count: number;
//   thread_count: number;
//   today_count: number;
//   created_at: string;
//   updated_at: string;
//   parent?: Board;
//   children?: Board[];
// }

// // Moderator 版主
// // export interface Moderator {
// //   id: number;
// //   user_id: number;
// //   board_id: number;
// //   can_delete_post: boolean;
// //   can_pin_post: boolean;
// //   can_edit_any_post: boolean;
// //   can_manage_moderator: boolean;
// //   can_ban_user: boolean;
// //   created_at: string;
// //   user?: User;
// //   board?: Board;
// // }

// // Topic 专题
// // export interface Topic {
// //   id: number;
// //   title: string;
// //   description: string;
// //   cover: string;
// //   creator_id: number;
// //   is_public: boolean;
// //   post_count: number;
// //   follower_count: number;
// //   created_at: string;
// //   updated_at: string;
// //   creator?: User;
// // }

// // TopicPost 专题帖子
// export interface TopicPost {
//   id: number;
//   topic_id: number;
//   post_id: number;
//   sort_order: number;
//   added_by: number;
//   created_at: string;
//   post?: Post;
//   topic?: Topic;
// }

// TopicFollow 专题关注
export interface TopicFollow {
  id: number;
  user_id: number;
  topic_id: number;
  created_at: string;
  user?: User;
  topic?: Topic;
}

// TimelineEvent 时间线事件
export interface TimelineEvent {
  id: number;
  user_id: number;
  actor_id: number;
  action:
    | "create_post"
    | "create_comment"
    | "like_post"
    | "like_comment"
    | "follow_user"
    | "accept_answer"
    | "sign_in";
  target_id: number;
  target_type: string;
  payload: any;
  score: number;
  created_at: string;
  user?: User;
  actor?: User;
}

// Subscription 订阅
export interface Subscription {
  id: number;
  subscriber_id: number;
  target_user_id: number;
  target_type: string;
  target_id: number;
  is_active: boolean;
  created_at: string;
}

// // Question 问答
// export interface Question {
//   id: number;
//   post_id: number;
//   accepted_answer_id: number | null;
//   reward_score: number;
//   answer_count: number;
//   created_at: string;
//   updated_at: string;
//   post?: Post;
//   accepted_answer?: Comment;
// }

// AnswerVoteResult 投票结果
export interface AnswerVoteResult {
  vote_type: "up" | "down" | "";
  vote_count: number;
  action: "added" | "removed" | "updated";
}

// 添加版主申请类型
export interface ModeratorApplication {
  id: number;
  user_id: number;
  board_id: number;
  reason: string;
  status: "pending" | "approved" | "rejected";
  handled_by: number | null;
  handle_note: string;
  created_at: string;
  updated_at: string;
  user?: User;
  board?: Board;
  handler?: User;
}

// // 基础类型
// // export type SortBy = 'latest' | 'hot' | 'like' | 'random';

// // API 响应类型
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// export interface PageData<T> {
//   list: T[];
//   total: number;
//   page: number;
//   page_size: number;
// }

export interface AuthResult {
  token: string;
  user: User;
}

// // 投票结果
// export interface AnswerVoteResult {
//   vote_type: 'up' | 'down' | '';
//   vote_count: number;
//   action: 'added' | 'removed' | 'updated';
// }

// // 版主申请
// export interface ModeratorApplication {
//   id: number;
//   user_id: number;
//   board_id: number;
//   reason: string;
//   status: 'pending' | 'approved' | 'rejected';
//   handled_by: number | null;
//   handle_note: string;
//   created_at: string;
//   updated_at: string;
//   user?: User;
//   board?: Board;
//   handler?: User;
// }

export * from "./auth.type";
export * from "./board.type";
export * from "./comment.type";
export * from "./follow.type";
export * from "./like.type";
export * from "./notification.type";
export * from "./post.type";
export * from "./question.type";
export * from "./report.type";
export * from "./tags.type";
export * from "./timeline.type";
export * from "./topic.type";
export * from "./user.type";
export * from "./moderator.type";
