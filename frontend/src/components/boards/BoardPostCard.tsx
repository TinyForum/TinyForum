
// ─── PostCard ──────────────────────────────────────────────────────────────────

import { Post, postApi } from "@/lib/api";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";
import { Link, HeartIcon, EyeIcon, HeartCrack } from "lucide-react";
import { useState } from "react";

export function BoardPostCard({ post }: { post: Post }) {
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(post.like_count);

  const handleLike = async (e: React.MouseEvent) => {
    e.preventDefault();
    try {
      if (liked) {
        await postApi.unlike(post.id);
        setLikeCount(prev => prev - 1);
      } else {
        await postApi.like(post.id);
        setLikeCount(prev => prev + 1);
      }
      setLiked(!liked);
    } catch (error) {
      console.error('Failed to like:', error);
    }
  };

  const excerpt = post.summary || post.content.replace(/<[^>]*>/g, '').slice(0, 180);

  return (
    <div className="group bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 hover:border-blue-200 dark:hover:border-blue-700 hover:shadow-md transition-all duration-200 p-5">
      <Link href={`/posts/${post.id}`} className="block mb-3">
        <h3 className="text-base font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 line-clamp-1 transition-colors">
          {post.title}
        </h3>
        {excerpt && (
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400 line-clamp-2 leading-relaxed">
            {excerpt}
          </p>
        )}
      </Link>

      <div className="flex items-center justify-between text-xs text-gray-400">
        <div className="flex items-center gap-3">
          <Link href={`/users/${post.author_id}`} className="flex items-center gap-1.5 hover:text-blue-500 transition-colors">
            <img
              src={post.author?.avatar || '/default-avatar.png'}
              alt={post.author?.username}
              className="w-5 h-5 rounded-full object-cover"
            />
            <span>{post.author?.username}</span>
          </Link>
          <span>{new Date(post.created_at).toLocaleDateString('zh-CN')}</span>
        </div>

        <div className="flex items-center gap-3">
          <button onClick={handleLike} className="flex items-center gap-1 hover:text-red-500 transition-colors">
            {liked
              ? <HeartCrack className="w-4 h-4 text-red-500" />
              : <HeartIcon className="w-4 h-4" />}
            <span>{likeCount}</span>
          </button>
          <div className="flex items-center gap-1">
            <ChatBubbleLeftRightIcon className="w-4 h-4" />
            <span>{post.like_count}</span>
          </div>
          <div className="flex items-center gap-1">
            <EyeIcon className="w-4 h-4" />
            <span>{post.view_count}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
