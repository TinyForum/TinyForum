"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  User,
  Settings,
  HelpCircle,
  LogOut,
  LayoutDashboard,
  PenSquare,
  X,
  Sparkles,
  Bookmark,
  MessageCircleQuestion,
  LayoutGrid,
} from "lucide-react";
import Image from "next/image";
import Avatar from "../user/Avatar";

interface MobileMenuProps {
  isOpen: boolean;
  onClose: () => void;
  navItems: any[];
  isAuthenticated: boolean;
  user: any;
  onLogout: () => void;
  unreadCount: number;
}

export default function MobileMenu({
  isOpen,
  onClose,
  navItems,
  isAuthenticated,
  user,
  onLogout,
  unreadCount,
}: MobileMenuProps) {
  const pathname = usePathname();

  if (!isOpen) return null;

  const isActive = (href: string) => {
    if (href === "/") return pathname === href;
    return pathname.startsWith(href);
  };

  return (
    <>
      {/* 背景遮罩 */}
      <div
        className="fixed inset-0 bg-black/50 z-40 lg:hidden animate-fadeIn"
        onClick={onClose}
      />
      
      {/* 侧边菜单 */}
      <div className="mobile-menu fixed left-0 top-0 bottom-0 w-80 bg-base-100 z-50 shadow-2xl transform transition-transform duration-300 ease-out animate-slideRight flex flex-col">
        {/* 头部 */}
        <div className="p-4 border-b border-base-200 bg-gradient-to-r from-primary/5 to-secondary/5">
          <div className="flex items-center justify-between">
            <Link
              href="/"
              className="flex items-center gap-2 text-xl font-bold text-primary"
              onClick={onClose}
            >
              <Image src="/assets/brand/logo.svg" width={32} height={32} alt="logo" />
              <span>微论坛</span>
            </Link>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
          
          {/* 用户信息（如果已登录） */}
          {isAuthenticated && user && (
            <Link
              href={`/users/${user.id}`}
              onClick={onClose}
              className="flex items-center gap-3 mt-4 p-2 rounded-lg hover:bg-base-200 transition-colors"
            >
              <Avatar username={user.username} avatarUrl={user.avatar} size="md" />
              <div className="flex-1">
                <div className="font-medium">{user.username}</div>
                <div className="text-xs text-base-content/50">
                  积分: {user.score || 0}
                </div>
              </div>
            </Link>
          )}
        </div>

        {/* 导航菜单 */}
        <div className="flex-1 overflow-y-auto p-4">
          <div className="space-y-1">
            <div className="text-xs font-semibold text-base-content/50 mb-2 px-3">
              导航
            </div>
            {navItems.map((item) => {
              const active = isActive(item.href);
              const Icon = item.icon;
              return (
                <Link
                  key={item.key}
                  href={item.href}
                  onClick={onClose}
                  className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-all ${
                    active
                      ? "bg-primary/10 text-primary font-medium"
                      : "hover:bg-base-200"
                  }`}
                >
                  <Icon className={`w-5 h-5 ${active ? "text-primary" : ""}`} />
                  <span>{item.key.charAt(0).toUpperCase() + item.key.slice(1)}</span>
                  {active && (
                    <span className="ml-auto w-1.5 h-1.5 bg-primary rounded-full" />
                  )}
                </Link>
              );
            })}
            
            <div className="divider my-3"></div>
            
            {/* 快捷操作 */}
            <div className="text-xs font-semibold text-base-content/50 mb-2 px-3">
              快捷操作
            </div>
            
            {isAuthenticated ? (
              <>
                <Link
                  href="/posts/new"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors text-primary"
                >
                  <PenSquare className="w-5 h-5" />
                  <span>写帖子</span>
                </Link>
                <Link
                  href="/questions/ask"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <MessageCircleQuestion className="w-5 h-5" />
                  <span>提问</span>
                </Link>
                <Link
                  href="/timeline"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <Sparkles className="w-5 h-5" />
                  <span>时间线</span>
                  {unreadCount > 0 && (
                    <span className="ml-auto badge badge-primary badge-xs">
                      {unreadCount}
                    </span>
                  )}
                </Link>
                <Link
                  href="/topics"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <Bookmark className="w-5 h-5" />
                  <span>专题</span>
                </Link>
                <Link
                  href="/boards"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <LayoutGrid className="w-5 h-5" />
                  <span>板块</span>
                </Link>
                
                <div className="divider my-3"></div>
                
                <Link
                  href="/settings"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <Settings className="w-5 h-5" />
                  <span>设置</span>
                </Link>
                <Link
                  href="/help"
                  onClick={onClose}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <HelpCircle className="w-5 h-5" />
                  <span>帮助</span>
                </Link>
                
                {user?.role === "admin" && (
                  <Link
                    href="/admin"
                    onClick={onClose}
                    className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors text-primary"
                  >
                    <LayoutDashboard className="w-5 h-5" />
                    <span>管理后台</span>
                  </Link>
                )}
                
                <button
                  onClick={onLogout}
                  className="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors text-error mt-2"
                >
                  <LogOut className="w-5 h-5" />
                  <span>退出登录</span>
                </button>
              </>
            ) : (
              <div className="space-y-2 p-3">
                <Link
                  href="/auth/login"
                  onClick={onClose}
                  className="btn btn-ghost w-full"
                >
                  登录
                </Link>
                <Link
                  href="/auth/register"
                  onClick={onClose}
                  className="btn btn-primary w-full"
                >
                  注册
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}