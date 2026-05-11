"use client";

import { useState, useCallback } from "react";
import toast from "react-hot-toast";

export interface Feature {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  category: string;
}

const INITIAL_FEATURES: Feature[] = [
  {
    id: "1",
    name: "用户注册",
    description: "允许新用户注册账号",
    enabled: true,
    category: "用户管理",
  },
  {
    id: "2",
    name: "评论系统",
    description: "允许用户在内容下发表评论",
    enabled: true,
    category: "互动功能",
  },
  {
    id: "3",
    name: "内容审核",
    description: "新发布内容需要管理员审核",
    enabled: true,
    category: "安全策略",
  },
  {
    id: "4",
    name: "API 接口",
    description: "开放 RESTful API 接口供第三方调用",
    enabled: false,
    category: "开发者工具",
  },
  {
    id: "5",
    name: "邮件通知",
    description: "发送系统邮件通知（如注册确认、密码重置）",
    enabled: true,
    category: "通知系统",
  },
  {
    id: "6",
    name: "多语言支持",
    description: "启用多语言界面切换功能",
    enabled: true,
    category: "国际化",
  },
];

export function useFeatureFlags() {
  const [features, setFeatures] = useState<Feature[]>(INITIAL_FEATURES);
  const [togglingId, setTogglingId] = useState<string | null>(null);

  const toggle = useCallback(
    async (id: string, enabled: boolean) => {
      const feature = features.find((f) => f.id === id);
      if (!feature) return;
      setTogglingId(id);
      try {
        // TODO: await apiClient.patch(`/system/features/${id}`, { enabled });
        await new Promise((r) => setTimeout(r, 300));
        setFeatures((prev) =>
          prev.map((f) => (f.id === id ? { ...f, enabled } : f)),
        );
        toast.success(`${feature.name} 已${enabled ? "开启" : "关闭"}`);
      } catch {
        toast.error("操作失败，请重试");
      } finally {
        setTogglingId(null);
      }
    },
    [features],
  );

  const enableAll = useCallback(() => {
    setFeatures((prev) => prev.map((f) => ({ ...f, enabled: true })));
    toast.success("已启用所有功能");
  }, []);

  const grouped = features.reduce<Record<string, Feature[]>>((acc, f) => {
    (acc[f.category] ??= []).push(f);
    return acc;
  }, {});

  const enabledCount = features.filter((f) => f.enabled).length;

  return { features, grouped, enabledCount, toggle, enableAll, togglingId };
}
