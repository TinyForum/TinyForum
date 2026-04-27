"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  postApi,
  tagApi,
  userApi,
  boardApi,
  timelineApi,
  notificationApi,
  questionApi,
  PostType,
} from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { useTranslations } from "next-intl";
import PostFilterBar from "@/components/home/PostFilterBar";
import PostList from "@/components/home/PostList";
import LeftSidebar, { FilterType } from "@/components/home/LeftSidebar";
import RightSidebar from "@/components/home/RightSidebar";
import { SortBy } from "@/type/posts.types";
// import { PostType, SortBy } from "@/types";
import { useLeaderboard } from "@/hooks/useLeaderboard";
// export type FilterType = "all" | PostType;
export default function HomePage() {
  const { isAuthenticated, user } = useAuthStore();
  const [sortBy, setSortBy] = useState<SortBy>("random");
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [selectedBoard, setSelectedBoard] = useState<number | null>(null);
  const [filterType, setFilterType] = useState<FilterType>("all");
  const [page, setPage] = useState(1);
  const t = useTranslations("post");
  const { data: leaderboard } = useLeaderboard({ limit: 10 });
  // 帖子列表 - 根据 postType 选择不同的 API
  const {
    data: postsData,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ["posts", sortBy, selectedTag, selectedBoard, filterType, page],
    queryFn: async () => {
      // 如果是问答类型，使用专门的 questionApi.getSimple API
      if (filterType === "question") {
        const params: any = {
          page,
          page_size: 15,
        };

        // 问答支持的过滤参数
        if (selectedBoard) {
          params.board_id = selectedBoard;
        }

        // 注意：问答 API 目前不支持标签过滤和排序
        // 如果后端支持，可以添加：
        // if (selectedTag) params.tag_id = selectedTag;
        // if (sortBy === "latest") params.sort = "latest";

        const res = await questionApi.getSimple(params);
        // 将 QuestionSimple 转换为 Post 格式以兼容 PostList 组件
        const questionList = res.data.data.list;
        const transformedPosts = questionList.map((q: any) => ({
          id: q.id,
          title: q.title,
          summary: q.summary,
          content: "",
          type: "question",
          status: "published",
          author_id: q.author_id,
          view_count: 0,
          like_count: 0,
          pin_top: false,
          board_id: q.board_id,
          reward_score: q.reward_score,
          created_at: q.created_at,
          updated_at: q.updated_at,
          question: {
            answer_count: q.answer_count,
            reward_score: q.reward_score,
          },
          author: q.author,
          board: q.board,
          tags: q.tags || [],
        }));

        return {
          list: transformedPosts,
          total: res.data.data.total,
          page: res.data.data.page,
          page_size: res.data.data.page_size,
        };
      }

      // 普通帖子、文章使用 list API
      const params: any = {
        page,
        page_size: 15,
        sort_by: sortBy === "latest" ? "latest" : sortBy,
      };

      if (selectedTag) {
        params.tag_id = selectedTag;
      }
      if (selectedBoard) {
        params.board_id = selectedBoard;
      }
      if (filterType !== "all") {
        params.type = filterType;
      }

      const res = await postApi.list(params);
      return res.data.data;
    },
  });

  // 标签列表
  const { data: tags } = useQuery({
    queryKey: ["tags"],
    queryFn: () => tagApi.list().then((r) => r.data.data),
  });

  // 板块列表
  const { data: boards } = useQuery({
    queryKey: ["boards-tree"],
    queryFn: () => boardApi.getTree().then((r) => r.data.data),
  });

  // 用户排行榜
  // const { data: leaderboard } = useQuery({
  //   queryKey: ["leaderboard"],
  //   queryFn: () => userApi.leaderboard(10).then((r) => r.data.data),
  // });

  // 用户信息（已登录时）
  const { data: userProfile } = useQuery({
    queryKey: ["user-profile", user?.id],
    queryFn: () => userApi.getProfile(user!.id).then((r) => r.data.data),
    enabled: isAuthenticated && !!user?.id,
  });

  // 未读通知数
  const { data: unreadCount } = useQuery({
    queryKey: ["unread-count"],
    queryFn: () => notificationApi.unreadCount().then((r) => r.data.data.count),
    enabled: isAuthenticated,
    refetchInterval: 30000,
  });

  // 时间线事件（已登录时）
  const { data: timelineEvents } = useQuery({
    queryKey: ["timeline-events"],
    queryFn: () =>
      timelineApi
        .getFollowing({ page: 1, page_size: 5 })
        .then((r) => r.data.data.list),
    enabled: isAuthenticated,
  });

  const posts = postsData?.list ?? [];
  const total = postsData?.total ?? 0;
  const totalPages = Math.ceil(total / 15);

  const handleSortChange = (newSortBy: SortBy) => {
    setSortBy(newSortBy);
    setPage(1);
  };

  const handleTagChange = (tagId: number | null) => {
    setSelectedTag(tagId);
    setSelectedBoard(null);
    setPage(1);
  };

  const handleBoardChange = (boardId: number | null) => {
    setSelectedBoard(boardId);
    setSelectedTag(null);
    setPage(1);
  };

  const handlePostTypeChange = (type: FilterType) => {
    setFilterType(type);
    setPage(1);
  };

  return (
    <div className="h-full">
      <div className="container mx-auto max-w-7xl px-4 h-full">
        <div className="flex gap-6 h-full">
          {/* 左侧边栏 */}
          <div className="lg:w-64 xl:w-72 flex-none overflow-y-auto custom-scrollbar sticky top-6 max-h-[calc(100vh-6rem)]">
            <LeftSidebar
              boards={boards ?? []}
              tags={tags ?? []}
              selectedBoard={selectedBoard}
              selectedTag={selectedTag}
              filterType={filterType}
              onBoardChange={handleBoardChange}
              onTagChange={handleTagChange}
              onPostTypeChange={handlePostTypeChange}
            />
          </div>

          {/* 中间内容区域 */}
          <div className="flex-1 min-w-0 flex flex-col h-full">
            <div className="flex-shrink-0 sticky top-0 bg-base-200 z-10 pb-4">
              <PostFilterBar
                sortBy={sortBy}
                onSortChange={handleSortChange}
                isAuthenticated={isAuthenticated}
                onRefetch={refetch}
              />
            </div>

            <div className="flex-1 overflow-y-auto custom-scrollbar pb-6">
              <PostList
                posts={posts}
                isLoading={isLoading}
                totalPages={totalPages}
                currentPage={page}
                onPageChange={setPage}
              />
            </div>
          </div>

          {/* 右侧边栏 */}
          {leaderboard !== undefined && (
            <div className="lg:w-64 xl:w-72 flex-none overflow-y-auto custom-scrollbar sticky top-6 max-h-[calc(100vh-6rem)]">
              <RightSidebar
                isAuthenticated={isAuthenticated}
                user={user}
                userProfile={userProfile}
                leaderboard={leaderboard}
                unreadCount={unreadCount ?? 0}
                timelineEvents={timelineEvents ?? []}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
