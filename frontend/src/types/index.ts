import { Board } from "./board.type";
import { User } from "./user.type";
// 基础类型
export type SortBy = 'latest' | 'hot' | 'like' | 'random';

// API 响应类型
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export interface AuthResult {
  token: string;
  user: User;
}

// 投票结果
export interface AnswerVoteResult {
  vote_type: 'up' | 'down' | '';
  vote_count: number;
  action: 'added' | 'removed' | 'updated';
}

// 版主申请
export interface ModeratorApplication {
  id: number;
  user_id: number;
  board_id: number;
  reason: string;
  status: 'pending' | 'approved' | 'rejected';
  handled_by: number | null;
  handle_note: string;
  created_at: string;
  updated_at: string;
  user?: User;
  board?: Board;
  handler?: User;
}

export * from './auth.type'
export * from './board.type';
export * from './comment.type';
export * from './follow.type'
export * from './like.type';
export * from './notification.type';
export * from './post.type';
export * from './question.type';
export * from './report.type';
export * from './tags.type';
export * from './timeline.type';
export * from './topic.type';
export * from './user.type';