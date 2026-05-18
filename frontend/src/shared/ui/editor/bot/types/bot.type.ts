// // ---------- 类型定义（与后端 Go 模型对齐）----------
// export type TriggerType =
//   | "on_schedule"
//   | "on_new_post"
//   | "on_new_comment"
//   | "on_user_register"
//   | "on_keyword"
//   | "on_manual";

// export type CondType =
//   | "post_title_contains"
//   | "post_content_contains"
//   | "user_role_is"
//   | "user_post_count_gte"
//   | "board_id_in"
//   | "time_range"
//   | "custom_expr";

// export type ActionType =
//   | "reply_post"
//   | "delete_post"
//   | "hide_post"
//   | "pin_post"
//   | "lock_post"
//   | "create_post"
//   | "delete_comment"
//   | "ban_user"
//   | "send_message"
//   | "webhook"
//   | "notify_admin"
//   | "wait"
//   | "set_variable"
//   | "stop_if";

// export interface TriggerNode {
//   type: TriggerType;
//   params?: Record<string, unknown>;
// }

// export interface CondNode {
//   type: CondType;
//   negate?: boolean;
//   params: Record<string, unknown>;
// }

// export interface ActionNode {
//   type: ActionType;
//   params: Record<string, unknown>;
// }

// export interface Flow {
//   version: string;
//   trigger: TriggerNode;
//   conditions?: CondNode[];
//   actions: ActionNode[];
// }
