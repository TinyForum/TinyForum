export interface User {
  id: number;
  username: string;
  email: string;
  avatar: string;
  bio: string;
  role: 'user' | 'admin';
  score: number;
  is_active: boolean;
  last_login: string | null;
  created_at: string;
  follower_count?: number;
  following_count?: number;
  is_following?: boolean;
}

export  type SortBy = 'latest'  | 'hot' | "like"| 'random';
export interface Tag {
  id: number;
  name: string;
  description: string;
  color: string;
  post_count: number;
}

export type PostType = 'post' | 'article' | 'topic' | "question" | "all";
export type PostStatus = 'draft' | 'published' | 'hidden';

export interface Post {
  id: number;
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: PostType;
  status: PostStatus;
  author_id: number;
  author: User;
  tags: Tag[];
  view_count: number;
  like_count: number;
  pin_top: boolean;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  content: string;
  post_id: number;
  author_id: number;
  author: User;
  parent_id: number | null;
  like_count: number;
  replies?: Comment[];
  created_at: string;
}

export interface Notification {
  id: number;
  user_id: number;
  sender_id: number | null;
  sender?: User;
  type: 'comment' | 'like' | 'follow' | 'reply' | 'system';
  content: string;
  target_id: number | null;
  target_type: string;
  is_read: boolean;
  created_at: string;
}

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

// 添加以下类型定义

// Board 板块
export interface Board {
  id: number;
  name: string;
  slug: string;
  description: string;
  icon: string;
  cover: string;
  parent_id: number | null;
  sort_order: number;
  view_role: 'user' | 'admin';
  post_role: 'user' | 'admin';
  reply_role: 'user' | 'admin';
  post_count: number;
  thread_count: number;
  today_count: number;
  created_at: string;
  updated_at: string;
  parent?: Board;
  children?: Board[];
}

// Moderator 版主
export interface Moderator {
  id: number;
  user_id: number;
  board_id: number;
  can_delete_post: boolean;
  can_pin_post: boolean;
  can_edit_any_post: boolean;
  can_manage_moderator: boolean;
  can_ban_user: boolean;
  created_at: string;
  user?: User;
  board?: Board;
}

// Topic 专题
export interface Topic {
  id: number;
  title: string;
  description: string;
  cover: string;
  creator_id: number;
  is_public: boolean;
  post_count: number;
  follower_count: number;
  created_at: string;
  updated_at: string;
  creator?: User;
}

// TopicPost 专题帖子
export interface TopicPost {
  id: number;
  topic_id: number;
  post_id: number;
  sort_order: number;
  added_by: number;
  created_at: string;
  post?: Post;
  topic?: Topic;
}

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
  action: 'create_post' | 'create_comment' | 'like_post' | 'like_comment' | 'follow_user' | 'accept_answer' | 'sign_in';
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

// Question 问答
export interface Question {
  id: number;
  post_id: number;
  accepted_answer_id: number | null;
  reward_score: number;
  answer_count: number;
  created_at: string;
  updated_at: string;
  post?: Post;
  accepted_answer?: Comment;
}

// AnswerVoteResult 投票结果
export interface AnswerVoteResult {
  vote_type: 'up' | 'down' | '';
  vote_count: number;
  action: 'added' | 'removed' | 'updated';
}

// 扩展 Post 类型
export interface ExtendedPost extends Post {
  board_id?: number;
  board?: Board;
  is_question?: boolean;
  question?: Question;
  pin_in_board?: boolean;
}

// 扩展 Comment 类型
export interface ExtendedComment extends Comment {
  is_answer?: boolean;
  is_accepted?: boolean;
  vote_count?: number;
}