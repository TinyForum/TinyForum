// components/moderator/ModeratorSidebar.tsx
"use client";

import { 
  LayoutDashboard, 
  FileText, 
  Flag, 
  Ban,
  ChevronLeft,
  ChevronRight,
  Shield,
  ChevronDown
} from "lucide-react";
import { useState } from "react";

interface ModeratorBoard {
  id: number;
  name: string;
  slug: string;
  permissions?: {
    can_delete_post: boolean;
    can_pin_post: boolean;
    can_ban_user: boolean;
    can_manage_moderator: boolean;
  };
}

interface ModeratorSidebarProps {
  activeMenu: string;
  onMenuChange: (menu: string) => void;
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
  boards: ModeratorBoard[];  // 板块列表
  currentBoardId: number | null;
  onBoardChange: (boardId: number) => void;
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
  boards,
  currentBoardId,
  onBoardChange,
  permissions,
  t,
}: ModeratorSidebarProps) {
  const [showBoardDropdown, setShowBoardDropdown] = useState(false);
  const currentBoard = boards.find(b => b.id === currentBoardId);

  const menuItems = [
    { id: "dashboard", label: t("dashboard"), icon: LayoutDashboard, always: true },
    { id: "posts", label: t("posts_management"), icon: FileText, always: true },
    { id: "reports", label: t("reports_management"), icon: Flag, always: true },
    { id: "bans", label: t("bans_management"), icon: Ban, require: permissions.canBanUser },
  ];

  return (
    <aside className={`bg-base-200 border-r border-base-300 transition-all duration-300 flex flex-col ${
      collapsed ? "w-20" : "w-64"
    }`}>
      <div className="flex-1 p-4">
        {/* 折叠按钮 */}
        <button
          onClick={() => onCollapsedChange(!collapsed)}
          className="btn btn-ghost btn-sm w-full mb-4"
        >
          {collapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
          {!collapsed && <span className="ml-2">{t("collapse")}</span>}
        </button>

        {/* 板块切换器 */}
        {!collapsed && boards.length > 0 && (
          <div className="mb-4">
            <div className="text-xs text-base-content/60 mb-2 px-2">
              {t("managed_boards")}
            </div>
            <div className="relative">
              <button
                onClick={() => setShowBoardDropdown(!showBoardDropdown)}
                className="btn btn-outline btn-sm w-full justify-between"
              >
                <div className="flex items-center gap-2 truncate">
                  <Shield className="w-4 h-4 text-primary" />
                  <span className="truncate">{currentBoard?.name || t("select_board")}</span>
                </div>
                <ChevronDown className={`w-4 h-4 transition-transform ${showBoardDropdown ? "rotate-180" : ""}`} />
              </button>
              
              {showBoardDropdown && (
                <div className="absolute left-0 right-0 top-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg z-50">
                  {boards.map((board) => (
                    <button
                      key={board.id}
                      onClick={() => {
                        onBoardChange(board.id);
                        setShowBoardDropdown(false);
                      }}
                      className={`w-full text-left px-3 py-2 hover:bg-base-200 transition-colors ${
                        currentBoardId === board.id ? "bg-primary/10 text-primary" : ""
                      } first:rounded-t-lg last:rounded-b-lg`}
                    >
                      <div className="flex items-center gap-2">
                        <Shield className="w-4 h-4" />
                        <span className="text-sm truncate">{board.name}</span>
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        {/* 折叠状态下的板块指示器 */}
        {collapsed && boards.length > 0 && currentBoard && (
          <div className="mb-4 flex justify-center">
            <div className="tooltip tooltip-right" data-tip={currentBoard.name}>
              <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
                <Shield className="w-5 h-5 text-primary" />
              </div>
            </div>
          </div>
        )}

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

      {/* 底部信息（可选） */}
      {!collapsed && (
        <div className="p-4 border-t border-base-300">
          <div className="text-xs text-base-content/60">
            <p>{t("moderator_panel")}</p>
            <p className="mt-1">v1.0.0</p>
          </div>
        </div>
      )}
    </aside>
  );
}