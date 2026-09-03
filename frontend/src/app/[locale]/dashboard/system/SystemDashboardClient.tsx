"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  SystemSidebar,
  type SystemMenuId,
} from "@/features/system/components/SystemSidebar";
import { SiteConfigPanel } from "@/features/system/components/SiteConfigPanel";
import { SystemPluginsPanel } from "@/features/system/components/SystemPluginsPanel";
import { FeatureFlagsPanel } from "@/features/system/components/FeatureFlagsPanel";
import { useSiteConfig } from "@/features/system/hooks/useSiteConfig";
import { useFeatureFlags } from "@/features/system/hooks/useFeatureFlags";

const PAGE_META: Record<SystemMenuId, { titleKey: string; subtitleKey: string }> = {
  website_config: {
    titleKey: "website_config",
    subtitleKey: "auto_load_from_backend",
  },
  plugins_center: {
    titleKey: "plugins_center",
    subtitleKey: "plugin_install_tip",
  },
  features_flags: {
    titleKey: "features_flags",
    subtitleKey: "feature_control_tip",
  },
};

export default function SystemDashboardClient() {
  const t = useTranslations("System");
  const [activeMenu, setActiveMenu] = useState<SystemMenuId>("website_config");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const {
    grouped,
    groupOrder,
    configFile,
    isLoading,
    isSaving,
    isReloading,
    updateField,
    save,
    reloadConfig,
    switchConfig,
  } = useSiteConfig();

  const { grouped: ffGrouped, enabledCount, features, toggle, enableAll, togglingId } =
    useFeatureFlags();

  const meta = PAGE_META[activeMenu];

  return (
    <div className="flex h-[calc(100vh-64px)] bg-base-200 overflow-hidden rounded-xl border border-base-300 shadow-sm">
      <SystemSidebar
        active={activeMenu}
        collapsed={sidebarCollapsed}
        onSelect={(id) => setActiveMenu(id)}
        onToggleCollapse={() => setSidebarCollapsed((v) => !v)}
      />

      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="shrink-0 px-6 py-4 border-b border-base-300 bg-base-100 flex items-end gap-3">
          <div>
            <h1 className="text-lg font-bold leading-none">{t(meta.titleKey)}</h1>
            <p className="text-xs text-base-content/40 mt-1">{t(meta.subtitleKey)}</p>
          </div>
        </header>

        <div className="flex-1 overflow-y-auto custom-scrollbar">
          <div className="p-6 max-w-3xl">
            {activeMenu === "website_config" && (
              <SiteConfigPanel
                groups={grouped}
                groupOrder={groupOrder}
                isSaving={isSaving}
                isReloading={isReloading}
                isLoading={isLoading}
                configFile={configFile}
                onChangeFile={switchConfig}
                onUpdateField={updateField}
                onSave={save}
                onReload={reloadConfig}
              />
            )}
            {activeMenu === "plugins_center" && <SystemPluginsPanel />}
            {activeMenu === "features_flags" && (
              <FeatureFlagsPanel
                grouped={ffGrouped}
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
