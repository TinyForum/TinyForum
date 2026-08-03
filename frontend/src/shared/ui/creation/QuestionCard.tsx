"use client";

import Link from "next/link";
import { timeAgo, truncate } from "@/shared/lib/utils";
import { MessageSquare, HelpCircle } from "lucide-react";
import Avatar from "@/shared/ui/common/Avatar";
import { UserDO } from "@/shared/api/types/user.model.do";

/** 问答 */
export function QuestionCard({
  title,
  plainContent,
  postId,
  author,
  createdAt,
  commentCount,
}: {
  title: string;
  plainContent: string;
  postId: number;
  author?: UserDO;
  createdAt: string;
  commentCount?: number;
}) {
  return (
    <div className="p-4">
      <div className="flex items-start gap-2">
        <HelpCircle className="w-5 h-5 text-amber-500 flex-shrink-0 mt-0.5" />
        <div className="flex-1 min-w-0">
          <Link href={`/posts/${postId}`}>
            <h3 className="font-semibold text-sm line-clamp-2 hover:text-primary transition-colors">
              {title}
            </h3>
          </Link>
          {plainContent && (
            <p className="text-xs text-base-content/50 mt-1 line-clamp-2">
              {truncate(plainContent, 120)}
            </p>
          )}
          <div className="flex items-center gap-3 mt-2 text-xs text-base-content/40">
            {author && (
              <span className="flex items-center gap-1">
                <Avatar
                  username={author.username}
                  avatarUrl={author.avatar_url}
                  size="sm"
                />
                {author.username}
              </span>
            )}
            <span>{timeAgo(createdAt)}</span>
            <span className="flex items-center gap-1">
              <MessageSquare className="w-3 h-3" />
              {commentCount ?? 0} 回答
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
