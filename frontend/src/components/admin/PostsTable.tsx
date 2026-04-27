// import { Post } from "@/lib/api/types";

import { Post } from "@/lib/api";
import { formatDate } from "@/lib/utils";
// import { Post } from "@/type/admin.types";
// import { formatDate } from "date-fns";
import { Pin, PinOff } from "lucide-react";

// 帖子表格组件
export function PostsTable({
  posts,
  onTogglePin,
  isToggling,
  t,
}: {
  posts: Post[];
  onTogglePin: (id: number) => void;
  isToggling: boolean;
  t: (key: string) => string;
}) {
  if (posts.length === 0) {
    return (
      <div className="text-center py-8 text-base-content/60">
        {t("no_data")}
      </div>
    );
  }

  const getTypeBadge = (type: string) => {
    if (type === "article") return "badge-secondary";
    if (type === "topic") return "badge-accent";
    return "badge-ghost";
  };

  const getTypeText = (type: string) => {
    if (type === "article") return t("article");
    if (type === "topic") return t("topic");
    return t("post");
  };

  const getStatusBadge = (status: string) => {
    return status === "published" ? "badge-success" : "badge-warning";
  };

  const getStatusText = (status: string) => {
    if (status === "published") return t("published");
    if (status === "draft") return t("draft");
    return t("hidden");
  };

  return (
    <div className="overflow-x-auto">
      <table className="table table-zebra">
        <thead>
          <tr>
            <th>{t("article_title")}</th>
            <th>{t("author")}</th>
            <th>{t("type")}</th>
            <th>{t("status")}</th>
            <th>{t("views_likes")}</th>
            <th>{t("publish_time")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {posts.map((post) => (
            <tr key={post.id}>
              <td className="max-w-xs">
                <a
                  href={`/posts/${post.id}`}
                  target="_blank"
                  rel="noreferrer"
                  className="text-sm hover:text-primary transition-colors line-clamp-1"
                >
                  {post.pin_top && (
                    <Pin className="w-3 h-3 inline mr-1 text-primary" />
                  )}
                  {post.title}
                </a>
              </td>
              <td className="text-sm text-base-content/60">
                {post.author?.username}
              </td>
              <td>
                <span className={`badge badge-sm ${getTypeBadge(post.type)}`}>
                  {getTypeText(post.type)}
                </span>
              </td>
              <td>
                <span
                  className={`badge badge-sm ${getStatusBadge(post.status)}`}
                >
                  {getStatusText(post.status)}
                </span>
              </td>
              <td className="text-xs text-base-content/50">
                {post.view_count} / {post.like_count}
              </td>
              <td className="text-xs text-base-content/50">
                {formatDate(post.created_at)}
              </td>
              <td>
                <button
                  className={`btn btn-xs gap-1 ${post.pin_top ? "btn-ghost" : "btn-primary btn-outline"}`}
                  onClick={() => onTogglePin(post.id)}
                  disabled={isToggling}
                >
                  {post.pin_top ? (
                    <>
                      <PinOff className="w-3 h-3" /> {t("cancel_pin")}
                    </>
                  ) : (
                    <>
                      <Pin className="w-3 h-3" /> {t("pin")}
                    </>
                  )}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
