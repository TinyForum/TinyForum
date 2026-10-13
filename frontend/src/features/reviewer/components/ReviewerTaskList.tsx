"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useReviewerTasks } from "../hooks/useReviewerTasks";

interface ReviewerTaskListProps {
  // 按触发类型过滤（如 report_aggregate 用于"管理举报"）；不传则展示全部待审核
  triggerTypeFilter?: string;
  emptyMessageKey?: string;
}

// 待审核任务列表：审核员可对 post/comment 任务执行通过或拒绝
export function ReviewerTaskList({
  triggerTypeFilter,
  emptyMessageKey = "no_pending_content",
}: ReviewerTaskListProps) {
  const t = useTranslations("Reviewer");
  const {
    list,
    isLoading,
    error,
    refresh,
    approveTask,
    rejectTask,
    isReviewing,
    reviewingTaskId,
  } = useReviewerTasks();

  const [rejectingId, setRejectingId] = useState<number | null>(null);
  const [rejectReason, setRejectReason] = useState("");

  const tasks = triggerTypeFilter
    ? list.filter((task) => task.trigger_type === triggerTypeFilter)
    : list;

  const triggerLabel = (type: string) => {
    switch (type) {
      case "sensitive_word":
        return t("trigger_sensitive_word");
      case "report_aggregate":
        return t("trigger_report_aggregate");
      case "manual":
        return t("trigger_manual");
      default:
        return type;
    }
  };

  const targetLabel = (type: string) => {
    switch (type) {
      case "post":
        return t("target_post");
      case "comment":
        return t("target_comment");
      case "user":
        return t("target_user");
      default:
        return type;
    }
  };

  if (isLoading && tasks.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <span className="loading loading-spinner loading-lg mx-auto" />
          <p className="text-base-content/50 mt-2">{t("loading")}</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-error">{error}</p>
          <div className="mt-4">
            <button onClick={() => refresh()} className="btn btn-sm btn-primary">
              {t("retry")}
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (tasks.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t(emptyMessageKey)}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>{t("type")}</th>
            <th>{t("target_id")}</th>
            <th>{t("trigger_type")}</th>
            <th>{t("submitted_at")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((task) => {
            const busy = isReviewing && reviewingTaskId === task.id;
            return (
              <tr key={task.id}>
                <td>
                  <span className="badge">
                    {targetLabel(task.target_type)}
                  </span>
                </td>
                <td>#{task.target_id}</td>
                <td>
                  <span className="badge badge-ghost">
                    {triggerLabel(task.trigger_type)}
                  </span>
                </td>
                <td>{new Date(task.created_at).toLocaleString()}</td>
                <td>
                  {rejectingId === task.id ? (
                    <div className="flex gap-2">
                      <input
                        type="text"
                        placeholder={t("reject_reason")}
                        className="input input-xs input-bordered w-32"
                        value={rejectReason}
                        onChange={(e) => setRejectReason(e.target.value)}
                      />
                      <button
                        onClick={() => {
                          rejectTask(task.id, { note: rejectReason });
                          setRejectingId(null);
                          setRejectReason("");
                        }}
                        disabled={busy}
                        className="btn btn-xs btn-warning"
                      >
                        {t("confirm")}
                      </button>
                      <button
                        onClick={() => {
                          setRejectingId(null);
                          setRejectReason("");
                        }}
                        className="btn btn-xs btn-ghost"
                      >
                        {t("cancel")}
                      </button>
                    </div>
                  ) : (
                    <div className="flex gap-2">
                      <button
                        onClick={() => approveTask(task.id)}
                        disabled={busy}
                        className="btn btn-xs btn-success"
                      >
                        {busy ? t("reviewing") : t("approve")}
                      </button>
                      <button
                        onClick={() => setRejectingId(task.id)}
                        disabled={busy}
                        className="btn btn-xs btn-warning"
                      >
                        {t("reject")}
                      </button>
                    </div>
                  )}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
