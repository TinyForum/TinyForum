import {
  useAdminBoardPosts,
  useAdminDeletePost,
  useAdminPinPost,
} from "@/hooks/admin/useAdminModerator";
import { Search, Pin, Trash2 } from "lucide-react";
import { useState, useEffect } from "react";

// 类型定义
interface Author {
  id: number;
  username: string;
  avatar?: string;
}

interface Post {
  id: number;
  title: string;
  content: string;
  author_id: number;
  author?: Author;
  reply_count: number;
  view_count: number;
  is_pinned: boolean;
  status: string;
  created_at: string;
  updated_at: string;
}

interface PostsResponse {
  list: Post[];
  total: number;
  page: number;
  page_size: number;
}

export function PostManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [keyword, setKeyword] = useState<string>("");
  const [debouncedKeyword, setDebouncedKeyword] = useState<string>("");

  const {
    data: posts,
    isLoading,
    refetch,
  } = useAdminBoardPosts(selectedBoardId, {
    page: 1,
    page_size: 20,
    keyword: debouncedKeyword,
  });

  const deletePost = useAdminDeletePost();
  const pinPost = useAdminPinPost();

  // 防抖搜索
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedKeyword(keyword);
    }, 500);

    return () => clearTimeout(timer);
  }, [keyword]);

  const handleDeletePost = (boardId: number, postId: number) => {
    if (confirm("确定要删除该帖子吗？此操作不可恢复。")) {
      deletePost.mutate({ boardId, postId });
    }
  };

  const handlePinPost = (
    boardId: number,
    postId: number,
    pinInBoard: boolean,
  ) => {
    pinPost.mutate({ boardId, postId, pinInBoard });
  };

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  // 安全获取帖子列表
  const postList = (posts as unknown as PostsResponse)?.list || [];

  return (
    <div className="space-y-4">
      <div className="flex gap-4 flex-wrap">
        <div className="form-control">
          <label className="label">
            <span className="label-text">板块 ID</span>
          </label>
          <input
            type="number"
            placeholder="板块 ID"
            value={selectedBoardId}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setSelectedBoardId(parseInt(e.target.value) || 1)
            }
            className="input input-bordered w-32"
          />
        </div>
        <div className="flex-1">
          <div className="form-control">
            <label className="label">
              <span className="label-text">搜索</span>
            </label>
            <div className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  placeholder="搜索帖子..."
                  value={keyword}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    setKeyword(e.target.value)
                  }
                  onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) =>
                    e.key === "Enter" && setDebouncedKeyword(keyword)
                  }
                  className="input input-bordered w-full pl-9"
                />
              </div>
              <button
                className="btn btn-primary"
                onClick={() => setDebouncedKeyword(keyword)}
              >
                搜索
              </button>
              <button
                className="btn btn-outline"
                onClick={() => {
                  setKeyword("");
                  setDebouncedKeyword("");
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
        {postList.map((post: Post) => (
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
                      {post.is_pinned && (
                        <Pin className="w-4 h-4 text-primary" />
                      )}
                    </h3>
                    <div className="flex items-center gap-3 mt-1 text-xs text-base-content/50">
                      <span>
                        作者: {post.author?.username || `用户${post.author_id}`}
                      </span>
                      <span>回复: {post.reply_count || 0}</span>
                      <span>浏览: {post.view_count || 0}</span>
                      <span>
                        {new Date(post.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-1">
                    <button
                      className="btn btn-ghost btn-sm"
                      onClick={() =>
                        handlePinPost(selectedBoardId, post.id, !post.is_pinned)
                      }
                      disabled={pinPost.isPending}
                      title={post.is_pinned ? "取消置顶" : "置顶"}
                    >
                      <Pin className="w-4 h-4" />
                    </button>
                    <button
                      className="btn btn-ghost btn-sm text-error hover:text-error"
                      onClick={() => handleDeletePost(selectedBoardId, post.id)}
                      disabled={deletePost.isPending}
                      title="删除"
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
                  {post.is_pinned && (
                    <span className="badge badge-primary badge-sm">置顶</span>
                  )}
                  {post.status === "deleted" && (
                    <span className="badge badge-error badge-sm">已删除</span>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))}

        {postList.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            <Search className="w-12 h-12 mx-auto mb-3 opacity-30" />
            <p>暂无帖子</p>
            {keyword && (
              <p className="text-sm mt-2">
                没有找到 &quot;{keyword}&quot; 相关的帖子
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
