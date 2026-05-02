import { useState, useEffect, useCallback } from "react";
import { Shield, ShieldAlert, ShieldCheck, Trash2, Edit } from "lucide-react";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import {
  useAdminModeratorList,
  useRemoveModerator,
  useUpdateModeratorPermissions,
} from "../hooks/useAdminModerator";
import { Moderator } from "@/shared/api/modules/moderator";

// 类型定义
interface ModeratorPermission {
  can_delete_post: boolean;
  can_pin_post: boolean;
  can_edit_any_post: boolean;
  can_manage_moderator: boolean;
  can_ban_user: boolean;
}

interface ModeratorsTableProps {
  boardId: number;
}

// 权限配置
const PERMISSIONS_CONFIG = [
  {
    key: "can_delete_post" as const,
    labelKey: "delete_posts",
    descriptionKey: "can_delete_posts_desc",
  },
  {
    key: "can_pin_post" as const,
    labelKey: "pin_posts",
    descriptionKey: "can_pin_posts_desc",
  },
  {
    key: "can_edit_any_post" as const,
    labelKey: "edit_posts",
    descriptionKey: "can_edit_posts_desc",
  },
  {
    key: "can_manage_moderator" as const,
    labelKey: "manage_moderators",
    descriptionKey: "can_manage_moderators_desc",
  },
  {
    key: "can_ban_user" as const,
    labelKey: "ban_users",
    descriptionKey: "can_ban_users_desc",
  },
];

