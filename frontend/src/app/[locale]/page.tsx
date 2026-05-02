"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  postApi,
  tagApi,
  boardApi,
  timelineApi,
  notificationApi,
  questionApi,
  Tag,
  Post,
  TimelineEvent,
} from "@/shared/api";
import { useAuthStore } from "@/store/auth";
import LeftSidebar, { FilterType } from "@/layout/home/LeftSidebar";
import { SortBy } from "@/shared/type/posts.types";
import { useLeaderboard } from "@/features/leader/hooks/useLeaderboard";
import RightSidebar from "@/layout/home/RightSidebar";
import PostFilterBar from "@/layout/home/mid/PostFilterBar";
import PostList from "@/layout/home/mid/PostList";
import { Board } from "@/shared/api/types/board.model";

export default function HomePage() {
  const { isAuthenticated, user } = useAuthStore();
  const [sortBy, setSortBy] = useState<SortBy>("random");
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [selectedBoard, setSelectedBoard] = useState<number | null>(null);
  const [filterType, setFilterType] = useState<FilterType>("all");
  const [page, setPage] = useState(1);
  const { data: leaderboard } = useLeaderboard({ limit: 10 });

  // 帖子列表
  const {
    data: postsData,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ["posts", sortBy, selectedTag, selectedBoard, filterType, page],
    queryFn: async () => {
      const params: Record<string, unknown> = {
        page,
        page_size: 15,
      };

      if (filterType === "question") {
        if (selectedBoard) {
          params.board_id = selectedBoard;
        }
        const res = await questionApi.getSimple(params);
        const data = res.data.data as
          | { list: unknown[]; total: number; page: number; page_size: number }
          | undefined;
        return {
          list: data?.list || [],
          total: data?.total || 0,
          page: data?.page || 1,
          page_size: data?.page_size || 15,
        };
      }

      if (selectedTag) {
        params.tag_id = selectedTag;
      }
      if (selectedBoard) {
        params.board_id = selectedBoard;
      }
      if (filterType !== "all") {
        params.type = filterType;
      }
      if (sortBy !== "random") {
        params.sort_by = sortBy === "latest" ? "latest" : sortBy;
      }

      const res = await postApi.list(params);
      const data = res.data.data as
        | { list: unknown[]; total: number; page: number; page_size: number }
        | undefined;
      return {
        list: data?.list || [],
        total: data?.total || 0,
        page: data?.page || 1,
        page_size: data?.page_size || 15,
      };
    },
  });

  // 标签列表
  const { data: tags } = useQuery({
    queryKey: ["tags"],
    queryFn: () => tagApi.list().then((r) => (r.data.data as Tag[]) || []),
  });

  // 板块列表
  const { data: boards } = useQuery({
    queryKey: ["boards-tree"],
    queryFn: () =>
      boardApi.getTree().then((r) => (r.data.data as Board[]) || []),
  });

  // 用户信息（已登录时）
  // const { data: userProfile } = useQuery({
  //   queryKey: ["user-profile", user?.id],
  //   queryFn: () =>
  //     userAPI.getProfile(user!.id).then((r) => r.data.data as UserProfile),
  //   enabled: isAuthenticated && !!user?.id,
  // });
  // const { profile: userProfile, refresh: fetchUserProfile } = useUserProfile(
  //   user?.id,
  //   false,
  // );
  // // 导入自定义 hooks
  // const { user  } = useProfile();

  // 未读通知数
  const { data: unreadCount } = useQuery({
    queryKey: ["unread-count"],
    queryFn: () =>
      notificationApi.unreadCount().then((r) => {
        const data = r.data.data as { count: number } | undefined;
        return data?.count || 0;
      }),
    enabled: isAuthenticated,
    refetchInterval: 30000,
  });

  // 时间线事件
  const { data: timelineEvents } = useQuery({
    queryKey: ["timeline-events"],
    queryFn: () =>
      timelineApi
        .getFollowing({ page: 1, page_size: 5 })
        .then((r) => (r.data.data?.list as TimelineEvent[]) || []),
    enabled: isAuthenticated,
  });

  const posts = (postsData?.list as unknown[] as Post[]) ?? [];
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
            <div className="flex-shrink-0 sticky top-0 bg-base-200 pb-4 z-[10]">
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
                userProfile={user}
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
