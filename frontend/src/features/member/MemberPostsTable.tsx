// components/member/MemberPostsTable.tsx
"use client";

import { Post } from "@/shared/api";
import { useTranslations } from "next-intl";
import Link from "next/link";

interface MemberPostsTableProps {
  posts?: Post[];
  onDelete?: (id: number) => void;
}

export function MemberPostsTable({
  posts = [],
  onDelete,
}: MemberPostsTableProps) {
  const t = useTranslations("Member");
  if (posts.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_posts")}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>{t("title")}</th>
            <th>{t("board")}</th>
            <th>{t("likes")}</th>
            <th>{t("comments")}</th>
            <th>{t("created_at")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {posts.map((post: Post) => (
            <tr key={post.id}>
              <td>
                <Link href={`/post/${post.id}`} className="hover:link-hover">
                  {post.title}
                </Link>
              </td>
              <td>{post.title}</td>
              <td>{post.like_count}</td>
              <td>{post.view_count}</td>
              <td>{new Date(post.created_at).toLocaleDateString()}</td>
              <td>
                {onDelete && (
                  <button
                    onClick={() => onDelete(post.id)}
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
