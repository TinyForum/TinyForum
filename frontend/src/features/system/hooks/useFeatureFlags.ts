"use client";

import { useCallback, useState } from "react";
import { useTranslations } from "next-intl";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { configApi } from "@/shared/api/modules/config";

export interface Feature {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  category: string;
  // 后端是否支持持久化该开关；false 表示仅展示，不支持在线修改
  persistable: boolean;
}

interface FeatureDef {
  id: string;
  nameKey: string;
  descKey: string;
  categoryKey: string;
  // 后端无对应配置项的开关标记为不可持久化，避免制造"点了没用"的假开关
  persistable: boolean;
  defaultEnabled: boolean;
}

// 功能开关静态定义：文案经 next-intl 解析，禁止硬编码
const FEATURE_DEFS: FeatureDef[] = [
  {
    id: "user_registration",
    nameKey: "feature_user_registration",
    descKey: "feature_user_registration_desc",
    categoryKey: "category_user_management",
    persistable: false,
    defaultEnabled: true,
  },
  {
    id: "comments",
    nameKey: "feature_comments",
    descKey: "feature_comments_desc",
    categoryKey: "category_interaction",
    persistable: false,
    defaultEnabled: true,
  },
  {
    id: "content_review",
    nameKey: "feature_content_review",
    descKey: "feature_content_review_desc",
    categoryKey: "category_security",
    persistable: false,
    defaultEnabled: true,
  },
  {
    id: "api_access",
    nameKey: "feature_api_access",
    descKey: "feature_api_access_desc",
    categoryKey: "category_developer",
    persistable: true,
    defaultEnabled: false,
  },
  {
    id: "email_notification",
    nameKey: "feature_email_notification",
    descKey: "feature_email_notification_desc",
    categoryKey: "category_notification",
    persistable: false,
    defaultEnabled: true,
  },
  {
    id: "multilingual",
    nameKey: "feature_multilingual",
    descKey: "feature_multilingual_desc",
    categoryKey: "category_i18n",
    persistable: false,
    defaultEnabled: true,
  },
];

const basicConfigKey = ["config", "basic"] as const;

export function useFeatureFlags() {
  const t = useTranslations("System");
  const queryClient = useQueryClient();
  const [togglingId, setTogglingId] = useState<string | null>(null);

  // 通过 Query 获取 basic 配置（debug 字段映射 API 接口开关），替代在 useEffect 中命令式请求
  const { data } = useQuery({
    queryKey: basicConfigKey,
    queryFn: () => configApi.get("basic"),
    staleTime: 5 * 60 * 1000,
  });

  const debug = (data?.data?.data as Record<string, unknown> | undefined)?.debug;

  const features: Feature[] = FEATURE_DEFS.map((def) => ({
    id: def.id,
    name: t(def.nameKey),
    description: t(def.descKey),
    category: t(def.categoryKey),
    enabled:
      def.persistable && debug !== undefined
        ? !Boolean(debug)
        : def.defaultEnabled,
    persistable: def.persistable,
  }));

  const toggle = useCallback(
    async (id: string, enabled: boolean) => {
      const def = FEATURE_DEFS.find((d) => d.id === id);
      if (!def) return;
      if (!def.persistable) {
        toast.error(t("feature_toggle_unsupported"));
        return;
      }
      setTogglingId(id);
      try {
        await configApi.update("basic", { debug: !enabled });
        await queryClient.invalidateQueries({ queryKey: basicConfigKey });
        toast.success(
          enabled
            ? t("feature_enabled", { name: t(def.nameKey) })
            : t("feature_disabled", { name: t(def.nameKey) }),
        );
      } catch {
        toast.error(t("operation_failed_retry"));
      } finally {
        setTogglingId(null);
      }
    },
    [queryClient, t],
  );

  const enableAll = useCallback(async () => {
    const persistableDefs = FEATURE_DEFS.filter((d) => d.persistable);
    try {
      await Promise.all(
        persistableDefs.map(() => configApi.update("basic", { debug: false })),
      );
      await queryClient.invalidateQueries({ queryKey: basicConfigKey });
      toast.success(t("all_features_enabled"));
    } catch {
      toast.error(t("operation_failed_retry"));
    }
  }, [queryClient, t]);

  const grouped = features.reduce<Record<string, Feature[]>>((acc, f) => {
    (acc[f.category] ??= []).push(f);
    return acc;
  }, {});

  const enabledCount = features.filter((f) => f.enabled).length;

  return { features, grouped, enabledCount, toggle, enableAll, togglingId };
}
