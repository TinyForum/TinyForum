import {
  useAdminApplications,
  useReviewApplication,
} from "@/hooks/admin/useAdminModerator";
import { CheckCircle, XCircle, ClipboardList } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";

// 类型定义
interface User {
  id: number;
  username: string;
  avatar?: string;
  email?: string;
}

interface Application {
  id: number;
  user_id: number;
  board_id: number;
  user?: User;
  reason: string;
  status: "pending" | "approved" | "rejected";
  review_note?: string;
  created_at: string;
  updated_at: string;
}

interface ApplicationsResponse {
  list: Application[];
  total: number;
  page: number;
  page_size: number;
}

interface ModeratorPermissions {
  can_delete_post: boolean;
  can_pin_post: boolean;
  can_edit_any_post: boolean;
  can_manage_moderator: boolean;
  can_ban_user: boolean;
}

const DEFAULT_PERMISSIONS: ModeratorPermissions = {
  can_delete_post: true,
  can_pin_post: true,
  can_edit_any_post: true,
  can_manage_moderator: false,
  can_ban_user: true,
};

export function ApplicationManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number | undefined>();
  const [status, setStatus] = useState<"pending" | "approved" | "rejected">(
    "pending",
  );
  const [selectedApp, setSelectedApp] = useState<Application | null>(null);
  const [reviewNote, setReviewNote] = useState<string>("");
  const [permissions, setPermissions] = useState<ModeratorPermissions>(DEFAULT_PERMISSIONS);

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

  // 获取申请列表
  const applicationList = (applications as ApplicationsResponse)?.list || [];

  const handleReview = (applicationId: number, approve: boolean) => {
    reviewApplication.mutate(
      {
        applicationId,
        data: {
          approve,
          review_note: reviewNote || undefined,
          ...(approve && permissions),
        },
      },
      {
        onSuccess: () => {
          toast.success(approve ? "申请已通过" : "申请已拒绝");
          setSelectedApp(null);
          setReviewNote("");
          setPermissions(DEFAULT_PERMISSIONS);
          refetch();
        },
        onError: () => {
          toast.error("操作失败");
        },
      },
    );
  };

  const openReviewDialog = (app: Application) => {
    setSelectedApp(app);
    setReviewNote("");
    setPermissions(DEFAULT_PERMISSIONS);
  };

  const getStatusBadge = (status: Application["status"]) => {
    switch (status) {
      case "pending":
        return <span className="badge badge-warning badge-sm">待审批</span>;
      case "approved":
        return <span className="badge badge-success badge-sm">已通过</span>;
      case "rejected":
        return <span className="badge badge-error badge-sm">已拒绝</span>;
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex gap-4 flex-wrap">
        <div className="form-control">
          <label className="label">
            <span className="label-text">板块 ID</span>
          </label>
          <input
            type="number"
            placeholder="板块 ID (可选)"
            value={selectedBoardId ?? ""}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setSelectedBoardId(
                e.target.value ? parseInt(e.target.value) : undefined,
              )
            }
            className="input input-bordered w-40"
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
              setStatus(e.target.value as "pending" | "approved" | "rejected")
            }
          >
            <option value="pending">待审批</option>
            <option value="approved">已通过</option>
            <option value="rejected">已拒绝</option>
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
        {applicationList.map((app: Application) => (
          <div
            key={app.id}
            className="card bg-base-100 shadow-sm border border-base-200 hover:shadow-md transition-shadow"
          >
            <div className="card-body p-4">
              <div className="flex justify-between items-start flex-wrap gap-2">
                <div className="space-y-2 flex-1">
                  <div className="flex items-center gap-3">
                    <div className="avatar placeholder">
                      <div className="bg-primary/10 rounded-full w-10 h-10 flex items-center justify-center">
                        <span className="text-primary font-medium">
                          {app.user?.username?.[0]?.toUpperCase() || "U"}
                        </span>
                      </div>
                    </div>
                    <div>
                      <p className="font-medium">
                        {app.user?.username || `用户${app.user_id}`}
                      </p>
                      <p className="text-xs text-base-content/50">
                        申请版块 ID: {app.board_id}
                      </p>
                    </div>
                  </div>
                  <p className="text-sm text-base-content/70">{app.reason}</p>
                  {app.review_note && (
                    <p className="text-xs text-base-content/50">
                      审批备注: {app.review_note}
                    </p>
                  )}
                  <p className="text-xs text-base-content/40">
                    申请时间: {new Date(app.created_at).toLocaleString()}
                  </p>
                  {app.status === "pending" && (
                    <div className="flex gap-2 mt-2">
                      <button
                        className="btn btn-sm btn-primary"
                        onClick={() => openReviewDialog(app)}
                      >
                        <CheckCircle className="w-3 h-3 mr-1" />
                        审批
                      </button>
                    </div>
                  )}
                </div>
                <div>{getStatusBadge(app.status)}</div>
              </div>
            </div>
          </div>
        ))}
        
        {applicationList.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            <ClipboardList className="w-12 h-12 mx-auto mb-3 opacity-30" />
            <p>暂无申请记录</p>
          </div>
        )}
      </div>

      {/* 审批对话框 */}
      {selectedApp && (
        <dialog
          className="modal modal-open"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setSelectedApp(null);
            }
          }}
        >
          <div
            className="modal-box max-w-md"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="font-bold text-lg">审批版主申请</h3>
            <div className="space-y-4 mt-4">
              <div className="space-y-2 p-3 bg-base-200/50 rounded-lg">
                <p className="text-sm">
                  <strong className="text-base-content/70">申请人:</strong>{" "}
                  {selectedApp.user?.username || `用户${selectedApp.user_id}`}
                </p>
                <p className="text-sm">
                  <strong className="text-base-content/70">申请版块:</strong>{" "}
                  {selectedApp.board_id}
                </p>
                <p className="text-sm">
                  <strong className="text-base-content/70">申请理由:</strong>{" "}
                  {selectedApp.reason}
                </p>
                <p className="text-sm">
                  <strong className="text-base-content/70">申请时间:</strong>{" "}
                  {new Date(selectedApp.created_at).toLocaleString()}
                </p>
              </div>

              <div className="divider my-2">审批操作</div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">审批备注</span>
                  <span className="label-text-alt text-base-content/50">
                    可选
                  </span>
                </label>
                <textarea
                  placeholder="填写审批备注..."
                  value={reviewNote}
                  onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                    setReviewNote(e.target.value)
                  }
                  className="textarea textarea-bordered w-full"
                  rows={3}
                />
              </div>

              {/* 通过时显示权限设置 */}
              <div className="space-y-2">
                <p className="font-medium text-sm">版主权限设置</p>
                <p className="text-xs text-base-content/50 mb-2">
                  (仅通过审批时生效)
                </p>
                <div className="space-y-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={permissions.can_delete_post}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setPermissions({
                          ...permissions,
                          can_delete_post: e.target.checked,
                        })
                      }
                      className="checkbox checkbox-sm checkbox-primary"
                    />
                    <span className="text-sm">可以删除帖子</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={permissions.can_pin_post}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setPermissions({
                          ...permissions,
                          can_pin_post: e.target.checked,
                        })
                      }
                      className="checkbox checkbox-sm checkbox-primary"
                    />
                    <span className="text-sm">可以置顶帖子</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={permissions.can_edit_any_post}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setPermissions({
                          ...permissions,
                          can_edit_any_post: e.target.checked,
                        })
                      }
                      className="checkbox checkbox-sm checkbox-primary"
                    />
                    <span className="text-sm">可以编辑任何帖子</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={permissions.can_manage_moderator}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setPermissions({
                          ...permissions,
                          can_manage_moderator: e.target.checked,
                        })
                      }
                      className="checkbox checkbox-sm checkbox-primary"
                    />
                    <span className="text-sm">可以管理其他版主</span>
                    <span className="text-xs text-warning ml-1">
                      (高级权限)
                    </span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={permissions.can_ban_user}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                        setPermissions({
                          ...permissions,
                          can_ban_user: e.target.checked,
                        })
                      }
                      className="checkbox checkbox-sm checkbox-primary"
                    />
                    <span className="text-sm">可以禁言用户</span>
                  </label>
                </div>
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
                {reviewApplication.isPending ? (
                  <span className="loading loading-spinner loading-xs" />
                ) : (
                  <XCircle className="w-4 h-4 mr-1" />
                )}
                拒绝
              </button>
              <button
                className="btn btn-success"
                onClick={() => handleReview(selectedApp.id, true)}
                disabled={reviewApplication.isPending}
              >
                {reviewApplication.isPending ? (
                  <span className="loading loading-spinner loading-xs" />
                ) : (
                  <CheckCircle className="w-4 h-4 mr-1" />
                )}
                通过
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}