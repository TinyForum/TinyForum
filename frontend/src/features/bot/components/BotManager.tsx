import { useState, useRef, useEffect } from "react";
import { MarketBots } from "./MarketBots";
import { MyBots } from "./MyBots";
import { AdminBots } from "./AdminBots";
import { useAuthStore } from "@/store";
import { BotFlowEditor } from "@/shared/ui/editor/bot/BotFlowEditor";

type TabType = "market" | "my" | "admin" | "nocode";

export function BotManager() {
  const [activeTab, setActiveTab] = useState<TabType>("market");
  const [page, setPage] = useState(1);
  const pageSize = 10;
  const { user } = useAuthStore();
  const isAdmin = user?.role === "admin" || user?.role === "super_admin";

  // 用于滚动重置的 ref
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // 切换 tab 或翻页时，将滚动容器滚动到顶部
  useEffect(() => {
    scrollContainerRef.current?.scrollTo({ top: 0, behavior: "auto" });
  }, [activeTab, page]);

  const onPageChange = (newPage: number) => setPage(newPage);

  return (
    // 外层容器：占满视口高度，使用 flex 列布局

    <div className="flex flex-col gap-6  h-[calc(100vh-16rem)]">
      {/* ========= 固定区域：标题 & Tab ========= */}
      <div className="flex-shrink-0 border-b border-base-200 bg-base-100/95 backdrop-blur-sm z-10">
        <div className="container mx-auto px-4 py-4">
          {/* 页面标题（可自定义） */}
          <h1 className="text-2xl font-bold mb-2">机器人管理</h1>

          {/* Tabs 组件 */}
          <div className="tabs tabs-boxed">
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
            <button
              className={`tab ${activeTab === "nocode" ? "tab-active" : ""}`}
              onClick={() => {
                setActiveTab("nocode");
                setPage(1);
              }}
            >
              可视化创建
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
        </div>
      </div>

      {/* ========= 滚动内容区域 ========= */}
      <div ref={scrollContainerRef} className="flex-1 overflow-y-auto">
        <div className="container mx-auto px-4 py-6">
          {activeTab === "market" && (
            <MarketBots
              page={page}
              pageSize={pageSize}
              onPageChange={onPageChange}
            />
          )}
          {activeTab === "my" && (
            <MyBots
              page={page}
              pageSize={pageSize}
              onPageChange={onPageChange}
            />
          )}
          {activeTab === "nocode" && <BotFlowEditor />}
          {activeTab === "admin" && isAdmin && (
            <AdminBots
              page={page}
              pageSize={pageSize}
              onPageChange={onPageChange}
            />
          )}
        </div>
      </div>
    </div>
  );
}
