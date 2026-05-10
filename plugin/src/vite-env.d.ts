// src/vite-env.d.ts
/// <reference types="vite/client" />

declare const __PLUGIN_SLUG__: string;

interface Window {
  // 允许动态的 __plugin_xxx__ 属性，值为插件入口函数
  [key: `__plugin_${string}__`]: (api: import('./PLUGIN.type').PluginAPI) => Promise<void>;
}