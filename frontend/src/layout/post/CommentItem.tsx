"use client";
// 评论项
import { useState } from "react";
import Link from "next/link";
import { timeAgo } from "@/shared/lib/utils";
import { CornerDownRight, Trash2 } from "lucide-react";
import { useAuthStore } from "@/store/auth";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import Avatar from "@/shared/ui/common/Avatar";
import { commentApi } from "@/shared/api/modules/comments";
import { Reply } from "@/shared/api/types/comment.model";

interface CommentItemProps {
  reply: Reply;
  postId: number;
  onReply?: (parentId: number, username: string) => void;
}

export default function CommentItem({
  reply,
  postId,
  onReply,
}: CommentItemProps) {
  const { user, isAuthenticated } = useAuthStore();
  const queryClient = useQueryClient();
  const [showReplies, setShowReplies] = useState(true);
  const t = useTranslations("Comment");

  const deleteMutation = useMutation({
    mutationFn: () => commentApi.delete(reply.id),
    onSuccess: () => {
      toast.success(t("deleted"));
      queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    },
    onError: () => toast.error(t("delete_failed")),
  });

  const canDelete = user?.id === reply.author_id || user?.role === "admin";

  return (
    <div className="flex gap-3">
      <Link href={`/users/${reply.author_id}`} className="flex-none">
        <div className="avatar">
          <div className="w-8 h-8 rounded-full">
            <Avatar
              username={reply.author?.username}
              avatarUrl={reply.author?.avatar_url} // 数据库中的头像
              size="md"
            />
          </div>
        </div>
      </Link>

      <div className="flex-1">
        <div className="bg-base-200 rounded-xl p-3">
          <div className="flex items-center justify-between mb-1">
            <Link
              href={`/users/${reply.author?.id}`}
              className="text-sm font-semibold hover:text-primary transition-colors"
            >
              {reply.author?.username}
            </Link>
            <span className="text-xs text-base-content/40">
              {timeAgo(reply.created_at)}
            </span>
          </div>
          <p className="text-sm text-base-content/80 whitespace-pre-wrap">
            {reply.content}
          </p>
        </div>
        <div className="flex items-center gap-3 mt-1.5 px-1">
          {isAuthenticated && onReply && (
            <button
              className="text-xs text-base-content/50 hover:text-primary transition-colors flex items-center gap-1"
              onClick={() => {
                if (reply.author?.username) {
                  onReply(reply.id, reply.author.username);
                }
              }}
            >
              <CornerDownRight className="w-3 h-3" /> {t("reply")}
            </button>
          )}
          {reply.replies && reply.replies.length > 0 && (
            <button
              className="text-xs text-base-content/50 hover:text-primary transition-colors"
              onClick={() => setShowReplies(!showReplies)}
            >
              {showReplies
                ? t("collapse")
                : `${t("expand") + reply.replies.length + t("replies")}`}
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
        {showReplies && reply.replies && reply.replies.length > 0 && (
          <div className="mt-3 space-y-3 ml-2 border-l-2 border-base-300 pl-3">
            {reply.replies.map((reply) => (
              <CommentItem
                key={reply.id}
                reply={reply}
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
