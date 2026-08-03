// src/app/posts/[id]/PostDetailClient.tsx
"use client";

import { useAuthStore } from "@/store/auth";
import Image from "next/image";
import Link from "next/link";
import { timeAgo, formatDate } from "@/shared/lib/utils";
import {
  Eye,
  Heart,
  Share2,
  Pencil,
  Trash2,
  Clock,
  ArrowLeft,
  Tag,
} from "lucide-react";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import Avatar from "@/shared/ui/common/Avatar";
import {
  useDeletePost,
  useLikePost,
  usePost,
  useUnlikePost,
} from "@/features/post/hooks/usePosts";

function normalizeUrl(url: string): string {
  if (!url) return url;
  if (url.startsWith("http://") || url.startsWith("https://") || url.startsWith("/")) return url;
  return "/" + url;
}

function isYouTubeUrl(url: string): boolean {
  return /(youtube\.com|youtu\.be)/.test(url);
}

function getYouTubeEmbedUrl(url: string): string {
  const match = url.match(/(?:v=|\/)([\w-]{11})/);
  if (match) return `https://www.youtube.com/embed/${match[1]}`;
  return url.replace("watch?v=", "embed/").replace("youtu.be/", "youtube.com/embed/");
}

// 导入自定义 hooks

export default function PostDetailClient({ postId }: { postId: number }) {
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const t = useTranslations("Post");

  // 1. 获取帖子详情（含是否已点赞）
  const { data, isLoading, error } = usePost(postId);

  // 2. 点赞 / 取消点赞 mutation
  const likePost = useLikePost();
  const unlikePost = useUnlikePost();

  // 3. 删除 mutation
  const deletePost = useDeletePost();

  // 处理点赞点击
  const handleLikeClick = () => {
    if (!isAuthenticated) {
      toast.error(t("please_login_first"));
      return;
    }
    // 根据当前 liked 状态决定调用 like 或 unlike
    if (data?.liked) {
      unlikePost.mutate(postId, {
        onSuccess: () => {
          toast.success(t("like_has_been_removed"));
        },
        onError: () => {
          toast.error(t("operation_failed"));
        },
      });
    } else {
      likePost.mutate(postId, {
        onSuccess: () => {
          toast.success(t("like_successful"));
        },
        onError: () => {
          toast.error(t("operation_failed"));
        },
      });
    }
  };

  // 处理分享
  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href);
    toast.success(t("link_copied"));
  };

  // 处理删除
  const handleDelete = () => {
    if (!confirm(t("are_you_sure_to_delete_this_post"))) return;
    deletePost.mutate(postId, {
      onSuccess: () => {
        toast.success(t("post_deleted"));
        router.push("/");
      },
      onError: () => {
        toast.error(t("deletion_failed"));
      },
    });
  };

  // 加载状态
  if (isLoading) {
    return <PostDetailSkeleton />;
  }

  // 错误或数据不存在
  if (error || !data) {
    return (
      <div className="text-center py-20">
        <p className="text-base-content/40 mb-4">
          {t("the_post_does_not_exist_or_has_been_deleted")}
        </p>
        <Link href="/" className="btn btn-primary">
          {t("back_to_home")}
        </Link>
      </div>
    );
  }

  const { post, liked } = data;
  const isAuthor = user?.id === post.creation.author?.id;
  const isAdmin = user?.role === "admin";

  return (
    <>
      {/* 返回按钮 */}
      <button
        onClick={() => router.back()}
        className="btn btn-ghost btn-sm gap-1 mb-4"
      >
        <ArrowLeft className="w-4 h-4" /> {t("back")}
      </button>

      <article className="card bg-base-100 border border-base-200/60 shadow-sm rounded-2xl overflow-hidden mb-6">
        <div className="card-body p-0">
          {/* 视频首屏 */}
          {post.creation.video_url && (
            <div className="bg-black">
              {isYouTubeUrl(post.creation.video_url) ? (
                <div className="relative w-full aspect-video">
                  <iframe
                    src={getYouTubeEmbedUrl(post.creation.video_url)}
                    title={post.creation.title}
                    className="w-full h-full"
                    allowFullScreen
                  />
                </div>
              ) : (
                <video
                  src={normalizeUrl(post.creation.video_url)}
                  controls
                  className="w-full max-h-[70vh] object-contain"
                  preload="metadata"
                  poster={post.creation.cover_url ? normalizeUrl(post.creation.cover_url) : undefined}
                />
              )}
            </div>
          )}

          {/* 封面图 Hero */}
          {!post.creation.video_url && post.creation.cover_url && (
            <div className="relative w-full aspect-[2/1] lg:aspect-[3/1] max-h-96">
              <Image
                src={normalizeUrl(post.creation.cover_url)}
                alt={post.creation.title}
                fill
                className="object-cover"
                priority
                sizes="(max-width: 768px) 100vw, 768px"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/40 via-transparent to-transparent" />
            </div>
          )}

          {/* 正文区域 */}
          <div className="p-5 lg:p-8">
            {/* 作者行 */}
            <div className="flex items-center gap-3 mb-5">
              <Link href={`/users/${post.creation.author?.id}`}>
                <Avatar
                  username={post.creation.author?.username}
                  avatarUrl={post.creation.author?.avatar_url}
                  size="md"
                />
              </Link>
              <div className="flex-1 min-w-0">
                <Link href={`/users/${post.creation.author?.id}`} className="font-semibold text-sm hover:text-primary transition-colors">
                  {post.creation.author?.username}
                </Link>
                <div className="flex items-center gap-2 text-xs text-base-content/40 mt-0.5">
                  <Clock className="w-3 h-3" />
                  <span title={formatDate(post.created_at)}>{timeAgo(post.created_at)}</span>
                  <span>·</span>
                  <Eye className="w-3 h-3" />
                  <span>{post.creation.view_count} 阅读</span>
                </div>
              </div>
              {(isAuthor || isAdmin) && (
                <div className="flex gap-1">
                  <Link href={`/posts/${postId}/edit`} className="btn btn-ghost btn-sm btn-circle">
                    <Pencil className="w-4 h-4" />
                  </Link>
                  <button className="btn btn-ghost btn-sm btn-circle text-error" onClick={handleDelete} disabled={deletePost.isPending}>
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              )}
            </div>

            {/* 类型标签 */}
            <div className="flex items-center flex-wrap gap-1.5 mb-4">
              <span className="text-xs font-medium px-2.5 py-1 rounded-full bg-base-200 text-base-content/60">
                {post.creation.creation_type === "article" ? "文章" : post.creation.creation_type === "topic" ? t("the_topic") : post.creation.creation_type === "question" ? "问答" : t("the_post")}
              </span>
              {post.creation.tags?.map((tag) => (
                <Link key={tag.id} href={`/posts?tag_id=${tag.id}`}
                  className="text-xs px-2.5 py-1 rounded-full"
                  style={{ backgroundColor: tag.color + "18", color: tag.color, borderColor: tag.color + "30" }}
                >
                  <Tag className="w-2.5 h-2.5 inline mr-0.5" />{tag.name}
                </Link>
              ))}
            </div>

            {/* 标题 */}
            <h1 className="text-xl lg:text-2xl font-bold text-base-content leading-snug mb-4">
              {post.creation.title}
            </h1>

            {/* 图片矩阵 */}
            {post.creation.image_urls && post.creation.image_urls.length > 0 && (
              <div className="mb-5 -mx-5 lg:-mx-8">
                <ImageGrid images={post.creation.image_urls} />
              </div>
            )}

            {/* 正文内容 */}
            <div
              className="prose-content text-base-content/80 leading-relaxed text-[15px]"
              dangerouslySetInnerHTML={{ __html: post.creation.content }}
            />

            {/* 底部操作栏 */}
            <div className="flex items-center justify-center gap-6 mt-8 pt-5 border-t border-base-200">
              <button
                className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-all ${
                  liked
                    ? "bg-error/10 text-error"
                    : "bg-base-200 text-base-content/60 hover:bg-base-300"
                }`}
                onClick={handleLikeClick}
                disabled={likePost.isPending || unlikePost.isPending}
              >
                {liked ? <Heart className="w-4 h-4 fill-current" /> : <Heart className="w-4 h-4" />}
                {post.creation.like_count}
              </button>
              <button className="flex items-center gap-2 px-4 py-2 rounded-full text-sm bg-base-200 text-base-content/60 hover:bg-base-300 transition-all" onClick={handleShare}>
                <Share2 className="w-4 h-4" /> 分享
              </button>
            </div>
          </div>
        </div>
      </article>
    </>
  );
}

/** 图片矩阵展示 */
function ImageGrid({ images }: { images: string[] }) {
  const count = images.length;
  if (count === 0) return null;
  const cols = count === 1 ? 1 : count === 2 ? 2 : count === 4 ? 2 : 3;

  return (
    <div
      className="grid gap-1 rounded-xl overflow-hidden"
      style={{ gridTemplateColumns: `repeat(${cols}, 1fr)` }}
    >
      {images.map((url, idx) => (
        <div key={idx} className="relative aspect-square cursor-pointer">
          <Image
            src={normalizeUrl(url)}
            alt={`图片 ${idx + 1}`}
            fill
            className="object-cover hover:scale-105 transition-transform"
            sizes={cols === 1 ? "100vw" : cols === 2 ? "50vw" : "33vw"}
          />
        </div>
      ))}
    </div>
  );
}

// 骨架屏组件（保持不变）
function PostDetailSkeleton() {
  return (
    <div className="space-y-4">
      <div className="skeleton h-10 w-3/4" />
      <div className="skeleton h-4 w-1/2" />
      <div className="skeleton h-64 w-full" />
    </div>
  );
}
