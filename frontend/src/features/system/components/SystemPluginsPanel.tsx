"use client";

import { useState } from "react";
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
import { MyPluginsTab } from "./MyPluginsTab";

export function SystemPluginsPanel() {
  const [activeTab, setActiveTab] = useState<
    "manage" | "market" | "config" | "logs" | "dev" | "upload" | "my"
  >("manage");

  const tabs = [
    { id: "manage", label: "插件管理", icon: Puzzle },
    { id: "market", label: "插件市场", icon: Store },
    { id: "config", label: "插件配置", icon: Settings },
    { id: "logs", label: "运行日志", icon: Activity },
    { id: "dev", label: "开发工具", icon: Wrench },
    { id: "upload", label: "上传插件", icon: CloudUpload },
    { id: "my", label: "我的插件", icon: Puzzle },
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
              <span className="hidden sm:inline">{tab.label}</span>
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
