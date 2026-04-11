"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { postApi, tagApi, userApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { useTranslations } from "next-intl";
import PostFilterBar from "@/components/home/PostFilterBar";
import PostList from "@/components/home/PostList";
import Sidebar from "@/components/home/Sidebar";

export default function HomePage() {
  const { isAuthenticated } = useAuthStore();
  const [sortBy, setSortBy] = useState<"" | "hot">("");
  const [selectedTag, setSelectedTag] = useState<number | null>(null);
  const [page, setPage] = useState(1);
  const t = useTranslations("post");

  const { data: postsData, isLoading } = useQuery({
    queryKey: ["posts", sortBy, selectedTag, page],
    queryFn: () =>
      postApi
        .list({
          page,
          page_size: 15,
          sort_by: sortBy,
          tag_id: selectedTag ?? undefined,
        })
        .then((r) => r.data.data),
  });

  const { data: tags } = useQuery({
    queryKey: ["tags"],
    queryFn: () => tagApi.list().then((r) => r.data.data),
  });

  const { data: leaderboard } = useQuery({
    queryKey: ["leaderboard"],
    queryFn: () => userApi.leaderboard(10).then((r) => r.data.data),
  });

  const posts = postsData?.list ?? [];
console.log('API Response:', postsData); 
  const total = postsData?.total ?? 0;
  const totalPages = Math.ceil(total / 15);

  const handleSortChange = (newSortBy: "" | "hot") => {
    setSortBy(newSortBy);
    setPage(1);
  };

  const handleTagChange = (tagId: number | null) => {
    setSelectedTag(tagId);
    setPage(1);
  };

  return (
    <div className="flex flex-col lg:flex-row gap-6">
      {/* Main content */}
      <div className="flex-1 min-w-0">
        <PostFilterBar
          sortBy={sortBy}
          onSortChange={handleSortChange}
          isAuthenticated={isAuthenticated}
        />

        <PostList
          posts={posts}
          isLoading={isLoading}
          totalPages={totalPages}
          currentPage={page}
          onPageChange={setPage}
        />
      </div>

      {/* Sidebar */}
      <Sidebar
        tags={tags ?? []}
        selectedTag={selectedTag}
        onTagChange={handleTagChange}
        leaderboard={leaderboard ?? []}
      />
    </div>
  );
}