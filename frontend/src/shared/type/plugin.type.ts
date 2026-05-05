// ─── Plugin System Types ─────────────────────────────────────────────────────

// ── 运行状态 ──────────────────────────────────────────────────────────────────
export type PluginStatus = "active" | "inactive" | "error" | "loading";

// ── 插件类型 ──────────────────────────────────────────────────────────────────
export type PluginType =
  | "ui" // 前端组件注入（PluginSlot 渲染）
  | "backend" // 纯后端逻辑扩展（服务端加载，不注入前端）
  | "lib" // 通用工具库（被其他插件依赖）
  | "app" // 独立子应用（注册独立路由页面）
  | "miniapp"; // 小程序/H5 轻应用

// ── 价格模型 ──────────────────────────────────────────────────────────────────
export type PluginPricingModel =
  | "free" // 完全免费
  | "freemium" // 基础免费，高级付费
  | "paid"; // 纯付费

export interface PluginPricing {
  model: PluginPricingModel;
  /** 付费/高级版价格，单位：元 */
  price?: number;
  /** 计费周期 */
  cycle?: "once" | "monthly" | "yearly";
  /** 免费版功能限制说明（freemium 时填写） */
  freeLimit?: string;
  /** 购买/授权链接 */
  purchaseUrl?: string;
}

// ── 插件分类（行业标准）─────────────────────────────────────────────────────────
export const PLUGIN_CATEGORIES = {
  COMMERCE: "commerce", // 电商/O2O/分销/返利
  MARKETING: "marketing", // 营销推广/活动/投票/SEO
  CONTENT: "content", // 内容管理/知识付费/小说/漫画/问答
  COMMUNITY: "community", // 社区/论坛/社交/评论/私信
  MEDIA: "media", // 音视频/直播/音乐/图库
  EDUCATION: "education", // 教育培训/在线课程/考试
  SUPPORT: "support", // 客服/工单/聊天/问卷/IM
  PAYMENT: "payment", // 支付网关/订阅/发票
  SECURITY: "security", // 认证/登录/权限/验证/防刷
  STORAGE: "storage", // 存储/OSS/COS/CDN/备份
  ANALYTICS: "analytics", // 数据统计/分析/报表
  DESIGN: "design", // UI组件/主题/样式/布局
  DEVTOOLS: "devtools", // API扩展/后端逻辑/开发工具
  INTEGRATION: "integration", // 第三方集成/Webhook/通知
  UTILITY: "utility", // 系统工具/增强功能/任务调度
} as const;

export type PluginCategory =
  (typeof PLUGIN_CATEGORIES)[keyof typeof PLUGIN_CATEGORIES];

export const PLUGIN_CATEGORY_LABELS: Record<PluginCategory, string> = {
  [PLUGIN_CATEGORIES.COMMERCE]: "电商/O2O",
  [PLUGIN_CATEGORIES.MARKETING]: "营销推广",
  [PLUGIN_CATEGORIES.CONTENT]: "内容/知识付费",
  [PLUGIN_CATEGORIES.COMMUNITY]: "社区/互动",
  [PLUGIN_CATEGORIES.MEDIA]: "媒体/音视频",
  [PLUGIN_CATEGORIES.EDUCATION]: "教育培训",
  [PLUGIN_CATEGORIES.SUPPORT]: "客服/工单/IM",
  [PLUGIN_CATEGORIES.PAYMENT]: "支付/财务",
  [PLUGIN_CATEGORIES.SECURITY]: "安全/认证",
  [PLUGIN_CATEGORIES.STORAGE]: "存储/备份",
  [PLUGIN_CATEGORIES.ANALYTICS]: "统计/分析",
  [PLUGIN_CATEGORIES.DESIGN]: "UI/主题设计",
  [PLUGIN_CATEGORIES.DEVTOOLS]: "开发/API",
  [PLUGIN_CATEGORIES.INTEGRATION]: "外部集成",
  [PLUGIN_CATEGORIES.UTILITY]: "系统工具",
};

// ── 版本兼容性 ────────────────────────────────────────────────────────────────
export type CompatVersion = "v0" | "v1" | "v2" | "v3";

