"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { useAuthStore } from "@/store/auth";
import type { Post } from "@/shared/api/types";
import { PostListSkeleton } from "@/shared/ui/common/PostListSkeleton";
import { EmptyPostList } from "@/shared/ui/common/EmptyPostList";
import { Pagination } from "@/shared/ui/common/Pagination";
import PostCard from "@/shared/ui/post/PostCard";

interface PostListProps {
  posts: Post[];
  isLoading: boolean;
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}

export default function PostList({
  posts,
  isLoading,
  totalPages,
  currentPage,
  onPageChange,
}: PostListProps) {
  const { isAuthenticated } = useAuthStore();

  if (isLoading) {
    return <PostListSkeleton />;
  }

  if (posts.length === 0) {
    return <EmptyPostList isAuthenticated={isAuthenticated} />;
  }

  return (
    <>
      <div className="space-y-3 z-0">
        {posts.map((post: Post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </div>

      {totalPages > 1 && (
        <Pagination
          totalPages={totalPages}
          currentPage={currentPage}
          onPageChange={onPageChange}
        />
      )}
    </>
  );
}
