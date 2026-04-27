"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { postApi } from "@/lib/api";
import PostCard from "@/components/post/PostCard";
import { Search } from "lucide-react";
import { useTranslations } from "next-intl";

function PostsContent() {
  const searchParams = useSearchParams();
  const keyword = searchParams.get("keyword") || "";
  const tagId = searchParams.get("tag_id")
    ? Number(searchParams.get("tag_id"))
    : undefined;
  const [page, setPage] = useState(1);
  const t = useTranslations("posts");

  const { data, isLoading } = useQuery({
    queryKey: ["posts", "search", keyword, tagId, page],
    queryFn: () =>
      postApi
        .list({ page, page_size: 20, keyword, tag_id: tagId })
        .then((r) => r.data.data),
  });

  const posts = data?.list ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.ceil(total / 20);

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center gap-3 mb-6">
        <Search className="w-5 h-5 text-primary" />
        <h1 className="text-xl font-bold">
          {keyword ? `${t("search")}：${keyword}` : t("all_posts")}
          {total > 0 && (
            <span className="text-base-content/40 font-normal ml-2">
              ({total} {t("results")})
            </span>
          )}
        </h1>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="skeleton h-28 w-full rounded-xl" />
          ))}
        </div>
      ) : posts.length === 0 ? (
        <div className="text-center py-20 text-base-content/40">
          <Search className="w-12 h-12 mx-auto mb-4 opacity-30" />
          <p className="text-lg">{t("no_results")}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center mt-6">
          <div className="join">
            <button
              className="join-item btn btn-sm"
              disabled={page === 1}
              onClick={() => setPage((p) => p - 1)}
            >
              «
            </button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
              <button
                key={p}
                className={`join-item btn btn-sm ${page === p ? "btn-active btn-primary" : ""}`}
                onClick={() => setPage(p)}
              >
                {p}
              </button>
            ))}
            <button
              className="join-item btn btn-sm"
              disabled={page === totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              »
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default function PostsPage() {
  return (
    <Suspense
      fallback={
        <div className="flex justify-center py-20">
          <span className="loading loading-spinner loading-lg text-primary" />
        </div>
      }
    >
      <PostsContent />
    </Suspense>
  );
}
