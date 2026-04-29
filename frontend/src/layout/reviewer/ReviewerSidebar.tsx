// components/reviewer/ReviewerSidebar.tsx
"use client";

import { useTranslations } from "next-intl";

interface ReviewerSidebarProps {
  activeMenu: string;
  onMenuChange: (menu: string) => void;
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
  menus: Array<{ id: string; label: string; icon: string }>;
  stats: {
    pending: number;
    reported: number;
    reviewedToday: number;
  };
}

export function ReviewerSidebar({
  activeMenu,
  onMenuChange,
  collapsed,
  onCollapsedChange,
  menus,
  stats,
}: ReviewerSidebarProps) {
  const t = useTranslations("Reveiw");
  return (
    <div
      className={`bg-base-200 border-r border-base-300 transition-all duration-300 ${
        collapsed ? "w-20" : "w-64"
      }`}
    >
      <div className="p-4 flex justify-between items-center border-b border-base-300">
        {!collapsed && (
          <h2 className="font-bold text-lg">{t("reviewer_panel")}</h2>
        )}
        <button
          onClick={() => onCollapsedChange(!collapsed)}
          className="btn btn-ghost btn-sm"
        >
          {collapsed ? "→" : "←"}
        </button>
      </div>

      <nav className="p-2 space-y-1">
        {menus.map((menu) => (
          <button
            key={menu.id}
            onClick={() => onMenuChange(menu.id)}
            className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
              activeMenu === menu.id
                ? "bg-primary text-primary-content"
                : "hover:bg-base-300"
            }`}
          >
            <span className="text-xl">{menu.icon}</span>
            {!collapsed && (
              <span className="flex-1 text-left">{menu.label}</span>
            )}
            {!collapsed && menu.id === "pending" && stats.pending > 0 && (
              <span className="badge badge-sm badge-error">
                {stats.pending}
              </span>
            )}
            {!collapsed && menu.id === "reports" && stats.reported > 0 && (
              <span className="badge badge-sm badge-warning">
                {stats.reported}
              </span>
            )}
          </button>
        ))}
      </nav>

      {!collapsed && (
        <div className="absolute bottom-0 w-64 p-4 border-t border-base-300 text-sm text-base-content/60">
          <p>
            {t("reviewed_today")}: {stats.reviewedToday}
          </p>
        </div>
      )}
    </div>
  );
}
