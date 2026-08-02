// hooks/useBoardPostLike.ts
import { useMutation } from "@tanstack/react-query";
import { postApi } from "@/shared/api/modules/posts";

// 板块帖子卡片：点赞/取消点赞
export function useBoardPostLike() {
  return {
    like: useMutation({
      mutationFn: (id: number) => postApi.like(id),
    }),
    unlike: useMutation({
      mutationFn: (id: number) => postApi.unlike(id),
    }),
  };
}
