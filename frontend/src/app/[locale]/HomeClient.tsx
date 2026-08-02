// src/app/[locale]/page.tsx 客户端组件
"use client";

import { useState, useCallback, useEffect, useRef } from "react";
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
import { useQuestionList } from "@/features/qustion/hooks/useQuestions";
import type { Post } from "@/shared/api/types/post.model";
import type { QuestionSimple } from "@/shared/api/types/question.model";

export default function HomeClient() {
  const { isAuthenticated, user } = useAuthStore();

  const [sortBy, setSortBy] = useState<SortBy>("random");
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [selectedBoard, setSelectedBoard] = useState<number | null>(null);
  const [filterType, setFilterType] = useState<FilterType>("all");
  const [page, setPage] = useState(1);
  const [accumulatedPosts, setAccumulatedPosts] = useState<Post[]>([]);
  const [accumulatedQuestions, setAccumulatedQuestions] = useState<QuestionSimple[]>([]);
  const [total, setTotal] = useState(0);
  const pageRef = useRef(page);

  useEffect(() => {
    pageRef.current = page;
  }, [page]);

  const { data: boards = [] } = useBoardTree();
  const { tags = [] } = useTags();
  const { data: leaderboard } = useLeaderboard({ limit: 10 });
  const { unreadCount } = useUnreadCount(isAuthenticated);
  const { data: timelineEvents = [] } = useTimelineEvents(isAuthenticated);

  const pageSize = 15;

  let postParams = undefined;
  let questionParams = undefined;
  let usePostsEnabled = false;
  let useQuestionsEnabled = false;

  switch (filterType) {
    case "question":
      useQuestionsEnabled = true;
      questionParams = {
        page,
        page_size: pageSize,
        board_id: selectedBoard ?? undefined,
      };
      break;
    default:
      usePostsEnabled = true;
      postParams = {
        page,
        page_size: pageSize,
        sort_by: sortBy,
        type: filterType !== "all" ? filterType : undefined,
        board_id: selectedBoard ?? undefined,
        tag_id: selectedTag ?? undefined,
      };
      break;
  }

  const {
    data: postsData,
    isLoading: postsLoading,
    isFetching: postsFetching,
    refetch: refetchPosts,
  } = usePosts(postParams, { enabled: usePostsEnabled });

  const {
    data: questionsData,
    isLoading: questionsLoading,
    isFetching: questionsFetching,
    refetch: refetchQuestions,
  } = useQuestionList(questionParams, { enabled: useQuestionsEnabled });

  useEffect(() => {
    if (filterType === "question" && questionsData) {
      if (page === 1) {
        setAccumulatedQuestions(questionsData.list ?? []);
      } else {
        setAccumulatedQuestions((prev) => [...prev, ...(questionsData.list ?? [])]);
      }
      setTotal(questionsData.total ?? 0);
    } else if (postsData) {
      if (page === 1) {
        setAccumulatedPosts(postsData.list ?? []);
      } else {
        setAccumulatedPosts((prev) => [...prev, ...(postsData.list ?? [])]);
      }
      setTotal(postsData.total ?? 0);
    }
  }, [postsData, questionsData, filterType, page]);

  const handleLoadMore = useCallback(() => {
    if (filterType === "question") {
      if (accumulatedQuestions.length < total && !questionsFetching) {
        setPage((p) => p + 1);
      }
    } else {
      if (accumulatedPosts.length < total && !postsFetching) {
        setPage((p) => p + 1);
      }
    }
  }, [accumulatedPosts.length, accumulatedQuestions.length, total, questionsFetching, postsFetching, filterType]);

  const resetAndRefetch = useCallback(() => {
    setPage(1);
    setAccumulatedPosts([]);
    setAccumulatedQuestions([]);
    setTotal(0);
  }, []);

  const handleSortChange = (newSortBy: SortBy) => {
    resetAndRefetch();
    setSortBy(newSortBy);
    setTimeout(() => {
      if (filterType === "question") refetchQuestions();
      else refetchPosts();
    }, 0);
  };

  const handleTagChange = (tagId: number | null) => {
    resetAndRefetch();
    setSelectedTag(tagId);
    setSelectedBoard(null);
    setTimeout(() => {
      if (filterType === "question") refetchQuestions();
      else refetchPosts();
    }, 0);
  };

  const handleBoardChange = (boardId: number | null) => {
    resetAndRefetch();
    setSelectedBoard(boardId);
    setSelectedTag(null);
    setTimeout(() => {
      if (filterType === "question") refetchQuestions();
      else refetchPosts();
    }, 0);
  };

  const handlePostTypeChange = (type: FilterType) => {
    resetAndRefetch();
    setFilterType(type);
    setTimeout(() => {
      if (type === "question") refetchQuestions();
      else refetchPosts();
    }, 0);
  };

  const isLoading = filterType === "question" ? questionsLoading : postsLoading;
  const isFetching = filterType === "question" ? questionsFetching : postsFetching;
  const hasMore = filterType === "question"
    ? accumulatedQuestions.length < total
    : accumulatedPosts.length < total;

  return (
    <div className="h-full">
      <div className="container mx-auto max-w-7xl px-4 h-full">
        <div className="flex gap-6 h-full">
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

          <div className="flex-1 min-w-0 flex flex-col h-full">
            <div className="flex-shrink-0 sticky top-0 bg-base-200 pb-4 z-[10]">
              <PostFilterBar
                sortBy={sortBy}
                onSortChange={handleSortChange}
                isAuthenticated={isAuthenticated}
                onRefetch={() => {
                  resetAndRefetch();
                  setTimeout(() => {
                    if (filterType === "question") refetchQuestions();
                    else refetchPosts();
                  }, 0);
                }}
              />
            </div>

            <div className="flex-1 overflow-y-auto custom-scrollbar pb-6">
              {(() => {
                switch (filterType) {
                  case "question":
                    return (
                      <QuestionList
                        questions={accumulatedQuestions}
                        isLoading={isLoading}
                        totalPages={0}
                        currentPage={page}
                        onPageChange={setPage}
                      />
                    );
                  default:
                    return (
                      <PostList
                        posts={accumulatedPosts}
                        isLoading={isLoading || isFetching}
                        hasMore={hasMore}
                        onLoadMore={handleLoadMore}
                        layout="waterfall"
                      />
                    );
                }
              })()}
            </div>
          </div>

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
