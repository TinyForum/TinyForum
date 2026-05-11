
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
