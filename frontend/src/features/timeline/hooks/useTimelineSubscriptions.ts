// hooks/useTimelineSubscriptions.ts
"use client";

import { useCallback } from "react";
import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { timelineApi } from "@/shared/api/modules/timeline";
import { timelineKeys } from "./useTimelineKeys";
import { Subscription } from "@/shared/api/types/timeline.model";

// 时间线订阅列表查询 Hook
export function useTimelineSubscriptions(enabled: boolean) {
  const t = useTranslations("Timeline");

  const query = useQuery<Subscription[]>({
    queryKey: [...timelineKeys.all, "subscriptions"],
    queryFn: async () => {
      const response = await timelineApi.getSubscriptions();
      if (response.data.code !== 0) {
        throw new Error(response.data.message || t("loading_failed"));
      }
      if (!response.data.data) {
        throw new Error(t("subscription_data_empty"));
      }
      return response.data.data;
    },
    enabled,
  });

  // 取消关注
  const unsubscribe = useCallback(async (userId: number) => {
    const response = await timelineApi.unsubscribe(userId);
    return response.data;
  }, []);

  return {
    data: query.data,
    unsubscribe,
  };
}
