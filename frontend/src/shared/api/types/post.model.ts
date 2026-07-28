// import { Board } from "./board.model";
import { Creation } from "./creation.model";
// import { Question } from "./question.model";
// import { Tag } from "./tag.model";
// import { UserDO } from "./user.model.do";

// ─── 普通帖子 ─────────────────────────────────────────────────────────────────
export type PostType = "post" | "article" | "topic" | "questions";
export type PostStatus = "draft" | "published" | "pending" | "hidden";
export interface Post {
  id: number;
  creations_id: number;
  creation: Creation;

  created_at: string;
  updated_at: string;
}

export interface PostListParams {
  page?: number;
  page_size?: number;
  keyword?: string;
  sort_by?: string;
  type?: PostType;
  author_id?: number;
  tag_id?: number;
  board_id?: number;
}

export interface CreatePostPayload {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  type?: PostType;
  board_id?: number;
  tag_ids?: number[];
  status?: PostStatus;
}

export interface UpdatePostPayload {
  title?: string;
  content?: string;
  summary?: string;
  cover?: string;
  tag_ids?: number[];
}

export interface PostDetailResponse {
  post: Post;
  liked: boolean;
}
