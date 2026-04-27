import { useAdminReports } from "@/hooks/admin/useAdminModerator";
import { Flag, CheckCircle, XCircle, Trash2 } from "lucide-react";
import { useState } from "react";

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

  if (isLoading)
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <input
          type="number"
          placeholder="板块 ID"
          value={selectedBoardId}
          onChange={(e) => setSelectedBoardId(parseInt(e.target.value))}
          className="input input-bordered w-32"
        />
        <select
          className="select select-bordered w-32"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
        >
          <option value="pending">待处理</option>
          <option value="resolved">已处理</option>
          <option value="rejected">已驳回</option>
        </select>
        <button className="btn btn-outline" onClick={() => refetch()}>
          刷新
        </button>
      </div>

      <div className="space-y-3">
        {reports?.list?.map((report: any) => (
          <div
            key={report.id}
            className="card bg-base-100 shadow-sm border border-base-200"
          >
            <div className="card-body">
              <div className="space-y-2">
                <div className="flex justify-between">
                  <div className="badge badge-error gap-1">
                    <Flag className="w-3 h-3" />
                    举报
                  </div>
                  <span className="text-sm text-gray-500">
                    {new Date(report.created_at).toLocaleString()}
                  </span>
                </div>
                <p className="text-sm">
                  <strong>举报人:</strong> {report.reporter?.username}
                </p>
                <p className="text-sm">
                  <strong>被举报内容:</strong> {report.content_preview}
                </p>
                <p className="text-sm">
                  <strong>原因:</strong> {report.reason}
                </p>
                {report.status === "pending" && (
                  <div className="flex gap-2 mt-2">
                    <button className="btn btn-sm btn-success">
                      <CheckCircle className="w-4 h-4 mr-1" />
                      通过
                    </button>
                    <button className="btn btn-sm btn-outline">
                      <XCircle className="w-4 h-4 mr-1" />
                      驳回
                    </button>
                    <button className="btn btn-sm btn-error">
                      <Trash2 className="w-4 h-4 mr-1" />
                      删除内容
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
        {reports?.list?.length === 0 && (
          <div className="text-center py-8 text-gray-500">暂无举报</div>
        )}
      </div>
    </div>
  );
}
