// components/question/AnswerCard.tsx
"use client";

import Link from "next/link";
import { ArrowUpIcon, ArrowDownIcon } from "lucide-react";
import {
  CalendarIcon,
  CheckBadgeIcon,
  ChatBubbleLeftRightIcon,
  ShareIcon,
  FlagIcon,
  TrophyIcon,
} from "@heroicons/react/24/outline";
import type { Comment } from "@/shared/api/types";
import { formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";
import { toast } from "react-hot-toast";
import { useEffect, useState } from "react";
import CommentSection from "../post/CommentSection";
import { useAnswerVote } from "@/features/answer/hooks/useAnswerVote";
// import CommentSection from '../CommentSection'; // 导入评论组件

interface AnswerCardProps {
  answer: Comment;
  isAccepted: boolean;
  isAuthor: boolean;
  canAccept: boolean;
  onAccept: () => void;
  currentUserId?: number;
  rewardScore?: number;
  answerNumber?: number;
}

export function AnswerCard({
  answer,
  isAccepted,
  canAccept,
  onAccept,
  currentUserId,
  rewardScore = 0,
  answerNumber,
}: AnswerCardProps) {
  const [showComments, setShowComments] = useState(false); // 控制评论区域显示

  const {
    userVote, // 'up' | 'down' | ''
    voteCount, // 总投票数
    loading, // 加载状态
    handleVote, // 投票/取消投票函数
  } = useAnswerVote(answer.id, currentUserId);

  // 调试：打印投票状态变化
  useEffect(() => {
    console.log(`Answer ${answer.id} - Vote Status:`, {
      userVote,
      voteCount,
      loading,
    });
  }, [answer.id, userVote, voteCount, loading]);

  const handleUpVote = async () => {
    if (loading) return;
    await handleVote("up");
  };

  const handleDownVote = async () => {
    if (loading) return;
    await handleVote("down");
  };

  const timeAgo = formatDistanceToNow(new Date(answer.created_at), {
    addSuffix: true,
    locale: zhCN,
  });

  const handleShare = () => {
    navigator.clipboard.writeText(
      `${window.location.origin}/answers/${answer.id}`,
    );
    toast.success("链接已复制到剪贴板");
  };

  const handleReport = () => {
    toast(
      (t) => (
        <div className="flex flex-col gap-2">
          <p className="text-sm">确认举报这个回答？</p>
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
                toast.dismiss(t.id);
                toast.success("举报已提交，我们会尽快处理");
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

  // 切换评论区域显示
  const toggleComments = () => {
    setShowComments(!showComments);
    if (!showComments) {
      // 延迟滚动，确保评论区域已渲染
      setTimeout(() => {
        const commentSection = document.getElementById(`comments-${answer.id}`);
        commentSection?.scrollIntoView({ behavior: "smooth", block: "start" });
      }, 100);
    }
  };

  useEffect(() => {
    console.log(
      `Answer ${answer.id} - voteCount: ${voteCount}, userVote: ${userVote}`,
    );
  }, [answer.id, voteCount, userVote]);

  return (
    <div
      className={`card bg-base-100 shadow-md border transition-all duration-300 ${
        isAccepted
          ? "border-success shadow-lg bg-gradient-to-r from-success/5 to-transparent"
          : "border-base-200 hover:shadow-lg"
      }`}
      id={`answer-${answer.id}`}
    >
      <div className="card-body p-5">
        <div className="flex gap-4">
          {/* 投票区域 */}
          <div className="flex flex-col items-center gap-1">
            <button
              onClick={handleUpVote}
              disabled={loading}
              className={`btn btn-sm btn-ghost p-1 min-h-0 h-auto ${
                userVote === "up" ? "text-primary" : "text-base-content/40"
              } hover:text-primary transition-colors disabled:opacity-50`}
              aria-label="赞同"
            >
              <ArrowUpIcon className="w-5 h-5" />
            </button>
            <span
              className={`text-sm font-semibold ${
                voteCount > 0
                  ? "text-primary"
                  : voteCount < 0
                    ? "text-error"
                    : "text-base-content/60"
              }`}
            >
              {voteCount}
            </span>
            <button
              onClick={handleDownVote}
              disabled={loading}
              className={`btn btn-sm btn-ghost p-1 min-h-0 h-auto ${
                userVote === "down" ? "text-error" : "text-base-content/40"
              } hover:text-error transition-colors disabled:opacity-50`}
              aria-label="反对"
            >
              <ArrowDownIcon className="w-5 h-5" />
            </button>
          </div>

          <div className="flex-1 min-w-0">
            {/* 头部信息 */}
            <div className="flex flex-wrap items-center justify-between gap-3 mb-3">
              <div className="flex flex-wrap items-center gap-3 text-sm">
                {/* 回答编号 */}
                {answerNumber && (
                  <span className="text-xs font-mono text-base-content/40">
                    #{answerNumber}
                  </span>
                )}

                {/* 作者信息 */}
                <div className="flex items-center gap-2">
                  <div className="avatar placeholder">
                    <div className="w-6 h-6 rounded-full bg-primary/10 text-primary">
                      <span className="text-xs">
                        {answer.author?.username?.[0]?.toUpperCase() || "U"}
                      </span>
                    </div>
                  </div>
                  <Link
                    href={`/users/${answer.author_id}`}
                    className="font-medium hover:text-primary transition-colors"
                  >
                    {answer.author?.username || `用户${answer.author_id}`}
                  </Link>
                </div>

                {/* 时间 */}
                <div className="flex items-center gap-1 text-base-content/60">
                  <CalendarIcon className="w-3.5 h-3.5" />
                  <span>{timeAgo}</span>
                </div>
              </div>

              {/* 状态标签 */}
              <div className="flex gap-2">
                {isAccepted && (
                  <div className="badge badge-success gap-1">
                    <TrophyIcon className="w-3 h-3" />
                    已采纳
                  </div>
                )}
                {rewardScore > 0 && isAccepted && (
                  <div className="badge badge-warning gap-1">
                    💰 +{rewardScore}
                  </div>
                )}
              </div>
            </div>

            {/* 回答内容 */}
            <div className="prose prose-sm max-w-none mb-4">
              <div
                className="text-base-content/80 leading-relaxed break-words"
                dangerouslySetInnerHTML={{ __html: answer.content }}
              />
            </div>

            {/* 底部操作按钮 */}
            <div className="flex flex-wrap gap-2 pt-2 border-t border-base-200">
              {/* 采纳按钮 */}
              {canAccept && !isAccepted && (
                <button
                  onClick={onAccept}
                  className="btn btn-xs btn-success btn-outline gap-1"
                >
                  <CheckBadgeIcon className="w-3.5 h-3.5" />
                  采纳为答案
                </button>
              )}

              {/* 分享按钮 */}
              <button
                onClick={handleShare}
                className="btn btn-xs btn-ghost gap-1"
              >
                <ShareIcon className="w-3.5 h-3.5" />
                分享
              </button>

              {/* 举报按钮 */}
              <button
                onClick={handleReport}
                className="btn btn-xs btn-ghost gap-1 hover:text-error"
              >
                <FlagIcon className="w-3.5 h-3.5" />
                举报
              </button>

              {/* 评论按钮 - 添加显示/隐藏评论区域功能 */}
              <button
                onClick={toggleComments}
                className="btn btn-xs btn-ghost gap-1"
              >
                <ChatBubbleLeftRightIcon className="w-3.5 h-3.5" />
                {/* 评论 {answer.replies?.length > 0 && `(${answer.comment_count})`} */}
                评论{" "}
                {answer.replies &&
                  answer.replies.length > 0 &&
                  `(${answer.replies.length})`}
              </button>
            </div>

            {/* 评论区域 - 可展开/收起 */}
            {showComments && (
              <div
                className="mt-4 pt-4 border-t border-base-200"
                id={`comments-${answer.id}`}
              >
                <CommentSection postId={answer.id} />
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
