// hooks/useProfileContent.ts
import { useQuery } from "@tanstack/react-query";
import { postApi } from "@/shared/api/modules/posts";
import { Post, PostType } from "@/shared/api/types/post.model";
import { PageData } from "@/shared/api/types/basic.model";

// 用户内容（帖子/文章/问答/话题）Hook
export function useProfileContent(userId: number, tab: PostType) {
  return useQuery<PageData<Post>>({
    queryKey: ["user-posts", userId, tab],
    queryFn: async () => {
      const res = await postApi.list({
        author_id: userId,
        type: ["article", "post", "question"].includes(tab) ? tab : "topic",
        page: 1,
        page_size: 20,
      });
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取帖子失败");
      }
      if (!res.data.data) {
        throw new Error("获取帖子失败");
      }
      return res.data.data;
    },
  });
}
