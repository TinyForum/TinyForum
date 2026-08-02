// hooks/useTimelineSubscriptions.ts
import { useCallback } from "react";
import { useQuery } from "@tanstack/react-query";
import { timelineApi } from "@/shared/api/modules/timeline";
import { timelineKeys } from "./useTimelineKeys";
import { Subscription } from "@/shared/api/types/timeline.model";

// 时间线订阅列表查询 Hook
export function useTimelineSubscriptions(enabled: boolean) {
  const query = useQuery<Subscription[]>({
    queryKey: [...timelineKeys.all, "subscriptions"],
    queryFn: async () => {
      const response = await timelineApi.getSubscriptions();
      if (response.data.code !== 0) {
        throw new Error(response.data.message || "加载失败");
      }
      if (!response.data.data) {
        throw new Error("订阅列表数据为空");
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
