// components/reviewer/PendingContentTable.tsx
"use client";

import { useTranslations } from "next-intl";
import Link from "next/link";

interface PendingContentTableProps {
  onApprove: (id: number, type: string) => void;
  onReject: (id: number, type: string, reason?: string) => void;
  onDelete: (id: number, type: string) => void;
  reviews: Review[];
}
interface Review {
  id: number;
  type: string;
  title: string;
  content: string;
  author_name: string;
  created_at: string;
}

export function PendingContentTable({
  onApprove,
  onReject,
  onDelete,
  reviews,
}: PendingContentTableProps) {
  // TODO: 从 props 获取数据
  const t = useTranslations("Reviewer");

  if (reviews.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_pending_content")}</p>
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
            <th>{t("content")}</th>
            <th>{t("author")}</th>
            <th>{t("submitted_at")}</th>
            <th>{t("actions")}</th>
          </tr>
        </thead>
        <tbody>
          {reviews.map((item: Review) => (
            <tr key={item.id}>
              <td>
                <span className="badge">{item.type}</span>
              </td>
              <td>
                <Link
                  href={`/${item.type}/${item.id}`}
                  className="hover:link-hover"
                >
                  {item.title || item.content?.substring(0, 50)}
                </Link>
              </td>
              <td>{item.author_name}</td>
              <td>{new Date(item.created_at).toLocaleString()}</td>
              <td>
                <div className="flex gap-2">
                  <button
                    onClick={() => onApprove(item.id, item.type)}
                    className="btn btn-xs btn-success"
                  >
                    {t("approve")}
                  </button>
                  <button
                    onClick={() => onReject(item.id, item.type)}
                    className="btn btn-xs btn-warning"
                  >
                    {t("reject")}
                  </button>
                  <button
                    onClick={() => onDelete(item.id, item.type)}
                    className="btn btn-xs btn-error"
                  >
                    {t("delete")}
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
