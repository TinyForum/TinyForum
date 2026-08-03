"use client";

import Link from "next/link";
import Image from "next/image";
import { Play, Video } from "lucide-react";
import Avatar from "@/shared/ui/common/Avatar";
import { UserDO } from "@/shared/api/types/user.model.do";
import { normalizeUrl, getYouTubeThumb, isYouTubeUrl } from "./WorkCard";

/** 短视频：小红书/Instagram 风格竖版封面 */
export function ShortVideoCard({
  coverUrl,
  videoUrl,
  title,
  postId,
  author,
}: {
  coverUrl?: string;
  videoUrl?: string;
  title: string;
  postId: number;
  author?: UserDO;
}) {
  const thumbUrl = coverUrl
    || (videoUrl && isYouTubeUrl(videoUrl) ? getYouTubeThumb(videoUrl) : undefined);

  return (
    <Link href={`/posts/${postId}`} className="block group">
      <div className="relative w-full aspect-[3/4] bg-base-200">
        {thumbUrl ? (
          <Image
            src={normalizeUrl(thumbUrl)}
            alt={title}
            fill
            className="object-cover"
            sizes="(max-width: 768px) 100vw, 33vw"
          />
        ) : videoUrl && !isYouTubeUrl(videoUrl) ? (
          <video
            src={normalizeUrl(videoUrl)}
            className="w-full h-full object-cover"
            preload="metadata"
            muted
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-base-content/20">
            <Video className="w-12 h-12" />
          </div>
        )}
        <div className="absolute inset-0 bg-black/20 flex items-center justify-center">
          <div className="w-14 h-14 rounded-full bg-black/60 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Play className="w-7 h-7 text-white ml-0.5" />
          </div>
        </div>
        <div className="absolute top-3 left-3 bg-black/50 backdrop-blur-sm text-white text-[10px] px-2 py-0.5 rounded-full">
          短视频
        </div>
      </div>
      <div className="p-3">
        <h3 className="font-medium text-sm line-clamp-2">{title}</h3>
        <div className="flex items-center gap-1.5 mt-2">
          {author && (
            <Avatar
              username={author.username}
              avatarUrl={author.avatar_url}
              size="sm"
            />
          )}
          <span className="text-xs text-base-content/50">
            {author?.username}
          </span>
        </div>
      </div>
    </Link>
  );
}
