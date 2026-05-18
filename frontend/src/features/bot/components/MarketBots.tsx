// features/bot/components/MarketBots.tsx
import { BotStatus } from "@/shared/api/types/bot.model.do";
import { useBots, useBotActions } from "../hooks/bot";
import { Pagination } from "@/shared/ui/common/Pagination";
import { useState } from "react";
import toast from "react-hot-toast";

// 状态徽章映射
const STATUS_BADGE_MAP: Record<BotStatus, string> = {
  active: "badge-success",
  inactive: "badge-ghost",
  error: "badge-error",
  loading: "badge-warning",
  stopped: "badge-secondary",
};

// 渲染状态徽章
const renderStatusBadge = (status: BotStatus) => (
  <span className={`badge ${STATUS_BADGE_MAP[status] || "badge-ghost"}`}>
    {status}
  </span>
);

interface MarketBotsProps {
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function MarketBots({ page, pageSize, onPageChange }: MarketBotsProps) {
  const { bots, loading, total, refresh } = useBots({
    page,
    pageSize,
    autoLoad: true,
  });
  const { runBot, loading: actionLoading } = useBotActions();
  const [runningBotId, setRunningBotId] = useState<number | null>(null);

  const totalPages = Math.ceil(total / pageSize);

  const handleRun = async (id: number) => {
    setRunningBotId(id);
    try {
      await runBot(id);
      toast.success("机器人运行成功");
      await refresh(); // 刷新列表以更新执行次数和状态
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "运行失败，请重试");
    } finally {
      setRunningBotId(null);
    }
  };

  // 加载中且无数据时显示全屏加载状态
  if (loading && bots.length === 0) {
    return <div className="flex justify-center p-8">加载机器人列表中...</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">机器人市场</h1>
      </div>

      {bots.length === 0 ? (
        <div className="card bg-base-100 shadow-xl p-8 text-center">
          暂无公开机器人
        </div>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {bots.map((bot) => (
              <div key={bot.id} className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <div className="flex justify-between items-start">
                    <h2 className="card-title">{bot.name}</h2>
                    {renderStatusBadge(bot.status)}
                  </div>
                  <p className="text-sm text-base-content/70 line-clamp-2">
                    {bot.description}
                  </p>
                  <div className="flex flex-wrap gap-2 text-xs text-base-content/60 mt-2">
                    <span>版本: {bot.version}</span>
                    <span>类型: {bot.type}</span>
                    <span>触发: {bot.triggerType}</span>
                    {bot.triggerType === "schedule" && (
                      <span>计划: {bot.cronExpr}</span>
                    )}
                    <span>执行次数: {bot.execCount}</span>
                    <span>作者: {bot.creatorName}</span>
                  </div>
                  <div className="card-actions justify-end mt-4">
                    <button
                      className="btn btn-xs btn-outline"
                      onClick={() => handleRun(bot.id)}
                      disabled={actionLoading || runningBotId === bot.id}
                    >
                      {runningBotId === bot.id ? "运行中..." : "运行"}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* 只有当存在多页时才显示分页组件 */}
          {totalPages > 1 && (
            <Pagination
              totalPages={totalPages}
              currentPage={page}
              onPageChange={onPageChange}
            />
          )}
        </>
      )}
    </div>
  );
}
