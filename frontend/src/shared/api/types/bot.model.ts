// ========== 机器人相关类型定义（与后端 vo.BotResponse 对应） ==========

/** 机器人状态 */
export type BotStatus = "active" | "inactive" | "error" | "loading" | "stopped";

/** 机器人类型 */
export type BotType =
  | "chat"
  | "moderate"
  | "notify"
  | "sync"
  | "task"
  | "webhook"
  | "analysis";

/** 触发类型 */
export type BotTriggerType = "schedule" | "event" | "webhook" | "manual";

/** 定价模型 */
export type BotPricingModel = "free" | "freemium" | "paid" | "subscription";

/** 资源限制 */
export interface ResourceLimit {
  maxMemoryMB: number;
  maxCPU: number;
}

/** 定价信息 */
export interface BotPricing {
  model: BotPricingModel;
  price?: number; // 单位：元
  cycle?: string; // once/monthly/yearly
  freeLimit?: string;
  purchaseUrl?: string;
}

/** 机器人配置字段定义（用于前端表单动态渲染） */
export interface BotConfigField {
  key: string;
  label: string;
  type: "text" | "number" | "boolean" | "select" | "textarea" | "secret";
  defaultValue?: unknown; // 替换 any
  placeholder?: string;
  description?: string;
  required?: boolean;
  options?: Array<{ label: string; value: unknown }>; // 替换 any
}

/** 机器人权限（与后端保持一致） */
export type BotPermission =
  | "read:user"
  | "read:posts"
  | "write:posts"
  | "read:comments"
  | "write:comments"
  | "send:message"
  | "manage:content"
  | "read:stats";

/** 机器人视图对象（后端 vo.BotResponse） */
export interface BotVO {
  // id: number;
  slug: string;
  name: string;
  version: string;
  description: string;
  summary: string;
  avatar_url: string;
  screenshots: string[];
  homepage_url: string;
  type: BotType;
  tags: string[];
  creator_id: number;
  creator_name: string;
  script_code: string; // Lua 脚本
  script_url?: string; // 可选外部 URL
  trigger_type: BotTriggerType;
  cron_expr: string;
  event_filter: string;
  timeout_sec: number;
  retry_times: number;
  envVars: Record<string, string>;
  resource_limit?: ResourceLimit;
  pricing: BotPricing;
  permissions: BotPermission[];
  enabled: boolean;
  status: BotStatus;
  exec_count: number;
  last_exec_at?: string;
  error_msg: string;
  config_schema: BotConfigField[];
  config_values: Record<string, unknown>; // 替换 Record<string, any>
  created_at: string;
  updated_at: string;
}

/** 机器人列表分页响应（后端返回格式） */
export interface BotListResponse {
  list: BotVO[];
  total: number;
  page: number;
}

/** 创建机器人请求（无 enabled 字段） */
export interface CreateBotRequest {
  name: string;
  version: string;
  description?: string;
  summary?: string;
  avatarUrl?: string;
  screenshots?: string[];
  homepageUrl?: string;
  type: BotType;
  tags?: string[];
  scriptCode: string;
  scriptUrl?: string;
  triggerType: BotTriggerType;
  cronExpr?: string;
  eventFilter?: string;
  timeoutSec?: number;
  retryTimes?: number;
  envVars?: Record<string, string>;
  resourceLimit?: ResourceLimit;
  pricing?: BotPricing;
  permissions?: BotPermission[];
  configSchema?: BotConfigField[];
  configValues?: Record<string, unknown>; // 替换 any
  enabled?: boolean;
}

/** 更新机器人请求（所有字段可选，包含 enabled） */
export interface UpdateBotRequest {
  name?: string;
  version?: string;
  description?: string;
  summary?: string;
  avatarUrl?: string;
  screenshots?: string[];
  homepageUrl?: string;
  type?: BotType;
  tags?: string[];
  scriptCode?: string;
  scriptUrl?: string;
  triggerType?: BotTriggerType;
  cronExpr?: string;
  eventFilter?: string;
  timeoutSec?: number;
  retryTimes?: number;
  envVars?: Record<string, string>;
  resourceLimit?: ResourceLimit;
  pricing?: BotPricing;
  permissions?: BotPermission[];
  configSchema?: BotConfigField[];
  configValues?: Record<string, unknown>; // 替换 any
  enabled?: boolean;
}

// ========== 机器人相关类型定义（与后端 vo.BotResponse 对应） ==========

export interface ResourceLimit {
  maxMemoryMB: number;
  maxCPU: number;
}

