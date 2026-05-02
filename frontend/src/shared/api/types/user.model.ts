import { UserRoleType } from "@/shared/type/roles.types";

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
