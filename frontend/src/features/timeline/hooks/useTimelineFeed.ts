// hooks/useTimelineFeed.ts
"use client";

import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { timelineApi } from "@/shared/api/modules/timeline";
import { timelineKeys } from "./useTimelineKeys";
import { PageData } from "@/shared/api/types/basic.model";
import { TimelineEvent } from "@/shared/api/types/timeline.model";

// 时间线动态查询 Hook：根据 Tab 类型与分页加载动态
export function useTimelineFeed(
  activeTab: "home" | "following",
  page: number,
  pageSize: number,
  enabled: boolean,
) {
  const t = useTranslations("Timeline");

  return useQuery<PageData<TimelineEvent>>({
    queryKey: [...timelineKeys.all, activeTab, page],
    queryFn: async () => {
      const response =
        activeTab === "home"
          ? await timelineApi.getHome({ page, page_size: pageSize })
          : await timelineApi.getFollowing({ page, page_size: pageSize });

      if (response.data.code !== 0) {
        throw new Error(response.data.message || t("loading_failed"));
      }
      if (!response.data.data) {
        throw new Error(t("timeline_data_empty"));
      }
      return response.data.data;
    },
    enabled,
  });
}
