"use client";

import { useState, useCallback, useEffect } from "react";
import toast from "react-hot-toast";
import { configApi } from "@/shared/api/modules/config";
import type { ConfigReloadResult } from "@/shared/api/modules/config";

export interface SiteConfig {
  siteName: string;
  siteDescription: string;
  siteKeywords: string;
  adminEmail: string;
  itemsPerPage: number;
  enableMaintenanceMode: boolean;
  defaultTheme: "light" | "dark" | "auto";
}

const DEFAULT_CONFIG: SiteConfig = {
  siteName: "TinyForum",
  siteDescription: "一个轻量级论坛平台",
  siteKeywords: "论坛,社区,讨论",
  adminEmail: "admin@example.com",
  itemsPerPage: 20,
  enableMaintenanceMode: false,
  defaultTheme: "auto",
};

export function useSiteConfig() {
  const [config, setConfig] = useState<SiteConfig>(DEFAULT_CONFIG);
  const [isSaving, setIsSaving] = useState(false);
  const [isReloading, setIsReloading] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    configApi
      .get("app")
      .then((res) => {
        if (res.data.data) {
          setConfig((prev) => ({ ...prev, ...(res.data.data as Partial<SiteConfig>) }));
        }
      })
      .catch(() => {
        // use defaults if config fetch fails
      })
      .finally(() => setIsLoading(false));
  }, []);

  const update = useCallback(
    <K extends keyof SiteConfig>(key: K, value: SiteConfig[K]) => {
      setConfig((prev) => ({ ...prev, [key]: value }));
    },
    [],
  );

  const save = useCallback(async () => {
    setIsSaving(true);
    try {
      await configApi.update("app", config as unknown as Record<string, unknown>);
      toast.success("配置保存成功");
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })
        ?.response?.data?.message || "保存失败，请重试";
      toast.error(msg);
    } finally {
      setIsSaving(false);
    }
  }, [config]);

  const reloadConfig = useCallback(async (): Promise<ConfigReloadResult | null> => {
    setIsReloading(true);
    try {
      const res = await configApi.reload();
      const result = res.data.data;
      if (result) {
        toast.success(`配置已重新加载 (${result.reloadedFiles.length} 个文件)`);
      }
      return result || null;
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })
        ?.response?.data?.message || "重新加载配置失败";
      toast.error(msg);
      return null;
    } finally {
      setIsReloading(false);
    }
  }, []);

  return { config, update, save, reloadConfig, isSaving, isReloading, isLoading };
}
