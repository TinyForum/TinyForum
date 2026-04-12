"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { postApi, tagApi, userApi, boardApi, timelineApi, notificationApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { useTranslations } from "next-intl";
import PostFilterBar from "@/components/home/PostFilterBar";
import PostList from "@/components/home/PostList";
import LeftSidebar from "@/components/home/LeftSidebar";
import RightSidebar from "@/components/home/RightSidebar";
import { PostType, SortBy } from "@/types";

export default function HomePage() {
  const { isAuthenticated, user } = useAuthStore();
  const [sortBy, setSortBy] = useState<SortBy>("random");
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [selectedBoard, setSelectedBoard] = useState<number | null>(null);
  const [postType, setPostType] = useState<PostType>("all");
  const [page, setPage] = useState(1);
  const t = useTranslations("post");

  // 帖子列表
  const { data: postsData, isLoading, refetch } = useQuery({
    queryKey: ["posts", sortBy, selectedTag, selectedBoard, postType, page],
    queryFn: () =>
      postApi
        .list({
          page,
          page_size: 15,
          sort_by: sortBy === "latest" ? "latest" : sortBy,
          tag_id: selectedTag ?? undefined,
          board_id: selectedBoard ?? undefined,
          type: postType === "all" ? undefined : postType,
        })
        .then((r) => r.data.data),
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
  const { data: leaderboard } = useQuery({
    queryKey: ["leaderboard"],
    queryFn: () => userApi.leaderboard(10).then((r) => r.data.data),
  });

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
    refetchInterval: 30000, // 30秒刷新一次
  });

  // 时间线事件（已登录时）
  const { data: timelineEvents } = useQuery({
    queryKey: ["timeline-events"],
    queryFn: () => timelineApi.getHomeTimeline({ page: 1, page_size: 5 }).then((r) => r.data.data.list),
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

  const handlePostTypeChange = (type: "all" | "question" | "article") => {
    setPostType(type);
    setPage(1);
  };

  return (
    <div className="flex flex-col lg:flex-row gap-6">
      {/* 左侧边栏 - 社区信息 */}
      <div className="lg:w-64 xl:w-72 flex-none">
        <LeftSidebar
          boards={boards ?? []}
          tags={tags ?? []}
          selectedBoard={selectedBoard}
          selectedTag={selectedTag}
          postType={postType}
          onBoardChange={handleBoardChange}
          onTagChange={handleTagChange}
          onPostTypeChange={handlePostTypeChange}
        />
      </div>

      {/* 中间内容区域 */}
      <div className="flex-1 min-w-0">
        <PostFilterBar
          sortBy={sortBy}
          onSortChange={handleSortChange}
          isAuthenticated={isAuthenticated}
          onRefetch={refetch}
        />

        <PostList
          posts={posts}
          isLoading={isLoading}
          totalPages={totalPages}
          currentPage={page}
          onPageChange={setPage}
        />
      </div>

      {/* 右侧边栏 - 个人信息 */}
      <div className="lg:w-64 xl:w-72 flex-none">
        <RightSidebar
          isAuthenticated={isAuthenticated}
          user={user}
          userProfile={userProfile}
          leaderboard={leaderboard ?? []}
          unreadCount={unreadCount ?? 0}
          timelineEvents={timelineEvents ?? []}
        />
      </div>
    </div>
  );
}