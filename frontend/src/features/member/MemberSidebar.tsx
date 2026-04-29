// components/member/MemberSidebar.tsx
"use client";

import { useTranslations } from "next-intl";

interface MemberSidebarProps {
  activeMenu: string;
  onMenuChange: (menu: string) => void;
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
  menus: Array<{ id: string; label: string; icon: string; badge?: number }>;
}

export function MemberSidebar({
  activeMenu,
  onMenuChange,
  collapsed,
  onCollapsedChange,
  menus,
}: MemberSidebarProps) {
  const t = useTranslations("Member");
  return (
    <div
      className={`bg-base-200 border-r border-base-300 transition-all duration-300 flex flex-col ${
        collapsed ? "w-20" : "w-64"
      }`}
    >
      <div className="p-4 flex justify-between items-center border-b border-base-300">
        {!collapsed && (
          <h2 className="font-bold text-lg">{t("member_center")}</h2>
        )}
        <button
          onClick={() => onCollapsedChange(!collapsed)}
          className="btn btn-ghost btn-sm"
        >
          {collapsed ? "→" : "←"}
        </button>
      </div>

      <div className="p-2 space-y-1 flex-1">
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
              <>
                <span className="flex-1 text-left">{menu.label}</span>
                {menu.badge !== undefined && menu.badge > 0 && (
                  <span className="badge badge-sm badge-primary">
                    {menu.badge}
                  </span>
                )}
              </>
            )}
          </button>
        ))}
      </div>

      {!collapsed && (
        <div className="p-4 border-t border-base-300">
          <div className="flex items-center gap-3">
            <div className="avatar placeholder">
              <div className="bg-neutral text-neutral-content rounded-full w-10">
                <span>👤</span>
              </div>
            </div>
            <div>
              <p className="text-sm font-medium">{t("member")}</p>
              {/* TODO: 显示用户名 */}
              <p className="text-xs text-base-content/60">member@example.com</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
