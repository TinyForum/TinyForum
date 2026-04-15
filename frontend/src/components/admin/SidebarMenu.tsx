import { Award, BarChart3, ChevronLeft, ChevronRight, FileText, HelpCircle, LayoutDashboard, LogOut, Megaphone, Settings, Shield, User2Icon, Users } from "lucide-react";

// ==================== 业务面板菜单组件 ====================
export function SidebarMenu({ 
  activeMenu, 
  onMenuChange, 
  collapsed, 
  onCollapsedChange,
  t 
}: { 
  activeMenu: string; 
  onMenuChange: (menu: string) => void;
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
  t: (key: string) => string;
}) {
  const menuItems = [
    { id: "dashboard", label: t("dashboard"), icon: LayoutDashboard, color: "text-primary" },
    { id: "announcements", label: t("announcements"), icon: Megaphone, color: "text-warning" },
    { id: "users", label: t("user_management"), icon: Users, color: "text-info" },
    { id: "moderators_management", label: t("moderators_management"), icon: User2Icon, color: "text-info" },
    { id: "posts", label: t("post_management"), icon: FileText, color: "text-success" },
    { id: "qa", label: t("qa_management"), icon: HelpCircle, color: "text-secondary" },
    { id: "points", label: t("points_management"), icon: Award, color: "text-accent" },
    { id: "system", label: t("system"), icon: BarChart3, color: "text-primary" },
    { id: "settings", label: t("settings"), icon: Settings, color: "text-base-content" },
  ];

  return (
    <div className={`bg-base-200 border-r border-base-300 transition-all duration-300 ${collapsed ? 'w-20' : 'w-64'} flex flex-col`}>
      {/* Logo 区域 */}
      <div className="h-16 flex items-center justify-between px-4 border-b border-base-300">
        {!collapsed && (
          <div className="flex items-center gap-2">
            <Shield className="w-6 h-6 text-primary" />
            <span className="font-bold text-lg">{t("admin_panel")}</span>
          </div>
        )}
        {collapsed && (
          <div className="flex justify-center w-full">
            <Shield className="w-6 h-6 text-primary" />
          </div>
        )}
        <button
          onClick={() => onCollapsedChange(!collapsed)}
          className="btn btn-ghost btn-sm"
        >
          {collapsed ? <ChevronRight className="w-4 h-4" /> : <ChevronLeft className="w-4 h-4" />}
        </button>
      </div>

      {/* 菜单列表 */}
      <div className="flex-1 overflow-y-auto py-4">
        {menuItems.map((item) => (
          <button
            key={item.id}
            onClick={() => onMenuChange(item.id)}
            className={`w-full flex items-center gap-3 px-4 py-3 transition-colors ${
              activeMenu === item.id
                ? "bg-primary/10 text-primary border-r-2 border-primary"
                : "hover:bg-base-300"
            } ${collapsed ? 'justify-center' : ''}`}
          >
            <item.icon className={`w-5 h-5 ${item.color}`} />
            {!collapsed && <span className="text-sm">{item.label}</span>}
          </button>
        ))}
      </div>

      {/* 底部退出按钮 */}
      <div className="border-t border-base-300 p-4">
        <button
          className={`w-full flex items-center gap-3 text-error hover:bg-error/10 rounded-lg transition-colors ${
            collapsed ? 'justify-center py-2' : 'px-4 py-2'
          }`}
        >
          <LogOut className="w-5 h-5" />
          {!collapsed && <span className="text-sm">{t("logout")}</span>}
        </button>
      </div>
    </div>
  );
}
