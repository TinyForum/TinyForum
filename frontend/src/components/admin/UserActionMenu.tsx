// components/admin/UserActionMenu.tsx
import { useState, useRef, useEffect } from "react";
import { 
  MoreVertical, 
  Trash2, 
  AlertTriangle,
  Mail,
  Key,
  Crown,
  UserX,
  Lock,
  Unlock,
  CheckCircle,
  XCircle
} from "lucide-react";
import toast from "react-hot-toast";
import { User } from "@/lib/api";

interface UserActionMenuProps {
  user: User;
  isCurrentUser: boolean;
  onToggleActive: (id: number, active: boolean) => void;
  onToggleBlock: (id: number, blocked: boolean) => void;
  onToggleRole?: (id: number, role: string) => void;
  onDeleteUser?: (id: number, username: string) => void;
  onResetPassword?: (id: number, username: string) => void;
  isTogglingActive: boolean;
  isTogglingBlock: boolean;
  isDeleting?: boolean;
  isUpdatingRole?: boolean;
  t: (key: string, params?: any) => string;
}

export function UserActionMenu({
  user,
  isCurrentUser,
  onToggleActive,
  onToggleBlock,
  onToggleRole,
  onDeleteUser,
  onResetPassword,
  isTogglingActive,
  isTogglingBlock,
  isDeleting,
  isUpdatingRole,
  t
}: UserActionMenuProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // 点击外部关闭菜单
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      // 防止 body 滚动
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  if (isCurrentUser) {
    return (
      <span className="text-xs text-base-content/40 px-2">
        {t("current_user")}
      </span>
    );
  }

  // 获取激活/停用配置
  const getActiveConfig = () => {
    if (user.is_active) {
      return {
        label: t("inactivate"),
        icon: <XCircle className="w-4 h-4" />,
        className: "text-warning",
        hoverClass: "hover:bg-warning/10",
        action: () => onToggleActive(user.id, false)
      };
    } else {
      return {
        label: t("activate"),
        icon: <CheckCircle className="w-4 h-4" />,
        className: "text-success",
        hoverClass: "hover:bg-success/10",
        action: () => onToggleActive(user.id, true)
      };
    }
  };

  // 获取封禁/解封配置
  const getBlockConfig = () => {
    if (user.is_blocked) {
      return {
        label: t("unblock"),
        icon: <Unlock className="w-4 h-4" />,
        className: "text-success",
        hoverClass: "hover:bg-success/10",
        action: () => onToggleBlock(user.id, false)
      };
    } else {
      return {
        label: t("blocked"),
        icon: <Lock className="w-4 h-4" />,
        className: "text-error",
        hoverClass: "hover:bg-error/10",
        action: () => onToggleBlock(user.id, true)
      };
    }
  };

  // 获取角色切换配置
  const getRoleConfig = () => {
    if (user.role === "admin") {
      return {
        label: t("remove_admin"),
        icon: <UserX className="w-4 h-4" />,
        className: "text-warning",
        action: () => onToggleRole?.(user.id, "user")
      };
    } else {
      return {
        label: t("set_as_admin"),
        icon: <Crown className="w-4 h-4" />,
        className: "text-primary",
        action: () => onToggleRole?.(user.id, "admin")
      };
    }
  };

  const activeConfig = getActiveConfig();
  const blockConfig = getBlockConfig();
  const roleConfig = getRoleConfig();
  const isBlocked = user.is_blocked;

  // 打开删除确认模态框
  const openDeleteModal = () => {
    const modal = document.getElementById(`delete_modal_${user.id}`) as HTMLDialogElement;
    modal?.showModal();
    setIsOpen(false);
  };

  // 打开重置密码确认模态框
  const openResetPasswordModal = () => {
    const modal = document.getElementById(`reset_pwd_modal_${user.id}`) as HTMLDialogElement;
    modal?.showModal();
    setIsOpen(false);
  };

  // 发送邮件
  const handleSendEmail = () => {
    toast.success(`${t("send_email_to")} ${user.username}`);
    setIsOpen(false);
  };

  // 执行操作并关闭菜单
  const executeAction = (action: () => void) => {
    action();
    setIsOpen(false);
  };

  return (
    <>
      <div className="relative" ref={dropdownRef}>
        {/* 触发器按钮 */}
        <button
          className={`btn btn-ghost btn-xs btn-square ${isOpen ? "bg-base-200" : ""}`}
          onClick={() => setIsOpen(!isOpen)}
          type="button"
        >
          <MoreVertical className="w-4 h-4" />
        </button>

        {/* 下拉菜单 - 使用 Portal 渲染到 body */}
        {isOpen && (
          <div 
            className="fixed z-[9999] min-w-48"
            style={{
              top: (() => {
                const rect = dropdownRef.current?.getBoundingClientRect();
                return rect ? rect.bottom + 4 : 0;
              })(),
              left: (() => {
                const rect = dropdownRef.current?.getBoundingClientRect();
                const viewportWidth = window.innerWidth;
                const menuWidth = 192; // w-48 = 12rem = 192px
                if (rect) {
                  // 如果右侧空间不足，向左对齐
                  if (rect.right + menuWidth > viewportWidth) {
                    return rect.right - menuWidth;
                  }
                  return rect.left;
                }
                return 0;
              })(),
            }}
          >
            <div className="menu menu-sm p-1 shadow-xl bg-base-100 rounded-box border border-base-200">
              {/* 封禁/解封 - 优先级最高的操作 */}
              <div className="menu-item">
                <button
                  className={`${blockConfig.className} ${blockConfig.hoverClass} gap-2 w-full text-left px-3 py-2 rounded-md`}
                  onClick={() => executeAction(blockConfig.action)}
                  disabled={isTogglingBlock}
                >
                  {blockConfig.icon}
                  <span>{blockConfig.label}</span>
                </button>
              </div>

              {/* 激活/停用 - 仅未封禁时显示 */}
              {!isBlocked && (
                <div className="menu-item">
                  <button
                    className={`${activeConfig.className} ${activeConfig.hoverClass} gap-2 w-full text-left px-3 py-2 rounded-md`}
                    onClick={() => executeAction(activeConfig.action)}
                    disabled={isTogglingActive}
                  >
                    {activeConfig.icon}
                    <span>{activeConfig.label}</span>
                  </button>
                </div>
              )}

              {/* 分隔线 */}
              <div className="divider my-1" />

              {/* 账户设置分组标题 */}
              <div className="menu-title text-xs opacity-50 px-3 py-1">
                {t("account_settings")}
              </div>

              {/* 角色管理 */}
              {onToggleRole && (
                <div className="menu-item">
                  <button
                    className={`${roleConfig.className} hover:bg-primary/10 gap-2 w-full text-left px-3 py-2 rounded-md`}
                    onClick={() => executeAction(roleConfig.action)}
                    disabled={isUpdatingRole}
                  >
                    {roleConfig.icon}
                    <span>{roleConfig.label}</span>
                  </button>
                </div>
              )}

              {/* 重置密码 */}
              {onResetPassword && (
                <div className="menu-item">
                  <button
                    className="text-info hover:bg-info/10 gap-2 w-full text-left px-3 py-2 rounded-md"
                    onClick={openResetPasswordModal}
                  >
                    <Key className="w-4 h-4" />
                    <span>{t("reset_password")}</span>
                  </button>
                </div>
              )}

              {/* 发送邮件 */}
              <div className="menu-item">
                <button
                  className="text-base-content/70 hover:bg-base-200 gap-2 w-full text-left px-3 py-2 rounded-md"
                  onClick={handleSendEmail}
                >
                  <Mail className="w-4 h-4" />
                  <span>{t("send_email")}</span>
                </button>
              </div>

              {/* 删除用户 - 危险操作 */}
              {onDeleteUser && (
                <>
                  <div className="divider my-1" />
                  <div className="menu-title text-xs opacity-50 px-3 py-1">
                    <span>{t("danger_zone")}</span>
                  </div>
                  <div className="menu-item">
                    <button
                      className="text-error hover:bg-error/10 gap-2 w-full text-left px-3 py-2 rounded-md"
                      onClick={openDeleteModal}
                      disabled={isDeleting}
                    >
                      <Trash2 className="w-4 h-4" />
                      <span>{t("delete_user")}</span>
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        )}
      </div>

      {/* 删除确认对话框 */}
      {onDeleteUser && (
        <dialog id={`delete_modal_${user.id}`} className="modal" onClick={(e) => e.stopPropagation()}>
          <div className="modal-box">
            <div className="flex items-center gap-3 mb-4">
              <div className="bg-error/10 p-2 rounded-full">
                <AlertTriangle className="w-6 h-6 text-error" />
              </div>
              <h3 className="font-bold text-lg">{t("confirm_delete")}</h3>
            </div>
            <p className="py-2 text-base-content/80">
              {t("delete_user_confirm")}
            </p>
            <div className="bg-base-200 rounded-lg p-3 my-3">
              <div className="flex items-center gap-2">
                <div className="avatar placeholder">
                  <div className="bg-neutral-focus text-neutral-content rounded-full w-8">
                    <span className="text-xs">{user.username[0]?.toUpperCase()}</span>
                  </div>
                </div>
                <div>
                  <p className="font-medium">{user.username}</p>
                  <p className="text-xs text-base-content/60">{user.email}</p>
                </div>
              </div>
            </div>
            <div className="alert alert-warning mt-4">
              <AlertTriangle className="w-4 h-4" />
              <span className="text-sm">{t("delete_warning")}</span>
            </div>
            <div className="modal-action">
              <form method="dialog" className="flex gap-2">
                <button className="btn btn-ghost">{t("cancel")}</button>
                <button 
                  className="btn btn-error"
                  onClick={() => onDeleteUser(user.id, user.username)}
                >
                  {t("confirm_delete")}
                </button>
              </form>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button>{t("close")}</button>
          </form>
        </dialog>
      )}

      {/* 重置密码确认对话框 */}
      {onResetPassword && (
        <dialog id={`reset_pwd_modal_${user.id}`} className="modal" onClick={(e) => e.stopPropagation()}>
          <div className="modal-box">
            <div className="flex items-center gap-3 mb-4">
              <div className="bg-info/10 p-2 rounded-full">
                <Key className="w-6 h-6 text-info" />
              </div>
              <h3 className="font-bold text-lg">{t("reset_password")}</h3>
            </div>
            <p className="py-2 text-base-content/80">
              {t("reset_password_confirm", { username: user.username })}
            </p>
            <div className="alert alert-info mt-4">
              <Mail className="w-4 h-4" />
              <span className="text-sm">{t("reset_password_notice")}</span>
            </div>
            <div className="modal-action">
              <form method="dialog" className="flex gap-2">
                <button className="btn btn-ghost">{t("cancel")}</button>
                <button 
                  className="btn btn-info"
                  onClick={() => onResetPassword(user.id, user.username)}
                >
                  {t("confirm_reset")}
                </button>
              </form>
            </div>
          </div>
          <form method="dialog" className="modal-backdrop">
            <button>{t("close")}</button>
          </form>
        </dialog>
      )}
    </>
  );
}