import { Post, postApi } from "@/shared/api";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";
import { Link, EyeIcon, HeartIcon, HeartCrack } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";

// 热门帖子卡片
export function HotPostCard({ post, rank }: { post: Post; rank: number }) {
  const [liked, setLiked] = useState(false);
  const [likesCount, setLikesCount] = useState(post.like_count || 0);

  const handleLike = async () => {
    if (!post.id) return;
    try {
      if (liked) {
        await postApi.unlike(post.id);
        setLiked(false);
        setLikesCount((prev) => prev - 1);
      } else {
        await postApi.like(post.id);
        setLiked(true);
        setLikesCount((prev) => prev + 1);
      }
    } catch (error) {
      console.error("Like failed:", error);
      toast.error("操作失败");
    }
  };

  const rankColors: Record<number, string> = {
    1: "text-yellow-500",
    2: "text-gray-400",
    3: "text-amber-600",
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
      <div className="flex items-start gap-3">
        {rank > 0 && (
          <div
            className={`text-2xl font-bold w-8 ${rankColors[rank] || "text-gray-300"}`}
          >
            {rank}
          </div>
        )}
        <div className="flex-1">
          <Link href={`/posts/${post.id}`}>
            <h3 className="font-semibold text-gray-900 hover:text-indigo-600 mb-2 line-clamp-1">
              {post.title}
            </h3>
          </Link>
          <div className="flex items-center gap-3 text-xs text-gray-400">
            <div className="flex items-center gap-1">
              <EyeIcon className="w-3 h-3" />
              {post.view_count}
            </div>
            <div className="flex items-center gap-1">
              <ChatBubbleLeftRightIcon className="w-3 h-3" />
              {post.question?.answer_count || 0}
            </div>
            <button
              onClick={handleLike}
              className="flex items-center gap-1 hover:text-red-500"
            >
              {liked ? (
                <HeartCrack className="w-3 h-3 text-red-500" />
              ) : (
                <HeartIcon className="w-3 h-3" />
              )}
              {likesCount}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
