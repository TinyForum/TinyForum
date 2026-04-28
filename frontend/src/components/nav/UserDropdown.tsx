// components/layout/UserDropdown.tsx
"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useEffect, useCallback, useRef } from "react";
import {
  User as UserIcon,
  Settings,
  HelpCircle,
  LogOut,
  Sparkles,
  Bookmark,
  MessageCircleQuestion,
  ShieldCheckIcon,
  Crown,
  Hammer,
  Eye,
} from "lucide-react";
import Avatar from "../user/Avatar";
import { useLogoutStore } from "@/store/logout";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { User } from "@/lib/api";

interface UserDropdownProps {
  user: User;
  isOpen?: boolean;
  onOpenChange?: (isOpen: boolean) => void;
}

export default function UserDropdown({
  user,
  isOpen: controlledIsOpen,
  onOpenChange,
}: UserDropdownProps) {
  const router = useRouter();
  const t = useTranslations("Common");
  const { logout, isLoading } = useLogoutStore();
  const dropdownRef = useRef<HTMLDivElement>(null);

  // 内部状态（非受控模式）
  const [internalIsOpen, setInternalIsOpen] = useState(false);

  // 判断是否为受控组件
  const isControlled = controlledIsOpen !== undefined;
  const isOpen = isControlled ? controlledIsOpen : internalIsOpen;

  // 使用 useCallback 稳定化 setIsOpen 函数
  const setIsOpen = useCallback(
    (newIsOpen: boolean) => {
      if (isControlled) {
        onOpenChange?.(newIsOpen);
      } else {
        setInternalIsOpen(newIsOpen);
      }
    },
    [isControlled, onOpenChange],
  );

  // 处理点击外部关闭
  useEffect(() => {
    if (!isOpen) return;

    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      if (dropdownRef.current && !dropdownRef.current.contains(target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isOpen, setIsOpen]); // 添加 setIsOpen 依赖

  // 处理 ESC 键关闭
  useEffect(() => {
    if (!isOpen) return;

    const handleEsc = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsOpen(false);
      }
    };

    document.addEventListener("keydown", handleEsc);
    return () => document.removeEventListener("keydown", handleEsc);
  }, [isOpen, setIsOpen]); // 添加 setIsOpen 依赖

  // 处理登出
  const handleLogout = async () => {
    try {
      await logout();
      toast.success(t("logout_success"));
      setIsOpen(false);
      router.push("/");
      router.refresh();
    } catch {
      // 移除未使用的 error 变量
      toast.error(t("logout_failed"));
    }
  };

  // 处理菜单项点击
  const handleMenuClick = () => {
    setIsOpen(false);
  };

  // 获取后台入口配置
  const getDashboardConfig = () => {
    const role = user?.role;

    if (role === "admin" || role === "super_admin") {
      return {
        icon: <Crown className="w-4 h-4" />,
        label: t("admin_dashboard"),
        path: "/dashboard/admin",
        className: "text-primary",
      };
    }

    if (role === "moderator") {
      return {
        icon: <Hammer className="w-4 h-4" />,
        label: t("moderator_dashboard"),
        path: "/dashboard/moderator",
        className: "text-secondary",
      };
    }

    if (role === "reviewer") {
      return {
        icon: <Eye className="w-4 h-4" />,
        label: t("reviewer_dashboard"),
        path: "/dashboard/reviewer",
        className: "text-accent",
      };
    }

    if (role === "member") {
      return {
        icon: <Eye className="w-4 h-4" />,
        label: t("member_dashboard"),
        path: "/dashboard/member",
        className: "text-accent",
      };
    }
    if (role === "user") {
      return {
        icon: <Eye className="w-4 h-4" />,
        label: t("user_dashboard"),
        path: "/dashboard/user",
        className: "text-accent",
      };
    }
    return null;
  };

  const dashboardConfig = getDashboardConfig();

  // 是否有管理权限
  const hasManagementAccess =
    user.role === "super_admin" ||
    user.role === "admin" ||
    user.role === "reviewer" ||
    user.role === "moderator" ||
    user.role === "member" ||
    user.role === "user";

  return (
    <div className="dropdown dropdown-end user-dropdown-container" ref={dropdownRef}>
      <div
        tabIndex={0}
        role="button"
        className="btn btn-ghost btn-circle avatar hover:ring-2 hover:ring-primary/20 transition-all"
        onClick={() => setIsOpen(!isOpen)}
      >
        <div className="w-9 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
          <Avatar username={user.username} avatarUrl={user.avatar} size="md" />
        </div>
      </div>

      {isOpen && (
        <ul className="dropdown-content menu bg-base-100 rounded-box z-10 w-64 p-2 shadow-xl border border-base-200 mt-2">
          {/* 用户信息 */}
          <li className="menu-title">
            <div className="flex items-center gap-2">
              <div className="avatar placeholder">
                <div className="w-10 rounded-full bg-primary/10">
                  <Avatar
                    username={user?.username}
                    avatarUrl={user?.avatar}
                    size="md"
                  />
                </div>
              </div>
              <div className="flex-1 min-w-0">
                <span className="text-base-content font-medium block truncate">
                  {user?.username}
                </span>
                <span className="text-xs text-base-content/50 truncate block">
                  {user?.email}
                </span>
                {/* 角色标签 */}
                {user.role !== "user" && (
                  <span
                    className={`badge badge-xs mt-1 ${
                      user?.role === "super_admin"
                        ? "badge-error"
                        : user?.role === "admin"
                          ? "badge-warning"
                          : user?.role === "moderator"
                            ? "badge-secondary"
                            : user?.role === "reviewer"
                              ? "badge-accent"
                              : "badge-ghost"
                    }`}
                  >
                    {t(`role.${user.role}`)}
                  </span>
                )}
              </div>
            </div>
          </li>

          <div className="divider my-1"></div>

          {/* 快速链接 */}
          <li onClick={handleMenuClick}>
            <Link href={`/users/${user.id}`} className="gap-2">
              <UserIcon className="w-4 h-4" />
              {t("profile")}
            </Link>
          </li>
          <li onClick={handleMenuClick}>
            <Link href="/timeline/me" className="gap-2">
              <Sparkles className="w-4 h-4" />
              {t("my_timeline")}
            </Link>
          </li>
          <li onClick={handleMenuClick}>
            <Link href="/topics/me" className="gap-2">
              <Bookmark className="w-4 h-4" />
              {t("my_topics")}
            </Link>
          </li>
          <li onClick={handleMenuClick}>
            <Link href="/questions/me" className="gap-2">
              <MessageCircleQuestion className="w-4 h-4" />
              {t("my_questions")}
            </Link>
          </li>

          <div className="divider my-1"></div>

          {/* 版主申请入口 */}
          {user.role && (
            <li onClick={handleMenuClick}>
              <Link href="/boards/applications" className="gap-2">
                <ShieldCheckIcon className="w-4 h-4" />
                {t("my_moderator_applications")}
              </Link>
            </li>
          )}

          {/* 管理后台入口 */}
          {hasManagementAccess && dashboardConfig && (
            <>
              <div className="divider my-1"></div>
              <li onClick={handleMenuClick}>
                <Link
                  href={dashboardConfig.path}
                  className={`gap-2 ${dashboardConfig.className}`}
                >
                  {dashboardConfig.icon}
                  {dashboardConfig.label}
                </Link>
              </li>
            </>
          )}

          <div className="divider my-1"></div>

          {/* 设置和帮助 */}
          <li onClick={handleMenuClick}>
            <Link href="/settings" className="gap-2">
              <Settings className="w-4 h-4" />
              {t("settings")}
            </Link>
          </li>
          <li onClick={handleMenuClick}>
            <Link href="/help" className="gap-2">
              <HelpCircle className="w-4 h-4" />
              {t("help_center")}
            </Link>
          </li>

          <div className="divider my-1"></div>

          {/* 登出按钮 */}
          <li>
            <button
              onClick={handleLogout}
              className="text-error gap-2"
              disabled={isLoading}
            >
              {isLoading ? (
                <span className="loading loading-spinner loading-xs" />
              ) : (
                <LogOut className="w-4 h-4" />
              )}
              {isLoading ? t("logging_out") : t("logout")}
            </button>
          </li>
        </ul>
      )}
    </div>
  );
}