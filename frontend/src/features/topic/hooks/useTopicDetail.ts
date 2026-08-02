// hooks/useTopicDetail.ts
import { useCallback } from "react";
import { topicApi } from "@/shared/api/modules/topics";

// 话题详情页：封装话题相关操作
export function useTopicDetail() {
  // 获取话题详情
  const fetchTopic = useCallback((topicId: number) => {
    return topicApi.getById(topicId);
  }, []);

  // 获取话题帖子列表
  const fetchTopicPosts = useCallback(
    (topicId: number, params?: { page?: number; page_size?: number }) => {
      return topicApi.getPosts(topicId, params);
    },
    [],
  );

  // 获取话题收藏者列表
  const fetchTopicFollowers = useCallback(
    (topicId: number, params?: { page?: number; page_size?: number }) => {
      return topicApi.getFollowers(topicId, params);
    },
    [],
  );

  // 关注话题
  const followTopic = useCallback((topicId: number) => {
    return topicApi.follow(topicId);
  }, []);

  // 取消关注话题
  const unfollowTopic = useCallback((topicId: number) => {
    return topicApi.unfollow(topicId);
  }, []);

  // 添加帖子到话题
  const addPostToTopic = useCallback((topicId: number, postId: number) => {
    return topicApi.addPost(topicId, { post_id: postId });
  }, []);

  return {
    fetchTopic,
    fetchTopicPosts,
    fetchTopicFollowers,
    followTopic,
    unfollowTopic,
    addPostToTopic,
  };
}
