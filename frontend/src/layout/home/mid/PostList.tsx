"use client";

import { useRef, useEffect, useCallback } from "react";
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

export default function PostList({
  posts,
  isLoading,
  hasMore,
  onLoadMore,
  layout = "waterfall",
}: PostListProps) {
  const { isAuthenticated } = useAuthStore();
  const sentinelRef = useRef<HTMLDivElement>(null);

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
      rootMargin: "200px",
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
        <div
          className="columns-1 sm:columns-2 xl:columns-3 gap-4 space-y-4"
        >
          {posts.map((post: Post) => (
            <div key={post.id} className="break-inside-avoid">
              <WorkCard post={post} />
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
