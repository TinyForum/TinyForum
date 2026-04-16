// components/layout/UserDropdown.tsx
"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  User,
  Settings,
  HelpCircle,
  LogOut,
  LayoutDashboard,
  Sparkles,
  Bookmark,
  MessageCircleQuestion,
  ShieldCheckIcon,
  Crown,
  Hammer,
  Eye,
} from "lucide-react";
import Avatar from "../user/Avatar";
import { useAuthStore } from "@/store/auth";
import { useLogoutStore } from "@/store/logout";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";

interface UserDropdownProps {
  user: any;
}

export default function UserDropdown({ user }: UserDropdownProps) {
  const router = useRouter();
  const t = useTranslations("common");
  const { logout, isLoading } = useLogoutStore();
  const { user: currentUser } = useAuthStore();

  // 使用 Store 中的最新用户数据
  const displayUser = currentUser || user;

  // 处理登出
  const handleLogout = async () => {
    try {
      await logout();
      toast.success(t("logout_success"));
      router.push("/");
      router.refresh();
    } catch (error) {
      toast.error(t("logout_failed"));
    }
  };

  // 获取后台入口配置
  const getDashboardConfig = () => {
    const role = displayUser?.role;
    
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
    
    return null;
  };

  const dashboardConfig = getDashboardConfig();

  // 是否有管理权限
  const hasManagementAccess = 
    displayUser?.role === "admin" || 
    displayUser?.role === "super_admin" || 
    displayUser?.role === "moderator" || 
    displayUser?.role === "reviewer";

  return (
    <div className="dropdown dropdown-end">
      <div
        tabIndex={0}
        role="button"
        className="btn btn-ghost btn-circle avatar hover:ring-2 hover:ring-primary/20 transition-all"
      >
        <div className="w-9 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
          <Avatar
            username={displayUser?.username}
            avatarUrl={displayUser?.avatar}
            size="md"
          />
        </div>
      </div>
      
      <ul
        tabIndex={0}
        className="dropdown-content menu bg-base-100 rounded-box z-10 w-64 p-2 shadow-xl border border-base-200 mt-2"
      >
        {/* 用户信息 */}
        <li className="menu-title">
          <div className="flex items-center gap-2">
            <div className="avatar placeholder">
              <div className="w-10 rounded-full bg-primary/10">
                {displayUser?.avatar ? (
                  <Avatar 
                    username={displayUser?.username} 
                    avatarUrl={displayUser?.avatar} 
                    size="sm" 
                  />
                ) : (
                  <span className="text-primary font-medium">
                    {displayUser?.username?.[0]?.toUpperCase()}
                  </span>
                )}
              </div>
            </div>
            <div className="flex-1 min-w-0">
              <span className="text-base-content font-medium block truncate">
                {displayUser?.username}
              </span>
              <span className="text-xs text-base-content/50 truncate block">
                {displayUser?.email}
              </span>
              {/* 角色标签 */}
              {displayUser?.role && displayUser?.role !== "user" && (
                <span className={`badge badge-xs mt-1 ${
                  displayUser?.role === "super_admin" ? "badge-error" :
                  displayUser?.role === "admin" ? "badge-warning" :
                  displayUser?.role === "moderator" ? "badge-secondary" :
                  displayUser?.role === "reviewer" ? "badge-accent" :
                  "badge-ghost"
                }`}>
                  {t(`role.${displayUser?.role}`)}
                </span>
              )}
            </div>
          </div>
        </li>

        <div className="divider my-1"></div>

        {/* 个人统计 */}
        <li className="px-2 py-1">
          <div className="flex justify-between text-sm">
            <span className="text-base-content/60">{t("score")}</span>
            <span className="font-bold text-primary">{displayUser?.score || 0}</span>
          </div>
          <div className="flex justify-between text-sm mt-1">
            <span className="text-base-content/60">{t("followers")}</span>
            <span className="font-bold">{displayUser?.followers_count || 0}</span>
          </div>
        </li>

        <div className="divider my-1"></div>

        {/* 快速链接 */}
        <li>
          <Link href={`/users/${displayUser?.id}`} className="gap-2">
            <User className="w-4 h-4" />
            {t("profile")}
          </Link>
        </li>
        <li>
          <Link href="/timeline" className="gap-2">
            <Sparkles className="w-4 h-4" />
            {t("my_timeline")}
          </Link>
        </li>
        <li>
          <Link href="/topics/my" className="gap-2">
            <Bookmark className="w-4 h-4" />
            {t("my_topics")}
          </Link>
        </li>
        <li>
          <Link href="/questions/my" className="gap-2">
            <MessageCircleQuestion className="w-4 h-4" />
            {t("my_questions")}
          </Link>
        </li>

        <div className="divider my-1"></div>

        {/* 版主申请入口 */}
        {displayUser && (
          <li>
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
            <li>
              <Link href={dashboardConfig.path} className={`gap-2 ${dashboardConfig.className}`}>
                {dashboardConfig.icon}
                {dashboardConfig.label}
              </Link>
            </li>
          </>
        )}

        <div className="divider my-1"></div>

        {/* 设置和帮助 */}
        <li>
          <Link href="/settings" className="gap-2">
            <Settings className="w-4 h-4" />
            {t("settings")}
          </Link>
        </li>
        <li>
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
    </div>
  );
}