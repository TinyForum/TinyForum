"use client";
// 整个评论区
import { useState, useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/store/auth";
import CommentItem from "./CommentItem";
import toast from "react-hot-toast";
import { Send, MessageSquare } from "lucide-react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import {
  useCommentsByPost,
  useCommentsTree,
  useCreateComment,
} from "@/features/comment/hooks/useComment";

interface CommentSectionProps {
  postId: number;
}

export default function CommentSection({ postId }: CommentSectionProps) {
  const { isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [content, setContent] = useState("");
  const [replyTo, setReplyTo] = useState<{
    id: number;
    username: string;
  } | null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const t = useTranslations("Comment");

  // 使用封装好的查询 Hook（固定分页参数，可后续扩展）
  const { data, isLoading } = useCommentsByPost(postId, {
    page: 1,
    page_size: 50,
  });
  const { data: commentsTree } = useCommentsTree(postId);

  // 使用封装好的创建评论 Mutation
  const createMutation = useCreateComment({
    onSuccess: (_, variables) => {
      toast.success(t("comment_success") || "评论成功");
      setContent("");
      setReplyTo(null);
      // 可选额外操作，失效已由 Hook 内部处理
    },
    onError: () => toast.error(t("comment_failed") || "评论失败"),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;
    createMutation.mutate({
      post_id: postId,
      content: content.trim(),
      parent_id: replyTo?.id,
    });
  };

  const handleReply = (parentId: number, username: string) => {
    setReplyTo({ id: parentId, username });
    textareaRef.current?.focus();
  };

  const comments = data?.list ?? [];
  const total = data?.total ?? 0;

  return (
    <div className="mt-8">
      <h3 className="text-lg font-bold flex items-center gap-2 mb-6">
        <MessageSquare className="w-5 h-5 text-primary" />
        {t("comments")}{" "}
        <span className="text-base-content/40 font-normal text-base">
          ({total})
        </span>
      </h3>
      {/* 评论输入框 */}
      {isAuthenticated ? (
        <form onSubmit={handleSubmit} className="mb-8">
          {replyTo && (
            <div className="flex items-center gap-2 mb-2 text-sm text-base-content/60 bg-base-200 px-3 py-1.5 rounded-lg">
              <span>
                {t("reply_to")}{" "}
                <strong className="text-primary">@{replyTo.username}</strong>
              </span>
              <button
                type="button"
                className="ml-auto text-error hover:text-error/80 text-xs"
                onClick={() => setReplyTo(null)}
              >
                {t("cancel")}
              </button>
            </div>
          )}
          <div className="flex gap-3">
            <textarea
              ref={textareaRef}
              className="textarea textarea-bordered flex-1 resize-none focus:outline-none focus:border-primary"
              rows={3}
              placeholder={
                replyTo
                  ? `${t("reply_to")} @${replyTo.username}...`
                  : t("write_comment")
              }
              value={content}
              onChange={(e) => setContent(e.target.value)}
              maxLength={2000}
            />
          </div>
          <div className="flex items-center justify-between mt-2">
            <span className="text-xs text-base-content/40">
              {content.length}/2000
            </span>
            <button
              type="submit"
              className="btn btn-primary btn-sm gap-1"
              disabled={!content.trim() || createMutation.isPending}
            >
              {createMutation.isPending ? (
                <span className="loading loading-spinner loading-xs" />
              ) : (
                <Send className="w-4 h-4" />
              )}
              {t("send")}
            </button>
          </div>
        </form>
      ) : (
        <div className="alert mb-6">
          <span>
            {t("please")}{" "}
            <Link href="/auth/login" className="link link-primary">
              {t("login")}
            </Link>{" "}
            {t("to_comment")}
          </span>
        </div>
      )}

      {/* 评论列表 */}
      {isLoading ? (
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner loading-md text-primary" />
        </div>
      ) : comments.length === 0 ? (
        <div className="text-center py-12 text-base-content/40">
          <MessageSquare className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p>{t("no_comments")}</p>
        </div>
      ) : (
        <div className="space-y-6">
          {commentsTree?.map((comment) => (
            <CommentItem
              key={comment.id}
              reply={comment.reply} // 注意：CommentItem 期望 Reply 类型
              postId={postId}
              onReply={handleReply}
            />
          ))}
        </div>
      )}
    </div>
  );
}
