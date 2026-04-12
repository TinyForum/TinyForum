"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import { notificationApi, timelineApi } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import {
  Bell,
  PenSquare,
  Search,
  LogOut,
  User,
  LayoutDashboard,
  Trophy,
  Menu,
  Home,
  TrendingUp,
  Bookmark,
  Settings,
  HelpCircle,
  TelescopeIcon,
  Sparkles,
  MessageCircleQuestion,
  LayoutGrid,
  Compass,
} from "lucide-react";
import { useState, useEffect, useRef } from "react";
import Image from "next/image";
import Avatar from "../user/Avatar";
import { useTranslations } from "next-intl";
import LanguageSwitcher from "../LanguageSwitcher";
import SearchBar from "../nav/SearchBar";
import MobileMenu from "../nav/MobileMenu";
import UserDropdown from "../nav/UserDropdown";
import NotificationBell from "../nav/NotificationBell";
import NavLinks from "../nav/NavLinks";
import QuickActions from "../nav/QuickActions";

// 导航标签配置
export const NAV_ITEMS = [
  { key: "home", href: "/", icon: Home, requiresAuth: false },
  { key: "explore", href: "/explore", icon: Compass, requiresAuth: false },
  { key: "boards", href: "/boards", icon: LayoutGrid, requiresAuth: false },
  { key: "questions", href: "/questions", icon: MessageCircleQuestion, requiresAuth: false },
  { key: "topics", href: "/topics", icon: Bookmark, requiresAuth: false },
  { key: "timeline", href: "/timeline", icon: Sparkles, requiresAuth: true },
  { key: "leaderboard", href: "/leaderboard", icon: Trophy, requiresAuth: false },
] as const;

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const router = useRouter();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const searchInputRef = useRef<HTMLInputElement>(null);

  const t = useTranslations("nav");

  // 获取未读通知数
  const { data: unreadData } = useQuery({
    queryKey: ["notifications", "unread"],
    queryFn: () => notificationApi.unreadCount().then((r) => r.data.data),
    enabled: isAuthenticated,
    refetchInterval: 30000,
  });

  // 获取未读时间线更新数
  const { data: timelineData } = useQuery({
    queryKey: ["timeline", "unread"],
    queryFn: () => timelineApi.getFollowing({ page: 1, page_size: 1 }).then((r) => r.data.data),
    enabled: isAuthenticated,
    refetchInterval: 60000,
  });

  const unreadCount = unreadData?.count ?? 0;
  const timelineUpdateCount = timelineData?.total ?? 0;

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery("");
    }
  };

  const handleLogout = () => {
    logout();
    router.push("/");
    setIsMobileMenuOpen(false);
  };

  // 过滤显示的导航项
  const visibleNavItems = NAV_ITEMS.filter(
    (item) => !item.requiresAuth || (item.requiresAuth && isAuthenticated)
  );

  // 点击外部关闭移动端菜单
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      if (!target.closest('.mobile-menu')) {
        setIsMobileMenuOpen(false);
      }
    };

    if (isMobileMenuOpen) {
      document.addEventListener("mousedown", handleClickOutside);
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "unset";
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.body.style.overflow = "unset";
    };
  }, [isMobileMenuOpen]);

  return (
  <>
  <nav className="navbar bg-base-100/95 backdrop-blur-sm shadow-sm sticky top-0 z-50 border-b border-base-300 transition-all duration-200">
    <div className="container mx-auto max-w-8xl px-4 w-full">
      {/* 左侧区域：Logo + 汉堡菜单 */}
      <div className="flex items-center gap-2 flex-shrink-0">
        {/* 移动端菜单按钮 */}
        <button
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          className="btn btn-ghost btn-sm btn-square lg:hidden"
          aria-label="菜单"
        >
          <Menu className="w-5 h-5" />
        </button>

        {/* Logo */}
        <Link
          href="/"
          className="mx-2 flex items-center gap-2 text-xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent shrink-0"
        >
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
            <Image src="/assets/brand/logo.svg" width={32} height={32} alt="logo" className="brightness-0 invert" />
          </div>
          <span className="hidden sm:block text-sm">{t("brand")}</span>
        </Link>
      </div>

      {/* 桌面端导航标签 - 独立区域 */}
      <div className="hidden lg:flex items-center gap-1 flex-shrink-0">
        <NavLinks items={visibleNavItems} />
      </div>

      {/* 搜索区域 - 弹性占据剩余空间 */}
      <div className="flex-1 min-w-0 px-4">
        <div className="max-w-md mx-auto">
          <SearchBar
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            onSearch={handleSearch}
          />
        </div>
      </div>

      {/* 右侧区域 - 紧凑排列 */}
      <div className="flex items-center gap-1 flex-shrink-0">
        <LanguageSwitcher />
        
        {/* 快捷操作 */}
        <QuickActions isAuthenticated={isAuthenticated} />

        {isAuthenticated ? (
          <>
            {/* 通知中心 */}
            <NotificationBell unreadCount={unreadCount} />

            {/* 时间线更新提示 */}
            {timelineUpdateCount > 0 && (
              <Link
                href="/timeline"
                className="btn btn-ghost btn-sm btn-circle relative hidden sm:flex flex-shrink-0"
              >
                <Sparkles className="w-5 h-5" />
                <span className="absolute -top-1 -right-1 w-2 h-2 bg-primary rounded-full animate-pulse" />
              </Link>
            )}

            {/* 用户下拉菜单 */}
            <UserDropdown user={user} onLogout={handleLogout} />
          </>
        ) : (
          <div className="flex gap-1 flex-shrink-0">
            <Link href="/auth/login" className="btn btn-ghost btn-sm">
              {t("login")}
            </Link>
            <Link href="/auth/register" className="btn btn-primary btn-sm">
              {t("register")}
            </Link>
          </div>
        )}
      </div>
    </div>
  </nav>

  {/* 移动端侧边菜单 */}
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