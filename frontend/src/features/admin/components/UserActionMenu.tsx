// components/admin/UserActionMenu.tsx
import { useState, useRef, useEffect, JSX, useCallback } from "react";
import {
  MoreVertical,
  Trash2,
  AlertTriangle,
  Mail,
  Key,
  Crown,
  Shield,
  Eye,
  Hammer,
  Lock,
  Unlock,
  CheckCircle,
  XCircle,
  User as UserIcon,
  UserPlus,
  Bot,
  User2,
  EyeIcon,
  MonitorUpIcon,
} from "lucide-react";
import toast from "react-hot-toast";
import { UsersIcon } from "@heroicons/react/24/solid";
import { useTranslations } from "next-intl";
import { UserRoleType } from "@/shared/type/roles.types";
import { UserDO } from "@/shared/api/types/user.model";

// 类型定义

interface UserActionMenuProps {
  user: UserDO;
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
}

// 角色配置类型
interface RoleOption {
  label: string;
  icon: JSX.Element;
  className: string;
  nextRole: UserRoleType;
}

// 菜单位置类型
interface MenuPosition {
  top: number;
  left: number;
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
  isDeleting = false,
  isUpdatingRole = false,
}: UserActionMenuProps) {
  const t = useTranslations("Admin");
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [showRoleMenu, setShowRoleMenu] = useState<boolean>(false);
  const [menuPosition, setMenuPosition] = useState<MenuPosition>({
    top: 0,
    left: 0,
  });
  const dropdownRef = useRef<HTMLDivElement>(null);

  // 计算菜单位置
  const calculateMenuPosition = useCallback((): void => {
    if (!dropdownRef.current) return;

    const rect = dropdownRef.current.getBoundingClientRect();
    const viewportWidth = window.innerWidth;
    const menuWidth = 208; // min-w-52 = 208px
    let left = rect.left;

    // 如果菜单会超出右边界，则向左对齐
    if (rect.right + menuWidth > viewportWidth) {
      left = rect.right - menuWidth;
    }

    setMenuPosition({
      top: rect.bottom + 4,
      left: left,
    });
  }, []);

  // 点击外部关闭菜单
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
        setShowRoleMenu(false);
      }
    };

    if (isOpen) {
      calculateMenuPosition();
      document.addEventListener("mousedown", handleClickOutside);
      // 阻止 body 滚动
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.body.style.overflow = "";
    };
  }, [isOpen, calculateMenuPosition]);

  // 监听窗口大小变化重新计算菜单位置
  useEffect(() => {
    if (!isOpen) return;

    const handleResize = () => {
      calculateMenuPosition();
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [isOpen, calculateMenuPosition]);

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
        label: t("set_inactive"),
        icon: <XCircle className="w-4 h-4" />,
        className: "text-warning",
        hoverClass: "hover:bg-warning/10",
        action: () => onToggleActive(user.id, false),
      };
    } else {
      return {
        label: t("set_active"),
        icon: <CheckCircle className="w-4 h-4" />,
        className: "text-success",
        hoverClass: "hover:bg-success/10",
        action: () => onToggleActive(user.id, true),
      };
    }
  };

  // 获取封禁/解封配置
  const getBlockConfig = () => {
    if (user.is_blocked) {
      return {
        label: t("set_unblock"),
        icon: <Unlock className="w-4 h-4" />,
        className: "text-success",
        hoverClass: "hover:bg-success/10",
        action: () => onToggleBlock(user.id, false),
      };
    } else {
      return {
        label: t("set_block"),
        icon: <Lock className="w-4 h-4" />,
        className: "text-error",
        hoverClass: "hover:bg-error/10",
        action: () => onToggleBlock(user.id, true),
      };
    }
  };

  // 获取可用的角色切换选项
  const getAvailableRoleOptions = (): RoleOption[] => {
    const currentRole = user.role as UserRoleType;
    const options: RoleOption[] = [];

    switch (currentRole) {
      case "super_admin":
        options.push({
          label: t("role.admin"),
          icon: <Crown className="w-4 h-4" />,
          className: "text-error",
          nextRole: "admin",
        });
        options.push({
          label: t("role.user"),
          icon: <UserIcon className="w-4 h-4" />,
          className: "text-info",
          nextRole: "user",
        });
        break;
      case "admin":
        options.push({
          label: t("role.reviewer"),
          icon: <EyeIcon className="w-4 h-4" />,
          className: "text-warning",
          nextRole: "reviewer",
        });
        options.push({
          label: t("role.moderator"),
          icon: <Hammer className="w-4 h-4" />,
          className: "text-purple-500",
          nextRole: "moderator",
        });
        options.push({
          label: t("role.user"),
          icon: <UserIcon className="w-4 h-4" />,
          className: "text-info",
          nextRole: "user",
        });
        break;
      case "moderator":
        options.push({
          label: t("role.admin"),
          icon: <Shield className="w-4 h-4" />,
          className: "text-warning",
          nextRole: "admin",
        });
        options.push({
          label: t("role.reviewer"),
          icon: <Eye className="w-4 h-4" />,
          className: "text-indigo-500",
          nextRole: "reviewer",
        });
        options.push({
          label: t("role.user"),
          icon: <UserIcon className="w-4 h-4" />,
          className: "text-info",
          nextRole: "user",
        });
        break;
      case "reviewer":
        options.push({
          label: t("role.admin"),
          icon: <Shield className="w-4 h-4" />,
          className: "text-warning",
          nextRole: "admin",
        });
        options.push({
          label: t("role.moderator"),
          icon: <Hammer className="w-4 h-4" />,
          className: "text-purple-500",
          nextRole: "moderator",
        });
        options.push({
          label: t("role.user"),
          icon: <UserIcon className="w-4 h-4" />,
          className: "text-info",
          nextRole: "user",
        });
        break;
      case "member":
        options.push({
          label: t("role.reviewer"),
          icon: <EyeIcon className="w-4 h-4" />,
          className: "text-warning",
          nextRole: "reviewer",
        });
        options.push({
          label: t("role.admin"),
          icon: <Shield className="w-4 h-4" />,
          className: "text-warning",
          nextRole: "admin",
        });
        options.push({
          label: t("role.user"),
          icon: <UserIcon className="w-4 h-4" />,
          className: "text-info",
          nextRole: "user",
        });
        break;
      case "user":
        options.push({
          label: t("role.member"),
          icon: <UserPlus className="w-4 h-4" />,
          className: "text-success",
          nextRole: "member",
        });
        options.push({
          label: t("role.guest"),
          icon: <User2 className="w-4 h-4" />,
          className: "text-info",
          nextRole: "guest",
        });
        break;
      case "guest":
        options.push({
          label: t("role.user"),
          icon: <UserIcon className="w-4 h-4" />,
          className: "text-info",
          nextRole: "user",
        });
        break;
      case "system":
        // 系统管理员没有社区权限
        break;
    }

    return options;
  };

  // 获取当前角色显示配置
  const getCurrentRoleConfig = () => {
    const roleIcons: Partial<Record<UserRoleType, JSX.Element>> = {
      super_admin: <Crown className="w-4 h-4" />,
      admin: <Shield className="w-4 h-4" />,
      moderator: <Hammer className="w-4 h-4" />,
      reviewer: <Eye className="w-4 h-4" />,
      member: <UserPlus className="w-4 h-4" />,
      user: <UsersIcon className="w-4 h-4" />,
      system: <MonitorUpIcon className="w-4 h-4" />,
    };

    return {
      icon: roleIcons[user.role as UserRoleType] || (
        <UserIcon className="w-4 h-4" />
      ),
      label: t(`role.${user.role}`),
    };
  };

  const activeConfig = getActiveConfig();
  const blockConfig = getBlockConfig();
  const roleOptions = getAvailableRoleOptions();
  const currentRoleConfig = getCurrentRoleConfig();
  const isBlocked = user.is_blocked;

  // 打开删除确认模态框
  const openDeleteModal = () => {
    const modal = document.getElementById(
      `delete_modal_${user.id}`,
    ) as HTMLDialogElement;
    modal?.showModal();
    setIsOpen(false);
  };

  // 打开重置密码确认模态框
  const openResetPasswordModal = () => {
    const modal = document.getElementById(
      `reset_pwd_modal_${user.id}`,
    ) as HTMLDialogElement;
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
    setShowRoleMenu(false);
  };

  // 切换角色
  const handleRoleChange = (role: UserRoleType) => {
    onToggleRole?.(user.id, role);
    setIsOpen(false);
    setShowRoleMenu(false);
  };

  return (
    <>
      <div className="relative" ref={dropdownRef}>
        <button
          className={`btn btn-ghost btn-xs btn-square ${isOpen ? "bg-base-200" : ""}`}
          onClick={() => {
            setIsOpen(!isOpen);
            setShowRoleMenu(false);
          }}
          type="button"
        >
          <MoreVertical className="w-4 h-4" />
        </button>

        {isOpen && (
          <div
            className="fixed z-[9999] min-w-52"
            style={{
              top: menuPosition.top,
              left: menuPosition.left,
            }}
          >
            <div className="menu menu-sm p-1 shadow-xl bg-base-100 rounded-box border border-base-200">
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

              <div className="divider my-1" />

              <div className="menu-title text-xs opacity-50 px-3 py-1">
                {t("account_settings")}
              </div>

              {onToggleRole && roleOptions.length > 0 && (
                <div className="menu-item relative">
                  {showRoleMenu ? (
                    <div className="pl-2 border-l-2 border-base-300 ml-2">
                      <div className="menu-title text-xs opacity-50 px-2 py-1">
                        {t("change_role")}
                      </div>
                      {roleOptions.map((option) => (
                        <button
                          key={option.nextRole}
                          className={`${option.className} hover:bg-base-200 gap-2 w-full text-left px-3 py-2 rounded-md text-sm`}
                          onClick={() => handleRoleChange(option.nextRole)}
                          disabled={isUpdatingRole}
                        >
                          {option.icon}
                          <span>{option.label}</span>
                        </button>
                      ))}
                      <button
                        className="text-base-content/50 hover:bg-base-200 gap-2 w-full text-left px-3 py-2 rounded-md text-sm"
                        onClick={() => setShowRoleMenu(false)}
                      >
                        ← {t("back")}
                      </button>
                    </div>
                  ) : (
                    <button
                      className="text-base-content/70 hover:bg-base-200 gap-2 w-full text-left px-3 py-2 rounded-md"
                      onClick={() => setShowRoleMenu(true)}
                      disabled={isUpdatingRole}
                    >
                      {currentRoleConfig.icon}
                      <span>{currentRoleConfig.label}</span>
                      <span className="ml-auto text-xs opacity-50">→</span>
                    </button>
                  )}
                </div>
              )}

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

              <div className="menu-item">
                <button
                  className="text-base-content/70 hover:bg-base-200 gap-2 w-full text-left px-3 py-2 rounded-md"
                  onClick={handleSendEmail}
                >
                  <Mail className="w-4 h-4" />
                  <span>{t("send_email")}</span>
                </button>
              </div>

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
        <dialog
          id={`delete_modal_${user.id}`}
          className="modal"
          onClick={(e) => e.stopPropagation()}
        >
          <div className="modal-box">
            <div className="flex items-center gap-3 mb-4">
              <div className="bg-error/10 p-2 rounded-full">
                <AlertTriangle className="w-6 h-6 text-error" />
              </div>
              <h3 className="font-bold text-lg">{t("confirm_delete")}</h3>
            </div>
            <p className="py-2 text-base-content/80">
              {t("delete_user_confirm", { username: user.username })}
            </p>
            <div className="bg-base-200 rounded-lg p-3 my-3">
              <div className="flex items-center gap-2">
                <div className="avatar placeholder">
                  <div className="bg-neutral-focus text-neutral-content rounded-full w-8">
                    <span className="text-xs">
                      {user.username[0]?.toUpperCase()}
                    </span>
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
        <dialog
          id={`reset_pwd_modal_${user.id}`}
          className="modal"
          onClick={(e) => e.stopPropagation()}
        >
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
