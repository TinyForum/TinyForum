"use client";

import { Eye, Heart, MessageSquare, Pin } from "lucide-react";
import { useTranslations } from "next-intl";
import { PostType } from "@/shared/api/types/post.model";

/** 通用统计底部栏 */
export function CardFooter({
  creationType,
  pinTop,
  viewCount,
  likeCount,
  commentCount,
}: {
  creationType: PostType;
  pinTop: boolean;
  viewCount: number;
  likeCount: number;
  commentCount?: number;
}) {
  const t = useTranslations("Common");

  return (
    <div className="flex items-center gap-3 px-4 py-2.5 border-t border-base-200/60 text-xs text-base-content/50">
      {pinTop && (
        <span className="text-primary font-medium flex items-center gap-1">
          <Pin className="w-3 h-3" />
          {t("pinned")}
        </span>
      )}
      <span className="bg-base-200/80 px-1.5 py-0.5 rounded text-[10px]">
        {t(`creation_type.${creationType}`)}
      </span>
      <span className="flex items-center gap-1 ml-auto">
        <Eye className="w-3 h-3" />
        {viewCount}
      </span>
      <span className="flex items-center gap-1">
        <Heart className="w-3 h-3" />
        {likeCount}
      </span>
      {commentCount !== undefined && (
        <span className="flex items-center gap-1">
          <MessageSquare className="w-3 h-3" />
          {commentCount}
        </span>
      )}
    </div>
  );
}
