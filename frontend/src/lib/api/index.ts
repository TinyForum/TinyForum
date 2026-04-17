/**
 * api/index.ts
 * 统一出口 —— 业务层只需 import from '@/api'
 *
 * 用法示例：
 *   import { postApi, commentApi } from '@/api'
 *   import type { Post, Comment } from '@/api'
 */

// ─── Client ───────────────────────────────────────────────────────────────────
export { default as apiClient } from "./client";

// ─── Types ────────────────────────────────────────────────────────────────────
export type {
  ApiResponse,
  PageData,
  PostType,
  PostStatus,
  // UserRole,
  VoteType,
  NotificationType,
  User,
  Tag,
  Post,
  Comment,
  Question,
  AnswerVoteResult,
  Notification,
  Board,
  Moderator,
  ModeratorApplication,
  TimelineEvent,
  Subscription,
  Topic,
  TopicPost,
  TopicFollow,
  AuthResult,
} from "./types";

// ─── Modules ──────────────────────────────────────────────────────────────────
export { authApi } from "./modules/auth";
export { postApi } from "./modules/posts";
export { commentApi } from "./modules/comments";
export { userApi } from "./modules/users";
export { tagApi } from "./modules/tags";
export { notificationApi } from "./modules/notifications";
export { boardApi } from "./modules/boards";
export { timelineApi } from "./modules/timeline";
export { topicApi } from "./modules/topics";
export { adminApi } from "./modules/admin";
export { announcementApi } from "./modules/announcements";
export { questionApi } from "./modules/questions";


// ─── Payload types（按需 re-export） ──────────────────────────────────────────
export type { RegisterPayload, LoginPayload } from "./modules/auth";
export type {
  PostListParams,
  CreatePostPayload,
  UpdatePostPayload,
  PostDetailResult,
} from "./modules/posts";
export type {
  CreateCommentPayload,
  VoteStatusResult,
} from "./modules/comments";
export type { UpdateProfilePayload } from "./modules/users";
export type { CreateTagPayload, UpdateTagPayload } from "./modules/tags";
export type {
  CreateBoardPayload,
  UpdateBoardPayload,
  AddModeratorPayload,
  BanUserPayload,
  PinPostPayload,
} from "./modules/boards";
export type { CreateTopicPayload, UpdateTopicPayload } from "./modules/topics";