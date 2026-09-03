"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  Puzzle,
  Store,
  Settings,
  Wrench,
  Activity,
  CloudUpload,
} from "lucide-react";
import { DeveloperToolsTab } from "@/features/plugin/components/DeveloperToolsTab";
import { PluginConfigTab } from "@/features/plugin/components/PluginConfigTab";
import { PluginLogsTab } from "@/features/plugin/components/PluginLogsTab";
import { PluginManagementTab } from "@/features/plugin/components/PluginManagementTab";
import { PluginMarketTab } from "@/features/plugin/components/PluginMarketTab";
import { UploadPluginTab } from "@/features/plugin/components/UploadPluginTab";
import { MyPluginsTab } from "@/features/plugin/components/MyPluginsTab";

export function SystemPluginsPanel() {
  const t = useTranslations("System");
  const [activeTab, setActiveTab] = useState<
    "manage" | "market" | "config" | "logs" | "dev" | "upload" | "my"
  >("manage");

  const tabs = [
    { id: "manage", labelKey: "plugin_management", icon: Puzzle },
    { id: "market", labelKey: "plugin_market", icon: Store },
    { id: "config", labelKey: "plugin_config", icon: Settings },
    { id: "logs", labelKey: "runtime_logs", icon: Activity },
    { id: "dev", labelKey: "developer_tools", icon: Wrench },
    { id: "upload", labelKey: "upload_plugin", icon: CloudUpload },
    { id: "my", labelKey: "my_plugins", icon: Puzzle },
  ] as const;

  return (
    <div className="space-y-5">
      <div className="tabs tabs-boxed bg-base-200 p-1 w-fit flex-wrap">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              className={`tab gap-2 ${activeTab === tab.id ? "tab-active" : ""}`}
              onClick={() => setActiveTab(tab.id)}
            >
              <Icon className="w-4 h-4" />
              <span className="hidden sm:inline">{t(tab.labelKey)}</span>
            </button>
          );
        })}
      </div>

      <div>
        {activeTab === "manage" && <PluginManagementTab />}
        {activeTab === "market" && <PluginMarketTab />}
        {activeTab === "config" && <PluginConfigTab />}
        {activeTab === "logs" && <PluginLogsTab />}
        {activeTab === "dev" && <DeveloperToolsTab />}
        {activeTab === "upload" && <UploadPluginTab />}
        {activeTab === "my" && <MyPluginsTab />}
      </div>
    </div>
  );
}
