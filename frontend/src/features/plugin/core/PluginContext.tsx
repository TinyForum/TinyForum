"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  useRef,
} from "react";
import { loadPlugins } from "../PluginLoader";
import { pluginRegistry } from "../PluginRegistry";
import { pluginApi } from "@/shared/api/modules/plugin/plugins";
import { PluginMeta, RegisteredPlugin } from "@/shared/api/types/plugin.model";
import { useAuthStore } from "@/store";

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
// 统一走 shared/api Client 层，复用 axios 实例（withCredentials、401 处理）
async function fetchEnabledPlugins(): Promise<PluginMeta[]> {
  try {
    const res = await pluginApi.listEnabled();
    return res.data.data?.list ?? [];
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
  const { isAuthenticated, isHydrated } = useAuthStore();

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
    if (!isHydrated) return;
    initializedRef.current = true;
    if (!isAuthenticated) {
      setIsInitialized(true);
      return;
    }
    load();
  }, [load, isAuthenticated, isHydrated]);

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
