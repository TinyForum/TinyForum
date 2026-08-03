"use client";

import Link from "next/link";
import Image from "next/image";
import { timeAgo, truncate } from "@/shared/lib/utils";
import Avatar from "@/shared/ui/common/Avatar";
import { UserDO } from "@/shared/api/types/user.model.do";
import { normalizeUrl } from "./WorkCard";

/** 图片矩阵：Instagram 风格 */
export function ImageCard({
  images,
  title,
  postId,
  author,
  createdAt,
  plainContent,
}: {
  images: string[];
  title: string;
  postId: number;
  author?: UserDO;
  createdAt: string;
  plainContent: string;
}) {
  const count = images.length;
  const cols = count === 1 ? 1 : count === 2 ? 2 : count === 4 ? 2 : 3;
  const maxDisplay = 9;
  const isSingle = count === 1;

  return (
    <div className="group">
      <Link href={`/posts/${postId}`}>
        <div
          className="grid gap-0.5"
          style={{ gridTemplateColumns: `repeat(${cols}, 1fr)` }}
        >
          {images.slice(0, maxDisplay).map((url, idx) => (
            <div
              key={idx}
              className={`relative ${isSingle ? "aspect-[4/3]" : "aspect-square"} overflow-hidden`}
            >
              <Image
                src={normalizeUrl(url)}
                alt=""
                fill
                className="object-cover group-hover:scale-[1.02] transition-transform duration-300"
                sizes={cols === 1 ? "100vw" : cols === 2 ? "50vw" : "33vw"}
              />
              {count > maxDisplay && idx === maxDisplay - 1 && (
                <div className="absolute inset-0 bg-black/60 flex items-center justify-center text-white text-2xl font-bold">
                  +{count - maxDisplay}
                </div>
              )}
            </div>
          ))}
        </div>
      </Link>
      <div className="p-3">
        <Link href={`/posts/${postId}`}>
          <h3 className="font-medium text-sm line-clamp-2 hover:text-primary transition-colors">
            {title}
          </h3>
        </Link>
        {plainContent && (
          <p className="text-xs text-base-content/50 mt-1 line-clamp-2">
            {truncate(plainContent, 100)}
          </p>
        )}
        <div className="flex items-center gap-1.5 mt-2">
          {author && (
            <Avatar
              username={author.username}
              avatarUrl={author.avatar_url}
              size="sm"
            />
          )}
          <span className="text-[11px] text-base-content/50">
            {author?.username} · {timeAgo(createdAt)}
          </span>
        </div>
      </div>
    </div>
  );
}
