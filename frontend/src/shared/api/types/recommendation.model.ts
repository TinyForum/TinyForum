export interface RecommendationItem {
  creation_id: number;
  title: string;
  summary: string;
  cover_url: string;
  author_id: number;
  author_name: string;
  author_avatar: string;
  view_count: number;
  like_count: number;
  comment_count: number;
  board_id: number;
  board_name: string;
  score: number;
  reason: string;
  created_at: string;
}

export interface RecommendationResponse {
  items: RecommendationItem[];
  total: number;
  page: number;
  page_size: number;
  has_more: boolean;
  session_id: string;
  strategy: string;
  generated_at: string;
}

export interface BehaviorEvent {
  target_id: number;
  target_type: string;
  behavior_type: string;
  value?: number;
  session_id?: string;
  context_json?: string;
}

export interface RecommendationFeedback {
  creation_id: number;
  feedback_type: string;
  source_type?: string;
  position?: number;
  session_id?: string;
}

export interface BatchFeedback {
  feedbacks: RecommendationFeedback[];
}

export interface UserInterestProfile {
  active_tags: string[];
  tag_weights: Record<string, number>;
  updated_at: string;
}
