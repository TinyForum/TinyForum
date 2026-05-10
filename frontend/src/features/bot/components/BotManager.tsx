// features/bot/components/BotManager.tsx
import { useState } from "react";
import { MarketBots } from "./MarketBots";
import { MyBots } from "./MyBots";
import { AdminBots } from "./AdminBots";
import { useAuthStore } from "@/store";

type TabType = "market" | "my" | "admin";

export function BotManager() {
  const [activeTab, setActiveTab] = useState<TabType>("market");
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { user } = useAuthStore();
  const isAdmin = user?.role === "admin" || user?.role === "super_admin";

  const onPageChange = (newPage: number) => setPage(newPage);

  return (
    <div className="container mx-auto p-4">
      {/* Tabs */}
      <div className="tabs tabs-boxed mb-6">
        <button
          className={`tab ${activeTab === "market" ? "tab-active" : ""}`}
          onClick={() => {
            setActiveTab("market");
            setPage(1);
          }}
        >
          机器人市场
        </button>
        <button
          className={`tab ${activeTab === "my" ? "tab-active" : ""}`}
          onClick={() => {
            setActiveTab("my");
            setPage(1);
          }}
        >
          我的机器人
        </button>
        {isAdmin && (
          <button
            className={`tab ${activeTab === "admin" ? "tab-active" : ""}`}
            onClick={() => {
              setActiveTab("admin");
              setPage(1);
            }}
          >
            管理机器人
          </button>
        )}
      </div>

      {/* 内容区域，根据 activeTab 渲染对应组件，切换时会重新创建，避免状态交叉 */}
      {activeTab === "market" && (
        <MarketBots
          page={page}
          pageSize={pageSize}
          onPageChange={onPageChange}
        />
      )}
      {activeTab === "my" && (
        <MyBots page={page} pageSize={pageSize} onPageChange={onPageChange} />
      )}
      {activeTab === "admin" && isAdmin && (
        <AdminBots
          page={page}
          pageSize={pageSize}
          onPageChange={onPageChange}
        />
      )}
    </div>
  );
}
