import { UserRoleType } from "@/shared/type/roles.types";
import { PageRequest } from "./basic.model";

export interface UserStatsVO {
  total_post: number;
  total_comment: number;
  total_favorite: number;
  total_like: number;
  total_follower: number;
  total_following: number;
  total_report: number;
  total_violation: number;
  total_question: number;
  total_answer: number;
  total_upload: number;
  total_score: number;
}

export interface UserDO {
  id: number;
  username: string;
  email: string;
  avatar: string;
  bio: string;
  role: UserRoleType;
  score: number;
  is_active: boolean;
  is_blocked: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}
export interface UserPostsVO {
  id: number;
  title: string;
  summary: string;
  cover: string;
  type: string;
  post_status: string;
  moderation_status: string;
  view_count: number;
  likes_count: number;
  comment_count: number;
  pin_top: boolean;
  tags: string[];
  board_name: string;
  pin_in_board: boolean;
  created_at: string;
  updated_at: string;
}

export interface GetUserPostsRequest extends PageRequest {
  keyword?: string;
  status?: string;
  moderation_status?: string;
  tag?: string;
  board_name?: string;
}
