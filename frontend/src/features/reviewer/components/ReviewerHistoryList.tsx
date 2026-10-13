"use client";

import { useTranslations } from "next-intl";
import { useReviewerHistory } from "../hooks/useReviewerHistory";

// 审核历史列表：展示已处理（通过/拒绝）的审核任务
export function ReviewerHistoryList() {
  const t = useTranslations("Reviewer");
  const { list, isLoading, error, refresh } = useReviewerHistory();

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

  const statusBadge = (status: string) => {
    switch (status) {
      case "approved":
        return <span className="badge badge-success">{t("approved")}</span>;
      case "rejected":
        return <span className="badge badge-error">{t("rejected")}</span>;
      default:
        return <span className="badge">{status}</span>;
    }
  };

  if (isLoading && list.length === 0) {
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

  if (list.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_review_history")}</p>
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
            <th>{t("reviewer")}</th>
            <th>{t("reviewed_at")}</th>
            <th>{t("action")}</th>
          </tr>
        </thead>
        <tbody>
          {list.map((task) => (
            <tr key={task.id}>
              <td>
                <span className="badge">{targetLabel(task.target_type)}</span>
              </td>
              <td>#{task.target_id}</td>
              <td>
                <span className="badge badge-ghost">
                  {triggerLabel(task.trigger_type)}
                </span>
              </td>
              <td>{task.reviewer_name || "-"}</td>
              <td>
                {task.reviewed_at
                  ? new Date(task.reviewed_at).toLocaleString()
                  : "-"}
              </td>
              <td>{statusBadge(task.status)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
