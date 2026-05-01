import { UserRoleType } from "@/shared/type/roles.types";

export interface Board {
  id: number;
  name: string;
  slug: string;
  description: string;
  icon: string;
  cover: string;
  parent_id?: number;
  parent?: Board;
  children?: Board[];
  sort_order: number;
  view_role: UserRoleType;
  post_role: UserRoleType;
  reply_role: UserRoleType;
  post_count: number;
  thread_count: number;
  today_count: number;
}
