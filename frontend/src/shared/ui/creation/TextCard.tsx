"use client";

import Link from "next/link";
import { timeAgo, truncate } from "@/shared/lib/utils";
import Avatar from "@/shared/ui/common/Avatar";
import { UserDO } from "@/shared/api/types/user.model.do";

/** 纯文本：Twitter 风格 */
export function TextCard({
  title,
  plainContent,
  postId,
  author,
  createdAt,
}: {
  title: string;
  plainContent: string;
  postId: number;
  author?: UserDO;
  createdAt: string;
}) {
  return (
    <div className="p-4">
      <div className="flex items-start gap-3">
        <Link
          href={author ? `/users/${author.id}` : "#"}
          className="flex-shrink-0"
        >
          <Avatar
            username={author?.username}
            avatarUrl={author?.avatar_url}
            size="sm"
          />
        </Link>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-1 flex-wrap">
            <Link
              href={author ? `/users/${author.id}` : "#"}
              className="font-semibold text-sm hover:text-primary transition-colors"
            >
              {author?.username}
            </Link>
            <span className="text-xs text-base-content/40">
              @{author?.username}
            </span>
            <span className="text-xs text-base-content/40">·</span>
            <span className="text-xs text-base-content/40">
              {timeAgo(createdAt)}
            </span>
          </div>
          <Link href={`/posts/${postId}`} className="block mt-1">
            <h3 className="font-bold text-sm">{title}</h3>
            {plainContent && (
              <p className="text-sm text-base-content/70 mt-1 line-clamp-6">
                {truncate(plainContent, 280)}
              </p>
            )}
          </Link>
        </div>
      </div>
    </div>
  );
}
