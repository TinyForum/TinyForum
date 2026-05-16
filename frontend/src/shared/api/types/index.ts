/**
 * api/types/index.ts
 * 所有 API 共用的请求 / 响应类型
 * 从这里统一 re-export，业务层只需 import from '@/api/types'
 */

// ─── 通用包装 ─────────────────────────────────────────────────────────────────

// api/types/index.ts

// ─── 枚举 ─────────────────────────────────────────────────────────────────────

// export type UserRole =
//   | "guest"
//   | "user"
//   | "member"
//   | "moderator"
//   | "reviewer"
//   | "admin"
//   | "super_admin"
//   | "bot";

// ─── 实体类型 ─────────────────────────────────────────────────────────────────

// export interface AuthResult {
//   user: UserDO;
// }

// lib/api/types/index.ts 中添加
// export interface Announcement {
//   id: number;
//   title: string;
//   content: string;
//   type: "system" | "feature" | "maintenance" | "policy";
//   is_pinned: boolean;
//   view_count: number;
//   author_id: number;
//   author?: User;
//   created_at: string;
//   updated_at: string;
// }
