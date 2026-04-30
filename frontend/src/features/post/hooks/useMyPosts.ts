import { useState, useCallback, useEffect } from "react";
import { toast } from "react-hot-toast";

interface UseMyPostsOptions {
  userId?: number;
  page: number;
  keyword?: string;
  enabled?: boolean;
}

export function useMyPosts({
  userId,
  page,
  keyword = "",
  enabled = true,
}: UseMyPostsOptions) {
  const [posts, setPosts] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const fetchPosts = useCallback(async () => {
    if (!userId || !enabled) return;

    setLoading(true);
    try {
      // 调用 API
      const response = await fetch(
        `/api/users/${userId}/posts?page=${page}&keyword=${keyword}`,
      );
      const data = await response.json();
      if (data.code === 0) {
        setPosts(data.data.list);
        setTotal(data.data.total);
      }
    } catch {
      toast.error("加载帖子失败");
    } finally {
      setLoading(false);
    }
  }, [userId, page, keyword, enabled]);

  const deletePost = useCallback(
    async (postId: number) => {
      try {
        const response = await fetch(`/api/posts/${postId}`, {
          method: "DELETE",
        });
        const data = await response.json();
        if (data.code === 0) {
          toast.success("删除成功");
          await fetchPosts();
          return true;
        }
      } catch {
        toast.error("删除失败");
        return false;
      }
    },
    [fetchPosts],
  );

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  return { posts, total, loading, deletePost, refresh: fetchPosts };
}
