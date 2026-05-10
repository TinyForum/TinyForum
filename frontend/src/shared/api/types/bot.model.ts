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
  defaultValue?: any;
  placeholder?: string;
  description?: string;
  required?: boolean;
  options?: Array<{ label: string; value: any }>;
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
  triggerType: BotTriggerType;
  cronExpr: string;
  eventFilter: string;
  timeoutSec: number;
  retryTimes: number;
  resourceLimit?: ResourceLimit;
  pricing: BotPricing;
  permissions: BotPermission[];
  enabled: boolean;
  status: BotStatus;
  execCount: number;
  lastExecAt?: string;
  errorMsg: string;
  configSchema: BotConfigField[];
  configValues: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

// ========== 请求/响应类型 ==========

/** 创建机器人请求体（与 request.CreateBotRequest 对应） */
// export interface CreateBotRequest {
//   name: string;
//   version: string;
//   description?: string;
//   summary?: string;
//   avatarUrl?: string;
//   screenshots?: string[];
//   homepageUrl?: string;
//   type: BotType;
//   tags?: string[];
//   scriptCode: string; // Lua 源码
//   scriptUrl?: string; // 外部 URL（与 scriptCode 二选一）
//   triggerType: BotTriggerType;
//   cronExpr?: string; // 当 triggerType = "schedule" 时必填
//   eventFilter?: string; // 当 triggerType = "event" 时使用
//   timeoutSec?: number; // 默认 10
//   retryTimes?: number; // 默认 0
//   envVars?: Record<string, string>;
//   resourceLimit?: ResourceLimit;
//   pricing?: BotPricing; // 可选，默认 free
//   permissions?: BotPermission[];
//   configSchema?: BotConfigField[];
//   configValues?: Record<string, any>;
// }

/** 更新机器人请求体（部分字段可选） */
// export interface UpdateBotRequest {
//   name?: string;
//   description?: string;
//   summary?: string;
//   avatarUrl?: string;
//   scriptCode?: string;
//   scriptUrl?: string;
//   triggerType?: BotTriggerType;
//   cronExpr?: string;
//   eventFilter?: string;
//   timeoutSec?: number;
//   retryTimes?: number;
//   envVars?: Record<string, string>;
//   resourceLimit?: ResourceLimit;
//   pricing?: BotPricing;
//   permissions?: BotPermission[];
//   configSchema?: BotConfigField[];
//   configValues?: Record<string, any>;
//   enabled?: boolean;
// }

/** 机器人列表分页响应（后端返回格式） */
export interface BotListResponse {
  list: BotVO[];
  total: number;
  page: number;
}
// 补充 BotVO 接口缺失字段
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
  envVars: Record<string, string>; // 环境变量
  resourceLimit?: ResourceLimit;
  pricing: BotPricing;
  permissions: BotPermission[];
  enabled: boolean; // 启用状态
  status: BotStatus;
  execCount: number;
  lastExecAt?: string;
  errorMsg: string;
  configSchema: BotConfigField[];
  configValues: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

// 创建机器人请求（无 enabled 字段）
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
  configValues?: Record<string, any>;
  enabled?: boolean;
}

// 更新机器人请求（所有字段可选，包含 enabled）
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
  configValues?: Record<string, any>;
  enabled?: boolean;
}
