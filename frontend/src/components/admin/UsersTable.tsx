// components/admin/UsersTable.tsx
import Avatar from "@/components/user/Avatar";
import { RoleBadge } from "@/components/common/RoleBadge";
import { UserRoleType } from "@/type/roles.types";
import { Ban, CheckCircle, XCircle } from "lucide-react";
import { formatDate } from "@/lib/utils";
import { User } from "@/shared/api";
import { UserActionMenu } from "./UserActionMenu";
import { useTranslations } from "next-intl";

// 翻译函数的类型定义

// 用户表格组件 Props 类型
interface UsersTableProps {
  users: User[];
  currentUserId?: number;
  onToggleActive: (id: number, active: boolean) => void;
  onToggleBlock: (id: number, blocked: boolean) => void;
  onToggleRole?: (id: number, role: string) => void;
  onDeleteUser?: (id: number, username: string) => void;
  onResetPassword?: (id: number, username: string) => void;
  isTogglingActive: boolean;
  isTogglingBlock: boolean;
  isDeleting?: boolean;
  isUpdatingRole?: boolean;
}

// 用户表格组件
export function UsersTable({
  users,
  currentUserId,
  onToggleActive,
  onToggleBlock,
  onToggleRole,
  onDeleteUser,
  onResetPassword,
  isTogglingActive,
  isTogglingBlock,
  isDeleting = false,
  isUpdatingRole = false,
}: UsersTableProps) {
  const t = useTranslations("Admin");
  if (users.length === 0) {
    return (
      <div className="text-center py-8 text-base-content/60">
        {t("no_data")}
      </div>
    );
  }

  // 获取用户状态显示
  const getUserStatusBadges = (user: User): React.ReactNode[] => {
    const badges: React.ReactNode[] = [];

    // IsBlocked 优先级最高，如果被封禁，只显示封禁状态
    if (user.is_blocked) {
      badges.push(
        <span key="blocked" className="badge badge-sm badge-error gap-1">
          <Ban className="w-3 h-3" /> {t("blocked")}
        </span>,
      );
    } else {
      // 未封禁时显示激活状态
      badges.push(
        <span
          key="active"
          className={`badge badge-sm ${user.is_active ? "badge-success" : "badge-warning"} gap-1`}
        >
          {user.is_active ? (
            <CheckCircle className="w-3 h-3" />
          ) : (
            <XCircle className="w-3 h-3" />
          )}
          {user.is_active ? t("activated") : t("inactive")}
        </span>,
      );
    }

    return badges;
  };

  return (
    <div className="overflow-x-auto">
      <table className="table table-zebra" style={{ overflow: "visible" }}>
        <thead>
          <tr>
            <th>{t("user")}</th>
            <th>{t("email")}</th>
            <th>{t("roles")}</th>
            <th>{t("score")}</th>
            <th>{t("registration_at")}</th>
            <th>{t("last_login")}</th>
            <th>{t("status")}</th>
            <th>{t("operation")}</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user: User) => {
            const isBlocked = user.is_blocked;

            return (
              <tr
                key={user.id}
                className={isBlocked ? "bg-base-200/50 opacity-75" : ""}
                style={{ overflow: "visible" }}
              >
                <td>
                  <div className="flex items-center gap-2">
                    <div className="avatar">
                      <div className="w-8 h-8 rounded-full">
                        <Avatar
                          username={user.username}
                          avatarUrl={user.avatar}
                          size="md"
                        />
                      </div>
                    </div>
                    <span
                      className={`font-medium text-sm ${isBlocked ? "line-through text-base-content/50" : ""}`}
                    >
                      {user.username}
                    </span>
                  </div>
                </td>
                <td className="text-sm text-base-content/60">{user.email}</td>
                <td>
                  <RoleBadge
                    role={user.role as UserRoleType}
                    showIcon
                    size="sm"
                  />
                </td>
                <td className="text-sm font-medium text-warning">
                  {user.score}
                </td>
                <td className="text-xs text-base-content/50">
                  {formatDate(user.created_at)}
                </td>
                <td className="text-xs text-base-content/50">
                  {user.last_login
                    ? formatDate(user.last_login)
                    : t("never_logged_in")}
                </td>
                <td>
                  <div className="flex gap-1 flex-wrap">
                    {getUserStatusBadges(user)}
                  </div>
                </td>
                <td>
                  <UserActionMenu
                    user={user}
                    isCurrentUser={user.id === currentUserId}
                    onToggleActive={onToggleActive}
                    onToggleBlock={onToggleBlock}
                    onToggleRole={onToggleRole}
                    onDeleteUser={onDeleteUser}
                    onResetPassword={onResetPassword}
                    isTogglingActive={isTogglingActive}
                    isTogglingBlock={isTogglingBlock}
                    isDeleting={isDeleting}
                    isUpdatingRole={isUpdatingRole}
                  />
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
