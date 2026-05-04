"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  useRef,
} from "react";
import { pluginRegistry } from "./PluginRegistry";
import { loadPlugins } from "./PluginLoader";
import type { PluginMeta, RegisteredPlugin } from "./types";

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

async function fetchEnabledPlugins(): Promise<PluginMeta[]> {
  try {
    const res = await fetch("/api/v1/plugins?enabled=true");
    if (!res.ok) return [];
    const json = await res.json();
    return json?.data ?? [];
  } catch {
    return [];
  }
}

export function PluginProvider({
  children,
  getUser,
  getLocale,
}: {
  children: React.ReactNode;
  getUser: () => { id: string; username: string; role: string } | null;
  getLocale: () => string;
}) {
  const [plugins, setPlugins] = useState<RegisteredPlugin[]>([]);
  const [isInitialized, setIsInitialized] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const initializedRef = useRef(false);

  // 订阅 registry 变更
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
      await loadPlugins(metas, { getUser, getLocale });
    } finally {
      setIsLoading(false);
      setIsInitialized(true);
    }
  }, [getUser, getLocale]);

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

export function usePlugins() {
  return useContext(PluginContext);
}
