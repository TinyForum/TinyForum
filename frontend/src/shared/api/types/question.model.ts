import { Post } from "./post.model";
import { Comment } from "./comment.model";
import { BaseModel } from "./basic.model";

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
