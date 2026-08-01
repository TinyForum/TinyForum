// components/user/MyCommentsTable.tsx
"use client";

import { CommentResponse } from "@/shared/api/types/comment.model";
import { useTranslations } from "next-intl";
import Link from "next/link";

interface MyCommentsTableProps {
  comments: CommentResponse[];
}

export function MyCommentsTable({ comments }: MyCommentsTableProps) {
  // const comments = []; // TODO: 从 props 传入

  const t = useTranslations("User");
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
          {comments.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center text-base-content/60 py-8">
                {t("no_comments")}
              </td>
            </tr>
          ) : (
            comments.map((comment: CommentResponse) => (
              <tr key={comment.id}>
                <td className="max-w-md truncate">{comment.reply.content}</td>
                <td>
                  <Link
                    href={`/post/${comment.works_id}`}
                    className="hover:link-hover"
                  >
                    {comment.reply.content}
                  </Link>
                </td>
                <td>{comment.reply.like_count}</td>
                <td>{new Date(comment.created_at).toLocaleDateString()}</td>
                <td>
                  <button className="btn btn-xs btn-error">
                    {t("delete")}
                  </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
