import { pluginRegistry } from "./PluginRegistry";
import type {
  PluginAPI,
  PluginEvent,
  PluginEventHandler,
  SlotName,
} from "./types";

interface UserContext {
  id: string;
  username: string;
  role: string;
}

interface PluginAPIOptions {
  pluginId: string;
  pluginName: string;
  getUser: () => UserContext | null;
  getLocale: () => string;
}

/**
 * 为每个插件创建独立的 API 沙箱
 * 插件只能访问这里暴露的能力，无法直接操作 store 或内部状态
 */
export function createPluginAPI(options: PluginAPIOptions): PluginAPI {
  const { pluginId, pluginName, getUser, getLocale } = options;
  const registeredHandlers: Array<{
    event: PluginEvent;
    handler: PluginEventHandler;
  }> = [];

  const api: PluginAPI = {
    registerSlot(slotName, component, slotOptions) {
      pluginRegistry.registerSlotComponent(slotName as SlotName, {
        pluginId,
        pluginName,
        component,
        order: slotOptions?.order ?? 0,
      });
    },

    on(event, handler) {
      registeredHandlers.push({ event, handler });
      pluginRegistry.on(event, handler);
    },

    off(event, handler) {
      pluginRegistry.off(event, handler);
    },

    getUser,
    getLocale,

    log(level, message) {
      const prefix = `[Plugin:${pluginName}]`;
      if (level === "error") console.error(prefix, message);
      else if (level === "warn") console.warn(prefix, message);
      else console.info(prefix, message);
    },
  };

  return api;
}
