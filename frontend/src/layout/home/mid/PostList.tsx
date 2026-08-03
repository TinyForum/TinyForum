"use client";

import { useRef, useEffect, useCallback, useMemo } from "react";
import { useAuthStore } from "@/store/auth";
import type { Post } from "@/shared/api/types/post.model";
import { PostListSkeleton } from "@/shared/ui/common/PostListSkeleton";
import { EmptyPostList } from "@/shared/ui/common/EmptyPostList";
import WorkCard from "@/shared/ui/creation/WorkCard";

interface PostListProps {
  posts: Post[];
  isLoading: boolean;
  hasMore: boolean;
  onLoadMore: () => void;
  layout?: "waterfall" | "list";
}

const COLUMN_COUNT = 3;

function estimatePostHeight(post: Post): number {
  const type = post.creation?.creation_type;
  const hasImages = (post.creation?.image_urls?.length ?? 0) > 0;
  switch (type) {
    case "short_video":
      return 380;
    case "long_video":
      return 260;
    case "article":
      return post.creation?.cover_url ? 300 : 180;
    case "image_text":
    case "image":
      return hasImages ? 340 : 180;
    case "question":
      return 150;
    case "topic":
    case "post":
      return hasImages ? 280 : 160;
    default:
      return 200;
  }
}

export default function PostList({
  posts,
  isLoading,
  hasMore,
  onLoadMore,
  layout = "waterfall",
}: PostListProps) {
  const { isAuthenticated } = useAuthStore();
  const sentinelRef = useRef<HTMLDivElement>(null);

  const columns = useMemo(() => {
    const cols: Post[][] = Array.from({ length: COLUMN_COUNT }, () => []);
    const heights = new Array(COLUMN_COUNT).fill(0);
    for (const post of posts) {
      const minIdx = heights.indexOf(Math.min(...heights));
      cols[minIdx].push(post);
      heights[minIdx] += estimatePostHeight(post);
    }
    return cols;
  }, [posts]);

  const handleIntersect = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      if (entries[0].isIntersecting && hasMore && !isLoading) {
        onLoadMore();
      }
    },
    [hasMore, isLoading, onLoadMore],
  );

  useEffect(() => {
    const sentinel = sentinelRef.current;
    if (!sentinel) return;

    const observer = new IntersectionObserver(handleIntersect, {
      rootMargin: "400px",
    });
    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [handleIntersect]);

  if (isLoading && posts.length === 0) {
    return <PostListSkeleton />;
  }

  if (posts.length === 0 && !isLoading) {
    return <EmptyPostList isAuthenticated={isAuthenticated} />;
  }

  return (
    <>
      {layout === "waterfall" ? (
        <div className="flex flex-row gap-4 w-full">
          {columns.map((col, colIdx) => (
            <div key={colIdx} className="flex-1 flex flex-col gap-4 min-w-0">
              {col.map((post: Post) => (
                <WorkCard key={post.id} post={post} />
              ))}
            </div>
          ))}
        </div>
      ) : (
        <div className="space-y-3 z-0">
          {posts.map((post: Post) => (
            <WorkCard key={post.id} post={post} />
          ))}
        </div>
      )}

      <div ref={sentinelRef} className="h-4" />

      {isLoading && posts.length > 0 && (
        <div className="flex justify-center py-6">
          <span className="loading loading-spinner loading-md text-primary" />
        </div>
      )}

      {!hasMore && posts.length > 0 && (
        <p className="text-center text-sm text-base-content/40 py-6">
          已经到底了
        </p>
      )}
    </>
  );
}
