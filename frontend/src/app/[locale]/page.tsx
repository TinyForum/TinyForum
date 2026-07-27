"use client";

import { useState } from "react";
import { useAuthStore } from "@/store/auth";
import LeftSidebar, { FilterType } from "@/layout/home/LeftSidebar";
import { useLeaderboard } from "@/features/leader/hooks/useLeaderboard";
import { useUnreadCount } from "@/features/notification/hooks/useUnreadCount";
import { RightSidebar } from "@/layout/home/RightSidebar";
import PostFilterBar from "@/layout/home/mid/PostFilterBar";
import PostList from "@/layout/home/mid/PostList";
import QuestionList from "@/layout/home/mid/QuestionList";
import { SortBy } from "@/shared/ui/type/home.type";
import { useTimelineEvents } from "@/features/timeline/hooks/useTimelineEvents";
import { useBoardTree } from "@/features/boards/hooks/useBoardTree";
import { usePosts } from "@/features/post/hooks/usePosts";
import { useTags } from "@/features/tag/hooks/useTags";
import { useQuestions } from "@/features/qustion/hooks/useQuestions";

export default function HomePage() {
  const { isAuthenticated, user } = useAuthStore();

  const [sortBy, setSortBy] = useState<SortBy>("random");
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [selectedBoard, setSelectedBoard] = useState<number | null>(null);
  const [filterType, setFilterType] = useState<FilterType>("all");
  const [page, setPage] = useState(1);

  const { data: boards = [] } = useBoardTree();
  const { tags = [] } = useTags();
  const { data: leaderboard } = useLeaderboard({ limit: 10 });
  const { unreadCount } = useUnreadCount();
  const { data: timelineEvents = [] } = useTimelineEvents(isAuthenticated);

  // ----- 根据 filterType 配置参数和启用状态 -----
  let postParams = undefined;
  let questionParams = undefined;
  let usePostsEnabled = false;
  let useQuestionsEnabled = false;

  switch (filterType) {
    case "questions":
      useQuestionsEnabled = true;
      questionParams = {
        page,
        page_size: 15,
        board_id: selectedBoard ?? undefined,
      };
      break;
    default:
      // 处理 all, posts, articles, topic 等
      usePostsEnabled = true;
      postParams = {
        page,
        page_size: 15,
        sort_by: sortBy === "latest" ? "latest" : undefined,
        type: filterType !== "all" ? filterType : undefined,
        board_id: selectedBoard ?? undefined,
        tag_id: selectedTag ?? undefined,
      };
      break;
  }

  // ----- 调用 Hooks（始终在顶层） -----
  const {
    data: postsData,
    isLoading: postsLoading,
    refetch: refetchPosts,
  } = usePosts(postParams, { enabled: usePostsEnabled });

  const {
    data: questionsData,
    isLoading: questionsLoading,
    refetch: refetchQuestions,
  } = useQuestions(questionParams, { enabled: useQuestionsEnabled });

  // ----- 根据 filterType 选择数据、加载状态、刷新函数、总数 -----
  let rawData: any[] = [];
  let isLoading = false;
  let refetch: () => void = () => {};
  let total = 0;

  switch (filterType) {
    case "questions":
      rawData = questionsData?.list ?? [];
      isLoading = questionsLoading;
      refetch = refetchQuestions;
      total = questionsData?.total ?? 0;
      break;
    default:
      rawData = postsData?.list ?? [];
      isLoading = postsLoading;
      refetch = refetchPosts;
      total = postsData?.total ?? 0;
      break;
  }

  const totalPages = Math.ceil(total / 15);

  // ----- 事件处理 -----
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
              boards={boards}
              tags={tags}
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
              {(() => {
                switch (filterType) {
                  case "questions":
                    return (
                      <QuestionList
                        questions={rawData}
                        isLoading={isLoading}
                        totalPages={totalPages}
                        currentPage={page}
                        onPageChange={setPage}
                      />
                    );
                  default:
                    return (
                      <PostList
                        posts={rawData}
                        isLoading={isLoading}
                        totalPages={totalPages}
                        currentPage={page}
                        onPageChange={setPage}
                      />
                    );
                }
              })()}
            </div>
          </div>

          {/* 右侧边栏 */}
          <div className="lg:w-64 xl:w-72 flex-none overflow-y-auto custom-scrollbar sticky top-6 max-h-[calc(100vh-6rem)]">
            <RightSidebar
              isAuthenticated={isAuthenticated}
              userProfile={user}
              leaderboard={leaderboard ?? []}
              unreadCount={unreadCount ?? 0}
              timelineEvents={timelineEvents}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
