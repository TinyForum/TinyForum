import { Post } from "./post.model";
import { UserDO } from "./user.model.do";

export interface Topic {
  id: number;
  title: string;
  description: string;
  cover: string;
  creator_id: number;
  creator?: UserDO;
  is_public: boolean;
  post_count: number;
  follower_count: number;
  created_at: string;
}
export interface TopicPost {
  id: number;
  topic_id: number;
  post_id: number;
  post?: Post;
  sort_order: number;
  added_by: number;
}

export interface TopicFollow {
  id: number;
  user_id: number;
  topic_id: number;
}
