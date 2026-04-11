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
} from "lucide-react";
import { useState } from "react";
import Image from "next/image";
import Avatar from "../user/Avatar";
import { useTranslations } from "next-intl";
import LanguageSwitcher from "../LanguageSwitcher";

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuthStore();
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");

  const t = useTranslations("nav");
  const { data: unreadData } = useQuery({
    queryKey: ["notifications", "unread"],
    queryFn: () => notificationApi.unreadCount().then((r) => r.data.data),
    enabled: isAuthenticated,
    refetchInterval: 30000,
  });

  const unreadCount = unreadData?.count ?? 0;

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/posts?keyword=${encodeURIComponent(searchQuery.trim())}`);
    }
  };

  const handleLogout = () => {
    logout();
    router.push("/");
  };

  return (
    <div className="navbar bg-base-100 shadow-sm sticky top-0 z-50 border-b border-base-300">
      <div className="container mx-auto max-w-6xl px-4 w-full">
        {/* Logo */}
        <div className="flex-none mr-4">
          <Link
            href="/"
            className="flex items-center gap-2 text-xl font-bold text-primary"
          >
            <div className="w-8 h-8 rounded-lg flex items-center justify-center text-white text-sm font-black">
              <Image src="/logo.svg" width={500} height={500} alt="logo" />
            </div>
            <span className="hidden sm:block">{t("brand")}</span>
          </Link>
        </div>

        {/* Search */}
        <div className="flex-1 max-w-md">
          <form onSubmit={handleSearch}>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-base-content/40" />
              <input
                type="text"
                placeholder= {t('search')+"..."}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="input input-bordered input-sm w-full pl-9 focus:outline-none focus:border-primary"
              />
            </div>
          </form>
        </div>
  <LanguageSwitcher/>
        {/* Right side */}
        <div className="flex-none flex items-center gap-2 ml-4">
          {isAuthenticated ? (
            <>
              {/* Write post button */}
              <Link
                href="/posts/new"
                className="btn btn-primary btn-sm gap-1 hidden sm:flex"
              >
                <PenSquare className="w-4 h-4" />
                {t("create_post")}
              </Link>

              {/* Notifications */}
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
            

              {/* User dropdown */}
              <div className="dropdown dropdown-end">
                <div
                  tabIndex={0}
                  role="button"
                  className="btn btn-ghost btn-circle avatar"
                >
                  <div className="w-9 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                    <Avatar
                      username={user?.username}
                      avatarUrl={user?.avatar} // 数据库中的头像
                      size="md"
                    />
                  </div>
                </div>
                <ul
                  tabIndex={0}
                  className="dropdown-content menu bg-base-100 rounded-box z-10 w-52 p-2 shadow-lg border border-base-300 mt-2"
                >
                  <li className="menu-title">
                    <span className="text-base-content font-medium">
                      {user?.username}
                    </span>
                    <span className="text-xs text-base-content/50">
                      {user?.email}
                    </span>
                  </li>
                  <div className="divider my-1"></div>
                  <li>
                    <Link href={`/users/${user?.id}`}>
                      <User className="w-4 h-4" />{t('profile')}
                    </Link>
                  </li>
                  <li>
                    <Link href="/leaderboard">
                      <Trophy className="w-4 h-4" />{t('leaderboard')}
                    </Link>
                  </li>
                  {user?.role === "admin" && (
                    <li>
                      <Link href="/admin">
                        <LayoutDashboard className="w-4 h-4" /> {t('admin')}
                      </Link>
                    </li>
                  )}
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
            <div className="flex gap-2">
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
    </div>
  );
}
