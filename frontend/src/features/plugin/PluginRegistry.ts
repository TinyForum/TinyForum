import {
  RegisteredPlugin,
  SlotComponent,
  PluginEventHandler,
  PluginMeta,
  SlotName,
  PluginEvent,
} from "@/shared/type/plugin.type";

type Listener = () => void;

class PluginRegistry {
  private static instance: PluginRegistry;

  private plugins: Map<string, RegisteredPlugin> = new Map();
  private slots: Map<string, SlotComponent[]> = new Map();
  private eventListeners: Map<string, PluginEventHandler[]> = new Map();
  private changeListeners: Set<Listener> = new Set();

  static getInstance(): PluginRegistry {
    if (!PluginRegistry.instance) {
      PluginRegistry.instance = new PluginRegistry();
    }
    return PluginRegistry.instance;
  }

  // ── 插件注册 ─────────────────────────────────
  registerPlugin(meta: PluginMeta): void {
    this.plugins.set(meta.id, { meta, status: "loading" });
    this.notify();
  }

  updatePluginStatus(
    id: string,
    status: RegisteredPlugin["status"],
    error?: string,
  ): void {
    const existing = this.plugins.get(id);
    if (existing) {
      this.plugins.set(id, { ...existing, status, error });
      this.notify();
    }
  }

  getPlugin(id: string): RegisteredPlugin | undefined {
    return this.plugins.get(id);
  }

  getAllPlugins(): RegisteredPlugin[] {
    return Array.from(this.plugins.values());
  }
  getPluginConfig(pluginId: string): Record<string, unknown> {
    return this.plugins.get(pluginId)?.meta.config ?? {};
  }

  // ── 插槽管理 ─────────────────────────────────
  registerSlotComponent(
    slotName: SlotName | string,
    component: SlotComponent,
  ): void {
    const existing = this.slots.get(slotName) ?? [];
    const sorted = [...existing, component].sort(
      (a, b) => (a.order ?? 0) - (b.order ?? 0),
    );
    this.slots.set(slotName, sorted);
    this.notify();
  }

  getSlotComponents(slotName: string): SlotComponent[] {
    return this.slots.get(slotName) ?? [];
  }

  removePluginSlots(pluginId: string): void {
    for (const [key, components] of this.slots.entries()) {
      this.slots.set(
        key,
        components.filter((c) => c.pluginId !== pluginId),
      );
    }
    this.notify();
  }

  // ── 事件系统 ─────────────────────────────────
  on(event: PluginEvent, handler: PluginEventHandler): void {
    const handlers = this.eventListeners.get(event) ?? [];
    this.eventListeners.set(event, [...handlers, handler]);
  }

  off(event: PluginEvent, handler: PluginEventHandler): void {
    const handlers = this.eventListeners.get(event) ?? [];
    this.eventListeners.set(
      event,
      handlers.filter((h) => h !== handler),
    );
  }

  emit(event: PluginEvent, data: unknown): void {
    const handlers = this.eventListeners.get(event) ?? [];
    handlers.forEach((h) => {
      try {
        h(data);
      } catch (err) {
        console.error(
          `[PluginRegistry] Event handler error for ${event}:`,
          err,
        );
      }
    });
  }

  // ── 订阅变更 ─────────────────────────────────
  subscribe(listener: Listener): () => void {
    this.changeListeners.add(listener);
    return () => this.changeListeners.delete(listener);
  }

  private notify(): void {
    this.changeListeners.forEach((l) => l());
  }

  // ── 清理 ─────────────────────────────────────
  unregisterPlugin(id: string): void {
    this.removePluginSlots(id);
    this.plugins.delete(id);
    this.notify();
  }

  reset(): void {
    this.plugins.clear();
    this.slots.clear();
    this.eventListeners.clear();
    this.notify();
  }
}

export const pluginRegistry = PluginRegistry.getInstance();
