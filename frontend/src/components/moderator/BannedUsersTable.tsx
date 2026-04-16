// components/moderator/BannedUsersTable.tsx
interface BannedUsersTableProps {
  users: any[];
  onUnban?: (userId: number) => void;
  isUnbanning: boolean;
  t: (key: string) => string;
}

export function BannedUsersTable({ users, onUnban, isUnbanning, t }: BannedUsersTableProps) {
  if (!users?.length) {
    return <div className="text-center py-8 text-base-content/50">{t("no_banned_users")}</div>;
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
          {users.map((user) => (
            <tr key={user.id}>
              <td>{user.username}</td>
              <td>{user.reason}</td>
              <td>{new Date(user.created_at).toLocaleDateString()}</td>
              <td>{user.expires_at ? new Date(user.expires_at).toLocaleDateString() : t("permanent")}</td>
              <td>
                {onUnban && (
                  <button
                    className="btn btn-xs btn-success"
                    onClick={() => onUnban(user.user_id)}
                    disabled={isUnbanning}
                  >
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