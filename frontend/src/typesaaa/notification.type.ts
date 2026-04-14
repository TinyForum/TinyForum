export type NotificationType = 'comment' | 'like' | 'follow' | 'reply' | 'system';

import { User } from "./user.type";
export interface Notification extends BaseModel {
  user_id: number;
  sender_id: number | null;
  type: NotificationType;
  content: string;
  target_id: number | null;
  target_type: string;
  is_read: boolean;

  user?: User;
  sender?: User;
}