"use client";

import Link from "next/link";
import Image from "next/image";
import { Menu } from "lucide-react";
import { useTranslations } from "next-intl";

import SearchBar from "@/features/moderator/components/SearchBar";
import NotificationBell from "@/shared/ui/nav/NotificationBell";
import QuickActions from "@/shared/ui/nav/QuickActions";
import UserDropdown from "@/shared/ui/nav/UserDropdown";
import LanguageSwitcher from "./LanguageSwitcher";
import { UserDO } from "@/shared/api/types/user.model";
// User 类型直接从项目 store 引入，与 UserDropdown 期望的类型保持一致

// ─────────────────────────────────────────────
// 仅在 lg 以下断点渲染（由父级 Navbar 通过 isDesktop 判断控制）
// 结构：[汉堡] [Logo 图标] [··弹性搜索··] [右侧操作]
// 导航标签通过汉堡菜单（MobileMenu 侧边栏）展示
// ─────────────────────────────────────────────

interface MobileNavbarProps {
  isAuthenticated: boolean;
  // 修复：使用项目已有的 User 类型，不再使用自造的 AuthUser
  user: UserDO | null;
  searchQuery: string;
  onSearchQueryChange: (value: string) => void;
  unreadCount: number;
  isMobileMenuOpen: boolean;
  onMobileMenuToggle: () => void;
  isUserDropdownOpen: boolean;
  onUserDropdownOpenChange: (open: boolean) => void;
}

export default function MobileNavbar({
  isAuthenticated,
  user,
  searchQuery,
  onSearchQueryChange,
  unreadCount,
  isMobileMenuOpen,
  onMobileMenuToggle,
  onUserDropdownOpenChange,
}: MobileNavbarProps) {
  const t = useTranslations("Nav");

  return (
    // 父级 Navbar 已通过 isDesktop 保证只在移动端渲染此组件，无需断点 class
    <div className="flex items-center w-full gap-2">
      {/* ── 汉堡按钮（触发侧边菜单） ── */}
      <button
        onClick={onMobileMenuToggle}
        className="btn btn-ghost btn-sm btn-square flex-shrink-0"
        aria-label="菜单"
        aria-expanded={isMobileMenuOpen}
      >
        <Menu className="w-5 h-5" />
      </button>

      {/* ── Logo（仅图标，节省横向空间） ── */}
      <Link href="/" className="flex items-center flex-shrink-0">
        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
          <Image
            src="/assets/brand/logo.svg"
            width={32}
            height={32}
            alt="logo"
            className="brightness-0 invert"
          />
        </div>
      </Link>

      {/* ── 搜索框（弹性占据剩余空间） ── */}
      <div className="flex-1 min-w-0">
        <SearchBar
          keyword={searchQuery}
          onKeywordChange={onSearchQueryChange}
        />
      </div>

      {/* ── 右侧操作（紧凑排列） ── */}
      <div className="flex items-center gap-1 flex-shrink-0">
        <LanguageSwitcher />
        <QuickActions isAuthenticated={isAuthenticated} />

        {isAuthenticated ? (
          <>
            <NotificationBell unreadCount={unreadCount} />

            {/* timeline 气泡收入侧边菜单，顶栏不占位 */}

            {user && (
              <UserDropdown
                user={user}
                // isOpen={isUserDropdownOpen}
                onOpenChange={onUserDropdownOpenChange}
              />
            )}
          </>
        ) : (
          // 移动端仅保留登录按钮，注册收入侧边菜单
          <Link href="/auth/login" className="btn btn-ghost btn-sm">
            {t("login")}
          </Link>
        )}
      </div>
    </div>
  );
}
