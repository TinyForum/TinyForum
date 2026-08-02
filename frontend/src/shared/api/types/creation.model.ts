import { Board } from "./board.model";
import { PostStatus, PostType } from "./post.model";
import { Tag } from "./tag.model";
import { UserDO } from "./user.model.do";

export interface Creation {
  id: number;
  created_at: string;
  updated_at: string;
  title: string;
  content: string;
  summary: string;
  cover_url: string;
  video_url: string;
  slug: string;
  image_urls: string[];
  creation_type: PostType;
  creation_status: PostStatus;
  moderation_status: string;
  view_count: number;
  like_count: number;
  pin_top: boolean;
  is_original: boolean;
  pin_in_board: boolean;
  board_id: number;
  author_id: number;
  author?: UserDO;
  board?: Board;
  tags?: Tag[];
}
