"use client";

import Link from "next/link";
import Image from "next/image";
import { useQuery } from "@tanstack/react-query";
import { timeAgo, truncate } from "@/shared/lib/utils";
import {
  Eye,
  Heart,
  MessageSquare,
  Pin,
  HelpCircle,
  Play,
  Image as ImageIcon,
  FileText,
  Video,
} from "lucide-react";
import Avatar from "@/shared/ui/common/Avatar";
import { userApi } from "@/shared/api/modules/user";
import { UserDO } from "@/shared/api/types/user.model.do";
import { Post, PostType } from "@/shared/api/types/post.model";

function normalizeUrl(url: string): string {
  if (!url) return url;
  if (url.startsWith("http://") || url.startsWith("https://") || url.startsWith("/")) return url;
  return "/" + url;
}

const TYPE_META: Record<PostType, { label: string; color: string; icon: React.ReactNode }> = {
  image_text:  { label: "图文", color: "bg-sky-100 text-sky-700", icon: <ImageIcon className="w-3 h-3" /> },
  short_video: { label: "短视频", color: "bg-rose-100 text-rose-700", icon: <Video className="w-3 h-3" /> },
  long_video:  { label: "长视频", color: "bg-red-100 text-red-700", icon: <Video className="w-3 h-3" /> },
  image:       { label: "图片", color: "bg-emerald-100 text-emerald-700", icon: <ImageIcon className="w-3 h-3" /> },
  article:     { label: "文章", color: "bg-violet-100 text-violet-700", icon: <FileText className="w-3 h-3" /> },
  question:    { label: "问答", color: "bg-amber-100 text-amber-700", icon: <HelpCircle className="w-3 h-3" /> },
  topic:       { label: "话题", color: "bg-cyan-100 text-cyan-700", icon: <MessageSquare className="w-3 h-3" /> },
  post:        { label: "帖子", color: "bg-blue-100 text-blue-700", icon: <FileText className="w-3 h-3" /> },
};

interface WorkCardProps {
  post: Post;
  commentCount?: number;
}

