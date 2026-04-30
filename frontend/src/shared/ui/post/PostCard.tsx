"use client";

import Link from "next/link";
import Image from "next/image";
import { useQuery } from "@tanstack/react-query";
import { timeAgo, truncate } from "@/shared/lib/utils";
import { Eye, Heart, MessageSquare, Pin, Tag, HelpCircle } from "lucide-react";
import { Post, userApi, type User } from "@/shared/api";
import Avatar from "@/features/user/components/Avatar";

interface PostCardProps {
  post: Post;
  commentCount?: number;
}

export default function PostCard({ post, commentCount }: PostCardProps) {
  const isQuestion = post?.type === "question";

  const { data: fetchedAuthor } = useQuery({
    queryKey: ["user", post?.author_id],
    queryFn: () => userApi.getProfile(post.author_id).then((r) => r.data.data),
    enabled: !!(post && !post.author && post.author_id),
  });

  if (!post) return null;
  const author = (post.author || fetchedAuthor) as User | undefined;
  const rewardScore = post.question?.reward_score || 0;
  const answerCount = post.question?.answer_count || 0;
  const isAccepted = post.question?.accepted_answer_id != null;

  const getPostTypeLabel = () => {
    switch (post.type) {
      case "question":
        return "问答";
      case "article":
        return "文章";
      case "topic":
        return "话题";
      default:
        return "帖子";
    }
  };

  const getPostTypeColor = () => {
    switch (post.type) {
      case "question":
        return "bg-primary/10 text-primary";
      case "article":
        return "bg-secondary/10 text-secondary";
      case "topic":
        return "bg-accent/10 text-accent";
      default:
        return "bg-base-300 text-base-content";
    }
  };

  return (
    <div
      className={`bg-base-100 shadow-sm hover:shadow-md transition-all duration-200 border border-base-300 rounded-xl ${
        post.pin_top ? "border-l-4 border-l-primary" : ""
      }`}
    >
      <div className="p-4">
        {/* 头像和标题区域 */}
        <div className="flex items-start gap-3">
          <Link href={`/users/${post.author_id}`} className="flex-none">
            <Avatar
              username={author?.username || `用户${post.author_id}`}
              avatarUrl={author?.avatar}
              size="md"
            />
          </Link>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              {post.pin_top && (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full bg-primary/10 text-primary">
                  <Pin className="w-3 h-3" /> 置顶
                </span>
              )}
              <span
                className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full ${getPostTypeColor()}`}
              >
                {isQuestion && <HelpCircle className="w-3 h-3" />}
                {getPostTypeLabel()}
              </span>
              {isQuestion && rewardScore > 0 && (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full bg-warning/10 text-warning">
                  💰 {rewardScore} 积分悬赏
                </span>
              )}
              {isQuestion && isAccepted && (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full bg-success/10 text-success">
                  ✓ 已采纳
                </span>
              )}
            </div>

            <Link href={`/posts/${post.id}`} className="group">
              <h2 className="text-base font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-2 mt-1">
                {post.title}
              </h2>
            </Link>

            {isQuestion && answerCount > 0 && (
              <div className="flex items-center gap-3 mt-1 text-xs text-base-content/50">
                <span className="flex items-center gap-1">
                  <MessageSquare className="w-3 h-3" />
                  {answerCount} 个回答
                </span>
              </div>
            )}

            {post.summary && !isQuestion && (
              <p className="text-sm text-base-content/60 mt-1 line-clamp-2">
                {truncate(post.summary, 120)}
              </p>
            )}
          </div>

          {post.cover && (
            <Link
              href={`/posts/${post.id}`}
              className="flex-none hidden sm:block"
            >
              <div className="w-20 h-16 rounded-lg overflow-hidden relative">
                <Image
                  src={post.cover}
                  alt={post.title}
                  fill
                  className="object-cover"
                  sizes="80px"
                />
              </div>
            </Link>
          )}
        </div>

        {/* 底部信息 */}
        <div className="flex items-center justify-between mt-3 flex-wrap gap-2">
          <div className="flex items-center gap-3 text-xs text-base-content/50">
            <Link
              href={`/users/${post.author_id}`}
              className="hover:text-primary transition-colors font-medium"
            >
              {author?.username || `用户${post.author_id}`}
            </Link>
            <span>{timeAgo(post.created_at)}</span>
            {post.tags && post.tags.length > 0 && (
              <div className="flex items-center gap-1">
                <Tag className="w-3 h-3" />
                {post.tags.slice(0, 2).map((tag) => (
                  <Link
                    key={tag.id}
                    href={`/posts?tag_id=${tag.id}`}
                    className="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full"
                    style={{
                      backgroundColor: tag.color + "20",
                      color: tag.color,
                      borderColor: tag.color + "40",
                    }}
                  >
                    {tag.name}
                  </Link>
                ))}
              </div>
            )}
          </div>

          <div className="flex items-center gap-3 text-xs text-base-content/50">
            <span className="flex items-center gap-1">
              <Eye className="w-3.5 h-3.5" /> {post.view_count}
            </span>
            <span className="flex items-center gap-1">
              <Heart className="w-3.5 h-3.5" /> {post.like_count}
            </span>
            {commentCount !== undefined && (
              <span className="flex items-center gap-1">
                <MessageSquare className="w-3.5 h-3.5" /> {commentCount}
              </span>
            )}
            {isQuestion && answerCount > 0 && commentCount === undefined && (
              <span className="flex items-center gap-1">
                <MessageSquare className="w-3.5 h-3.5" /> {answerCount}
              </span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
