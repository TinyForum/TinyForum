"use client";

import {
  Settings,
  Puzzle,
  ToggleLeft,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";

export type SystemMenuId = "config" | "plugins" | "features";

interface MenuItem {
  id: SystemMenuId;
  label: string;
  sublabel: string;
  icon: React.ReactNode;
  accent: string;
}

const MENUS: MenuItem[] = [
  {
    id: "config",
    label: "网站配置",
    sublabel: "Site Config",
    icon: <Settings className="w-5 h-5" />,
    accent: "text-primary",
  },
  {
    id: "plugins",
    label: "插件管理",
    sublabel: "Plugins",
    icon: <Puzzle className="w-5 h-5" />,
    accent: "text-secondary",
  },
  {
    id: "features",
    label: "功能开关",
    sublabel: "Feature Flags",
    icon: <ToggleLeft className="w-5 h-5" />,
    accent: "text-accent",
  },
];

interface SystemSidebarProps {
  active: SystemMenuId;
  collapsed: boolean;
  onSelect: (id: SystemMenuId) => void;
  onToggleCollapse: () => void;
}

export function SystemSidebar({
  active,
  collapsed,
  onSelect,
  onToggleCollapse,
}: SystemSidebarProps) {
  return (
    <aside
      className={`${
        collapsed ? "w-16" : "w-56"
      } flex flex-col border-r border-base-300 bg-base-100 transition-all duration-300 shrink-0`}
    >
      {/* Header */}
      <div className="flex items-center justify-between border-b border-base-300 px-3 py-4 min-h-[64px]">
        {!collapsed && (
          <div className="overflow-hidden">
            <p className="font-bold text-base leading-tight truncate">
              系统管理
            </p>
            <p className="text-[11px] text-base-content/40 tracking-widest uppercase">
              System
            </p>
          </div>
        )}
        <button
          onClick={onToggleCollapse}
          className="btn btn-ghost btn-xs btn-square shrink-0 ml-auto"
          title={collapsed ? "展开" : "收起"}
        >
          {collapsed ? (
            <ChevronRight className="w-4 h-4" />
          ) : (
            <ChevronLeft className="w-4 h-4" />
          )}
        </button>
      </div>

      {/* Nav */}
      <nav className="flex-1 p-2 space-y-1">
        {MENUS.map((menu) => {
          const isActive = active === menu.id;
          return (
            <button
              key={menu.id}
              onClick={() => onSelect(menu.id)}
              title={collapsed ? menu.label : undefined}
              className={`
                w-full flex items-center gap-3 rounded-lg px-3 py-2.5 text-left
                transition-all duration-150 group
                ${
                  isActive
                    ? "bg-primary text-primary-content shadow-sm"
                    : "hover:bg-base-200 text-base-content/70 hover:text-base-content"
                }
              `}
            >
              <span className={`shrink-0 ${isActive ? "" : menu.accent}`}>
                {menu.icon}
              </span>
              {!collapsed && (
                <div className="overflow-hidden">
                  <p className="text-sm font-medium leading-none truncate">
                    {menu.label}
                  </p>
                  <p
                    className={`text-[10px] mt-0.5 truncate ${isActive ? "text-primary-content/60" : "text-base-content/30"}`}
                  >
                    {menu.sublabel}
                  </p>
                </div>
              )}
              {!collapsed && isActive && (
                <span className="ml-auto w-1.5 h-1.5 rounded-full bg-primary-content/60 shrink-0" />
              )}
            </button>
          );
        })}
      </nav>

      {/* Footer hint */}
      {!collapsed && (
        <div className="p-3 border-t border-base-300">
          <p className="text-[10px] text-base-content/25 text-center leading-relaxed">
            修改配置后立即生效
          </p>
        </div>
      )}
    </aside>
  );
}
