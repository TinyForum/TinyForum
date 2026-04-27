import {
  useAdminBoardPosts,
  useAdminDeletePost,
  useAdminPinPost,
} from "@/hooks/admin/useAdminModerator";
import { Search, Pin, Trash2 } from "lucide-react";
import { useState } from "react";

export function PostManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [keyword, setKeyword] = useState("");
  const [debouncedKeyword, setDebouncedKeyword] = useState("");

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

  const handleSearch = () => {
    setDebouncedKeyword(keyword);
  };

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

  if (isLoading)
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <input
          type="number"
          placeholder="板块 ID"
          value={selectedBoardId}
          onChange={(e) => setSelectedBoardId(parseInt(e.target.value))}
          className="input input-bordered w-32"
        />
        <div className="flex-1 flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              placeholder="搜索帖子..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              className="input input-bordered w-full pl-9"
            />
          </div>
          <button className="btn btn-primary" onClick={handleSearch}>
            搜索
          </button>
          <button className="btn btn-outline" onClick={() => refetch()}>
            刷新
          </button>
        </div>
      </div>

      <div className="space-y-3">
        {posts?.list?.map((post: any) => (
          <div
            key={post.id}
            className="card bg-base-100 shadow-sm border border-base-200"
          >
            <div className="card-body">
              <div className="space-y-2">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <h3 className="font-medium flex items-center gap-2">
                      {post.title}
                      {post.is_pinned && (
                        <Pin className="w-4 h-4 text-primary" />
                      )}
                    </h3>
                    <p className="text-sm text-gray-500 mt-1">
                      作者: {post.author?.username} | 回复: {post.reply_count} |
                      浏览: {post.view_count}
                    </p>
                  </div>
                  <div className="flex gap-1">
                    <button
                      className="btn btn-ghost btn-sm"
                      onClick={() =>
                        handlePinPost(selectedBoardId, post.id, !post.is_pinned)
                      }
                      disabled={pinPost.isPending}
                    >
                      <Pin className="w-4 h-4" />
                    </button>
                    <button
                      className="btn btn-ghost btn-sm text-error"
                      onClick={() => handleDeletePost(selectedBoardId, post.id)}
                      disabled={deletePost.isPending}
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <p className="text-sm text-gray-600 line-clamp-2">
                  {post.content}
                </p>
                <div className="flex gap-2">
                  {post.is_pinned && (
                    <span className="badge badge-primary">置顶</span>
                  )}
                  {post.status === "deleted" && (
                    <span className="badge badge-error">已删除</span>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))}
        {posts?.list?.length === 0 && (
          <div className="text-center py-8 text-gray-500">暂无帖子</div>
        )}
      </div>
    </div>
  );
}
