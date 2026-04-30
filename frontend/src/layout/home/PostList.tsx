"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { useAuthStore } from "@/store/auth";
import type { Post } from "@/shared/api/types";
import PostCard from "../post/PostCard";

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
      <div className="space-y-3">
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

// 骨架屏组件
function PostListSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_: unknown, i: number) => (
        <div key={i} className="skeleton h-28 w-full rounded-xl" />
      ))}
    </div>
  );
}

// 空状态组件
function EmptyPostList({ isAuthenticated }: { isAuthenticated: boolean }) {
  const t = useTranslations("Post");

  return (
    <div className="text-center py-20 text-base-content/40">
      <p className="text-lg">{t("no_posts")}</p>
      {isAuthenticated && (
        <Link href="/posts/new" className="btn btn-primary mt-4">
          {t("post_your_first_post")}
        </Link>
      )}
    </div>
  );
}

// 分页组件
function Pagination({
  totalPages,
  currentPage,
  onPageChange,
}: {
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}) {
  const getPageNumbers = (): number[] => {
    const pages: number[] = [];
    const maxVisible = 7;
    const start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
    const end = Math.min(totalPages, start + maxVisible - 1);

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    return pages;
  };

  return (
    <div className="flex justify-center mt-6">
      <div className="join">
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === 1}
          onClick={() => onPageChange(currentPage - 1)}
          aria-label="Previous page"
        >
          «
        </button>
        {getPageNumbers().map((p: number) => (
          <button
            key={p}
            className={`join-item btn btn-sm ${currentPage === p ? "btn-active btn-primary" : ""}`}
            onClick={() => onPageChange(p)}
            aria-label={`Go to page ${p}`}
            aria-current={currentPage === p ? "page" : undefined}
          >
            {p}
          </button>
        ))}
        <button
          className="join-item btn btn-sm"
          disabled={currentPage === totalPages}
          onClick={() => onPageChange(currentPage + 1)}
          aria-label="Next page"
        >
          »
        </button>
      </div>
    </div>
  );
}
