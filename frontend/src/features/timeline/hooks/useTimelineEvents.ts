// hooks/useTimelineEvents.ts
import { useQuery } from "@tanstack/react-query";
import { timelineApi } from "@/shared/api/modules/timeline";
import { TimelineEvent } from "@/shared/api/types/timeline.model";

export const useTimelineEvents = (enabled: boolean) => {
  return useQuery({
    queryKey: ["timeline-events"],
    queryFn: async () => {
      const res = await timelineApi.getFollowing({ page: 1, page_size: 5 });
      return (res.data.data?.list as TimelineEvent[]) || [];
    },
    enabled,
  });
};
