import { Flag, CheckCircle, XCircle, Trash2 } from "lucide-react";
import { useState } from "react";

import type { ReportResponse } from "@/shared/api/modules/admin/report";
import {
  useAdminGetReports,
  useAdminHandleReport,
} from "../hooks/useAdminReport";

// 扩展类型（如果目标对象信息需要）
interface ReportWithTarget extends ReportResponse {
  target_title?: string;
  content_preview?: string;
}

export function ReportManagement() {
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

  // 处理通过举报
  const handleResolve = (reportId: number) => {
    handleReport({
      id: reportId,
      data: {
        status: "resolved",
        handle_note: "内容违规，已处理",
      },
    });
  };

  // 处理驳回举报
  const handleReject = (reportId: number, reason?: string) => {
    const rejectReason = reason || "举报不成立";
    handleReport({
      id: reportId,
      data: {
        status: "rejected",
        handle_note: rejectReason,
        reject_reason: rejectReason,
      },
    });
  };

  // 删除被举报内容（需单独实现 API，暂留空）
  const handleDeleteContent = (reportId: number) => {
    console.log("删除内容:", reportId);
    // TODO: 调用删除内容 API
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
            <span className="label-text">状态</span>
          </label>
          <select
            className="select select-bordered w-32"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
          >
            <option value="pending">待处理</option>
            <option value="resolved">已处理</option>
            <option value="rejected">已驳回</option>
          </select>
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text">&nbsp;</span>
          </label>
          <button className="btn btn-outline" onClick={() => refetch()}>
            刷新
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
                      {report.status === "pending" && "待处理"}
                      {report.status === "resolved" && "已处理"}
                      {report.status === "rejected" && "已驳回"}
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
                  <strong className="text-base-content/70">举报人:</strong>{" "}
                  {report.is_anonymous
                    ? "匿名用户"
                    : report.reporter?.nickname || `用户${report.reporter?.id}`}
                </p>

                {/* 目标标题如果有，可以从 target_info 获取 */}
                {(report as ReportWithTarget).target_title && (
                  <p className="text-sm">
                    <strong className="text-base-content/70">标题:</strong>{" "}
                    {(report as ReportWithTarget).target_title}
                  </p>
                )}

                <p className="text-sm">
                  <strong className="text-base-content/70">被举报内容:</strong>{" "}
                  <span className="line-clamp-2">
                    {report.content_snapshot || "无内容预览"}
                  </span>
                </p>

                <p className="text-sm">
                  <strong className="text-base-content/70">原因:</strong>{" "}
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
                      通过
                    </button>
                    <button
                      className="btn btn-sm btn-outline gap-1"
                      onClick={() => handleReject(report.id)}
                      disabled={isHandling}
                    >
                      <XCircle className="w-4 h-4" />
                      驳回
                    </button>
                    <button
                      className="btn btn-sm btn-error gap-1"
                      onClick={() => handleDeleteContent(report.id)}
                      disabled={isHandling}
                    >
                      <Trash2 className="w-4 h-4" />
                      删除内容
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
            <p>暂无举报</p>
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
              第 {page} / {Math.ceil(pageData.total / pageSize)} 页
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
