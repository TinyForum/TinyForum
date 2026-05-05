"use client";

import { useState } from "react";
import {
  SystemSidebar,
  type SystemMenuId,
} from "@/features/system/components/SystemSidebar";
import { SiteConfigPanel } from "@/features/system/components/SiteConfigPanel";
import { SystemPluginsPanel } from "@/features/system/components/SystemPluginsPanel";
import { FeatureFlagsPanel } from "@/features/system/components/FeatureFlagsPanel";
import { useSiteConfig } from "@/features/system/hooks/useSiteConfig";
import { useFeatureFlags } from "@/features/system/hooks/useFeatureFlags";

const PAGE_META: Record<SystemMenuId, { title: string; subtitle: string }> = {
  config: {
    title: "网站配置",
    subtitle: "管理站点基础信息、显示偏好与高级选项",
  },
  plugins: {
    title: "插件管理",
    subtitle: "安装、启用或移除扩展插件，变更在下次页面加载时生效",
  },
  features: {
    title: "功能开关",
    subtitle: "实时控制各功能模块的启用状态，修改后立即生效",
  },
};

export default function SystemPage() {
  const [activeMenu, setActiveMenu] = useState<SystemMenuId>("config");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const { config, update, save, isSaving } = useSiteConfig();
  const { grouped, enabledCount, features, toggle, enableAll, togglingId } =
    useFeatureFlags();

  const meta = PAGE_META[activeMenu];

  return (
    <div className="flex h-[calc(100vh-64px)] bg-base-200 overflow-hidden rounded-xl border border-base-300 shadow-sm">
      {/* Sidebar */}
      <SystemSidebar
        active={activeMenu}
        collapsed={sidebarCollapsed}
        onSelect={(id) => setActiveMenu(id)}
        onToggleCollapse={() => setSidebarCollapsed((v) => !v)}
      />

      {/* Main content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Page header */}
        <header className="shrink-0 px-6 py-4 border-b border-base-300 bg-base-100 flex items-end gap-3">
          <div>
            <h1 className="text-lg font-bold leading-none">{meta.title}</h1>
            <p className="text-xs text-base-content/40 mt-1">{meta.subtitle}</p>
          </div>
        </header>

        {/* Scrollable content */}
        <div className="flex-1 overflow-y-auto custom-scrollbar">
          <div className="p-6 max-w-3xl">
            {activeMenu === "config" && (
              <SiteConfigPanel
                config={config}
                isSaving={isSaving}
                update={update}
                onSave={save}
              />
            )}
            {activeMenu === "plugins" && <SystemPluginsPanel />}
            {activeMenu === "features" && (
              <FeatureFlagsPanel
                grouped={grouped}
                enabledCount={enabledCount}
                total={features.length}
                togglingId={togglingId}
                onToggle={toggle}
                onEnableAll={enableAll}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
