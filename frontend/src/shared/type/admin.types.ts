// type/admin.types.ts
export type TabType =
  | "users"
  | "posts"
  | "qa"
  | "announcements"
  | "points"
  | "statistics";

export interface User {
  id: number;
  username: string;
  email: string;
  role: "user" | "moderator" | "reviewer" | "bot" | "admin" | "super_admin";
  is_active: boolean;
  score: number;
  created_at: string;
  avatar?: string;
  last_login?: string;
  post_count?: number;
}

export interface Post {
  id: number;
  title: string;
  type: "article" | "topic" | "post";
  status: "published" | "draft" | "hidden";
  pin_top: boolean;
  view_count: number;
  like_count: number;
  comment_count: number;
  created_at: string;
  author?: { id: number; username: string };
}

// export interface Announcement {
//   id: number;
//   title: string;
//   content: string;
//   type: "global" | "board";
//   board_id?: number;
//   is_pinned: boolean;
//   created_at: string;
//   created_by: string;
//   expires_at?: string;
// }

export interface QAItem {
  id: number;
  question: string;
  answer: string;
  status: "pending" | "answered" | "closed";
  created_at: string;
  author: string;
  category: string;
}

export interface PointsRecord {
  id: number;
  username: string;
  amount: number;
  type: "award" | "deduct";
  reason: string;
  created_at: string;
}

export interface Stats {
  totalUsers: number;
  totalPosts: number;
  todayActive: number;
  onlineNow: number;
  totalPoints: number;
  userGrowth: number;
  postGrowth: number;
  exchangeRate: number;
  todayAwarded: number;
}
