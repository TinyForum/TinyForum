import { User } from "./user.type";

type ActionType = 'create_post' | 'create_comment' | 'like_post' | 'like_comment' | 'follow_user' | 'accept_answer' | 'sign_in';

interface TimelineEvent extends BaseModel {
  user_id: number;
  actor_id: number;
  action: ActionType;
  target_id: number;
  target_type: string;
  payload: string;
  score: number;

  user?: User;
  actor?: User;
}

interface UserTimeline extends BaseModel {
  user_id: number;
  timeline_type: string;
  last_read_at: Date;
}

interface TimelineSubscription extends BaseModel {
  subscriber_id: number;
  target_user_id: number;
  target_type: string;
  target_id: number;
  is_active: boolean;
}