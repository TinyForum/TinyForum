"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import { notificationApi } from "@/lib/api";
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
  SearchCheckIcon,
  SearchIcon,
  TelescopeIcon,
} from "lucide-react";
import { useState, useEffect, useRef } from "react";
import Image from "next/image";
import Avatar from "../user/Avatar";
import { useTranslations } from "next-intl";
import LanguageSwitcher from "../LanguageSwitcher";

// 导航标签配置
const NAV_ITEMS = [
  { key: "home", href: "/", icon: Home, requiresAuth: false },
  { key: "explore", href: "/explore", icon: TelescopeIcon, requiresAuth: false },
  { key: "track", href: "/track", icon: TrendingUp, requiresAuth: true },
  { key: "leaderboard", href: "/leaderboard", icon: Trophy, requiresAuth: false },
] as const;

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const searchInputRef = useRef<HTMLInputElement>(null);
  const mobileMenuRef = useRef<HTMLDivElement>(null);

  const t = useTranslations("nav");
  const { data: unreadData } = useQuery({
    queryKey: ["notifications", "unread"],
    queryFn: () => notificationApi.unreadCount().then((r) => r.data.data),
    enabled: isAuthenticated,
    refetchInterval: 30000,
  });

  const unreadCount = unreadData?.count ?? 0;

  // 过滤显示的导航项
  const visibleNavItems = NAV_ITEMS.filter(
    (item) => !item.requiresAuth || (item.requiresAuth && isAuthenticated)
  );

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/posts?keyword=${encodeURIComponent(searchQuery.trim())}`);
      setIsSearchExpanded(false);
      setSearchQuery("");
    }
  };

  const handleLogout = () => {
    logout();
    router.push("/");
    setIsMobileMenuOpen(false);
  };

  // 点击外部关闭移动端菜单
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (mobileMenuRef.current && !mobileMenuRef.current.contains(event.target as Node)) {
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

  // 移动端搜索展开时自动聚焦
  useEffect(() => {
    if (isSearchExpanded && searchInputRef.current) {
      searchInputRef.current.focus();
    }
  }, [isSearchExpanded]);

  return (
    <>
      <nav className="navbar bg-base-100 shadow-sm sticky top-0 z-50 border-b border-base-300">
        <div className="container mx-auto max-w-7xl px-4 w-full">
          {/* 左侧区域：Logo + 汉堡菜单 */}
          <div className="flex items-center gap-2">
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
              className="flex items-center gap-2 text-xl font-bold text-primary shrink-0"
            >
              <div className="w-8 h-8 rounded-lg flex items-center justify-center text-white text-sm font-black">
                <Image src="/logo.svg" width={500} height={500} alt="logo" />
              </div>
              <span className="hidden sm:block">{t("brand")}</span>
            </Link>
          </div>

          {/* 桌面端导航标签（隐藏在小屏幕上） */}
          <div className="hidden lg:flex items-center gap-1 ml-4">
            {visibleNavItems.map((item) => (
              <Link
                key={item.key}
                href={item.href}
                className="btn btn-ghost btn-sm gap-2"
              >
                <item.icon className="w-4 h-4" />
                <span>{t(item.key)}</span>
              </Link>
            ))}
          </div>

          {/* 搜索区域 - 响应式设计 */}
          <div className="flex-1 max-w-md mx-4">
            {/* 移动端：搜索图标 + 展开输入框 */}
            <div className="lg:hidden">
              {!isSearchExpanded ? (
                <button
                  onClick={() => setIsSearchExpanded(true)}
                  className="btn btn-ghost btn-sm btn-circle"
                  aria-label="搜索"
                >
                  <Search className="w-5 h-5" />
                </button>
              ) : (
                <form onSubmit={handleSearch} className="fixed inset-x-0 top-0 z-50 p-4 bg-base-100 shadow-lg animate-slideDown">
                  <div className="flex gap-2">
                    <div className="relative flex-1">
                      <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                      <input
                        ref={searchInputRef}
                        type="text"
                        placeholder={t("search") + "..."}
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="input input-bordered w-full pl-9"
                        autoFocus
                      />
                    </div>
                    <button
                      type="button"
                      onClick={() => setIsSearchExpanded(false)}
                      className="btn btn-ghost btn-sm"
                    >
                      取消
                    </button>
                  </div>
                </form>
              )}
            </div>

            {/* 桌面端：完整搜索框 */}
            <form onSubmit={handleSearch} className="hidden lg:block">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
                <input
                  type="text"
                  placeholder={t("search") + "..."}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="input input-bordered input-sm w-full pl-9 focus:outline-none focus:border-primary"
                />
              </div>
            </form>
          </div>

          {/* 右侧区域：语言切换 + 用户操作 */}
          <div className="flex items-center gap-1 shrink-0">
            <LanguageSwitcher />

            {isAuthenticated ? (
              <>
                {/* 发帖按钮 */}
                <Link
                  href="/posts/new"
                  className="btn btn-primary btn-sm gap-1 hidden md:flex"
                >
                  <PenSquare className="w-4 h-4" />
                  <span className="hidden sm:inline">{t("create_post")}</span>
                </Link>

                {/* 通知 */}
                <Link
                  href="/notifications"
                  className="btn btn-ghost btn-sm btn-circle relative"
                >
                  <Bell className="w-5 h-5" />
                  {unreadCount > 0 && (
                    <span className="badge badge-error badge-xs absolute -top-1 -right-1 min-w-[16px] h-4 text-[10px]">
                      {unreadCount > 99 ? "99+" : unreadCount}
                    </span>
                  )}
                </Link>

                {/* 用户头像下拉菜单 */}
                <div className="dropdown dropdown-end">
                  <div
                    tabIndex={0}
                    role="button"
                    className="btn btn-ghost btn-circle avatar"
                  >
                    <div className="w-9 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                      <Avatar
                        username={user?.username}
                        avatarUrl={user?.avatar}
                        size="md"
                      />
                    </div>
                  </div>
                  <ul
                    tabIndex={0}
                    className="dropdown-content menu bg-base-100 rounded-box z-10 w-56 p-2 shadow-lg border border-base-300 mt-2"
                  >
                    <li className="menu-title">
                      <span className="text-base-content font-medium">
                        {user?.username}
                      </span>
                      <span className="text-xs text-base-content/50 truncate">
                        {user?.email}
                      </span>
                    </li>
                    <div className="divider my-1"></div>
                    <li>
                      <Link href={`/users/${user?.id}`}>
                        <User className="w-4 h-4" />
                        {t("profile")}
                      </Link>
                    </li>
                    <li>
                      <Link href="/settings">
                        <Settings className="w-4 h-4" />
                        {t("settings")}
                      </Link>
                    </li>
                  
                  
                     {user?.role === "admin" && (
                    <li>
                      <Link href="/admin">
                        <LayoutDashboard className="w-4 h-4" /> {t("admin")}
                      </Link>
                    </li>
                  )}
                  
                    <li>
                      <Link href="/help">
                        <HelpCircle className="w-4 h-4" />
                        {t("help")}
                      </Link>
                    </li>
                    <div className="divider my-1"></div>
                    <li>
                      <button onClick={handleLogout} className="text-error">
                        <LogOut className="w-4 h-4" /> {t("logout")}
                      </button>
                    </li>
                  </ul>
                </div>
              </>
            ) : (
              <div className="flex gap-1">
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
      {isMobileMenuOpen && (
        <>
          <div 
            className="fixed inset-0 bg-black/50 z-40 lg:hidden animate-fadeIn"
            onClick={() => setIsMobileMenuOpen(false)}
          />
          <div
            ref={mobileMenuRef}
            className="fixed left-0 top-0 bottom-0 w-72 bg-base-100 z-50 shadow-xl transform transition-transform duration-300 ease-in-out lg:hidden animate-slideRight"
          >
            <div className="p-4 border-b border-base-200">
              <div className="flex items-center justify-between">
                <Link
                  href="/"
                  className="flex items-center gap-2 text-xl font-bold text-primary"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <Image src="/logo.svg" width={32} height={32} alt="logo" />
                  <span>{t("brand")}</span>
                </Link>
                <button
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="btn btn-ghost btn-sm btn-circle"
                >
                  ✕
                </button>
              </div>
            </div>

            <div className="p-4">
              {/* 移动端导航菜单 */}
              <div className="space-y-1">
                {visibleNavItems.map((item) => (
                  <Link
                    key={item.key}
                    href={item.href}
                    onClick={() => setIsMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                  >
                    <item.icon className="w-5 h-5" />
                    <span>{t(item.key)}</span>
                  </Link>
                ))}
                
                <div className="divider my-2"></div>
                
                {isAuthenticated ? (
                  <>
                    <Link
                      href="/posts/new"
                      onClick={() => setIsMobileMenuOpen(false)}
                      className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors text-primary"
                    >
                      <PenSquare className="w-5 h-5" />
                      <span>{t("create_post")}</span>
                    </Link>
                    <Link
                      href="/settings"
                      onClick={() => setIsMobileMenuOpen(false)}
                      className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                    >
                      <Settings className="w-5 h-5" />
                      <span>{t("settings")}</span>
                    </Link>
                    <Link
                      href="/help"
                      onClick={() => setIsMobileMenuOpen(false)}
                      className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                    >
                      <HelpCircle className="w-5 h-5" />
                      <span>{t("help")}</span>
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors text-error"
                    >
                      <LogOut className="w-5 h-5" />
                      <span>{t("logout")}</span>
                    </button>
                  </>
                ) : (
                  <div className="space-y-2 p-3">
                    <Link
                      href="/auth/login"
                      onClick={() => setIsMobileMenuOpen(false)}
                      className="btn btn-ghost w-full"
                    >
                      {t("login")}
                    </Link>
                    <Link
                      href="/auth/register"
                      onClick={() => setIsMobileMenuOpen(false)}
                      className="btn btn-primary w-full"
                    >
                      {t("register")}
                    </Link>
                  </div>
                )}
              </div>
            </div>
          </div>
        </>
      )}

  
    </>
  );
}