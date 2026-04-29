"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store";
import { notificationApi, timelineApi } from "@/shared/api";
import { useQuery } from "@tanstack/react-query";
import {
  Trophy,
  Home,
  Bookmark,
  Sparkles,
  MessageCircleQuestion,
  LayoutGrid,
  Compass,
} from "lucide-react";
import { useState, useEffect } from "react";
import { useTranslations } from "next-intl";
import MobileMenu from "../nav/MobileMenu";
import DesktopNavbar from "./DesktopNavbar";
import MobileNavbar from "./MobileNavbar";
import { useMediaQuery } from "@/features/common/hooks/useMediaQuery";

// ─────────────────────────────────────────────
// Navbar — 入口组件，只负责：
//   1. 共享状态管理（auth、查询、路由）
//   2. 将 props 分发给 DesktopNavbar / MobileNavbar
//   3. 渲染移动端侧边菜单（MobileMenu）
//
// 不含任何断点 class，断点由子组件自己声明
// ─────────────────────────────────────────────

export default function Navbar() {
  const isDesktop = useMediaQuery("(min-width: 1024px)"); // lg = 1024px
  const { user, isAuthenticated, logout } = useAuthStore();
  const router = useRouter();
  const t = useTranslations("Nav");

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [isUserDropdownOpen, setIsUserDropdownOpen] = useState(false);
  useEffect(() => {
    console.log("isUserDropdownOpen", isUserDropdownOpen);
  }, [isUserDropdownOpen]);

  // ── 数据查询 ──────────────────────────────────
  const { data: unreadData } = useQuery({
    queryKey: ["notifications", "unread"],
    queryFn: () => notificationApi.unreadCount().then((r) => r.data.data),
    enabled: isAuthenticated,
    refetchInterval: 30000,
  });

  const { data: timelineData } = useQuery({
    queryKey: ["timeline", "unread"],
    queryFn: () =>
      timelineApi
        .getFollowing({ page: 1, page_size: 1 })
        .then((r) => r.data.data),
    enabled: isAuthenticated,
    refetchInterval: 60000,
  });

  const unreadCount = unreadData?.count ?? 0;
  const timelineUpdateCount = timelineData?.total ?? 0;

  // ── 导航项配置 ────────────────────────────────
  const NAV_ITEMS = [
    {
      key: "home",
      name: t("home"),
      href: "/",
      icon: Home,
      requiresAuth: false,
    },
    {
      key: "explore",
      name: t("explore"),
      href: "/explore",
      icon: Compass,
      requiresAuth: false,
    },
    {
      key: "boards",
      name: t("boards"),
      href: "/boards",
      icon: LayoutGrid,
      requiresAuth: false,
    },
    {
      key: "questions",
      name: t("questions"),
      href: "/questions",
      icon: MessageCircleQuestion,
      requiresAuth: false,
    },
    {
      key: "topics",
      name: t("topics"),
      href: "/topics",
      icon: Bookmark,
      requiresAuth: false,
    },
    {
      key: "timeline",
      name: t("timeline"),
      href: "/timeline",
      icon: Sparkles,
      requiresAuth: true,
    },
    {
      key: "leaderboard",
      name: t("leaderboard"),
      href: "/leaderboard",
      icon: Trophy,
      requiresAuth: false,
    },
  ] as const;

  const visibleNavItems = NAV_ITEMS.filter(
    (item) => !item.requiresAuth || isAuthenticated,
  );

  // ── 事件处理 ──────────────────────────────────
  const handleLogout = () => {
    logout();
    router.push("/");
    setIsMobileMenuOpen(false);
  };

  // 移动端菜单打开时锁定 body 滚动，点击外部关闭
  useEffect(() => {
    if (!isMobileMenuOpen) return;

    const handleClickOutside = (e: MouseEvent) => {
      if (!(e.target as HTMLElement).closest(".mobile-menu")) {
        setIsMobileMenuOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    document.body.style.overflow = "hidden";

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.body.style.overflow = "unset";
    };
  }, [isMobileMenuOpen]);

  // ── 共享 props（两端复用部分） ─────────────────
  const sharedProps = {
    isAuthenticated,
    user,
    searchQuery,
    onSearchQueryChange: setSearchQuery,
    unreadCount,
    isUserDropdownOpen,
    onUserDropdownOpenChange: setIsUserDropdownOpen,
  };

  return (
    <>
      <nav className="navbar bg-base-100/95 backdrop-blur-sm shadow-sm sticky top-0 border-b border-base-300 transition-all duration-200">
        <div className="container mx-auto max-w-none px-4 w-full flex items-center">
          {/* JS 级别二选一，不靠 CSS display 切换 */}
          {isDesktop ? (
            <DesktopNavbar
              {...sharedProps}
              navItems={visibleNavItems}
              timelineUpdateCount={timelineUpdateCount}
            />
          ) : (
            <MobileNavbar
              {...sharedProps}
              isMobileMenuOpen={isMobileMenuOpen}
              onMobileMenuToggle={() => setIsMobileMenuOpen((v) => !v)}
            />
          )}
        </div>
      </nav>

      <MobileMenu
        isOpen={isMobileMenuOpen}
        onClose={() => setIsMobileMenuOpen(false)}
        navItems={visibleNavItems}
        isAuthenticated={isAuthenticated}
        user={user}
        onLogout={handleLogout}
        unreadCount={unreadCount}
      />
    </>
  );
}