export default function WorkCard({ post, commentCount }: WorkCardProps) {
  const creation = post?.creation;
  const creationType = (creation?.creation_type || "image_text") as PostType;
  const meta = TYPE_META[creationType] || TYPE_META.image_text;
  const coverUrl = creation?.cover_url;
  const videoUrl = creation?.video_url;
  const imageUrls = creation?.image_urls || [];

  const { data: fetchedAuthor } = useQuery({
    queryKey: ["user", creation?.author?.id],
    queryFn: () => userApi.getProfile(creation!.author_id).then((r) => r.data.data),
    enabled: !!(post && creation && !creation.author && creation.author_id),
  });

  if (!post || !creation) return null;

  const author = (creation.author || fetchedAuthor) as UserDO | undefined;
  const stripHtml = (html: string) => html.replace(/<[^>]*>/g, "").trim();
  const plainContent = stripHtml(creation.content || "");

  return (
    <div
      className={`bg-base-100 shadow-sm hover:shadow-md transition-all duration-200 border border-base-300 rounded-xl overflow-hidden h-fit ${
        creation.pin_top ? "ring-2 ring-primary/30" : ""
      }`}
    >
      {/* 文章：封面 Hero */}
      {creationType === "article" && coverUrl && (
        <ArticleCover coverUrl={coverUrl} title={creation.title} postId={post.id} />
      )}

      {/* 视频：内嵌视频播放器 */}
      {(creationType === "short_video" || creationType === "long_video") && (
        <VideoPreview
          coverUrl={coverUrl}
          videoUrl={videoUrl}
          title={creation.title}
          postId={post.id}
        />
      )}

      {/* 图片矩阵：image_text / image / topic / post */}
      {["image_text", "image", "topic", "post"].includes(creationType) && imageUrls.length > 0 && (
        <ImageGridDisplay images={imageUrls} postId={post.id} />
      )}

      {/* 图片类型无图但有封面时用小图 */}
      {["image_text", "image", "topic", "post"].includes(creationType) && imageUrls.length === 0 && coverUrl && (
        <Link href={`/posts/${post.id}`} className="block">
          <div className="relative w-full aspect-[2/1]">
            <Image src={normalizeUrl(coverUrl)} alt={creation.title} fill className="object-cover" sizes="(max-width: 768px) 100vw, 50vw" />
          </div>
        </Link>
      )}

      <div className="p-4">
        {/* 类型标签 + 置顶 */}
        <div className="flex items-center gap-2 flex-wrap mb-2">
          {creation.pin_top && (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full bg-primary/10 text-primary">
              <Pin className="w-3 h-3" /> 置顶
            </span>
          )}
          <span className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full ${meta.color}`}>
            {meta.icon}
            {meta.label}
          </span>
        </div>

        {/* 标题 */}
        <Link href={`/posts/${post.id}`} className="group">
          <h2 className="text-base font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-2">
            {creation.title}
          </h2>
        </Link>

        {/* 内容摘要 */}
        {plainContent && (
          <p className="text-sm text-base-content/60 mt-1 line-clamp-3">
            {truncate(plainContent, 150)}
          </p>
        )}

        {/* 问答：回答数 */}
        {creationType === "question" && (
          <div className="flex items-center gap-3 mt-2 text-xs text-base-content/50">
            <span className="flex items-center gap-1">
              <MessageSquare className="w-3 h-3" />{commentCount ?? 0} 个回答
            </span>
          </div>
        )}

        {/* 底部信息栏 */}
        <div className="flex items-center justify-between mt-3 pt-3 border-t border-base-200 flex-wrap gap-2">
          <div className="flex items-center gap-2">
            <Link href={`/users/${creation.author_id}`}>
              <Avatar username={author?.username || `用户${creation.author_id}`} avatarUrl={author?.avatar_url} size="sm" />
            </Link>
            <div className="flex flex-col">
              <Link href={`/users/${creation.author_id}`} className="text-xs font-medium text-base-content hover:text-primary transition-colors">
                {author?.username || `用户${creation.author_id}`}
              </Link>
              <span className="text-[10px] text-base-content/40">{timeAgo(post.created_at)}</span>
            </div>
          </div>

          <div className="flex items-center gap-3 text-xs text-base-content/50">
            <span className="flex items-center gap-1"><Eye className="w-3.5 h-3.5" /> {creation.view_count}</span>
            <span className="flex items-center gap-1"><Heart className="w-3.5 h-3.5" /> {creation.like_count}</span>
            {commentCount !== undefined && (
              <span className="flex items-center gap-1"><MessageSquare className="w-3.5 h-3.5" /> {commentCount}</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

/** 文章封面 Hero */
function ArticleCover({ coverUrl, title, postId }: { coverUrl: string; title: string; postId: number }) {
  return (
    <Link href={`/posts/${postId}`} className="block">
      <div className="relative w-full aspect-[2/1]">
        <Image src={normalizeUrl(coverUrl)} alt={title} fill className="object-cover" sizes="(max-width: 768px) 100vw, 50vw" />
      </div>
    </Link>
  );
}

/** 视频预览（封面 + 播放按钮 / 内嵌播放器） */
function VideoPreview({ coverUrl, videoUrl, title, postId }: { coverUrl?: string; videoUrl?: string; title: string; postId: number }) {
  if (videoUrl) {
    return (
      <Link href={`/posts/${postId}`} className="block relative">
        {videoUrl.includes("youtube.com") || videoUrl.includes("youtu.be") ? (
          <div className="relative w-full aspect-video">
            <iframe
              src={videoUrl.replace("watch?v=", "embed/").replace("youtu.be/", "youtube.com/embed/")}
              title={title}
              className="w-full h-full"
              allowFullScreen
            />
          </div>
        ) : (
          <div className="relative w-full aspect-video bg-base-200">
            <video src={normalizeUrl(videoUrl)} className="w-full h-full object-contain" controls preload="metadata" />
          </div>
        )}
      </Link>
    );
  }

  if (coverUrl) {
    return (
      <Link href={`/posts/${postId}`} className="block relative">
        <div className="relative w-full aspect-video">
          <Image src={normalizeUrl(coverUrl)} alt={title} fill className="object-cover" sizes="(max-width: 768px) 100vw, 50vw" />
          <div className="absolute inset-0 flex items-center justify-center bg-black/20">
            <div className="w-14 h-14 rounded-full bg-black/60 flex items-center justify-center">
              <Play className="w-7 h-7 text-white ml-0.5" />
            </div>
          </div>
        </div>
      </Link>
    );
  }

  return null;
}

/** 九宫格图片矩阵 */
function ImageGridDisplay({ images, postId }: { images: string[]; postId: number }) {
  const count = images.length;
  if (count === 0) return null;

  const cols = count === 1 ? 1 : count === 4 ? 2 : 3;
  const maxDisplay = 9;

  return (
    <Link href={`/posts/${postId}`} className="block">
      <div
        className="grid gap-0.5"
        style={{ gridTemplateColumns: `repeat(${cols}, 1fr)` }}
      >
        {images.slice(0, maxDisplay).map((url, idx) => (
          <div key={idx} className="relative aspect-square">
            <Image
              src={normalizeUrl(url)}
              alt=""
              fill
              className="object-cover"
              sizes={cols === 1 ? "100vw" : cols === 2 ? "50vw" : "33vw"}
            />
            {count > maxDisplay && idx === maxDisplay - 1 && (
              <div className="absolute inset-0 bg-black/60 flex items-center justify-center text-white text-lg font-bold">
                +{count - maxDisplay}
              </div>
            )}
          </div>
        ))}
      </div>
    </Link>
  );
}
