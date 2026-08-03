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
  Video,
  Clock,
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

function isYouTubeUrl(url: string): boolean {
  return /(youtube\.com|youtu\.be)/.test(url);
}

function getYouTubeThumb(url: string): string {
  const match = url.match(/(?:v=|\/)([\w-]{11})/);
  if (match) return `https://img.youtube.com/vi/${match[1]}/hqdefault.jpg`;
  return "";
}

interface WorkCardProps {
  post: Post;
  commentCount?: number;
}

export default function WorkCard({ post, commentCount }: WorkCardProps) {
  const creation = post?.creation;
  const creationType = (creation?.creation_type || "image_text") as PostType;
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
    <div className="bg-base-100 rounded-2xl overflow-hidden shadow-sm hover:shadow-lg transition-shadow duration-300 border border-base-200/60 h-fit break-inside-avoid">
      {/* 文章：小红书风格封面 Hero */}
      {creationType === "article" && coverUrl && (
        <ArticleCard coverUrl={coverUrl} title={creation.title} postId={post.id} author={author} createdAt={post.created_at} />
      )}

      {/* 长视频：YouTube 风格横向卡片 */}
      {creationType === "long_video" && (
        <LongVideoCard coverUrl={coverUrl} videoUrl={videoUrl} title={creation.title} postId={post.id} author={author} createdAt={post.created_at} />
      )}

      {/* 短视频：小红书风格封面 + 播放 */}
      {creationType === "short_video" && (
        <ShortVideoCard coverUrl={coverUrl} title={creation.title} postId={post.id} author={author} />
      )}

      {/* 图文/图片/话题/帖子：Instagram 风格图片矩阵 */}
      {(creationType === "image_text" || creationType === "image" || creationType === "topic" || creationType === "post") && imageUrls.length > 0 && (
        <ImageCard images={imageUrls} title={creation.title} postId={post.id} author={author} createdAt={post.created_at} plainContent={plainContent} />
      )}

      {/* 无图片的帖子/话题：Twitter 风格纯文本 */}
      {(creationType === "post" || creationType === "topic" || creationType === "image_text" || creationType === "image") && imageUrls.length === 0 && !coverUrl && (
        <TextCard title={creation.title} plainContent={plainContent} postId={post.id} author={author} createdAt={post.created_at} />
      )}

      {/* 问答：简洁卡片 */}
      {creationType === "question" && (
        <QuestionCard title={creation.title} plainContent={plainContent} postId={post.id} author={author} createdAt={post.created_at} commentCount={commentCount} />
      )}

      {/* 图片类型无图但有封面时 */}
      {(creationType === "image_text" || creationType === "image" || creationType === "topic" || creationType === "post") && imageUrls.length === 0 && coverUrl && (
        <ImageCard images={[coverUrl]} title={creation.title} postId={post.id} author={author} createdAt={post.created_at} plainContent={plainContent} />
      )}

      {/* 通用底部栏：置顶标识 + 类型标签 + 统计 */}
      <CardFooter creationType={creationType} pinTop={creation.pin_top} viewCount={creation.view_count} likeCount={creation.like_count} commentCount={commentCount} />
    </div>
  );
}

/** 通用统计底部栏 */
function CardFooter({ creationType, pinTop, viewCount, likeCount, commentCount }: { creationType: PostType; pinTop: boolean; viewCount: number; likeCount: number; commentCount?: number }) {
  const typeLabels: Record<string, string> = {
    image_text: "图文", short_video: "短视频", long_video: "长视频",
    image: "图片", article: "文章", question: "问答", topic: "话题", post: "帖子",
  };

  return (
    <div className="flex items-center gap-3 px-4 py-2.5 border-t border-base-200/60 text-xs text-base-content/50">
      {pinTop && <span className="text-primary font-medium flex items-center gap-1"><Pin className="w-3 h-3" />置顶</span>}
      <span className="bg-base-200/80 px-1.5 py-0.5 rounded text-[10px]">{typeLabels[creationType]}</span>
      <span className="flex items-center gap-1 ml-auto"><Eye className="w-3 h-3" />{viewCount}</span>
      <span className="flex items-center gap-1"><Heart className="w-3 h-3" />{likeCount}</span>
      {commentCount !== undefined && <span className="flex items-center gap-1"><MessageSquare className="w-3 h-3" />{commentCount}</span>}
    </div>
  );
}

