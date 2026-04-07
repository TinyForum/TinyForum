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

export interface Tag {
  id: number;
  name: string;
  description: string;
  color: string;
  post_count: number;
}

export type PostType = 'post' | 'article' | 'topic';
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
