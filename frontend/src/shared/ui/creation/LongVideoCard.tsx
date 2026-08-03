"use client";

import Link from "next/link";
import Image from "next/image";
import { timeAgo } from "@/shared/lib/utils";
import { Play, Video, Clock } from "lucide-react";
import { UserDO } from "@/shared/api/types/user.model.do";
import { getYouTubeThumb, isYouTubeUrl, normalizeUrl } from "./WorkCard";

/** 长视频：YouTube 风格横向布局 */
export function LongVideoCard({
  coverUrl,
  videoUrl,
  title,
  postId,
  author,
  createdAt,
}: {
  coverUrl?: string;
  videoUrl?: string;
  title: string;
  postId: number;
  author?: UserDO;
  createdAt: string;
}) {
  const thumbUrl =
    videoUrl && isYouTubeUrl(videoUrl) ? getYouTubeThumb(videoUrl) : coverUrl;
  return (
    <Link href={`/posts/${postId}`} className="flex gap-3 p-3 group">
      <div className="relative w-40 lg:w-48 aspect-video rounded-xl overflow-hidden flex-shrink-0 bg-base-200">
        {thumbUrl ? (
          <Image
            src={normalizeUrl(thumbUrl)}
            alt={title}
            fill
            className="object-cover group-hover:scale-105 transition-transform duration-300"
            sizes="192px"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-base-content/20">
            <Video className="w-8 h-8" />
          </div>
        )}
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="w-10 h-10 rounded-full bg-black/70 flex items-center justify-center">
            <Play className="w-5 h-5 text-white ml-0.5" />
          </div>
        </div>
      </div>
      <div className="flex-1 min-w-0 py-1">
        <h3 className="font-semibold text-sm line-clamp-2 group-hover:text-primary transition-colors">
          {title}
        </h3>
        <p className="text-xs text-base-content/50 mt-1">{author?.username}</p>
        <p className="text-xs text-base-content/40 mt-1 flex items-center gap-1">
          <Clock className="w-3 h-3" />
          {timeAgo(createdAt)}
        </p>
      </div>
    </Link>
  );
}
