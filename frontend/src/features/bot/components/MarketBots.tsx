// features/bot/components/MarketBots.tsx
import { useTranslations } from "next-intl";
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
  const t = useTranslations("Bot");
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
      toast.success(t("run_bot_success_msg"));
      await refresh(); // 刷新列表以更新执行次数和状态
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "运行失败，请重试");
    } finally {
      setRunningBotId(null);
    }
  };

  // 加载中且无数据时显示全屏加载状态
  if (loading && bots.length === 0) {
    return <div className="flex justify-center p-8">{t("loading_bots")}</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">{t("bot_market")}</h1>
      </div>

      {bots.length === 0 ? (
        <div className="card bg-base-100 shadow-xl p-8 text-center">
          {t("no_public_bots")}
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
                    <span>{t("version_label")} {bot.version}</span>
                    <span>{t("type_label")} {bot.type}</span>
                    <span>{t("trigger_label")} {bot.triggerType}</span>
                    {bot.triggerType === "schedule" && (
                      <span>{t("schedule_label")} {bot.cronExpr}</span>
                    )}
                    <span>{t("exec_count_label")} {bot.execCount}</span>
                    <span>{t("author_label")} {bot.creatorName}</span>
                  </div>
                  <div className="card-actions justify-end mt-4">
                    <button
                      className="btn btn-xs btn-outline"
                      onClick={() => handleRun(bot.id)}
                      disabled={actionLoading || runningBotId === bot.id}
                    >
                      {runningBotId === bot.id ? t("running") : t("run")}
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
