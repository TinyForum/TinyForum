// components/user/MyCommentsTable.tsx
"use client";

import { Comment } from "@/shared/api";
import { useTranslations } from "next-intl";
import Link from "next/link";

interface MyCommentsTableProps {
  comments: Comment[];
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
            comments.map((comment: Comment) => (
              <tr key={comment.id}>
                <td className="max-w-md truncate">{comment.content}</td>
                <td>
                  <Link
                    href={`/post/${comment.post_id}`}
                    className="hover:link-hover"
                  >
                    {comment.content}
                  </Link>
                </td>
                <td>{comment.like_count}</td>
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
