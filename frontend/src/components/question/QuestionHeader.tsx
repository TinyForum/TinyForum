// components/question/QuestionHeader.tsx
"use client";

import Link from "next/link";
import {
  EyeIcon,
  ChatBubbleLeftRightIcon,
  CheckBadgeIcon,
  CalendarIcon,
  HeartIcon,
  ShareIcon,
  FlagIcon,
  TagIcon,
  CurrencyDollarIcon,
  SparklesIcon,
} from "@heroicons/react/24/outline";
import { HeartIcon as HeartSolidIcon } from "@heroicons/react/24/solid";
import type { Post, Tag } from "@/lib/api/types";
import { formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";
import { toast } from "react-hot-toast";

interface QuestionHeaderProps {
  question: Post;
  answersCount: number;
  liked: boolean;
  likesCount: number;
  hasAccepted: boolean;
  rewardScore?: number;
  onLike: () => void;
  onShare?: () => void;
  onReport?: () => void;
}

export function QuestionHeader({
  question,
  answersCount,
  liked,
  likesCount,
  hasAccepted,
  rewardScore,
  onLike,
  onShare,
  onReport,
}: QuestionHeaderProps) {
  const timeAgo = formatDistanceToNow(new Date(question.created_at), {
    addSuffix: true,
    locale: zhCN,
  });

  const handleShare = async () => {
    if (onShare) {
      onShare();
      return;
    }

    // 默认分享行为
    const url = window.location.href;
    try {
      await navigator.clipboard.writeText(url);
      toast.success("链接已复制到剪贴板");
    } catch {
      // 捕获错误但不使用 err 变量
      toast.error("复制失败，请手动复制");
    }
  };

  const handleReport = () => {
    if (onReport) {
      onReport();
      return;
    }

    // 默认举报行为
    toast(
      (t) => (
        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium">确定要举报这个问题吗？</p>
          <div className="flex gap-2 justify-end">
            <button
              className="btn btn-xs btn-ghost"
              onClick={() => toast.dismiss(t.id)}
            >
              取消
            </button>
            <button
              className="btn btn-xs btn-error"
              onClick={() => {
                // 这里调用举报 API
                toast.success("已提交举报，我们会尽快处理");
                toast.dismiss(t.id);
              }}
            >
              确认举报
            </button>
          </div>
        </div>
      ),
      { duration: 5000 },
    );
  };

  return (
    <>
      {rewardScore !== undefined && rewardScore > 0 && !hasAccepted && (
        <div className="bg-gradient-to-r from-primary to-secondary px-4 py-2.5">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <CurrencyDollarIcon className="w-4 h-4 text-primary-content" />
              <span className="text-sm font-medium text-primary-content">
                悬赏 {rewardScore} 积分
              </span>
            </div>
            <SparklesIcon className="w-3.5 h-3.5 text-primary-content/80 animate-pulse" />
          </div>
        </div>
      )}

      <div className="card bg-base-100 shadow-lg border border-base-200 mb-6 overflow-hidden hover:shadow-xl transition-shadow duration-300">
        <div className="card-body p-6">
          {/* 标题 */}
          <h1 className="text-xl md:text-2xl lg:text-3xl font-bold text-base-content leading-tight">
            {question.title}
          </h1>

          {/* 元信息区域 */}
          <div className="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-base-content/60">
            {/* 作者 - 使用 primary 主题色 */}
            <div className="flex items-center gap-1.5">
              <div className="avatar placeholder">
                <div className="w-6 h-6 rounded-full bg-primary/10 text-primary ring-1 ring-primary/20">
                  <span className="text-xs font-medium">
                    {question.author?.username?.[0]?.toUpperCase() || "U"}
                  </span>
                </div>
              </div>
              <Link
                href={`/users/${question.author_id}`}
                className="hover:text-primary transition-colors duration-200 font-medium"
              >
                {question.author?.username || `用户${question.author_id}`}
              </Link>
            </div>

            {/* 时间 */}
            <div className="flex items-center gap-1">
              <CalendarIcon className="w-3.5 h-3.5" />
              <span>{timeAgo}</span>
            </div>

            {/* 浏览 */}
            <div className="flex items-center gap-1">
              <EyeIcon className="w-3.5 h-3.5" />
              <span>{question.view_count || 0} 浏览</span>
            </div>

            {/* 回答数 */}
            <div className="flex items-center gap-1">
              <ChatBubbleLeftRightIcon className="w-3.5 h-3.5" />
              <span>{answersCount} 回答</span>
            </div>

            {/* 点赞按钮 - 使用红色主题 */}
            <button
              onClick={onLike}
              className={`flex items-center gap-1.5 px-2 py-0.5 rounded-full transition-all duration-200 ${
                liked
                  ? "text-primary bg-primary/10"
                  : "hover:text-primary hover:bg-primary/5"
              }`}
            >
              {liked ? (
                <HeartSolidIcon className="w-4 h-4" />
              ) : (
                <HeartIcon className="w-4 h-4" />
              )}
              <span className="text-sm font-medium">{likesCount}</span>
            </button>

            {/* 状态标签 - 使用主题色 */}
            <div className="flex gap-2">
              {hasAccepted && (
                <div className="badge badge-success gap-1 badge-sm">
                  <CheckBadgeIcon className="w-3 h-3" />
                  已解决
                </div>
              )}
              {!hasAccepted && rewardScore && rewardScore > 0 && (
                <div className="badge badge-warning gap-1 badge-sm">
                  <CurrencyDollarIcon className="w-3 h-3" />
                  待采纳
                </div>
              )}
              {!hasAccepted && (!rewardScore || rewardScore === 0) && (
                <div className="badge badge-ghost gap-1 badge-sm">待解决</div>
              )}
            </div>
          </div>

          {/* 标签区域 - 使用主题色边框 */}
          {question.tags && question.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-1">
              {question.tags.map((tag: Tag) => (
                <Link
                  key={tag.id}
                  href={`/questions?tag_id=${tag.id}`}
                  className="badge badge-ghost hover:badge-primary gap-1 transition-all duration-200 cursor-pointer group"
                  style={{
                    borderLeftColor:
                      tag.color || "var(--fallback-p,oklch(var(--p)/1))",
                    borderLeftWidth: "3px",
                  }}
                >
                  <TagIcon className="w-3 h-3 group-hover:scale-110 transition-transform" />
                  {tag.name}
                </Link>
              ))}
            </div>
          )}

          {/* 内容区域 - 优化富文本显示 */}
          <div className="prose prose-sm sm:prose-base max-w-none pt-3">
            <div
              className="text-base-content/80 leading-relaxed [&_p]:mb-3 [&_ul]:list-disc [&_ul]:pl-5 [&_ol]:list-decimal [&_ol]:pl-5 [&_h1]:text-2xl [&_h1]:font-bold [&_h2]:text-xl [&_h2]:font-bold [&_h3]:text-lg [&_h3]:font-bold [&_pre]:bg-base-200 [&_pre]:p-3 [&_pre]:rounded-lg [&_code]:text-primary [&_code]:bg-primary/5 [&_code]:px-1 [&_code]:rounded"
              dangerouslySetInnerHTML={{ __html: question.content }}
            />
          </div>
        </div>

        {/* 操作按钮区域 */}
        <div className="card-actions justify-start px-6 pb-6 pt-0 gap-2 border-t border-base-200">
          <button
            onClick={handleShare}
            className="btn btn-ghost btn-sm gap-1.5 hover:bg-primary/5 transition-all duration-200"
          >
            <ShareIcon className="w-4 h-4" />
            分享
          </button>
          <button
            onClick={handleReport}
            className="btn btn-ghost btn-sm gap-1.5 hover:bg-error/10 hover:text-error transition-all duration-200"
          >
            <FlagIcon className="w-4 h-4" />
            举报
          </button>
        </div>
      </div>
    </>
  );
}
