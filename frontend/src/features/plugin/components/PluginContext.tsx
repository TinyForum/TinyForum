"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  useRef,
} from "react";
import { PluginMeta, RegisteredPlugin } from "@/shared/type/plugin.type";
import { loadPlugins } from "../PluginLoader";
import { pluginRegistry } from "../PluginRegistry";

interface PluginContextValue {
  plugins: RegisteredPlugin[];
  isInitialized: boolean;
  isLoading: boolean;
  reload: () => Promise<void>;
}

const PluginContext = createContext<PluginContextValue>({
  plugins: [],
  isInitialized: false,
  isLoading: false,
  reload: async () => {},
});

// ── 直接请求 API，不依赖 hook ────────────────────────────────────────────────
// PluginContext 是全局 Provider，不能用 useQuery（它在 QueryClientProvider 内部）
// 用 fetch 直接请求，避免循环依赖和 Hook 规则问题
async function fetchEnabledPlugins(): Promise<PluginMeta[]> {
  try {
    const res = await fetch("/api/v1/plugins?enabled=true");
    console.log("[PluginContext] fetchEnabledPlugins: HTTP", res.status);

    if (!res.ok) {
      console.warn("[PluginContext] fetchEnabledPlugins: HTTP", res.status);
      return [];
    }
    const json = await res.json();
    // 兼容两种后端返回结构：
    // 1. { data: { list: [...] } }  — 分页结构
    // 2. { data: [...] }            — 数组结构
    const payload = json?.data;
    if (Array.isArray(payload)) return payload;
    if (Array.isArray(payload?.list)) return payload.list;
    if (Array.isArray(payload?.data)) return payload.data;
    if (Array.isArray(payload?.data?.list)) return payload.data.list;
    console.warn("[PluginContext] fetchEnabledPlugins: unexpected shape", json);
    return [];
  } catch (err) {
    console.error("[PluginContext] fetchEnabledPlugins error:", err);
    return [];
  }
}

export function PluginProvider({
  children,
  getUser,
  // getLocale,
}: {
  children: React.ReactNode;
  getUser: () => { id: string; username: string; role: string } | null;
  // getLocale: () => string;
}) {
  const [plugins, setPlugins] = useState<RegisteredPlugin[]>([]);
  const [isInitialized, setIsInitialized] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const initializedRef = useRef(false);

  // 订阅 registry 变更，驱动 UI 更新
  useEffect(() => {
    const unsubscribe = pluginRegistry.subscribe(() => {
      setPlugins([...pluginRegistry.getAllPlugins()]);
    });
    return unsubscribe;
  }, []);

  const load = useCallback(async () => {
    setIsLoading(true);
    try {
      const metas = await fetchEnabledPlugins();
      console.info(`[PluginContext] Loading ${metas.length} enabled plugins`);
      if (metas.length === 0) {
        console.warn(
          "[PluginContext] No enabled plugins returned from API. " +
            "Check: 1) API /api/v1/plugins?enabled=true is implemented, " +
            "2) Plugins exist in DB with enabled=true",
        );
      }
      await loadPlugins(metas, { getUser });
    } finally {
      setIsLoading(false);
      setIsInitialized(true);
    }
  }, [getUser]);

  useEffect(() => {
    if (initializedRef.current) return;
    initializedRef.current = true;
    load();
  }, [load]);

  return (
    <PluginContext.Provider
      value={{ plugins, isInitialized, isLoading, reload: load }}
    >
      {children}
    </PluginContext.Provider>
  );
}

export function usePluginsContext() {
  return useContext(PluginContext);
}
