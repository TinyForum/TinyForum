import { UserRoleType } from "@/shared/api/types/roles.model";
export interface UserDO {
  id: number;
  username: string;
  email: string;
  avatar_url: string;
  bio: string;
  role: UserRoleType;
  score: number;
  is_active: boolean;
  is_blocked: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}
export interface ProfileResponse {
  id: number;
  created_at: string; // ISO 8601 时间字符串
  updated_at: string;
  deleted_at: string | null;
  username: string;
  email: string;
  avatar_url: string;
  bio: string;
  role: UserRoleType;
  score: number;
  is_active: boolean;
  is_blocked: boolean;
  last_login: string | null;
  invited_by_id: number | null;
  follower_count: number;
  following_count: number;
  is_following: boolean;
}
