// components/member/MemberCommentsTable.tsx
"use client";

import { Comment } from "@/shared/api";
import { useTranslations } from "next-intl";
import Link from "next/link";

interface MemberCommentsTableProps {
  comments?: Comment[];
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
          {comments.map((comment: Comment) => (
            <tr key={comment.id}>
              <td className="max-w-md truncate">{comment.content}</td>
              <td>
                <Link
                  href={`/post/${comment.post_id}`}
                  className="hover:link-hover"
                >
                  {comment.parent?.post_id}
                </Link>
              </td>
              <td>{comment.like_count}</td>
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
