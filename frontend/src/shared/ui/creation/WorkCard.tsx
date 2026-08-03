"use client";

import { useQuery } from "@tanstack/react-query";
import { userApi } from "@/shared/api/modules/user";
import { UserDO } from "@/shared/api/types/user.model.do";
import { Post, PostType } from "@/shared/api/types/post.model";
import { ArticleCard } from "./ArticleCard";
import { CardFooter } from "./CardFooter";
import { ImageCard } from "./ImageCard";
import { QuestionCard } from "./QuestionCard";
import { ShortVideoCard } from "./ShortVideoCard";
import { TextCard } from "./TextCard";
import { LongVideoCard } from "./LongVideoCard";

export function normalizeUrl(url: string): string {
  if (!url) return url;
  if (
    url.startsWith("http://") ||
    url.startsWith("https://") ||
    url.startsWith("/")
  )
    return url;
  return "/" + url;
}

export function isYouTubeUrl(url: string): boolean {
  return /(youtube\.com|youtu\.be)/.test(url);
}

export function getYouTubeThumb(url: string): string {
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
    queryFn: () =>
      userApi.getProfile(creation!.author_id).then((r) => r.data.data),
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
        <ArticleCard
          coverUrl={coverUrl}
          title={creation.title}
          postId={post.id}
          author={author}
          createdAt={post.created_at}
        />
      )}

      {/* 长视频：YouTube 风格横向卡片 */}
      {creationType === "long_video" && (
        <LongVideoCard
          coverUrl={coverUrl}
          videoUrl={videoUrl}
          title={creation.title}
          postId={post.id}
          author={author}
          createdAt={post.created_at}
        />
      )}

      {/* 短视频：小红书风格封面 + 播放 */}
      {creationType === "short_video" && (
        <ShortVideoCard
          coverUrl={coverUrl}
          videoUrl={videoUrl}
          title={creation.title}
          postId={post.id}
          author={author}
        />
      )}

      {/* 图文/图片/话题/帖子：Instagram 风格图片矩阵 */}
      {(creationType === "image_text" ||
        creationType === "image" ||
        creationType === "topic" ||
        creationType === "post") &&
        imageUrls.length > 0 && (
          <ImageCard
            images={imageUrls}
            title={creation.title}
            postId={post.id}
            author={author}
            createdAt={post.created_at}
            plainContent={plainContent}
          />
        )}

      {/* 无图片的帖子/话题：Twitter 风格纯文本 */}
      {(creationType === "post" ||
        creationType === "topic" ||
        creationType === "image_text" ||
        creationType === "image") &&
        imageUrls.length === 0 &&
        !coverUrl && (
          <TextCard
            title={creation.title}
            plainContent={plainContent}
            postId={post.id}
            author={author}
            createdAt={post.created_at}
          />
        )}

      {/* 问答：简洁卡片 */}
      {creationType === "question" && (
        <QuestionCard
          title={creation.title}
          plainContent={plainContent}
          postId={post.id}
          author={author}
          createdAt={post.created_at}
          commentCount={commentCount}
        />
      )}

      {/* 图片类型无图但有封面时 */}
      {(creationType === "image_text" ||
        creationType === "image" ||
        creationType === "topic" ||
        creationType === "post") &&
        imageUrls.length === 0 &&
        coverUrl && (
          <ImageCard
            images={[coverUrl]}
            title={creation.title}
            postId={post.id}
            author={author}
            createdAt={post.created_at}
            plainContent={plainContent}
          />
        )}

      {/* 通用底部栏：置顶标识 + 类型标签 + 统计 */}
      <CardFooter
        creationType={creationType}
        pinTop={creation.pin_top}
        viewCount={creation.view_count}
        likeCount={creation.like_count}
        commentCount={commentCount}
      />
    </div>
  );
}
