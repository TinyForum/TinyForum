// components/moderator/ModeratorSidebar.tsx
"use client";

import { 
  LayoutDashboard, 
  FileText, 
  Flag, 
  Ban,
  ChevronLeft,
  ChevronRight,
  Shield
} from "lucide-react";

interface ModeratorSidebarProps {
  activeMenu: string;
  onMenuChange: (menu: string) => void;
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
  boardName?: string;
  permissions: {
    isModerator: boolean;
    canDeletePost: boolean;
    canPinPost: boolean;
    canBanUser: boolean;
    canManageModerator: boolean;
  };
  t: (key: string) => string;
}

export function ModeratorSidebar({
  activeMenu,
  onMenuChange,
  collapsed,
  onCollapsedChange,
  boardName,
  permissions,
  t,
}: ModeratorSidebarProps) {
  const menuItems = [
    { id: "dashboard", label: t("dashboard"), icon: LayoutDashboard, always: true },
    { id: "posts", label: t("posts_management"), icon: FileText, always: true },
    { id: "reports", label: t("reports_management"), icon: Flag, always: true },
    { id: "bans", label: t("bans_management"), icon: Ban, require: permissions.canBanUser },
  ];

  return (
    <aside className={`bg-base-200 border-r border-base-300 transition-all duration-300 ${
      collapsed ? "w-20" : "w-64"
    }`}>
      <div className="p-4">
        {/* 板块信息 */}
        {!collapsed && boardName && (
          <div className="mb-4 p-3 bg-primary/10 rounded-lg">
            <div className="flex items-center gap-2 text-primary">
              <Shield className="w-4 h-4" />
              <span className="font-medium text-sm">{t("managed_board")}</span>
            </div>
            <p className="font-bold mt-1">{boardName}</p>
          </div>
        )}

        {/* 折叠按钮 */}
        <button
          onClick={() => onCollapsedChange(!collapsed)}
          className="btn btn-ghost btn-sm w-full mb-4"
        >
          {collapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
          {!collapsed && <span className="ml-2">{t("collapse")}</span>}
        </button>

        {/* 菜单 */}
        <ul className="menu">
          {menuItems.map((item) => {
            const Icon = item.icon;
            if (!item.always && !item.require) {
              return null;
            }
            return (
              <li key={item.id}>
                <button
                  onClick={() => onMenuChange(item.id)}
                  className={`${activeMenu === item.id ? "active" : ""}`}
                >
                  <Icon className="w-4 h-4" />
                  {!collapsed && <span>{item.label}</span>}
                </button>
              </li>
            );
          })}
        </ul>
      </div>
    </aside>
  );
}