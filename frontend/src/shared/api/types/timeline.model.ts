import { UserDO } from "./user.model";

export interface TimelineEvent {
  id: number;
  user_id: number;
  actor_id: number;
  actor?: UserDO;
  action: string;
  target_id: number;
  target_type: string;
  payload: string;
  score: number;
  created_at: string;
}

export interface Subscription {
  id: number;
  subscriber_id: number;
  target_user_id: number;
  target_type: string;
  is_active: boolean;
}
