import { Flag, CheckCircle, XCircle, Trash2 } from "lucide-react";
import { useState } from "react";
import { useTranslations } from "next-intl";

import type { ReportResponse } from "@/shared/api/modules/admin/report";
import {
  useAdminGetReports,
  useAdminHandleReport,
  useAdminDeleteReportedContent,
} from "../hooks/useAdminReport";

// 扩展类型（如果目标对象信息需要）
interface ReportWithTarget extends ReportResponse {
  target_title?: string;
  content_preview?: string;
}

export function ReportManagement() {
  const t = useTranslations("Admin");
  const [status, setStatus] = useState<string>("pending");
  const [page, setPage] = useState(1);
  const pageSize = 20;

  // 使用标准 hook 获取举报列表
  const {
    data: pageData,
    isLoading,
    refetch,
  } = useAdminGetReports({
    page,
    pageSize,
    status: status,
  });

  const { mutate: handleReport, isPending: isHandling } =
    useAdminHandleReport();

  const { mutate: deleteContent, isPending: isDeletingContent } =
    useAdminDeleteReportedContent();

    // 处理通过举报
  const handleResolve = (reportId: number) => {
    handleReport({
      id: reportId,
      data: {
        status: "resolved",
        handle_note: t("content_violation_handled"),
      },
    });
  };

  // 处理驳回举报
  const handleReject = (reportId: number, reason?: string) => {
    const rejectReason = reason || t("report_not_established");
    handleReport({
      id: reportId,
      data: {
        status: "rejected",
        handle_note: rejectReason,
        reject_reason: rejectReason,
      },
    });
  };

  // 删除被举报内容（按 target_type 分派到帖子/评论/用户删除接口）
  const handleDeleteContent = (report: ReportResponse) => {
    if (!confirm(t("delete_content_confirm"))) return;
    deleteContent({ targetType: report.target_type, targetId: report.target_id });
  };

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  const reportList = pageData?.list || [];

  return (
    <div className="space-y-4">
      <div className="flex gap-4 flex-wrap">
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t("status")}</span>
          </label>
          <select
            className="select select-bordered w-32"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
          >
            <option value="pending">{t("pending")}</option>
            <option value="resolved">{t("resolved")}</option>
            <option value="rejected">{t("reject_success")}</option>
          </select>
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text">&nbsp;</span>
          </label>
          <button className="btn btn-outline" onClick={() => refetch()}>
            {t("refresh")}
          </button>
        </div>
      </div>

      <div className="space-y-3">
        {reportList.map((report: ReportResponse) => (
          <div
            key={report.id}
            className="card bg-base-100 shadow-sm border border-base-200 hover:shadow-md transition-shadow"
          >
            <div className="card-body p-4">
              <div className="space-y-2">
                <div className="flex justify-between items-center flex-wrap gap-2">
                  <div className="flex gap-2">
                    <div
                      className={`badge gap-1 ${
                        report.status === "pending"
                          ? "badge-error"
                          : report.status === "resolved"
                            ? "badge-success"
                            : "badge-warning"
                      }`}
                    >
                      <Flag className="w-3 h-3" />
                      {report.status === "pending" && t("pending")}
                      {report.status === "resolved" && t("resolved")}
                      {report.status === "rejected" && t("reject_success")}
                    </div>
                    <div className="badge badge-ghost badge-sm">
                      {report.target_type}
                    </div>
                  </div>
                  <span className="text-xs text-base-content/50">
                    {new Date(report.created_at).toLocaleString()}
                  </span>
                </div>

                <p className="text-sm">
                  <strong className="text-base-content/70">{t("reporter")}:</strong>{" "}
                  {report.is_anonymous
                    ? t("anonymous_user")
                    : report.reporter?.nickname || `用户${report.reporter?.id}`}
                </p>

                {/* 目标标题如果有，可以从 target_info 获取 */}
                {(report as ReportWithTarget).target_title && (
                  <p className="text-sm">
                    <strong className="text-base-content/70">{t("title")}:</strong>{" "}
                    {(report as ReportWithTarget).target_title}
                  </p>
                )}

                <p className="text-sm">
                  <strong className="text-base-content/70">{t("reported_content")}:</strong>{" "}
                  <span className="line-clamp-2">
                    {report.content_snapshot || t("no_content_preview")}
                  </span>
                </p>

                <p className="text-sm">
                  <strong className="text-base-content/70">{t("reason")}:</strong>{" "}
                  {report.reason}
                </p>

                {report.status === "pending" && (
                  <div className="flex gap-2 mt-2 flex-wrap">
                    <button
                      className="btn btn-sm btn-success gap-1"
                      onClick={() => handleResolve(report.id)}
                      disabled={isHandling}
                    >
                      <CheckCircle className="w-4 h-4" />
                      {t("approve")}
                    </button>
                    <button
                      className="btn btn-sm btn-outline gap-1"
                      onClick={() => handleReject(report.id)}
                      disabled={isHandling}
                    >
                      <XCircle className="w-4 h-4" />
                      {t("reject")}
                    </button>
                    <button
                      className="btn btn-sm btn-error gap-1"
                      onClick={() => handleDeleteContent(report)}
                      disabled={isHandling || isDeletingContent}
                    >
                      <Trash2 className="w-4 h-4" />
                      {t("delete_content")}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}

        {reportList.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            <Flag className="w-12 h-12 mx-auto mb-3 opacity-30" />
            <p>{t("no_reports")}</p>
          </div>
        )}
      </div>

      {/* 分页组件可选 */}
      {pageData && pageData.total > pageSize && (
        <div className="flex justify-center mt-4">
          <div className="join">
            <button
              className="join-item btn btn-sm"
              disabled={page <= 1}
              onClick={() => setPage((p) => p - 1)}
            >
              «
            </button>
            <span className="join-item btn btn-sm">
              {t("page_info", { page, totalPages: Math.ceil(pageData.total / pageSize) })}
            </span>
            <button
              className="join-item btn btn-sm"
              disabled={page >= Math.ceil(pageData.total / pageSize)}
              onClick={() => setPage((p) => p + 1)}
            >
              »
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
