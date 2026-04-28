// components/moderator/PendingPostsTable.tsx

export interface PendingPost {
  id: number;
  title: string;
  author_name: string;
  author_id: number;
  created_at: string;
  is_pinned: boolean;
  status: string;
  type: string;
  content?: string;
  summary?: string;
  cover?: string;
}

interface ModeratorPermissions {
  canDeletePost: boolean;
  canPinPost: boolean;
  canEditAnyPost: boolean;
  canManageModerator: boolean;
  canBanUser: boolean;
}

interface PendingPostsTableProps {
  posts: PendingPost[];
  onDelete?: (postId: number) => void;
  onPin?: (postId: number, pin: boolean) => void;
  isDeleting: boolean;
  isPinning: boolean;
  permissions: ModeratorPermissions;
  t: (key: string) => string;
}

export function PendingPostsTable({
  posts,
  onDelete,
  onPin,
  isDeleting,
  isPinning,
  permissions,
  t,
}: PendingPostsTableProps) {
  if (!posts?.length) {
    return (
      <div className="text-center py-8 text-base-content/50">
        {t("no_posts")}
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>{t("title")}</th>
            <th>{t("author")}</th>
            <th>{t("created_at")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {posts.map((post: PendingPost) => (
            <tr key={post.id}>
              <td>{post.title}</td>
              <td>{post.author_name}</td>
              <td>{new Date(post.created_at).toLocaleDateString()}</td>
              <td className="flex gap-2">
                {permissions.canPinPost && onPin && (
                  <button
                    className="btn btn-xs btn-ghost"
                    onClick={() => onPin(post.id, !post.is_pinned)}
                    disabled={isPinning}
                  >
                    {post.is_pinned ? t("unpin") : t("pin")}
                  </button>
                )}
                {permissions.canDeletePost && onDelete && (
                  <button
                    className="btn btn-xs btn-error"
                    onClick={() => onDelete(post.id)}
                    disabled={isDeleting}
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