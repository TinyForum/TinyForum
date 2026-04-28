import { useAdminReports } from "@/hooks/admin/useAdminModerator";
import { Flag, CheckCircle, XCircle, Trash2 } from "lucide-react";
import { useState } from "react";

// 类型定义
interface Reporter {
  id: number;
  username: string;
  avatar?: string;
}

interface Report {
  id: number;
  reporter_id: number;
  reporter?: Reporter;
  target_id: number;
  target_type: string;
  target_title?: string;
  content_preview: string;
  reason: string;
  status: "pending" | "resolved" | "rejected";
  created_at: string;
  updated_at: string;
}

interface ReportsResponse {
  list: Report[];
  total: number;
  page: number;
  page_size: number;
}

export function ReportManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [status, setStatus] = useState<string>("pending");

  const {
    data: reports,
    isLoading,
    refetch,
  } = useAdminReports(selectedBoardId, {
    page: 1,
    page_size: 20,
    status,
  });

  const handleResolve = (reportId: number) => {
    // TODO: 实现通过举报
    console.log("通过举报:", reportId);
  };

  const handleReject = (reportId: number) => {
    // TODO: 实现驳回举报
    console.log("驳回举报:", reportId);
  };

  const handleDeleteContent = (reportId: number) => {
    // TODO: 实现删除被举报内容
    console.log("删除内容:", reportId);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  // 安全获取报告列表
  const reportList = (reports as unknown as ReportsResponse)?.list || [];

  return (
    <div className="space-y-4">
      <div className="flex gap-4 flex-wrap">
        <div className="form-control">
          <label className="label">
            <span className="label-text">板块 ID</span>
          </label>
          <input
            type="number"
            placeholder="板块 ID"
            value={selectedBoardId}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setSelectedBoardId(parseInt(e.target.value) || 1)
            }
            className="input input-bordered w-32"
          />
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text">状态</span>
          </label>
          <select
            className="select select-bordered w-32"
            value={status}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
              setStatus(e.target.value)
            }
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
        {reportList.map((report: Report) => (
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
                  {report.reporter?.username || `用户${report.reporter_id}`}
                </p>

                {report.target_title && (
                  <p className="text-sm">
                    <strong className="text-base-content/70">标题:</strong>{" "}
                    {report.target_title}
                  </p>
                )}

                <p className="text-sm">
                  <strong className="text-base-content/70">被举报内容:</strong>{" "}
                  <span className="line-clamp-2">{report.content_preview}</span>
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
                    >
                      <CheckCircle className="w-4 h-4" />
                      通过
                    </button>
                    <button
                      className="btn btn-sm btn-outline gap-1"
                      onClick={() => handleReject(report.id)}
                    >
                      <XCircle className="w-4 h-4" />
                      驳回
                    </button>
                    <button
                      className="btn btn-sm btn-error gap-1"
                      onClick={() => handleDeleteContent(report.id)}
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
    </div>
  );
}
