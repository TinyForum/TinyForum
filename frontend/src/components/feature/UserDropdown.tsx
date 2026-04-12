"use client";

import Link from "next/link";
import {
  User,
  Settings,
  HelpCircle,
  LogOut,
  LayoutDashboard,
  Trophy,
  Bookmark,
  Sparkles,
  MessageCircleQuestion,
  LayoutGrid,
} from "lucide-react";
import Avatar from "../user/Avatar";

interface UserDropdownProps {
  user: any;
  onLogout: () => void;
}

export default function UserDropdown({ user, onLogout }: UserDropdownProps) {
  return (
    <div className="dropdown dropdown-end">
      <div
        tabIndex={0}
        role="button"
        className="btn btn-ghost btn-circle avatar hover:ring-2 hover:ring-primary/20 transition-all"
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
        className="dropdown-content menu bg-base-100 rounded-box z-10 w-64 p-2 shadow-xl border border-base-200 mt-2"
      >
        <li className="menu-title">
          <div className="flex items-center gap-2">
            <div className="avatar placeholder">
              <div className="w-10 rounded-full bg-primary/10">
                <span className="text-primary font-medium">
                  {user?.username?.[0]?.toUpperCase()}
                </span>
              </div>
            </div>
            <div className="flex-1">
              <span className="text-base-content font-medium block">
                {user?.username}
              </span>
              <span className="text-xs text-base-content/50 truncate block">
                {user?.email}
              </span>
            </div>
          </div>
        </li>
        
        <div className="divider my-1"></div>
        
        {/* 个人统计 */}
        <li className="px-2 py-1">
          <div className="flex justify-between text-sm">
            <span className="text-base-content/60">积分</span>
            <span className="font-bold text-primary">{user?.score || 0}</span>
          </div>
          <div className="flex justify-between text-sm mt-1">
            <span className="text-base-content/60">关注者</span>
            <span className="font-bold">{user?.followers_count || 0}</span>
          </div>
        </li>
        
        <div className="divider my-1"></div>
        
        {/* 快速链接 */}
        <li>
          <Link href={`/users/${user?.id}`} className="gap-2">
            <User className="w-4 h-4" />
            个人主页
          </Link>
        </li>
        <li>
          <Link href="/timeline" className="gap-2">
            <Sparkles className="w-4 h-4" />
            我的时间线
          </Link>
        </li>
        <li>
          <Link href="/topics/my" className="gap-2">
            <Bookmark className="w-4 h-4" />
            我的专题
          </Link>
        </li>
        <li>
          <Link href="/questions/my" className="gap-2">
            <MessageCircleQuestion className="w-4 h-4" />
            我的问答
          </Link>
        </li>
        
        <div className="divider my-1"></div>
        
        <li>
          <Link href="/settings" className="gap-2">
            <Settings className="w-4 h-4" />
            设置
          </Link>
        </li>
        <li>
          <Link href="/help" className="gap-2">
            <HelpCircle className="w-4 h-4" />
            帮助中心
          </Link>
        </li>
        
        {user?.role === "admin" && (
          <>
            <div className="divider my-1"></div>
            <li>
              <Link href="/admin" className="gap-2 text-primary">
                <LayoutDashboard className="w-4 h-4" />
                管理后台
              </Link>
            </li>
          </>
        )}
        
        <div className="divider my-1"></div>
        
        <li>
          <button onClick={onLogout} className="text-error gap-2">
            <LogOut className="w-4 h-4" />
            退出登录
          </button>
        </li>
      </ul>
    </div>
  );
}