/** 文章：小红书式封面 + 标题摘要 */
function ArticleCard({ coverUrl, title, postId, author, createdAt }: { coverUrl: string; title: string; postId: number; author?: UserDO; createdAt: string }) {
  return (
    <Link href={`/posts/${postId}`} className="block group">
      <div className="relative w-full aspect-[4/3]">
        <Image src={normalizeUrl(coverUrl)} alt={title} fill className="object-cover group-hover:scale-[1.02] transition-transform duration-300" sizes="(max-width: 768px) 100vw, 33vw" />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
        <div className="absolute bottom-0 left-0 right-0 p-4">
          <h3 className="text-white font-bold text-lg leading-tight line-clamp-2 drop-shadow-lg">{title}</h3>
          <div className="flex items-center gap-2 mt-2 text-white/80 text-xs">
            {author && <Avatar username={author.username} avatarUrl={author.avatar_url} size="sm" />}
            <span>{author?.username}</span>
            <span>·</span>
            <span>{timeAgo(createdAt)}</span>
          </div>
        </div>
      </div>
    </Link>
  );
}

/** 长视频：YouTube 风格横向布局 */
function LongVideoCard({ coverUrl, videoUrl, title, postId, author, createdAt }: { coverUrl?: string; videoUrl?: string; title: string; postId: number; author?: UserDO; createdAt: string }) {
  const thumbUrl = videoUrl && isYouTubeUrl(videoUrl) ? getYouTubeThumb(videoUrl) : coverUrl;
  return (
    <Link href={`/posts/${postId}`} className="flex gap-3 p-3 group">
      <div className="relative w-40 lg:w-48 aspect-video rounded-xl overflow-hidden flex-shrink-0 bg-base-200">
        {thumbUrl ? (
          <Image src={normalizeUrl(thumbUrl)} alt={title} fill className="object-cover group-hover:scale-105 transition-transform duration-300" sizes="192px" />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-base-content/20"><Video className="w-8 h-8" /></div>
        )}
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="w-10 h-10 rounded-full bg-black/70 flex items-center justify-center">
            <Play className="w-5 h-5 text-white ml-0.5" />
          </div>
        </div>
      </div>
      <div className="flex-1 min-w-0 py-1">
        <h3 className="font-semibold text-sm line-clamp-2 group-hover:text-primary transition-colors">{title}</h3>
        <p className="text-xs text-base-content/50 mt-1">{author?.username}</p>
        <p className="text-xs text-base-content/40 mt-1 flex items-center gap-1"><Clock className="w-3 h-3" />{timeAgo(createdAt)}</p>
      </div>
    </Link>
  );
}

/** 短视频：小红书/Instagram 风格竖版封面 */
function ShortVideoCard({ coverUrl, title, postId, author }: { coverUrl?: string; title: string; postId: number; author?: UserDO }) {
  return (
    <Link href={`/posts/${postId}`} className="block group">
      <div className="relative w-full aspect-[3/4] bg-base-200">
        {coverUrl ? (
          <Image src={normalizeUrl(coverUrl)} alt={title} fill className="object-cover" sizes="(max-width: 768px) 100vw, 33vw" />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-base-content/20"><Video className="w-12 h-12" /></div>
        )}
        <div className="absolute inset-0 bg-black/20 flex items-center justify-center">
          <div className="w-14 h-14 rounded-full bg-black/60 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Play className="w-7 h-7 text-white ml-0.5" />
          </div>
        </div>
        <div className="absolute top-3 left-3 bg-black/50 backdrop-blur-sm text-white text-[10px] px-2 py-0.5 rounded-full">短视频</div>
      </div>
      <div className="p-3">
        <h3 className="font-medium text-sm line-clamp-2">{title}</h3>
        <div className="flex items-center gap-1.5 mt-2">
          {author && <Avatar username={author.username} avatarUrl={author.avatar_url} size="sm" />}
          <span className="text-xs text-base-content/50">{author?.username}</span>
        </div>
      </div>
    </Link>
  );
}

