import { pluginRegistry } from "./PluginRegistry";
import { createPluginAPI } from "./PluginAPI";
import type { PluginEntryFn, PluginMeta } from "./types";

interface LoaderOptions {
  getUser: () => { id: string; username: string; role: string } | null;
  getLocale: () => string;
}

const LOAD_TIMEOUT = 10_000; // 10秒超时

/**
 * 动态加载单个插件脚本
 * 插件脚本必须将入口函数挂载到 window.__plugins__[pluginId]
 */
async function loadPluginScript(
  scriptUrl: string,
  pluginId: string,
): Promise<PluginEntryFn> {
  return new Promise((resolve, reject) => {
    // 防止重复加载
    const existingScript = document.getElementById(`plugin-script-${pluginId}`);
    if (existingScript) {
      const fn = (window as Record<string, unknown>)[
        `__plugin_${pluginId}__`
      ] as PluginEntryFn | undefined;
      if (fn) return resolve(fn);
    }

    const script = document.createElement("script");
    script.id = `plugin-script-${pluginId}`;
    script.src = scriptUrl;
    script.async = true;
    script.crossOrigin = "anonymous";

    const timer = setTimeout(() => {
      script.remove();
      reject(new Error(`Plugin "${pluginId}" load timeout`));
    }, LOAD_TIMEOUT);

    script.onload = () => {
      clearTimeout(timer);
      const fn = (window as Record<string, unknown>)[
        `__plugin_${pluginId}__`
      ] as PluginEntryFn | undefined;
      if (typeof fn === "function") {
        resolve(fn);
      } else {
        reject(
          new Error(
            `Plugin "${pluginId}" did not export a valid entry function. ` +
              `Expected window.__plugin_${pluginId}__ to be a function.`,
          ),
        );
      }
    };

    script.onerror = () => {
      clearTimeout(timer);
      script.remove();
      reject(new Error(`Failed to load plugin script: ${scriptUrl}`));
    };

    document.head.appendChild(script);
  });
}

/**
 * 加载并初始化一个插件
 */
export async function loadPlugin(
  meta: PluginMeta,
  options: LoaderOptions,
): Promise<void> {
  pluginRegistry.registerPlugin(meta);

  try {
    const entryFn = await loadPluginScript(meta.scriptUrl, meta.id);

    const api = createPluginAPI({
      pluginId: meta.id,
      pluginName: meta.name,
      getUser: options.getUser,
      getLocale: options.getLocale,
    });

    await entryFn(api);

    pluginRegistry.updatePluginStatus(meta.id, "active");
    console.info(
      `[PluginLoader] Plugin "${meta.name}" v${meta.version} loaded successfully`,
    );
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    pluginRegistry.updatePluginStatus(meta.id, "error", message);
    console.error(
      `[PluginLoader] Failed to load plugin "${meta.name}":`,
      message,
    );
  }
}

/**
 * 批量加载所有启用的插件
 */
export async function loadPlugins(
  plugins: PluginMeta[],
  options: LoaderOptions,
): Promise<void> {
  const enabledPlugins = plugins.filter((p) => p.enabled);

  // 并行加载，单个失败不影响其他
  await Promise.allSettled(
    enabledPlugins.map((meta) => loadPlugin(meta, options)),
  );
}

/**
 * 卸载插件（清理插槽 + 移除脚本标签）
 */
export function unloadPlugin(pluginId: string): void {
  pluginRegistry.unregisterPlugin(pluginId);

  const script = document.getElementById(`plugin-script-${pluginId}`);
  if (script) script.remove();

  // 清理全局命名空间
  delete (window as Record<string, unknown>)[`__plugin_${pluginId}__`];

  console.info(`[PluginLoader] Plugin "${pluginId}" unloaded`);
}
