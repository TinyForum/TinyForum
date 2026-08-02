// hooks/useFollowTopic.ts
import { useMutation } from "@tanstack/react-query";
import { topicApi } from "@/shared/api/modules/topics";

// 关注/取消关注话题
export function useFollowTopic() {
  return {
    follow: useMutation({
      mutationFn: (id: number) => topicApi.follow(id),
    }),
    unfollow: useMutation({
      mutationFn: (id: number) => topicApi.unfollow(id),
    }),
  };
}
