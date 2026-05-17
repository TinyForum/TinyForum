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

export interface BotDO {
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
  configValues?: Record<string, unknown>;
  enabled?: boolean;
}
