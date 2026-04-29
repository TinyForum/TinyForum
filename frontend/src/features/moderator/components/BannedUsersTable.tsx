// components/moderator/BannedUsersTable.tsx

import { BannedUser } from "@/shared/api/modules/moderator";

// export interface BannedUser {
//   id: number;
//   user_id: number;
//   username: string;
//   reason: string;
//   created_at: string;
//   expires_at: string | null;
//   board_id: number;
//   banned_by: number;
// }

interface BannedUsersTableProps {
  users: BannedUser[];
  onUnban?: (userId: number) => void;
  isUnbanning: boolean;
  t: (key: string) => string;
}

export function BannedUsersTable({
  users,
  onUnban,
  isUnbanning,
  t,
}: BannedUsersTableProps) {
  if (!users?.length) {
    return (
      <div className="text-center py-8 text-base-content/50">
        {t("no_banned_users")}
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>{t("user")}</th>
            <th>{t("reason")}</th>
            <th>{t("banned_at")}</th>
            <th>{t("expires_at")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user: BannedUser) => (
            <tr key={user.id}>
              <td>
                <div className="flex items-center gap-2">
                  <div className="avatar placeholder">
                    <div className="w-8 h-8 rounded-full bg-error/10 text-error">
                      <span className="text-xs font-medium">
                        {user.username?.[0]?.toUpperCase() || "U"}
                      </span>
                    </div>
                  </div>
                  <span className="font-medium">{user.username}</span>
                </div>
              </td>
              <td className="max-w-xs">
                <p className="line-clamp-2 text-sm">
                  {user.reason || t("no_reason")}
                </p>
              </td>
              <td>{new Date(user.created_at).toLocaleDateString()}</td>
              <td>
                {user.expires_at
                  ? new Date(user.expires_at).toLocaleDateString()
                  : t("permanent")}
              </td>
              <td>
                {onUnban && (
                  <button
                    className="btn btn-xs btn-success gap-1"
                    onClick={() => onUnban(user.user_id)}
                    disabled={isUnbanning}
                  >
                    {isUnbanning ? (
                      <span className="loading loading-spinner loading-xs" />
                    ) : null}
                    {t("unban")}
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
