"use client";

import { useState } from "react";
import { BotVO } from "@/shared/api/types/bot.model";

export function BotExecutionLog({ bot }: { bot: BotVO }) {
  const [open, setOpen] = useState(false);
  const logs = bot.lastExecLogs ?? [];
  const hasData = logs.length > 0 || !!bot.errorMsg || bot.lastExecDurationMs > 0;

  if (!hasData) return null;

  const logColor = (log: string) => {
    if (log.startsWith("✗") || log.includes("[error]")) return "text-red-600";
    if (log.startsWith("✓")) return "text-green-600";
    if (log.startsWith("→") || log.startsWith("⚡")) return "text-blue-600";
    if (log.startsWith("∥")) return "text-purple-600";
    if (log.startsWith("⟳")) return "text-indigo-600";
    return "text-gray-600";
  };

  return (
    <div className="mt-2 border-t pt-2">
      <button
        className="text-xs text-blue-600 hover:underline"
        onClick={() => setOpen(!open)}
      >
        {open ? "▲ 收起日志" : "▼ 展开执行日志"}
        {bot.errorMsg && (
          <span className="text-red-500 ml-1">(错误)</span>
        )}
      </button>
      {open && (
        <div className="mt-2 text-xs space-y-1">
          {bot.lastExecDurationMs > 0 && (
            <div className="text-gray-500">
              耗时: {bot.lastExecDurationMs} ms · 执行次数: {bot.execCount}
            </div>
          )}
          {bot.errorMsg && (
            <div className="bg-red-50 border border-red-200 rounded p-2 text-red-700 font-mono whitespace-pre-wrap break-all">
              {bot.errorMsg}
            </div>
          )}
          {logs.length > 0 && (
            <div className="bg-gray-50 border rounded p-2 font-mono max-h-48 overflow-y-auto">
              {logs.map((log, i) => (
                <div key={i} className={logColor(log)}>
                  {log}
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
