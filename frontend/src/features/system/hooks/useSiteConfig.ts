"use client";

import { useState, useCallback } from "react";
import toast from "react-hot-toast";

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
  siteName: "我的网站",
  siteDescription: "这是一个功能强大的网站平台",
  siteKeywords: "系统,管理,配置",
  adminEmail: "admin@example.com",
  itemsPerPage: 20,
  enableMaintenanceMode: false,
  defaultTheme: "auto",
};

export function useSiteConfig() {
  const [config, setConfig] = useState<SiteConfig>(DEFAULT_CONFIG);
  const [isSaving, setIsSaving] = useState(false);

  const update = useCallback(
    <K extends keyof SiteConfig>(key: K, value: SiteConfig[K]) => {
      setConfig((prev) => ({ ...prev, [key]: value }));
    },
    [],
  );

  const save = useCallback(async () => {
    setIsSaving(true);
    try {
      // TODO: await apiClient.put("/system/config", config);
      await new Promise((r) => setTimeout(r, 600));
      toast.success("网站配置保存成功");
    } catch {
      toast.error("保存失败，请重试");
    } finally {
      setIsSaving(false);
    }
  }, [config]);

  return { config, update, save, isSaving };
}
