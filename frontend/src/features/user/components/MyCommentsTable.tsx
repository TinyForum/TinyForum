// components/user/MyCommentsTable.tsx
"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";
import { toast } from "react-hot-toast";
import { useMeComments } from "@/features/user/hooks/useMeComments";
import { commentApi } from "@/shared/api/modules/comments";

export function MyCommentsTable() {
  const t = useTranslations("User");
  const {
    list: comments,
    isLoading,
    error,
    has_more: hasMore,
    loadMore,
    refresh,
  } = useMeComments();

  // 删除评论：确认后调用删除接口，成功后刷新列表
  const handleDelete = async (commentId: number) => {
    if (!confirm(t("delete_comment_confirm"))) return;
    try {
      await commentApi.delete(commentId);
      toast.success(t("delete_comment_success"));
      refresh();
    } catch {
      toast.error(t("delete_comment_failed"));
    }
  };

  if (isLoading && comments.length === 0) {
    return (
      <div className="overflow-x-auto rounded-xl border border-base-200 bg-base-100 shadow-sm">
        <table className="table">
          <tbody>
            {[...Array(3)].map((_, idx) => (
              <tr key={idx} className="animate-pulse">
                <td colSpan={5}>
                  <div className="skeleton h-12 w-full rounded-lg"></div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 rounded-2xl border border-error/20 bg-error/5 p-8 text-center shadow-sm">
        <p className="text-base-content/80">{error}</p>
        <button className="btn btn-primary btn-sm" onClick={refresh}>
          {t("retry")}
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-5">
      <div className="overflow-x-auto rounded-2xl border border-base-200 bg-base-100 shadow-md">
        <table className="table table-zebra table-md">
          <thead className="bg-base-200/60">
            <tr className="text-sm uppercase tracking-wide text-base-content/80">
              <th>{t("content")}</th>
              <th>{t("posts")}</th>
              <th className="text-center">{t("likes")}</th>
              <th>{t("created_at")}</th>
              <th className="text-center">{t("actions")}</th>
            </tr>
          </thead>
          <tbody>
            {comments.length === 0 ? (
              <tr>
                <td colSpan={5} className="py-12 text-center text-base-content/50">
                  {t("no_comments")}
                </td>
              </tr>
            ) : (
              comments.map((comment) => (
                <tr key={comment.id} className="hover:bg-base-200/30">
                  <td className="max-w-md truncate">{comment.content}</td>
                  <td>
                    <Link
                      href={`/posts/${comment.works_id}`}
                      className="hover:text-primary hover:underline"
                    >
                      {t("view_post")}
                    </Link>
                  </td>
                  <td className="text-center font-mono text-sm">
                    {comment.like_count}
                  </td>
                  <td className="text-sm text-base-content/70">
                    {new Date(comment.created_at).toLocaleDateString()}
                  </td>
                  <td className="text-center">
                    <button
                      className="btn btn-xs btn-error"
                      onClick={() => handleDelete(comment.id)}
                    >
                      {t("delete")}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {hasMore && (
        <div className="flex justify-center pt-2">
          <button
            className="btn btn-outline btn-primary rounded-full px-8"
            onClick={loadMore}
            disabled={isLoading}
          >
            {isLoading ? t("loading") : t("load_more")}
          </button>
        </div>
      )}
    </div>
  );
}
