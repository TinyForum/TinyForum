import { createPluginAPI } from "./PluginAPI";
import { pluginRegistry } from "./PluginRegistry";
import { PluginEntryFn, PluginMeta } from "@/shared/type/plugin.type";

interface LoaderOptions {
  getUser: () => { id: string; username: string; role: string } | null;
  // getLocale: () => string;
}
import * as React from "react";
import * as ReactDOM from "react-dom";
window.React = React;
window.ReactDOM = ReactDOM;

const LOAD_TIMEOUT = 10_000;

// ── window 类型扩展，允许动态键访问 ──────────────────────────────────────────
// 解决 TypeScript "Element implicitly has an 'any' type" 错误
declare global {
  interface Window {
    [key: string]: unknown;
  }
}

async function loadPluginScript(
  scriptUrl: string,
  pluginId: string,
): Promise<PluginEntryFn> {
  // window 的 key 规则：插件 id 中的 "-" 替换为 "_"
  const windowKey = `__plugin_${pluginId.replace(/-/g, "_")}__`;
  console.log(`Loading plugin "${pluginId}" from "${scriptUrl}"...`);

  return new Promise((resolve, reject) => {
    // 防止重复加载
    const existingScript = document.getElementById(`plugin-script-${pluginId}`);
    if (existingScript) {
      const fn = window[windowKey] as PluginEntryFn | undefined;
      if (typeof fn === "function") return resolve(fn);
      // 脚本存在但函数未挂载，说明上次加载失败，移除旧标签重试
      existingScript.remove();
    }

    const script = document.createElement("script");
    script.id = `plugin-script-${pluginId}`;
    script.src = scriptUrl;
    script.async = true;
    script.crossOrigin = "anonymous";

    const timer = setTimeout(() => {
      script.remove();
      reject(
        new Error(`Plugin "${pluginId}" load timeout (${LOAD_TIMEOUT}ms)`),
      );
    }, LOAD_TIMEOUT);

    script.onload = () => {
      clearTimeout(timer);
      const fn = window[windowKey] as PluginEntryFn | undefined;
      if (typeof fn === "function") {
        resolve(fn);
      } else {
        reject(
          new Error(
            `Plugin "${pluginId}" script loaded but entry function not found.\n` +
              `Expected: window.${windowKey} to be a function.\n` +
              `Got: ${typeof fn}\n` +
              `Make sure your plugin script contains:\n` +
              `  window.${windowKey} = function(api) { ... }`,
          ),
        );
      }
    };

    script.onerror = (event) => {
      clearTimeout(timer);
      script.remove();
      reject(
        new Error(
          `Failed to load plugin script: ${scriptUrl}\n` +
            `Check: 1) URL is accessible, 2) CORS allows this origin, 3) URL is HTTPS`,
        ),
      );
    };

    document.head.appendChild(script);
  });
}

export async function loadPlugin(
  meta: PluginMeta,
  options: LoaderOptions,
): Promise<void> {
  pluginRegistry.registerPlugin(meta);
  console.log(
    "[PluginLoader] Loading plugin: ",
    "ID: ",
    meta.id,
    " Name: ",
    meta.name,
    " Version: ",
    meta.version,
  );

  try {
    const entryFn = await loadPluginScript(meta.scriptUrl, meta.id.toString());

    const api = createPluginAPI({
      pluginId: meta.id,
      pluginName: meta.name,
      getUser: options.getUser,
    });

    await entryFn(api);

    pluginRegistry.updatePluginStatus(meta.id, "active");
    console.info(
      `[PluginLoader] ✅ Plugin "${meta.id}" v${meta.version} loaded`,
    );
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    pluginRegistry.updatePluginStatus(meta.id, "error", message);
    console.error(`[PluginLoader] ❌ Plugin "${meta.name}" failed:`, message);
  }
}

export async function loadPlugins(
  plugins: PluginMeta[],
  options: LoaderOptions,
): Promise<void> {
  // 防御：确保传入的是数组
  if (!Array.isArray(plugins)) {
    console.error(
      "[PluginLoader] loadPlugins: expected PluginMeta[], got:",
      typeof plugins,
      plugins,
    );
    return;
  }

  const enabledPlugins = plugins.filter((p) => p.enabled);
  console.log("plugin: ", plugins);
  console.info(
    `[PluginLoader] Loading ${enabledPlugins.length}/${plugins.length} enabled plugins`,
  );

  // 并行加载，单个失败不阻塞其他
  await Promise.allSettled(
    enabledPlugins.map((meta) => loadPlugin(meta, options)),
  );
}

export function unloadPlugin(pluginId: string): void {
  pluginRegistry.unregisterPlugin(pluginId);

  const script = document.getElementById(`plugin-script-${pluginId}`);
  if (script) script.remove();

  // 清理 window 挂载（使用扩展后的类型，无需类型断言）
  const windowKey = `__plugin_${pluginId.replace(/-/g, "_")}__`;
  delete window[windowKey];

  console.info(`[PluginLoader] Plugin "${pluginId}" unloaded`);
}
