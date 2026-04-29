// components/user/MyPostsTable.tsx
"use client";

import { Post } from "@/shared/api";
import { useTranslations } from "next-intl";
import Link from "next/link";

interface MyPostsTableProps {
  posts: Post[];
}

export function MyPostsTable({ posts }: MyPostsTableProps) {
  // const posts = []; // TODO: 从 props 传入
  const t = useTranslations("User");
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
          {posts.length === 0 ? (
            <tr>
              <td colSpan={6} className="text-center text-base-content/60 py-8">
                {t("no_posts")}
              </td>
            </tr>
          ) : (
            posts.map((post: Post) => (
              <tr key={post.id}>
                <td>
                  <Link href={`/post/${post.id}`} className="hover:link-hover">
                    {post.title}
                  </Link>
                </td>
                <td>{post.board?.name}</td>
                <td>{post.like_count}</td>
                <td>{post.view_count}</td>
                <td>{new Date(post.created_at).toLocaleDateString()}</td>
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
