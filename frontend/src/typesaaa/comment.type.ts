import { Like } from "./like.type";
import { Post } from "./post.type";
import { User } from "./user.type";

export interface Comment extends BaseModel {
  content: string;
  post_id: number;
  author_id: number;
  parent_id: number | null;
  like_count: number;
  is_answer: boolean;
  is_accepted: boolean;
  vote_count: number;

  post?: Post;
  author?: User;
  parent?: Comment;
  replies?: Comment[];
  likes?: Like[];
}
