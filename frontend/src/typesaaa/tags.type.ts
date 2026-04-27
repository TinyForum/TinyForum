import { Post } from "./post.type";

export interface Tag extends BaseModel {
  name: string;
  description: string;
  color: string;
  post_count: number;

  posts?: Post[];
}