export interface PluginCompatibility {
  /** 兼容的平台版本，空数组 = 全版本兼容 */
  versions: CompatVersion[];
  /** 最低 Node.js 版本（可选） */
  minNode?: string;
  /** 依赖的其他插件 ID */
  requires?: string[];
  /** 与哪些插件冲突 */
  conflicts?: string[];
}

// ── 配置字段 Schema ───────────────────────────────────────────────────────────
export interface PluginConfigField {
  key: string;
  label: string;
  type:
    | "text"
    | "number"
    | "boolean"
    | "select"
    | "textarea"
    | "color"
    | "url"
    | "secret";
  defaultValue?: unknown;
  placeholder?: string;
  description?: string;
  required?: boolean;
  options?: Array<{ label: string; value: string | number | boolean }>;
  min?: number;
  max?: number;
}

export type PluginConfigSchema = PluginConfigField[];

// ── 权限声明 ──────────────────────────────────────────────────────────────────
export type PluginPermission =
  | "read:user"
  | "read:posts"
  | "write:posts"
  | "read:comments"
  | "write:comments"
  | "read:settings"
  | "write:settings"
  | "send:email"
  | "send:sms"
  | "storage:read"
  | "storage:write"
  | "payment"
  | "network";

// ── 核心 PluginMeta ───────────────────────────────────────────────────────────
export interface PluginMeta {
  // 基础标识
  id: string;
  name: string;
  version: string;
  description: string;
  summary?: string;
  iconUrl?: string;
  screenshots?: string[];
  homepageUrl?: string;

  // 分类与类型
  type: PluginType;
  category: PluginCategory;
  tags?: string[];

  // 作者信息
  author: string;
  authorEmail?: string;
  authorUrl?: string;

  // 加载配置
  scriptUrl: string;
  serverEntry?: string;
  slots?: string[];
  routes?: string[];

  // 价格与兼容性
  pricing: PluginPricing;
  compatibility: PluginCompatibility;

  // 权限
  permissions?: PluginPermission[];

  // 运行时（服务端写入，前端只读）
  enabled: boolean;
  status?: PluginStatus;
  installCount?: number;
  rating?: number; // 0~5

  // 配置
  configSchema?: PluginConfigSchema;
  config?: Record<string, unknown>;

  // 时间戳
  createdAt?: string;
  updatedAt?: string;
}

// ── Runtime types ─────────────────────────────────────────────────────────────
export interface RegisteredPlugin {
  meta: PluginMeta;
  status: PluginStatus;
  error?: string;
}

export interface SlotComponent {
  pluginId: string;
  pluginName: string;
  component: React.ComponentType<Record<string, unknown>>;
  props?: Record<string, unknown>;
  order?: number;
}

// ── Plugin API Sandbox ────────────────────────────────────────────────────────
export interface PluginAPI {
  registerSlot(
    slotName: string,
    component: React.ComponentType<Record<string, unknown>>,
    options?: { order?: number },
  ): void;
  on(event: PluginEvent, handler: PluginEventHandler): void;
  off(event: PluginEvent, handler: PluginEventHandler): void;
  getUser(): { id: string; username: string; role: string } | null;
  getLocale(): string;
  getConfig(): Record<string, unknown>;
  log(level: "info" | "warn" | "error", message: string): void;
}

export type PluginEvent =
  | "post:view"
  | "post:create"
  | "post:delete"
  | "user:login"
  | "user:logout"
  | "comment:create"
  | "order:create"
  | "payment:success";

export type PluginEventHandler = (data: unknown) => void;
export type PluginEntryFn = (api: PluginAPI) => void | Promise<void>;

// ── Slot Names ────────────────────────────────────────────────────────────────
export const SLOT_NAMES = [
  "sidebar-top",
  "sidebar-bottom",
  "navbar-extra",
  "post-list-top",
  "post-list-bottom",
  "post-detail-bottom",
  "post-detail-sidebar",
  "dashboard-widget",
  "profile-extra",
  "profile-sidebar",
  "home-banner",
  "home-after-feed",
  "footer-extra",
  "payment-checkout",
  "user-card-extra",
] as const;

export type SlotName = (typeof SLOT_NAMES)[number];
