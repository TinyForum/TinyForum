import { Post } from "./post.model";
import { Comment } from "./comment.model";
import { BaseModel } from "./basic.model";
import { Creation } from "./creation.model";

// export interface Question {
//   id: number;
//   post_id: number;
//   accepted_answer_id?: number;
//   accepted_answer?: Comment;
//
//
// }
export interface QuestionDO {
  id: number;
  creations_id: number;
  creation: Creation;
  accepted_answer_id: number | null;
  created_at: string;
  updated_at: string;
  reward_score: number;
  answer_count: number;
  view_count: number;
  accepted_answer?: Comment;
}
export interface QuestionResponse {
  id: number;
  creations_id: number;
  creation: Creation;
  accepted_answer_id: number | null;
  created_at: string;
  updated_at: string;
  reward_score: number;
  answer_count: number;
  view_count: number;
  accepted_answer?: Comment;
}

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
    avatar_url?: string;
  };
  tags: Array<{
    id: number;
    name: string;
  }>;
}
