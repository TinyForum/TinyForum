import { Follow } from "./follow.type";
import { Post } from "./post.type";

export type UserRole = 'user' | 'admin' | 'moderator';

export interface User extends BaseModel {
  username: string;
  email: string;
  password: string;
  avatar: string;
  bio: string;
  role: UserRole;
  score: number;
  is_active: boolean;
  is_blocked: boolean;
  last_login: Date | null;

  posts?: Post[];
  comments?: Comment[];
  followers?: Follow[];
  following?: Follow[];
}