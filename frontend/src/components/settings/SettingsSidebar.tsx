"use client";

import { SettingsTab } from "@/app/[locale]/settings/page";
import {
  User,
  Shield,
  Palette,
  Bell,
  AlertTriangle,
  Settings,
} from "lucide-react";
import { useTranslations } from "next-intl";
// import { SettingsTab } from '@/app/settings/page';

interface SettingsSidebarProps {
  activeTab: SettingsTab;
  onTabChange: (tab: SettingsTab) => void;
}

interface MenuItem {
  id: SettingsTab;
  label: string;
  icon: typeof User;
  description: string;
  color?: string;
}

export default function SettingsSidebar({
  activeTab,
  onTabChange,
}: SettingsSidebarProps) {
  const t = useTranslations("Settings");

  const menuItems: MenuItem[] = [
    {
      id: "profile",
      label: "个人资料",
      icon: User,
      description: "管理您的个人信息和头像",
      color: "text-blue-500",
    },
    {
      id: "security",
      label: "安全设置",
      icon: Shield,
      description: "修改密码和账户安全",
      color: "text-green-500",
    },
    {
      id: "appearance",
      label: "外观主题",
      icon: Palette,
      description: "自定义界面主题和布局",
      color: "text-purple-500",
    },
    {
      id: "notifications",
      label: "通知设置",
      icon: Bell,
      description: "管理消息通知偏好",
      color: "text-yellow-500",
    },
    {
      id: "danger",
      label: "危险区域",
      icon: AlertTriangle,
      description: "删除账户等危险操作",
      color: "text-red-500",
    },
  ];

  return (
    <aside className="w-80 bg-base-100 border-r border-base-200 flex-shrink-0 hidden md:block">
      <div className="sticky top-0 h-screen overflow-y-auto">
        {/* 侧边栏头部 */}
        <div className="p-6 border-b border-base-200">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-xl">
              <Settings className="w-5 h-5 text-primary" />
            </div>
            <div>
              <h2 className="text-xl font-bold">设置</h2>
              <p className="text-xs text-base-content/60">个性化您的体验</p>
            </div>
          </div>
        </div>

        {/* 菜单列表 */}
        <nav className="p-4 space-y-2">
          {menuItems.map((item) => {
            const Icon = item.icon;
            const isActive = activeTab === item.id;

            return (
              <button
                key={item.id}
                onClick={() => onTabChange(item.id)}
                className={`
                  w-full flex items-start gap-3 p-3 rounded-xl transition-all duration-200
                  ${
                    isActive
                      ? "bg-primary/10 text-primary shadow-sm"
                      : "hover:bg-base-200 text-base-content/70 hover:text-base-content"
                  }
                `}
              >
                {/* 图标 */}
                <div
                  className={`
                  p-2 rounded-lg transition-all
                  ${isActive ? "bg-primary/20" : "bg-base-200/50"}
                `}
                >
                  <Icon
                    className={`w-5 h-5 ${isActive ? "text-primary" : item.color}`}
                  />
                </div>

                {/* 文字内容 */}
                <div className="flex-1 text-left">
                  <div className="font-medium text-sm">{item.label}</div>
                  <div className="text-xs text-base-content/40 mt-0.5">
                    {item.description}
                  </div>
                </div>

                {/* 激活指示器 */}
                {isActive && (
                  <div className="w-1 h-8 bg-primary rounded-full animate-pulse" />
                )}
              </button>
            );
          })}
        </nav>
      </div>
    </aside>
  );
}
