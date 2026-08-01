import { Reply } from "./comment.model";

export interface AnswerListParams {
  page?: number;
  page_size?: number;
  sort?: "vote" | "newest" | "oldest";
}
export interface AnswerResponse {
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
  is_accepted: boolean;
}
