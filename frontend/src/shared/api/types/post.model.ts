import { Board } from "./board.model";
import { Question } from "./question.model";
import { Tag } from "./tag.model";
import { UserDO } from "./user.model.do";

// ─── 普通帖子 ─────────────────────────────────────────────────────────────────
export type PostType = "post" | "article" | "topic" | "question";
export type PostStatus = "draft" | "published" | "pending" | "hidden";
export interface Post {
  id: number;
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: PostType;
  status: PostStatus;
  author_id: number;
  author?: UserDO;
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

export interface PostDetailResult {
  post: Post;
  liked: boolean;
}
