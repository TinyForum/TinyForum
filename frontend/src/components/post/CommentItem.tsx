"use client";

import { useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { Comment } from "@/lib/api/types";
import { timeAgo } from "@/lib/utils";
import { CornerDownRight, Trash2 } from "lucide-react";
import { useAuthStore } from "@/store/auth";
import { commentApi } from "@/lib/api";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Avatar from "../user/Avatar";
import { useTranslations } from "next-intl";

interface CommentItemProps {
  comment: Comment;
  postId: number;
  onReply?: (parentId: number, username: string) => void;
}

export default function CommentItem({
  comment,
  postId,
  onReply,
}: CommentItemProps) {
  const { user, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [showReplies, setShowReplies] = useState(true);
  const t = useTranslations("Comment");

  const deleteMutation = useMutation({
    mutationFn: () => commentApi.delete(comment.id),
    onSuccess: () => {
      toast.success(t("deleted"));
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    },
    onError: () => toast.error(t("delete_failed")),
  });

  const canDelete = user?.id === comment.author_id || user?.role === "admin";

  return (
    <div className="flex gap-3">
      <Link href={`/users/${comment.author_id}`} className="flex-none">
        <div className="avatar">
          <div className="w-8 h-8 rounded-full">
            <Avatar
              username={comment.author?.username}
              avatarUrl={comment.author?.avatar} // 数据库中的头像
              size="md"
            />
          </div>
        </div>
      </Link>

      <div className="flex-1">
        <div className="bg-base-200 rounded-xl p-3">
          <div className="flex items-center justify-between mb-1">
            <Link
              href={`/users/${comment.author?.id}`}
              className="text-sm font-semibold hover:text-primary transition-colors"
            >
              {comment.author?.username}
            </Link>
            <span className="text-xs text-base-content/40">
              {timeAgo(comment.created_at)}
            </span>
          </div>
          <p className="text-sm text-base-content/80 whitespace-pre-wrap">
            {comment.content}
          </p>
        </div>

        <div className="flex items-center gap-3 mt-1.5 px-1">
          {isAuthenticated && onReply && (
            <button
              className="text-xs text-base-content/50 hover:text-primary transition-colors flex items-center gap-1"
             onClick={() => {
  if (comment.author?.username) {
    onReply(comment.id, comment.author.username);
  }
}}

            >
              <CornerDownRight className="w-3 h-3" /> {t("reply")}
            </button>
          )}
          {comment.replies && comment.replies.length > 0 && (
            <button
              className="text-xs text-base-content/50 hover:text-primary transition-colors"
              onClick={() => setShowReplies(!showReplies)}
            >
              {showReplies ? t("collapse") : `${t("expand")+comment.replies.length+t("replies")}`}
            </button>
          )}
          {canDelete && (
            <button
              className="text-xs text-error/60 hover:text-error transition-colors flex items-center gap-1 ml-auto"
              onClick={() => deleteMutation.mutate()}
              disabled={deleteMutation.isPending}
            >
              <Trash2 className="w-3 h-3" /> {t("delete")}
            </button>
          )}
        </div>

        {/* Nested replies */}
        {showReplies && comment.replies && comment.replies.length > 0 && (
          <div className="mt-3 space-y-3 ml-2 border-l-2 border-base-300 pl-3">
            {comment.replies.map((reply) => (
              <CommentItem
                key={reply.id}
                comment={reply}
                postId={postId}
                onReply={onReply}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