export function ModeratorsTable({ boardId }: ModeratorsTableProps) {
  const t = useTranslations("Admin");
  const [moderators, setModerators] = useState<Moderator[]>([]);
  const [editingModerator, setEditingModerator] = useState<Moderator | null>(
    null,
  );
  const [permissions, setPermissions] = useState<ModeratorPermission>({
    can_delete_post: false,
    can_pin_post: false,
    can_edit_any_post: false,
    can_manage_moderator: false,
    can_ban_user: false,
  });

  const { data: moderatorsData, refetch } = useAdminModeratorList(
    (boardId = 1),
  );
  const removeModerator = useRemoveModerator(boardId);
  const updatePermissions = useUpdateModeratorPermissions(boardId);

  // 加载版主列表 - 使用 useCallback 稳定化函数
  const loadAllModerators = useCallback(() => {
    if (moderatorsData) {
      setModerators(moderatorsData);
    }
  }, [moderatorsData]);

  useEffect(() => {
    loadAllModerators();
  }, [loadAllModerators]);

  // 处理移除版主
  const handleRemoveModerator = async (userId: number, username: string) => {
    if (confirm(t("confirm_remove_moderator", { username }))) {
      try {
        await removeModerator.mutateAsync(userId);
        toast.success(t("moderator_removed"));
        refetch();
      } catch {
        toast.error(t("operation_failed"));
      }
    }
  };

  // 打开编辑权限对话框
  const openEditPermissions = (moderator: Moderator) => {
    setEditingModerator(moderator);
    setPermissions(moderator.permissions);
  };

  // 保存权限修改
  const handleSavePermissions = async () => {
    if (!editingModerator) return;

    try {
      await updatePermissions.mutateAsync({
        userId: editingModerator.user_id,
        data: permissions,
      });
      toast.success(t("permissions_updated"));
      setEditingModerator(null);
      refetch();
    } catch {
      toast.error(t("operation_failed"));
    }
  };

  // 处理权限变更
  const handlePermissionChange = (
    permissionKey: keyof ModeratorPermission,
    value: boolean,
  ) => {
    setPermissions((prev) => ({
      ...prev,
      [permissionKey]: value,
    }));
  };

  // 获取角色图标
  const getRoleIcon = (permissions: ModeratorPermission) => {
    if (permissions.can_manage_moderator) {
      return <ShieldAlert className="w-4 h-4 text-warning" />;
    }
    if (
      permissions.can_delete_post ||
      permissions.can_ban_user ||
      permissions.can_edit_any_post
    ) {
      return <ShieldCheck className="w-4 h-4 text-primary" />;
    }
    return <Shield className="w-4 h-4 text-base-content/50" />;
  };

  // 获取角色名称
  const getRoleName = (permissions: ModeratorPermission) => {
    if (permissions.can_manage_moderator) return t("role.super_moderator");
    if (permissions.can_ban_user && permissions.can_delete_post)
      return t("role.full_moderator");
    if (permissions.can_delete_post) return t("role.content_moderator");
    return t("role.junior_moderator");
  };

  if (!moderators.length) {
    console.log("moderators: ", moderators);
    return (
      <div className="text-center py-8 text-base-content/50">
        {t("no_moderators")}
      </div>
    );
  }

  return (
    <>
      <div className="overflow-x-auto">
        <table className="table">
          <thead>
            <tr>
              <th>{t("user")}</th>
              <th>{t("role")}</th>
              <th>{t("permissions")}</th>
              <th>{t("actions")}</th>
            </tr>
          </thead>
          <tbody>
            {moderators.map((moderator: Moderator) => (
              <tr key={moderator.id}>
                <td>
                  <div className="flex items-center gap-2">
                    <div className="avatar placeholder">
                      <div className="w-8 h-8 rounded-full bg-primary/10">
                        <span className="text-xs font-medium">
                          {moderator.user?.username?.[0]?.toUpperCase() || "U"}
                        </span>
                      </div>
                    </div>
                    <span>{moderator.user?.username}</span>
                  </div>
                </td>
                <td>
                  <div className="flex items-center gap-1">
                    {getRoleIcon(moderator.permissions)}
                    <span className="text-sm">
                      {getRoleName(moderator.permissions)}
                    </span>
                  </div>
                </td>
                <td>
                  <div className="flex gap-1 flex-wrap">
                    {moderator.permissions.can_delete_post && (
                      <span className="badge badge-sm">{t("delete")}</span>
                    )}
                    {moderator.permissions.can_pin_post && (
                      <span className="badge badge-sm">{t("pin")}</span>
                    )}
                    {moderator.permissions.can_edit_any_post && (
                      <span className="badge badge-sm">{t("edit")}</span>
                    )}
                    {moderator.permissions.can_manage_moderator && (
                      <span className="badge badge-sm badge-warning">
                        {t("manage")}
                      </span>
                    )}
                    {moderator.permissions.can_ban_user && (
                      <span className="badge badge-sm badge-error">
                        {t("ban")}
                      </span>
                    )}
                  </div>
                </td>
                <td>
                  <div className="flex gap-1">
                    <button
                      className="btn btn-ghost btn-xs"
                      onClick={() => openEditPermissions(moderator)}
                      title={t("edit_permissions")}
                    >
                      <Edit className="w-3 h-3" />
                    </button>
                    <button
                      className="btn btn-ghost btn-xs text-error"
                      onClick={() =>
                        handleRemoveModerator(
                          moderator.user_id,
                          moderator.user?.username || "",
                        )
                      }
                      disabled={removeModerator.isPending}
                      title={t("remove_moderator")}
                    >
                      <Trash2 className="w-3 h-3" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* 编辑权限对话框 */}
      {editingModerator && (
        <dialog
          className="modal modal-open"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setEditingModerator(null);
            }
          }}
        >
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">
              {t("edit_permissions_for")} {editingModerator.user?.username}
            </h3>
            <div className="space-y-3">
              {PERMISSIONS_CONFIG.map((config) => (
                <div key={config.key} className="form-control">
                  <label className="label cursor-pointer">
                    <div>
                      <span className="label-text font-medium">
                        {t(config.labelKey)}
                      </span>
                      <p className="text-xs text-base-content/50">
                        {t(config.descriptionKey)}
                      </p>
                    </div>
                    <input
                      type="checkbox"
                      className="toggle toggle-primary toggle-sm"
                      checked={permissions[config.key]}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        handlePermissionChange(config.key, e.target.checked)
                      }
                    />
                  </label>
                </div>
              ))}
            </div>
            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => setEditingModerator(null)}
              >
                {t("cancel")}
              </button>
              <button
                className="btn btn-primary"
                onClick={handleSavePermissions}
                disabled={updatePermissions.isPending}
              >
                {updatePermissions.isPending ? (
                  <span className="loading loading-spinner loading-xs" />
                ) : null}
                {t("save")}
              </button>
            </div>
          </div>
        </dialog>
      )}
    </>
  );
}
