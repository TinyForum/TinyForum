// src/components/post/PostCard.tsx
"use client";

import Link from "next/link";
import Image from "next/image";
import { useQuery } from "@tanstack/react-query";
import { timeAgo, truncate } from "@/lib/utils";
import { Eye, Heart, MessageSquare, Pin, Tag, HelpCircle } from "lucide-react";
import Avatar from "../user/Avatar";
import { Post, userApi } from "@/lib/api";

interface PostCardProps {
  post: Post;
  commentCount?: number;
}

export default function PostCard({ post, commentCount }: PostCardProps) {
  if (!post) return null;

  // 判断是否为问答帖
  const isQuestion = post.type === "question";

  // 如果 post.author 不存在，单独获取用户信息
  const { data: fetchedAuthor } = useQuery({
    queryKey: ["user", post.author_id],
    queryFn: () => userApi.getProfile(post.author_id).then((r) => r.data.data),
    enabled: !post.author && !!post.author_id, // 只在没有 author 时获取
  });

  // 使用已有的 author 或获取到的 author
  const author = post.author || fetchedAuthor;

  // 从 post.question 获取问答相关信息
  const rewardScore = post.question?.reward_score || 0;
  const answerCount = post.question?.answer_count || 0;
  const isAccepted = post.question?.accepted_answer_id != null;

  // 获取帖子类型显示文本
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

  // 获取帖子类型样式
  const getPostTypeClass = () => {
    switch (post.type) {
      case "question":
        return "badge-primary";
      case "article":
        return "badge-secondary";
      case "topic":
        return "badge-accent";
      default:
        return "badge-ghost";
    }
  };

  return (
    <div
      className={`card bg-base-100 shadow-sm hover:shadow-md transition-all duration-200 hover:-translate-y-0.5 border border-base-300 ${post.pin_top ? "border-l-4 border-l-primary" : ""}`}
    >
      <div className="card-body p-4">
        {/* Top row */}
        <div className="flex items-start gap-3">
          {/* Author avatar */}
          <Link href={`/users/${post.author_id}`} className="flex-none">
            <div className="avatar">
              <div className="w-10 h-10 rounded-full">
                <Avatar
                  username={author?.username || `用户${post.author_id}`}
                  avatarUrl={author?.avatar}
                  size="md"
                />
              </div>
            </div>
          </Link>

          <div className="flex-1 min-w-0">
            {/* Title and badges */}
            <div className="flex items-center gap-2 flex-wrap">
              {post.pin_top && (
                <span className="badge badge-primary badge-sm gap-1">
                  <Pin className="w-3 h-3" /> 置顶
                </span>
              )}
              <span className={`badge ${getPostTypeClass()} badge-sm gap-1`}>
                {isQuestion && <HelpCircle className="w-3 h-3" />}
                {getPostTypeLabel()}
              </span>
              {isQuestion && rewardScore > 0 && (
                <span className="badge badge-warning badge-sm gap-1">
                  💰 {rewardScore} 积分悬赏
                </span>
              )}
              {isQuestion && isAccepted && (
                <span className="badge badge-success badge-sm gap-1">
                  ✓ 已采纳
                </span>
              )}
            </div>

            <Link href={`/posts/${post.id}`} className="group">
              <h2 className="text-base font-semibold text-base-content group-hover:text-primary transition-colors line-clamp-2 mt-1">
                {post.title}
              </h2>
            </Link>

            {/* 问答帖子显示额外信息 */}
            {isQuestion && answerCount > 0 && (
              <div className="flex items-center gap-3 mt-1 text-xs text-base-content/50">
                <span className="flex items-center gap-1">
                  <MessageSquare className="w-3 h-3" />
                  {answerCount} 个回答
                </span>
              </div>
            )}

            {/* 普通帖子显示摘要 */}
            {post.summary && !isQuestion && (
              <p className="text-sm text-base-content/60 mt-1 line-clamp-2">
                {truncate(post.summary, 120)}
              </p>
            )}
          </div>

          {/* Cover image */}
          {post.cover && (
            <Link
              href={`/posts/${post.id}`}
              className="flex-none hidden sm:block"
            >
              <div className="w-20 h-16 rounded-lg overflow-hidden">
                <Image
                  src={post.cover}
                  alt={post.title}
                  width={80}
                  height={64}
                  className="object-cover w-full h-full"
                />
              </div>
            </Link>
          )}
        </div>

        {/* Bottom row */}
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
                    className="badge badge-sm"
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
            {/* 如果是问答帖，显示回答数 */}
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
