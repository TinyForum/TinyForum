import { Post } from "./post.type";

export interface Question extends BaseModel {
  post_id: number;
  accepted_answer_id: number | null;
  reward_score: number;
  answer_count: number;

  post?: Post;
  accepted_answer?: Comment;
}

export interface AnswerVote extends BaseModel {
  user_id: number;
  comment_id: number;
  vote_type: string; // up/down
}