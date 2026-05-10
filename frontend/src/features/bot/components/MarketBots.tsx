// features/bot/components/MarketBots.tsx
import {
  BotVO,
  CreateBotRequest,
  UpdateBotRequest,
} from "@/shared/api/types/bot.model";
import { useBots, useBotActions } from "../hooks/bot";
import { Pagination } from "@/shared/ui/common/Pagination";
import { BotFormModal } from "./BotFormModal";
import { useState } from "react";

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
  const totalPages = Math.ceil(total / pageSize);

  const [modalOpen, setModalOpen] = useState(false);
  const [editingBot, setEditingBot] = useState<BotVO | null>(null);

  const handleRun = async (id: number) => {
    await runBot(id);
    refresh();
  };

  const renderStatus = (status: string) => {
    const map: Record<string, string> = {
      active: "badge-success",
      inactive: "badge-ghost",
      error: "badge-error",
      loading: "badge-warning",
      stopped: "badge-secondary",
    };
    return (
      <span className={`badge ${map[status] || "badge-ghost"}`}>{status}</span>
    );
  };

  if (loading && bots.length === 0)
    return <div className="flex justify-center p-8">加载机器人列表中...</div>;

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
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {bots.map((bot) => (
            <div key={bot.id} className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <div className="flex justify-between items-start">
                  <h2 className="card-title">{bot.name}</h2>
                  {renderStatus(bot.status)}
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
                    disabled={actionLoading}
                  >
                    运行
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <Pagination
        totalPages={totalPages}
        currentPage={page}
        onPageChange={onPageChange}
      />

      <BotFormModal
        isOpen={modalOpen}
        editingBot={editingBot}
        onClose={() => setModalOpen(false)}
        onSave={async () => {}}
        isLoading={false}
      />
    </div>
  );
}
