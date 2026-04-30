// components/reviewer/ReportedContentTable.tsx
"use client";

import { useState } from "react";
import { ModeratorReport } from "@/shared/api/modules/moderator";
import Link from "next/link";
import { useTranslations } from "next-intl";

interface ReportedContentTableProps {
  reports: ModeratorReport[];
  onApprove?: (id: number, type: string) => void;
  onReject?: (id: number, type: string, reason?: string) => void;
  onDelete?: (id: number, type: string) => void;
  onBanUser?: (data: {
    userId: number;
    reason: string;
    expiresAt?: string;
  }) => void;
  isDeleting?: boolean;
  isBanning?: boolean;
  permissions?: {
    canDeletePost: boolean;
    canBanUser: boolean;
  };
}

export function ReportedContentTable({
  reports,
  onApprove,
  onReject,
  onDelete,
  onBanUser,
  isDeleting = false,
  isBanning = false,
  permissions,
}: ReportedContentTableProps) {
  const [banReason, setBanReason] = useState("");
  const [showBanModal, setShowBanModal] = useState<number | null>(null);
  const [rejectReason, setRejectReason] = useState("");
  const [showRejectModal, setShowRejectModal] = useState<number | null>(null);
  const t = useTranslations("Reviewer");

  if (reports.length === 0) {
    return (
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body text-center py-12">
          <p className="text-base-content/50">{t("no_reports")}</p>
        </div>
      </div>
    );
  }

  const handleBanUser = (userId: number) => {
    if (onBanUser && banReason) {
      onBanUser({
        userId,
        reason: banReason,
      });
      setShowBanModal(null);
      setBanReason("");
    }
  };

  const handleReject = (id: number, type: string) => {
    if (onReject) {
      onReject(id, type, rejectReason);
      setShowRejectModal(null);
      setRejectReason("");
    }
  };

  return (
    <div className="space-y-4">
      {reports.map((report) => (
        <div
          key={report.id}
          className="card bg-base-100 border border-base-300"
        >
          <div className="card-body">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-2">
                  <span className="badge badge-error">{t("report")}</span>
                  <span className="badge badge-sm">
                    {report.target_type === "post" ? t("post") : t("comment")}
                  </span>
                  <span className="text-sm text-base-content/60">
                    {t("reported_by")}: {report.reporter_name}
                  </span>
                  <span className="text-sm text-base-content/60">
                    {new Date(report.created_at).toLocaleString()}
                  </span>
                  <span
                    className={`badge badge-sm ${
                      report.status === "pending"
                        ? "badge-warning"
                        : report.status === "resolved"
                          ? "badge-success"
                          : "badge-secondary"
                    }`}
                  >
                    {t(report.status)}
                  </span>
                </div>
                <p className="font-medium">
                  {t("reason")}: {report.reason}
                </p>
              </div>
            </div>

            <div className="border-t border-base-300 pt-4 mt-2">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <p className="text-sm font-medium">{t("reported_content")}</p>
                  <Link
                    href={`/${report.target_type}/${report.target_id}`}
                    className="text-sm hover:link-hover"
                  >
                    {report.target_title ||
                      report.target_content?.substring(0, 100) ||
                      t("view_content")}
                  </Link>
                  {report.target_content && (
                    <p className="text-xs text-base-content/60 mt-1 line-clamp-2">
                      {report.target_content}
                    </p>
                  )}
                </div>
                <div className="flex gap-2 ml-4">
                  {onApprove && report.status === "pending" && (
                    <button
                      onClick={() =>
                        onApprove(report.target_id, report.target_type)
                      }
                      className="btn btn-xs btn-success"
                    >
                      {t("approve")}
                    </button>
                  )}

                  {onReject && report.status === "pending" && (
                    <>
                      {showRejectModal === report.id ? (
                        <div className="flex gap-2">
                          <input
                            type="text"
                            placeholder={t("reject_reason")}
                            className="input input-xs input-bordered w-32"
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                          />
                          <button
                            onClick={() =>
                              handleReject(report.target_id, report.target_type)
                            }
                            className="btn btn-xs btn-warning"
                          >
                            {t("confirm")}
                          </button>
                          <button
                            onClick={() => setShowRejectModal(null)}
                            className="btn btn-xs btn-ghost"
                          >
                            {t("cancel")}
                          </button>
                        </div>
                      ) : (
                        <button
                          onClick={() => setShowRejectModal(report.id)}
                          className="btn btn-xs btn-warning"
                        >
                          {t("reject")}
                        </button>
                      )}
                    </>
                  )}

                  {permissions?.canDeletePost &&
                    report.status === "pending" && (
                      <button
                        onClick={() =>
                          onDelete?.(report.target_id, report.target_type)
                        }
                        disabled={isDeleting}
                        className="btn btn-xs btn-error"
                      >
                        {isDeleting ? t("deleting") : t("delete")}
                      </button>
                    )}

                  {permissions?.canBanUser && report.status === "pending" && (
                    <>
                      {showBanModal === report.id ? (
                        <div className="flex gap-2">
                          <input
                            type="text"
                            placeholder={t("ban_reason")}
                            className="input input-xs input-bordered w-32"
                            value={banReason}
                            onChange={(e) => setBanReason(e.target.value)}
                          />
                          <button
                            onClick={() => handleBanUser(report.reporter_id)}
                            disabled={isBanning || !banReason}
                            className="btn btn-xs btn-error"
                          >
                            {isBanning ? t("banning") : t("confirm")}
                          </button>
                          <button
                            onClick={() => setShowBanModal(null)}
                            className="btn btn-xs btn-ghost"
                          >
                            {t("cancel")}
                          </button>
                        </div>
                      ) : (
                        <button
                          onClick={() => setShowBanModal(report.id)}
                          className="btn btn-xs btn-error"
                        >
                          {t("ban_user")}
                        </button>
                      )}
                    </>
                  )}

                  {report.status !== "pending" && (
                    <span className="text-xs text-base-content/40">
                      {t("already_processed")}
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
