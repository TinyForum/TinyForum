// features/bot/components/MyBots.tsx
import { useState } from "react";
import {
  BotVO,
  CreateBotRequest,
  UpdateBotRequest,
} from "@/shared/api/types/bot.model";
import { useMyBots, useBotActions } from "../hooks/bot";
import { Pagination } from "@/shared/ui/common/Pagination";
import { BotFormModal } from "./BotFormModal";

interface MyBotsProps {
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function MyBots({ page, pageSize, onPageChange }: MyBotsProps) {
  const { bots, loading, total, refresh } = useMyBots({
    page,
    pageSize,
    autoLoad: true,
  });
  const {
    createBot,
    updateBot,
    deleteBot,
    runBot,
    loading: actionLoading,
  } = useBotActions();
  const totalPages = Math.ceil(total / pageSize);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingBot, setEditingBot] = useState<BotVO | null>(null);

  const handleSave = async (
    data: CreateBotRequest | UpdateBotRequest,
    isEdit: boolean,
    botId?: number,
  ) => {
    let success = false;
    if (isEdit && botId)
      success = await updateBot(botId, data as UpdateBotRequest);
    else success = !!(await createBot(data as CreateBotRequest));
    if (success) {
      setModalOpen(false);
      refresh();
    }
  };

  const handleEdit = (bot: BotVO) => {
    setEditingBot(bot);
    setModalOpen(true);
  };
  const handleCreate = () => {
    setEditingBot(null);
    setModalOpen(true);
  };
  const handleDelete = async (id: number) => {
    if (window.confirm("确定删除该机器人吗？") && (await deleteBot(id)))
      refresh();
  };
  const handleRun = async (id: number) => {
    await runBot(id);
    refresh();
  };
  const handleToggleEnable = async (bot: BotVO) => {
    await updateBot(bot.id, { enabled: !bot.enabled });
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
        <h1 className="text-2xl font-bold">我的机器人</h1>
        <button className="btn btn-primary" onClick={handleCreate}>
          创建机器人
        </button>
      </div>

      {bots.length === 0 ? (
        <div className="card bg-base-100 shadow-xl p-8 text-center">
          还没有机器人，点击上方按钮创建第一个机器人吧！
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
                </div>
                <div className="card-actions justify-end mt-4">
                  <button
                    className="btn btn-xs btn-outline"
                    onClick={() => handleRun(bot.id)}
                    disabled={actionLoading}
                  >
                    运行
                  </button>
                  <button
                    className="btn btn-xs btn-outline"
                    onClick={() => handleToggleEnable(bot)}
                    disabled={actionLoading}
                  >
                    {bot.enabled ? "禁用" : "启用"}
                  </button>
                  <button
                    className="btn btn-xs btn-outline"
                    onClick={() => handleEdit(bot)}
                  >
                    编辑
                  </button>
                  <button
                    className="btn btn-xs btn-error"
                    onClick={() => handleDelete(bot.id)}
                    disabled={actionLoading}
                  >
                    删除
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
        onSave={handleSave}
        isLoading={actionLoading}
      />
    </div>
  );
}
