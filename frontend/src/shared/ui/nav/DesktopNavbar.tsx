"use client";

import Link from "next/link";
import Image from "next/image";
import { Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { User } from "@/shared/api";
import SearchBar from "@/features/moderator/components/SearchBar";
import NavLinks from "@/shared/ui/nav/NavLinks";
import NotificationBell from "@/shared/ui/nav/NotificationBell";
import QuickActions from "@/shared/ui/nav/QuickActions";
import { NavItem } from "@/shared/ui/nav/types";
import UserDropdown from "@/shared/ui/nav/UserDropdown";
import LanguageSwitcher from "./LanguageSwitcher";

interface DesktopNavbarProps {
  navItems: NavItem[];
  isAuthenticated: boolean;
  user: User | null;
  searchQuery: string;
  onSearchQueryChange: (value: string) => void;
  unreadCount: number;
  timelineUpdateCount: number;
  isUserDropdownOpen: boolean;
  onUserDropdownOpenChange: (open: boolean) => void;
}

export default function DesktopNavbar({
  navItems,
  isAuthenticated,
  user,
  searchQuery,
  onSearchQueryChange,
  unreadCount,
  timelineUpdateCount,
  onUserDropdownOpenChange,
}: DesktopNavbarProps) {
  const t = useTranslations("Nav");

  return (
    <div className="flex items-center w-full gap-2">
      {/* ── Logo + 品牌名 ── */}
      <Link
        href="/"
        className="flex items-center gap-2 text-xl font-bold bg-primary bg-clip-text text-transparent flex-shrink-0 mx-2"
      >
        <div className="w-8 h-8 rounded-lg flex items-center justify-center">
          <Image
            src="/assets/brand/logo.svg"
            width={32}
            height={32}
            alt="logo"
            className=""
          />
        </div>
        <span className="text-sm">{t("brand")}</span>
      </Link>

      {/* ── 导航标签 ── */}
      <div className="flex items-center gap-1 flex-shrink-0">
        <NavLinks items={navItems} />
      </div>

      {/* ── 搜索框（弹性居中） ── */}
      <div className="flex-1 min-w-0 px-4">
        <div className="max-w-md mx-auto">
          <SearchBar
            keyword={searchQuery}
            onKeywordChange={onSearchQueryChange}
          />
        </div>
      </div>

      {/* ── 右侧操作 ── */}
      <div className="flex items-center gap-1 flex-shrink-0">
        <LanguageSwitcher />
        <QuickActions isAuthenticated={isAuthenticated} />

        {isAuthenticated ? (
          <>
            <NotificationBell unreadCount={unreadCount} />

            {timelineUpdateCount > 0 && (
              <Link
                href="/timeline"
                className="btn btn-ghost btn-sm btn-circle relative"
              >
                <Sparkles className="w-5 h-5" />
                <span className="absolute -top-1 -right-1 w-2 h-2 bg-primary rounded-full animate-pulse" />
              </Link>
            )}

            {user && (
              <UserDropdown
                user={user}
                onOpenChange={onUserDropdownOpenChange}
              />
            )}
          </>
        ) : (
          <>
            <Link href="/auth/login" className="btn btn-ghost btn-sm">
              {t("login")}
            </Link>
            <Link href="/auth/register" className="btn btn-primary btn-sm">
              {t("register")}
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
