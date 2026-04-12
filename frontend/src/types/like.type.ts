import { Post } from "./post.type";
import { User } from "./user.type";

export interface Like extends BaseModel {
  user_id: number;
  post_id: number | null;
  comment_id: number | null;

  user?: User;
  post?: Post;
  comment?: Comment;
}