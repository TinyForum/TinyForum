// components/topic/TopicPostCard.tsx
"use client";

import Link from "next/link";
import {
  HeartIcon,
  ChatBubbleLeftRightIcon,
  EyeIcon,
} from "@heroicons/react/24/outline";
import type { Post } from "@/shared/api/types";
import { formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";

interface TopicPostCardProps {
  post: Post;
}

export function TopicPostCard({ post }: TopicPostCardProps) {
  const timeAgo = formatDistanceToNow(new Date(post.created_at), {
    addSuffix: true,
    locale: zhCN,
  });

  return (
    <Link href={`/questions/${post.id}`}>
      <div className="card bg-base-100 shadow-sm hover:shadow-md transition-all duration-300 cursor-pointer group">
        <div className="card-body p-4">
          <h3 className="card-title text-base font-bold text-base-content group-hover:text-primary transition-colors">
            {post.title}
          </h3>

          <p className="text-base-content/60 text-sm line-clamp-2">
            {post.content.replace(/<[^>]*>/g, "").substring(0, 100)}
          </p>

          <div className="flex flex-wrap items-center gap-4 text-xs text-base-content/40 mt-2">
            <div className="flex items-center gap-1">
              <div className="avatar placeholder">
                <div className="w-4 h-4 rounded-full bg-primary/10 text-primary">
                  <span className="text-[10px]">
                    {post.author?.username?.[0]?.toUpperCase() || "U"}
                  </span>
                </div>
              </div>
              <span>{post.author?.username || `用户${post.author_id}`}</span>
            </div>

            <div className="flex items-center gap-1">
              <EyeIcon className="w-3.5 h-3.5" />
              <span>{post.view_count || 0}</span>
            </div>

            <div className="flex items-center gap-1">
              <HeartIcon className="w-3.5 h-3.5" />
              <span>{post.like_count || 0}</span>
            </div>

            <div className="flex items-center gap-1">
              <ChatBubbleLeftRightIcon className="w-3.5 h-3.5" />
              {/* <span>{post.answer_count || 0}</span> */}
            </div>

            <span>{timeAgo}</span>

            {/* {post.is_accepted && (
              <span className="badge badge-success badge-xs">已解决</span>
            )} */}
          </div>
        </div>
      </div>
    </Link>
  );
}
