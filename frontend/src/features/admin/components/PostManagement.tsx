import { Search, Pin, Trash2 } from "lucide-react";
import { useState, useEffect } from "react";

import type { Post } from "@/shared/api/types";
import { useAdminGetPosts, useAdminTogglePin } from "../hooks/useAdminPosts";

export function PostManagement() {
  const [keyword, setKeyword] = useState<string>("");
  const [debouncedKeyword, setDebouncedKeyword] = useState<string>("");
  const [page, setPage] = useState(1);
  const pageSize = 20;

  const {
    data: postsData,
    isLoading,
    refetch,
  } = useAdminGetPosts({
    page,
    page_size: pageSize,
    keyword: debouncedKeyword || undefined,
  });

  const togglePin = useAdminTogglePin();

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedKeyword(keyword);
      setPage(1);
    }, 500);
    return () => clearTimeout(timer);
  }, [keyword]);

  const handleTogglePin = (postId: number) => {
    togglePin.mutate(postId);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  const postList = postsData?.list || [];
  const total = postsData?.total || 0;
  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="space-y-4">
      {/* 搜索框部分保持不变 */}
      <div className="flex gap-4 flex-wrap">
        <div className="flex-1">
          <div className="form-control">
            <label className="label">
              <span className="label-text">搜索</span>
            </label>
            <div className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  placeholder="搜索帖子（标题/内容）..."
                  value={keyword}
                  onChange={(e) => setKeyword(e.target.value)}
                  onKeyDown={(e) =>
                    e.key === "Enter" && setDebouncedKeyword(keyword)
                  }
                  className="input input-bordered w-full pl-9"
                />
              </div>
              <button
                className="btn btn-outline"
                onClick={() => {
                  setKeyword("");
                  setDebouncedKeyword("");
                  setPage(1);
                  refetch();
                }}
              >
                重置
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="space-y-3">
        {postList.map(
          (
            post: Post, // 现在 Post 类型已正确
          ) => (
            <div
              key={post.id}
              className="card bg-base-100 shadow-sm border border-base-200 hover:shadow-md transition-shadow"
            >
              <div className="card-body p-4">
                <div className="space-y-2">
                  <div className="flex justify-between items-start flex-wrap gap-2">
                    <div className="flex-1">
                      <h3 className="font-medium flex items-center gap-2 text-base">
                        {post.title}
                        {post.pin_top && (
                          <Pin className="w-4 h-4 text-primary" />
                        )}
                      </h3>
                      <div className="flex items-center gap-3 mt-1 text-xs text-base-content/50">
                        <span>
                          作者:{" "}
                          {post.author?.username || `用户${post.author_id}`}
                        </span>
                        {/* <span>回复: {post. || 0}</span> */}
                        <span>浏览: {post.view_count || 0}</span>
                        <span>
                          {new Date(post.created_at).toLocaleDateString()}
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-1">
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => handleTogglePin(post.id)}
                        disabled={togglePin.isPending}
                        title={post.pin_top ? "取消置顶" : "置顶"}
                      >
                        <Pin className="w-4 h-4" />
                      </button>
                      <button
                        className="btn btn-ghost btn-sm text-error hover:text-error"
                        disabled
                        title="删除功能暂未开放"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                  {post.content && (
                    <p className="text-sm text-base-content/70 line-clamp-2">
                      {post.content}
                    </p>
                  )}
                  <div className="flex gap-2 flex-wrap">
                    {post.pin_top && (
                      <span className="badge badge-primary badge-sm">置顶</span>
                    )}
                    {post.status === "published" && (
                      <span className="badge badge-error badge-sm">已发布</span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ),
        )}
        {postList.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            <Search className="w-12 h-12 mx-auto mb-3 opacity-30" />
            <p>暂无帖子</p>
            {debouncedKeyword && (
              <p className="text-sm mt-2">
                没有找到 &quot;{debouncedKeyword}&quot; 相关的帖子
              </p>
            )}
          </div>
        )}
      </div>

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="flex justify-center mt-4">
          <div className="join">
            <button
              className="join-item btn btn-sm"
              disabled={page <= 1}
              onClick={() => setPage((p) => p - 1)}
            >
              «
            </button>
            <span className="join-item btn btn-sm">
              第 {page} / {totalPages} 页
            </span>
            <button
              className="join-item btn btn-sm"
              disabled={page >= totalPages}
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
