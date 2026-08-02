// hooks/useExploreData.ts
import { useQuery } from "@tanstack/react-query";
import { toast } from "react-hot-toast";
import { postApi } from "@/shared/api/modules/posts";
import { tagApi } from "@/shared/api/modules/tags";
import { topicApi } from "@/shared/api/modules/topics";
import { userApi } from "@/shared/api/modules/user";
import { Post } from "@/shared/api/types/post.model";
import { Tag } from "@/shared/api/types/tag.model";
import { Topic } from "@/shared/api/types/topic.model";
import { LeaderboardItemResponse } from "@/shared/api/types/user.model";

// 探索页数据查询 Hook：聚合帖子、标签、话题与活跃用户
export function useExploreData(activeTab: string, sortBy?: string) {
  return useQuery<{
    posts: Post[];
    hotTags: Tag[];
    hotTopics: Topic[];
    activeUsers: LeaderboardItemResponse[];
  }>({
    queryKey: ["explore", activeTab],
    queryFn: async () => {
      try {
        const [postsResponse, tagsResponse, topicsResponse, usersResponse] =
          await Promise.all([
            postApi.list({
              page: 1,
              page_size: 10,
              sort_by: sortBy,
            }),
            tagApi.list(),
            topicApi.list({ page: 1, page_size: 8 }),
            userApi.getLeaderboardSimple({ limit: 10 }),
          ]);
        // 添加安全检查
        const posts =
          postsResponse.data.code === 0 && postsResponse.data.data
            ? postsResponse.data.data.list || []
            : [];
        const sortedTags =
          tagsResponse.data.code === 0
            ? [...(tagsResponse.data.data || [])].sort(
                (a, b) => (b.post_count || 0) - (a.post_count || 0),
              )
            : [];
        return {
          posts,
          hotTags: sortedTags.slice(0, 12),
          hotTopics:
            topicsResponse.data.code === 0
              ? topicsResponse.data.data?.list || []
              : [],
          activeUsers:
            usersResponse.data.code === 0 ? usersResponse.data.data || [] : [],
        };
      } catch (error) {
        console.error("Failed to load explore data:", error);
        toast.error("加载失败");
        throw error;
      }
    },
  });
}

// 探索页搜索 Hook：按关键词搜索帖子
export function useExploreSearch(submittedKeyword: string) {
  return useQuery<Post[]>({
    queryKey: ["explore-search", submittedKeyword],
    queryFn: async () => {
      try {
        const response = await postApi.list({
          keyword: submittedKeyword,
          page: 1,
          page_size: 20,
        });
        if (response.data.code !== 0) {
          throw new Error(response.data.message || "搜索失败");
        }
        const list = response.data.data?.list || [];
        if (list.length === 0) {
          toast("未找到相关结果");
        }
        return list;
      } catch (error) {
        console.error("Search failed:", error);
        toast.error("搜索失败");
        throw error;
      }
    },
    enabled: submittedKeyword.trim().length > 0,
  });
}
