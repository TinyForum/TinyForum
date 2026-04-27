import { Post } from "./post.type";
import { User } from "./user.type";

export interface Topic extends BaseModel {
  title: string;
  description: string;
  cover: string;
  creator_id: number;
  is_public: boolean;
  post_count: number;
  follower_count: number;

  creator?: User;
  posts?: TopicPost[];
  followers?: TopicFollow[];
}

export interface TopicPost extends BaseModel {
  topic_id: number;
  post_id: number;
  sort_order: number;
  added_by: number;

  topic?: Topic;
  post?: Post;
}

interface TopicFollow extends BaseModel {
  user_id: number;
  topic_id: number;

  user?: User;
  topic?: Topic;
}
