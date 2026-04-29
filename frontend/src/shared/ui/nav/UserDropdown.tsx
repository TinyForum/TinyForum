"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { Fragment, useEffect, useState } from "react";
import {
  Dialog,
  DialogPanel,
  DialogTitle,
  Menu,
  MenuButton,
  MenuItem,
  MenuItems,
  Transition,
} from "@headlessui/react";
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
import { useLogoutStore } from "@/store/logout";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { User } from "@/shared/api";
import Avatar from "@/features/user/components/Avatar";
import { createPortal } from "react-dom";

interface UserDropdownProps {
  user: User;
  onOpenChange?: (isOpen: boolean) => void;
}

export default function UserDropdown({
  user,
  onOpenChange,
}: UserDropdownProps) {
  const router = useRouter();
  const t = useTranslations("Common");
  const { logout, isLoading } = useLogoutStore();
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);

  // 点击登出按钮时打开模态框，而不是直接登出
  const handleLogoutClick = () => {
    if (isLoading) return;
    setShowLogoutConfirm(true);
  };

  // 实际执行登出的函数
  const performLogout = async () => {
    try {
      await logout();
      toast.success(t("logout_success"));
      router.push("/");
      router.refresh();
    } catch {
      toast.error(t("logout_failed"));
    } finally {
      setShowLogoutConfirm(false);
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      toast.success(t("logout_success"));
      router.push("/");
      router.refresh();
    } catch {
      toast.error(t("logout_failed"));
    }
  };

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
  const hasManagementAccess =
    user.role === "super_admin" ||
    user.role === "admin" ||
    user.role === "reviewer" ||
    user.role === "moderator" ||
    user.role === "member" ||
    user.role === "user";

  const menuSections = [
    { id: "user-info", type: "info" as const, condition: true },
    { id: "divider-1", type: "divider" as const, condition: true },
    {
      id: "profile",
      type: "link" as const,
      icon: <UserIcon className="w-4 h-4" />,
      label: t("profile"),
      href: `/users/${user.id}`,
      condition: true,
    },
    {
      id: "timeline",
      type: "link" as const,
      icon: <Sparkles className="w-4 h-4" />,
      label: t("my_timeline"),
      href: "/timeline/me",
      condition: true,
    },
    {
      id: "topics",
      type: "link" as const,
      icon: <Bookmark className="w-4 h-4" />,
      label: t("my_topics"),
      href: "/topics/me",
      condition: true,
    },
    {
      id: "questions",
      type: "link" as const,
      icon: <MessageCircleQuestion className="w-4 h-4" />,
      label: t("my_questions"),
      href: "/questions/me",
      condition: true,
    },
    { id: "divider-2", type: "divider" as const, condition: true },
    {
      id: "applications",
      type: "link" as const,
      icon: <ShieldCheckIcon className="w-4 h-4" />,
      label: t("my_moderator_applications"),
      href: "/boards/applications",
      condition: !!user.role,
    },
    {
      id: "divider-3",
      type: "divider" as const,
      condition: hasManagementAccess && !!dashboardConfig,
    },
    {
      id: "dashboard",
      type: "link" as const,
      icon: dashboardConfig?.icon,
      label: dashboardConfig?.label,
      href: dashboardConfig?.path,
      className: dashboardConfig?.className,
      condition: hasManagementAccess && !!dashboardConfig,
    },
    { id: "divider-4", type: "divider" as const, condition: true },
    {
      id: "settings",
      type: "link" as const,
      icon: <Settings className="w-4 h-4" />,
      label: t("settings"),
      href: "/settings",
      condition: true,
    },
    {
      id: "help",
      type: "link" as const,
      icon: <HelpCircle className="w-4 h-4" />,
      label: t("help_center"),
      href: "/help",
      condition: true,
    },
    { id: "divider-5", type: "divider" as const, condition: true },
    {
      id: "logout",
      type: "button" as const,
      icon: isLoading ? (
        <span className="loading loading-spinner loading-xs" />
      ) : (
        <LogOut className="w-4 h-4" />
      ),
      label: isLoading ? t("logging_out") : t("logout"),
      onClick: handleLogoutClick, // 改为打开模态框
      className: "text-error hover:bg-error/10",
      disabled: isLoading,
      condition: true,
    },
  ];

  return (
    <div className="relative bg-primary/10 rounded-lg">
      <Menu>
        {({ open }: { open: boolean }) => {
          useEffect(() => {
            onOpenChange?.(open);
          }, [open]);

          return (
            <>
              <MenuButton
                className="btn btn-ghost btn-circle avatar transition-all duration-200 hover:scale-105 hover:ring-2 hover:ring-primary/30 focus:ring-2 focus:ring-primary/40"
                aria-label={t("user_menu")}
              >
                <div className="w-9 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                  <Avatar
                    username={user.username}
                    avatarUrl={user.avatar}
                    size="md"
                  />
                </div>
              </MenuButton>

              <Transition
                as={Fragment}
                show={open}
                enter="transition ease-out duration-150"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-100"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <MenuItems className="absolute right-0 top-full mt-2 w-80 origin-top-right overflow-hidden rounded-2xl bg-white shadow-2xl border border-gray-100 focus:outline-none z-50">
                  <div className="max-h-[calc(100vh-6rem)] overflow-y-auto scrollbar-thin scrollbar-thumb-rounded-full scrollbar-thumb-gray-300">
                    <div className="p-2 space-y-1">
                      {menuSections.map((section) => {
                        if (!section.condition) return null;

                        switch (section.type) {
                          case "info":
                            return (
                              <div
                                key={section.id}
                                className="relative mb-2 rounded-xl bg-gradient-to-br from-primary/5 to-base-200/30 p-3"
                              >
                                <div className="flex items-center gap-3">
                                  <div className="avatar">
                                    <div className="w-12 rounded-full ring-2 ring-primary/20 ring-offset-2 ring-offset-base-100">
                                      <Avatar
                                        username={user?.username}
                                        avatarUrl={user?.avatar}
                                        size="md"
                                      />
                                    </div>
                                  </div>
                                  <div className="flex-1 min-w-0">
                                    <p className="text-base-content font-semibold truncate">
                                      {user?.username}
                                    </p>
                                    <p className="text-xs text-base-content/60 truncate">
                                      {user?.email}
                                    </p>
                                    {user.role !== "user" && (
                                      <span
                                        className={`badge badge-xs mt-1.5 ${
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
                              </div>
                            );

                          case "divider":
                            return (
                              <div
                                key={section.id}
                                className="my-2 border-t border-base-300/50"
                              />
                            );

                          // 链接
                          case "link":
                            return (
                              <MenuItem key={section.id}>
                                {({ focus }: { focus: boolean }) => (
                                  <li className="block">
                                    <Link
                                      href={section.href!}
                                      className={`
          group flex w-full items-center gap-3 rounded-xl px-4 py-2.5 
          text-sm font-medium transition-all duration-200
          ${
            focus
              ? "bg-primary/10 text-primary shadow-sm"
              : "text-base-content/80 hover:bg-gray-100 hover:text-base-content hover:shadow-sm"
          }
          ${section.className || ""}
        `}
                                    >
                                      <span className="shrink-0 transition-transform duration-200 group-hover:scale-105 group-hover:translate-x-0.5">
                                        {section.icon}
                                      </span>
                                      <span className="flex-1 text-left">
                                        {section.label}
                                      </span>
                                    </Link>
                                  </li>
                                )}
                              </MenuItem>
                            );

                          case "button":
                            return (
                              <MenuItem
                                key={section.id}
                                disabled={section.disabled}
                              >
                                {({ focus }: { focus: boolean }) => (
                                  <li className="block">
                                    <button
                                      type="button"
                                      onClick={section.onClick}
                                      disabled={section.disabled}
                                      className={`
          group flex w-full items-center gap-3 rounded-xl px-4 py-2.5 
          text-sm font-medium transition-all duration-200
          ${
            focus && !section.disabled
              ? "bg-red-50 text-red-700 shadow-sm ring-1 ring-red-200"
              : "text-red-600 hover:bg-red-50 hover:text-red-700 hover:shadow-sm"
          }
          ${
            section.disabled
              ? "cursor-not-allowed opacity-50"
              : "active:scale-[0.98]"
          }
          ${section.className || ""}
        `}
                                    >
                                      <span className="shrink-0 transition-transform duration-200 group-hover:scale-105 group-hover:translate-x-0.5">
                                        {section.icon}
                                      </span>
                                      <span className="flex-1 text-left">
                                        {section.label}
                                      </span>
                                    </button>
                                  </li>
                                )}
                              </MenuItem>
                            );

                          default:
                            return null;
                        }
                      })}
                    </div>
                  </div>
                </MenuItems>
              </Transition>
            </>
          );
        }}
      </Menu>

      {typeof window !== "undefined" &&
        createPortal(
          <Transition appear show={showLogoutConfirm} as={Fragment}>
            <Dialog
              className="relative z-[99999]"
              onClose={() => setShowLogoutConfirm(false)}
            >
              {/* 背景遮罩 */}
              <Transition.Child
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0"
                enterTo="opacity-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
              </Transition.Child>

              {/* 居中容器：用 items-center justify-center 确保垂直水平居中 */}
              <div className="fixed inset-0 flex items-center justify-center p-4">
                <Transition.Child
                  as={Fragment}
                  enter="ease-out duration-300"
                  enterFrom="opacity-0 scale-95"
                  enterTo="opacity-100 scale-100"
                  leave="ease-in duration-200"
                  leaveFrom="opacity-100 scale-100"
                  leaveTo="opacity-0 scale-95"
                >
                  <DialogPanel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                    <DialogTitle className="text-lg font-semibold leading-6 text-gray-900">
                      {t("confirm_logout")}
                    </DialogTitle>
                    <div className="mt-2">
                      <p className="text-sm text-gray-600">
                        {t("confirm_logout_message")}
                      </p>
                    </div>
                    <div className="mt-6 flex justify-end gap-3">
                      <button
                        type="button"
                        className="inline-flex justify-center rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary/20 transition"
                        onClick={() => setShowLogoutConfirm(false)}
                      >
                        {t("cancel")}
                      </button>
                      <button
                        type="button"
                        className="inline-flex justify-center rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500/20 transition"
                        onClick={performLogout}
                      >
                        {t("logout")}
                      </button>
                    </div>
                  </DialogPanel>
                </Transition.Child>
              </div>
            </Dialog>
          </Transition>,
          document.body, // ← 挂载到 body，完全脱离父级 stacking context
        )}
    </div>
  );
}