/** 图片矩阵：Instagram 风格 */
function ImageCard({ images, title, postId, author, createdAt, plainContent }: { images: string[]; title: string; postId: number; author?: UserDO; createdAt: string; plainContent: string }) {
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
            <div key={idx} className={`relative ${isSingle ? 'aspect-[4/3]' : 'aspect-square'} overflow-hidden`}>
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
        <Link href={`/posts/${postId}`}><h3 className="font-medium text-sm line-clamp-2 hover:text-primary transition-colors">{title}</h3></Link>
        {plainContent && <p className="text-xs text-base-content/50 mt-1 line-clamp-2">{truncate(plainContent, 100)}</p>}
        <div className="flex items-center gap-1.5 mt-2">
          {author && <Avatar username={author.username} avatarUrl={author.avatar_url} size="sm" />}
          <span className="text-[11px] text-base-content/50">{author?.username} · {timeAgo(createdAt)}</span>
        </div>
      </div>
    </div>
  );
}

/** 纯文本：Twitter 风格 */
function TextCard({ title, plainContent, postId, author, createdAt }: { title: string; plainContent: string; postId: number; author?: UserDO; createdAt: string }) {
  return (
    <div className="p-4">
      <div className="flex items-start gap-3">
        <Link href={author ? `/users/${author.id}` : "#"} className="flex-shrink-0">
          <Avatar username={author?.username} avatarUrl={author?.avatar_url} size="sm" />
        </Link>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-1 flex-wrap">
            <Link href={author ? `/users/${author.id}` : "#"} className="font-semibold text-sm hover:text-primary transition-colors">{author?.username}</Link>
            <span className="text-xs text-base-content/40">@{author?.username}</span>
            <span className="text-xs text-base-content/40">·</span>
            <span className="text-xs text-base-content/40">{timeAgo(createdAt)}</span>
          </div>
          <Link href={`/posts/${postId}`} className="block mt-1">
            <h3 className="font-bold text-sm">{title}</h3>
            {plainContent && <p className="text-sm text-base-content/70 mt-1 line-clamp-6">{truncate(plainContent, 280)}</p>}
          </Link>
        </div>
      </div>
    </div>
  );
}

/** 问答 */
function QuestionCard({ title, plainContent, postId, author, createdAt, commentCount }: { title: string; plainContent: string; postId: number; author?: UserDO; createdAt: string; commentCount?: number }) {
  return (
    <div className="p-4">
      <div className="flex items-start gap-2">
        <HelpCircle className="w-5 h-5 text-amber-500 flex-shrink-0 mt-0.5" />
        <div className="flex-1 min-w-0">
          <Link href={`/posts/${postId}`}>
            <h3 className="font-semibold text-sm line-clamp-2 hover:text-primary transition-colors">{title}</h3>
          </Link>
          {plainContent && <p className="text-xs text-base-content/50 mt-1 line-clamp-2">{truncate(plainContent, 120)}</p>}
          <div className="flex items-center gap-3 mt-2 text-xs text-base-content/40">
            {author && <span className="flex items-center gap-1"><Avatar username={author.username} avatarUrl={author.avatar_url} size="sm" />{author.username}</span>}
            <span>{timeAgo(createdAt)}</span>
            <span className="flex items-center gap-1"><MessageSquare className="w-3 h-3" />{commentCount ?? 0} 回答</span>
          </div>
        </div>
      </div>
    </div>
  );
}
