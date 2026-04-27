import {
  useAdminApplications,
  useReviewApplication,
} from "@/hooks/admin/useAdminModerator";
import { CheckCircle, XCircle } from "lucide-react";
import { useState } from "react";

export function ApplicationManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>();
  const [status, setStatus] = useState<"pending" | "approved" | "rejected">(
    "pending",
  );
  const [selectedApp, setSelectedApp] = useState<any>(null);
  const [reviewNote, setReviewNote] = useState("");
  const [permissions, setPermissions] = useState({
    can_delete_post: true,
    can_pin_post: true,
    can_edit_any_post: true,
    can_manage_moderator: false,
    can_ban_user: true,
  });

  const {
    data: applications,
    isLoading,
    refetch,
  } = useAdminApplications({
    board_id: selectedBoardId,
    status,
    page: 1,
    page_size: 20,
  });

  const reviewApplication = useReviewApplication();

  const handleReview = (applicationId: number, approve: boolean) => {
    reviewApplication.mutate({
      applicationId,
      data: {
        approve,
        review_note: reviewNote || undefined,
        ...(approve && permissions), // 只有通过时才设置权限
      },
    });
    setSelectedApp(null);
    setReviewNote("");
  };

  const openReviewDialog = (app: any) => {
    setSelectedApp(app);
    setReviewNote("");
    setPermissions({
      can_delete_post: true,
      can_pin_post: true,
      can_edit_any_post: true,
      can_manage_moderator: false,
      can_ban_user: true,
    });
  };

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
          placeholder="板块 ID (可选)"
          value={selectedBoardId || ""}
          onChange={(e) =>
            setSelectedBoardId(
              e.target.value ? parseInt(e.target.value) : undefined,
            )
          }
          className="input input-bordered w-40"
        />
        <select
          className="select select-bordered w-32"
          value={status}
          onChange={(e) => setStatus(e.target.value as any)}
        >
          <option value="pending">待审批</option>
          <option value="approved">已通过</option>
          <option value="rejected">已拒绝</option>
        </select>
        <button className="btn btn-outline" onClick={() => refetch()}>
          刷新
        </button>
      </div>

      <div className="space-y-3">
        {applications?.list?.map((app: any) => (
          <div
            key={app.id}
            className="card bg-base-100 shadow-sm border border-base-200"
          >
            <div className="card-body">
              <div className="flex justify-between items-start">
                <div className="space-y-2 flex-1">
                  <div className="flex items-center gap-3">
                    <div className="avatar placeholder">
                      <div className="bg-primary/10 rounded-full w-10">
                        <span className="text-primary">
                          {app.user?.username?.[0]?.toUpperCase()}
                        </span>
                      </div>
                    </div>
                    <div>
                      <p className="font-medium">{app.user?.username}</p>
                      <p className="text-sm text-gray-500">
                        申请版块 ID: {app.board_id}
                      </p>
                    </div>
                  </div>
                  <p className="text-sm text-gray-600">{app.reason}</p>
                  {app.review_note && (
                    <p className="text-sm text-gray-500">
                      审批备注: {app.review_note}
                    </p>
                  )}
                  {app.status === "pending" && (
                    <div className="flex gap-2 mt-2">
                      <button
                        className="btn btn-sm btn-success"
                        onClick={() => openReviewDialog(app)}
                      >
                        <CheckCircle className="w-4 h-4 mr-1" />
                        审批
                      </button>
                    </div>
                  )}
                </div>
                <div
                  className={`badge badge-lg ${
                    app.status === "pending"
                      ? "badge-warning"
                      : app.status === "approved"
                        ? "badge-success"
                        : "badge-error"
                  }`}
                >
                  {app.status === "pending"
                    ? "待审批"
                    : app.status === "approved"
                      ? "已通过"
                      : "已拒绝"}
                </div>
              </div>
            </div>
          </div>
        ))}
        {applications?.list?.length === 0 && (
          <div className="text-center py-8 text-gray-500">暂无申请记录</div>
        )}
      </div>

      {/* 审批对话框 */}
      {selectedApp && (
        <dialog
          className="modal modal-open"
          onClick={() => setSelectedApp(null)}
        >
          <div
            className="modal-box max-w-md"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="font-bold text-lg">审批版主申请</h3>
            <div className="space-y-4 mt-4">
              <div className="space-y-2">
                <p className="text-sm">
                  <strong>申请人:</strong> {selectedApp.user?.username}
                </p>
                <p className="text-sm">
                  <strong>申请版块:</strong> {selectedApp.board_id}
                </p>
                <p className="text-sm">
                  <strong>申请理由:</strong> {selectedApp.reason}
                </p>
              </div>

              <div className="divider">审批操作</div>

              <textarea
                placeholder="审批备注 (可选)"
                value={reviewNote}
                onChange={(e) => setReviewNote(e.target.value)}
                className="textarea textarea-bordered w-full"
                rows={3}
              />

              {/* 通过时显示权限设置 */}
              <div className="space-y-2">
                <p className="font-medium">版主权限设置 (通过时生效)</p>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_delete_post}
                    onChange={(e) =>
                      setPermissions({
                        ...permissions,
                        can_delete_post: e.target.checked,
                      })
                    }
                    className="checkbox checkbox-sm"
                  />
                  <span>可以删除帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_pin_post}
                    onChange={(e) =>
                      setPermissions({
                        ...permissions,
                        can_pin_post: e.target.checked,
                      })
                    }
                    className="checkbox checkbox-sm"
                  />
                  <span>可以置顶帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_edit_any_post}
                    onChange={(e) =>
                      setPermissions({
                        ...permissions,
                        can_edit_any_post: e.target.checked,
                      })
                    }
                    className="checkbox checkbox-sm"
                  />
                  <span>可以编辑任何帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_manage_moderator}
                    onChange={(e) =>
                      setPermissions({
                        ...permissions,
                        can_manage_moderator: e.target.checked,
                      })
                    }
                    className="checkbox checkbox-sm"
                  />
                  <span>可以管理其他版主</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_ban_user}
                    onChange={(e) =>
                      setPermissions({
                        ...permissions,
                        can_ban_user: e.target.checked,
                      })
                    }
                    className="checkbox checkbox-sm"
                  />
                  <span>可以禁言用户</span>
                </label>
              </div>
            </div>
            <div className="modal-action">
              <button className="btn" onClick={() => setSelectedApp(null)}>
                取消
              </button>
              <button
                className="btn btn-error"
                onClick={() => handleReview(selectedApp.id, false)}
                disabled={reviewApplication.isPending}
              >
                <XCircle className="w-4 h-4 mr-1" />
                拒绝
              </button>
              <button
                className="btn btn-success"
                onClick={() => handleReview(selectedApp.id, true)}
                disabled={reviewApplication.isPending}
              >
                <CheckCircle className="w-4 h-4 mr-1" />
                通过
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}
