// components/reviewer/ReviewHistoryTable.tsx
"use client";

import { useTranslations } from "next-intl";

interface ReviewHistoryTableProps {
  historys: History[];
}
interface History {
  id: number;
  content_type: string;
  content_title: string;
  action: string;
  reason: string;
  created_at: string;
}

export function ReviewHistoryTable({ historys }: ReviewHistoryTableProps) {
  // TODO: 从 propss 获取数据
  const t = useTranslations("Review");

  if (historys.length === 0) {
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
            <th>{t("content_type")}</th>
            <th>{t("content")}</th>
            <th>{t("action")}</th>
            <th>{t("reason")}</th>
            <th>{t("reviewed_at")}</th>
          </tr>
        </thead>
        <tbody>
          {historys.map((item: History) => (
            <tr key={item.id}>
              <td>
                <span className="badge">{item.content_type}</span>
              </td>
              <td>
                <span className="text-sm">{item.content_title}</span>
              </td>
              <td>
                <span
                  className={`badge ${
                    item.action === "approved"
                      ? "badge-success"
                      : item.action === "rejected"
                        ? "badge-warning"
                        : "badge-error"
                  }`}
                >
                  {t(item.action)}
                </span>
              </td>
              <td className="text-sm text-base-content/60">
                {item.reason || "-"}
              </td>
              <td className="text-sm">
                {new Date(item.created_at).toLocaleString()}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
