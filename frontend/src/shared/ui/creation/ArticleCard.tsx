"use client";

import Link from "next/link";
import Image from "next/image";
import { timeAgo } from "@/shared/lib/utils";

import Avatar from "@/shared/ui/common/Avatar";
import { UserDO } from "@/shared/api/types/user.model.do";
import { normalizeUrl } from "./WorkCard";

/** 文章：小红书式封面 + 标题摘要 */
export function ArticleCard({
  coverUrl,
  title,
  postId,
  author,
  createdAt,
}: {
  coverUrl: string;
  title: string;
  postId: number;
  author?: UserDO;
  createdAt: string;
}) {
  return (
    <Link href={`/posts/${postId}`} className="block group">
      <div className="relative w-full aspect-[4/3]">
        <Image
          src={normalizeUrl(coverUrl)}
          alt={title}
          fill
          className="object-cover group-hover:scale-[1.02] transition-transform duration-300"
          sizes="(max-width: 768px) 100vw, 33vw"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
        <div className="absolute bottom-0 left-0 right-0 p-4">
          <h3 className="text-white font-bold text-lg leading-tight line-clamp-2 drop-shadow-lg">
            {title}
          </h3>
          <div className="flex items-center gap-2 mt-2 text-white/80 text-xs">
            {author && (
              <Avatar
                username={author.username}
                avatarUrl={author.avatar_url}
                size="sm"
              />
            )}
            <span>{author?.username}</span>
            <span>·</span>
            <span>{timeAgo(createdAt)}</span>
          </div>
        </div>
      </div>
    </Link>
  );
}
