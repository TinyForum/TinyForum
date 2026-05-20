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
import { Post } from "@/shared/api/types/post.model";
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

  const isQuestionMode = filterType === "question";

  // 普通帖子参数
  const postParams = !isQuestionMode
    ? {
        page,
        page_size: 15,
        sort_by: sortBy === "latest" ? "latest" : undefined,
        type: filterType !== "all" ? filterType : undefined,

        board_id: selectedBoard ?? undefined,
        tag_id: selectedTag ?? undefined,
      }
    : undefined;

  const {
    data: postsData,
    isLoading: postsLoading,
    refetch: refetchPosts,
  } = usePosts(postParams, { enabled: !isQuestionMode });

  // 问答参数
  const questionParams = {
    page,
    page_size: 15,
    board_id: selectedBoard ?? undefined,
  };

  const {
    data: questionsData,
    isLoading: questionsLoading,
    refetch: refetchQuestions,
  } = useQuestions(questionParams, { enabled: isQuestionMode });

  const isLoading = isQuestionMode ? questionsLoading : postsLoading;
  const refetch = isQuestionMode ? refetchQuestions : refetchPosts;

  // 帖子模式使用的数据
  const rawPosts = postsData?.list ?? [];

  // 问答模式使用的数据（直接使用，无需断言）
  const rawQuestions = questionsData?.list ?? [];

  // 总条数
  const total = isQuestionMode
    ? (questionsData?.total ?? 0)
    : (postsData?.total ?? 0);
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
              {isQuestionMode ? (
                <QuestionList
                  questions={rawQuestions}
                  isLoading={isLoading}
                  totalPages={totalPages}
                  currentPage={page}
                  onPageChange={setPage}
                />
              ) : (
                <PostList
                  posts={rawPosts}
                  isLoading={isLoading}
                  totalPages={totalPages}
                  currentPage={page}
                  onPageChange={setPage}
                />
              )}
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
