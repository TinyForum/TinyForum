import { UserDO } from "./user.model.do";

export type NotificationType =
  | "comment"
  | "like"
  | "follow"
  | "reply"
  | "system";

export interface Notification {
  id: number;
  user_id: number;
  sender_id?: number;
  sender?: UserDO;
  type: NotificationType;
  content: string;
  target_id?: number;
  target_type: string;
  is_read: boolean;
  created_at: string;
}
