// components/member/MemberCommentsTable.tsx
"use client";

import { CommentResponse } from "@/shared/api/types/comment.model";
import { useTranslations } from "next-intl";
import Link from "next/link";

interface MemberCommentsTableProps {
  comments?: CommentResponse[];
  onDelete?: (id: number) => void;
}

export function MemberCommentsTable({
  comments = [],
  onDelete,
}: MemberCommentsTableProps) {
  const t = useTranslations("Member");
  if (comments.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_comments")}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>{t("content")}</th>
            <th>{t("posts")}</th>
            <th>{t("likes")}</th>
            <th>{t("created_at")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {comments.map((comment: CommentResponse) => (
            <tr key={comment.id}>
              <td className="max-w-md truncate">{comment.reply.content}</td>
              <td>
                <Link
                  href={`/post/${comment.works_id}`}
                  className="hover:link-hover"
                >
                  {comment.reply.parent_id}
                </Link>
              </td>
              <td>{comment.reply.like_count}</td>
              <td>{new Date(comment.created_at).toLocaleDateString()}</td>
              <td>
                {onDelete && (
                  <button
                    onClick={() => onDelete(comment.id)}
                    className="btn btn-xs btn-error"
                  >
                    {t("delete")}
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
