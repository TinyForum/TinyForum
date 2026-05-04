// ─── Plugin System Types ───────────────────────────────────────────────────

export type PluginStatus = "active" | "inactive" | "error" | "loading";

export interface PluginMeta {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  scriptUrl: string;
  enabled: boolean;
  status?: PluginStatus;
  config?: Record<string, unknown>;
  slots?: string[];
  createdAt?: string;
  updatedAt?: string;
}

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

// API沙箱接口 —— 只暴露受控能力给插件
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
  log(level: "info" | "warn" | "error", message: string): void;
}

export type PluginEvent =
  | "post:view"
  | "post:create"
  | "post:delete"
  | "user:login"
  | "user:logout"
  | "comment:create";

export type PluginEventHandler = (data: unknown) => void;

// 插件入口函数签名
export type PluginEntryFn = (api: PluginAPI) => void | Promise<void>;

// 挂载的插槽名
export const SLOT_NAMES = [
  "sidebar-top",
  "sidebar-bottom",
  "navbar-extra",
  "post-list-top",
  "post-list-bottom",
  "post-detail-bottom",
  "dashboard-widget",
  "profile-extra",
] as const;

export type SlotName = (typeof SLOT_NAMES)[number];
