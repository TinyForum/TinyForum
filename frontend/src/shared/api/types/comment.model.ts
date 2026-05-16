import { UserDO } from "./user.model.do";

export interface CreateCommentPayload {
  post_id: number;
  content: string;
  parent_id?: number;
}

// export interface VoteStatusResult {
//   has_voted: boolean;
//   vote_type: VoteType | "";
//   vote_count: number;
// }
export interface Comment {
  id: number;
  content: string;
  post_id: number;
  author_id: number;
  author?: UserDO;
  parent_id?: number;
  parent?: Comment;
  replies?: Comment[];
  like_count: number;
  is_answer: boolean;
  is_accepted: boolean;
  vote_count: number;
  created_at: string;
  updated_at: string;
}
