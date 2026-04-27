import { Board } from "./board.type";
import { Like } from "./like.type";
import { Question } from "./question.type";
import { Tag } from "./tags.type";

import { User } from "./user.type";

export type PostType = "all" | "post" | "question" | "article";
export type PostStatus = "draft" | "published" | "hidden";

export interface Post extends BaseModel {
  title: string;
  content: string;
  summary: string;
  cover: string;
  type: PostType;
  status: PostStatus;
  author_id: number;
  view_count: number;
  like_count: number;
  pin_top: boolean;
  board_id: number;
  pin_in_board: boolean;

  author?: User;
  tags?: Tag[];
  comments?: Comment[];
  likes?: Like[];
  board?: Board;
  question?: Question;
}
