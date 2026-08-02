// src/app/posts/[id]/PostDetailClient.tsx
"use client";

import { useAuthStore } from "@/store/auth";
import Image from "next/image";
import Link from "next/link";
import { timeAgo, formatDate } from "@/shared/lib/utils";
import {
  Eye,
  Heart,
  HeartOff,
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

      <article className="card bg-base-100 border border-base-300 shadow-sm mb-6">
        <div className="card-body p-6 lg:p-8">
          {/* 类型与标签 */}
          <div className="flex items-center flex-wrap gap-2 mb-3">
            <span
              className={`badge ${
                post.creation.creation_type === "article"
                  ? "badge-secondary"
                  : post.creation.creation_type === "topic"
                    ? "badge-accent"
                    : "badge-ghost"
              }`}
            >
              {post.creation.creation_type === "article"
                ? "文章"
                : post.creation.creation_type === "topic"
                  ? t("the_topic")
                  : t("the_post")}
            </span>
            {post.creation.tags?.map((tag) => (
              <Link
                key={tag.id}
                href={`/posts?tag_id=${tag.id}`}
                className="badge badge-sm gap-1"
                style={{
                  backgroundColor: tag.color + "20",
                  color: tag.color,
                  borderColor: tag.color + "40",
                }}
              >
                <Tag className="w-2.5 h-2.5" /> {tag.name}
              </Link>
            ))}
          </div>

          {/* 标题 */}
          <h1 className="text-2xl lg:text-3xl font-bold text-base-content leading-tight">
            {post.creation.title}
          </h1>

          {/* 作者信息 */}
          <div className="flex items-center gap-3 mt-4 pb-4 border-b border-base-300">
            <Link href={`/users/${post.creation.author?.id}`}>
              <div className="avatar">
                <div className="w-10 h-10 rounded-full">
                  <Avatar
                    username={post.creation.author?.username}
                    avatarUrl={post.creation.author?.avatar_url}
                    size="md"
                  />
                </div>
              </div>
            </Link>
            <div>
              <Link
                href={`/users/${post.creation.author?.id}`}
                className="font-medium hover:text-primary transition-colors text-sm"
              >
                {post.creation.author?.username}
              </Link>
              <div className="flex items-center gap-2 text-xs text-base-content/40">
                <Clock className="w-3 h-3" />
                <span title={formatDate(post.created_at)}>
                  {timeAgo(post.created_at)}
                </span>
                <span>·</span>
                <Eye className="w-3 h-3" />
                <span>
                  {post.creation.view_count} {t("read")}
                </span>
              </div>
            </div>

            {/* 操作按钮（编辑/删除） */}
            <div className="ml-auto flex items-center gap-1">
              {(isAuthor || isAdmin) && (
                <>
                  <Link
                    href={`/posts/${postId}/edit`}
                    className="btn btn-ghost btn-xs gap-1"
                  >
                    <Pencil className="w-3.5 h-3.5" /> {t("to_edit")}
                  </Link>
                  <button
                    className="btn btn-ghost btn-xs text-error gap-1"
                    onClick={handleDelete}
                    disabled={deletePost.isPending}
                  >
                    <Trash2 className="w-3.5 h-3.5" /> {t("delete")}
                  </button>
                </>
              )}
            </div>
          </div>

          {/* 封面图 */}
          {post.creation.cover_url && (
            <div className="my-4 rounded-xl overflow-hidden">
              <Image
                src={post.creation.cover_url}
                alt={post.creation.title}
                width={800}
                height={400}
                className="w-full object-cover max-h-72"
              />
            </div>
          )}

          {/* 正文内容 */}
          <div
            className="prose-content mt-4 text-base-content/80 leading-relaxed"
            dangerouslySetInnerHTML={{ __html: post.creation.content }}
          />

          {/* 底部操作栏 */}
          <div className="flex items-center gap-3 mt-8 pt-4 border-t border-base-300">
            <button
              className={`btn btn-sm gap-2 ${liked ? "btn-error" : "btn-ghost"}`}
              onClick={handleLikeClick}
              disabled={likePost.isPending || unlikePost.isPending}
            >
              {liked ? (
                <HeartOff className="w-4 h-4" />
              ) : (
                <Heart className="w-4 h-4" />
              )}
              {post.creation.like_count} {t("like")}
            </button>
            <button
              className="btn btn-ghost btn-sm gap-2"
              onClick={handleShare}
            >
              <Share2 className="w-4 h-4" /> {t("share")}
            </button>
          </div>
        </div>
      </article>
    </>
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