/** 定价信息 */
export interface BotPricing {
  model: BotPricingModel;
  price?: number; // 单位：元
  cycle?: string; // once/monthly/yearly
  freeLimit?: string;
  purchaseUrl?: string;
}

/** 机器人配置字段定义（用于前端表单动态渲染） */
export interface BotConfigField {
  key: string;
  label: string;
  type: "text" | "number" | "boolean" | "select" | "textarea" | "secret";
  defaultValue?: unknown; // 替换 any
  placeholder?: string;
  description?: string;
  required?: boolean;
  options?: Array<{ label: string; value: unknown }>; // 替换 any
}

/** 机器人视图对象（后端 vo.BotResponse） */
export interface BotVO {
  id: number;
  name: string;
  version: string;
  description: string;
  summary: string;
  avatarUrl: string;
  screenshots: string[];
  homepageUrl: string;
  type: BotType;
  tags: string[];
  creatorId: number;
  creatorName: string;
  scriptCode: string; // Lua 脚本
  scriptUrl?: string; // 可选外部 URL
  triggerType: BotTriggerType;
  cronExpr: string;
  eventFilter: string;
  timeoutSec: number;
  retryTimes: number;
  envVars: Record<string, string>;
  resourceLimit?: ResourceLimit;
  pricing: BotPricing;
  permissions: BotPermission[];
  enabled: boolean;
  status: BotStatus;
  execCount: number;
  lastExecAt?: string;
  errorMsg: string;
  configSchema: BotConfigField[];
  configValues: Record<string, unknown>; // 替换 Record<string, any>
  createdAt: string;
  updatedAt: string;
}

/** 机器人列表分页响应（后端返回格式） */
export interface BotListResponse {
  list: BotVO[];
  total: number;
  page: number;
}

/** 创建机器人请求（无 enabled 字段） */
export interface CreateBotRequest {
  name: string;
  version: string;
  description?: string;
  summary?: string;
  avatarUrl?: string;
  screenshots?: string[];
  homepageUrl?: string;
  type: BotType;
  tags?: string[];
  scriptCode: string;
  scriptUrl?: string;
  triggerType: BotTriggerType;
  cronExpr?: string;
  eventFilter?: string;
  timeoutSec?: number;
  retryTimes?: number;
  envVars?: Record<string, string>;
  resourceLimit?: ResourceLimit;
  pricing?: BotPricing;
  permissions?: BotPermission[];
  configSchema?: BotConfigField[];
  configValues?: Record<string, unknown>; // 替换 any
  enabled?: boolean;
}

/** 更新机器人请求（所有字段可选，包含 enabled） */
export interface UpdateBotRequest {
  name?: string;
  version?: string;
  description?: string;
  summary?: string;
  avatarUrl?: string;
  screenshots?: string[];
  homepageUrl?: string;
  type?: BotType;
  tags?: string[];
  scriptCode?: string;
  scriptUrl?: string;
  triggerType?: BotTriggerType;
  cronExpr?: string;
  eventFilter?: string;
  timeoutSec?: number;
  retryTimes?: number;
  envVars?: Record<string, string>;
  resourceLimit?: ResourceLimit;
  pricing?: BotPricing;
  permissions?: BotPermission[];
  configSchema?: BotConfigField[];
  configValues?: Record<string, unknown>; // 替换 any
  enabled?: boolean;
}
/** 手动触发时可携带的事件数据（任意 JSON 对象） */
export type RunEventData = Record<string, unknown>;

// ---------- 类型定义（与后端 Go 模型对齐）----------
export type TriggerType =
  | "on_schedule"
  | "on_new_post"
  | "on_new_comment"
  | "on_user_register"
  | "on_keyword"
  | "on_manual";

export type CondType =
  | "post_title_contains"
  | "post_content_contains"
  | "user_role_is"
  | "user_post_count_gte"
  | "board_id_in"
  | "time_range"
  | "custom_expr";

export type ActionType =
  | "reply_post"
  | "delete_post"
  | "hide_post"
  | "pin_post"
  | "lock_post"
  | "create_post"
  | "delete_comment"
  | "ban_user"
  | "send_message"
  | "webhook"
  | "notify_admin"
  | "wait"
  | "set_variable"
  | "stop_if";

export interface TriggerNode {
  type: TriggerType;
  params?: Record<string, unknown>;
}

export interface CondNode {
  type: CondType;
  negate?: boolean;
  params: Record<string, unknown>;
}

export interface ActionNode {
  type: ActionType;
  params: Record<string, unknown>;
}

export interface Flow {
  version: string;
  trigger: TriggerNode;
  conditions?: CondNode[];
  actions: ActionNode[];
}
