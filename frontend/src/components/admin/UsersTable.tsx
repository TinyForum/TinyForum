import Avatar from "@/components/user/Avatar";
// import { User } from "@/type/admin.types";
import { ShieldOff, ShieldCheck } from "lucide-react";
import { formatDate } from "@/lib/utils";
import { User } from "@/lib/api";
// 用户表格组件
export function UsersTable({
  users,
  currentUserId,
  onToggleActive,
  isToggling,
  t
}: {
  users: User[];
  currentUserId?: number;
  onToggleActive: (id: number, active: boolean) => void;
  isToggling: boolean;
  t: (key: string) => string;
}) {
  if (users.length === 0) {
    return <div className="text-center py-8 text-base-content/60">{t("no_data")}</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="table table-zebra">
        <thead>
          <tr>
            <th>{t("user")}</th>
            <th>{t("email")}</th>
            <th>{t("role")}</th>
            <th>{t("score")}</th>
            <th>{t("registration_at")}</th>
            <th>{t("status")}</th>
            <th>{t("operation")}</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id}>
              <td>
                <div className="flex items-center gap-2">
                  <div className="avatar">
                    <div className="w-8 h-8 rounded-full">
                      <Avatar username={u.username} avatarUrl={u.avatar} size="md" />
                    </div>
                  </div>
                  <span className="font-medium text-sm">{u.username}</span>
                </div>
              </td>
              <td className="text-sm text-base-content/60">{u.email}</td>
              <td>
                <span className={`badge badge-sm ${u.role === "admin" ? "badge-warning" : "badge-ghost"}`}>
                  {u.role === "admin" ? t("administrator") : t("user")}
                </span>
              </td>
              <td className="text-sm font-medium text-warning">{u.score}</td>
              <td className="text-xs text-base-content/50">{formatDate(u.created_at)}</td>
              <td>
                <span className={`badge badge-sm ${u.is_active ? "badge-success" : "badge-error"}`}>
                  {u.is_active ? t("active") : t("banned")}
                </span>
              </td>
              <td>
                {u.id !== currentUserId && (
                  <button
                    className={`btn btn-xs gap-1 ${u.is_active ? "btn-error btn-outline" : "btn-success btn-outline"}`}
                    onClick={() => onToggleActive(u.id, !u.is_active)}
                    disabled={isToggling}
                  >
                    {u.is_active ? (
                      <>
                        <ShieldOff className="w-3 h-3" /> {t("ban")}
                      </>
                    ) : (
                      <>
                        <ShieldCheck className="w-3 h-3" /> {t("unban")}
                      </>
                    )}
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