// app/admin/tasks/page.tsx
"use client";

import { useState } from "react";
import toast from "react-hot-toast";
import {
  Users,
  Shield,
  Flag,
  FileText,
  CheckCircle,
  XCircle,
  Eye,
  Search,
  UserPlus,
  UserMinus,
  Ban,
  Unlock,
  Pin,
  Trash2,
  Clock,
  MoreVertical,
} from "lucide-react";
import { useAddModerator, useAdminApplications, useAdminBannedUsers, useAdminBanUser, useAdminBoardPosts, useAdminDeletePost, useAdminModeratorList, useAdminPinPost, useAdminReports, useAdminUnbanUser, useRemoveModerator, useReviewApplication } from "@/hooks/admin/useAdminModerator";

// Hooks
// import {
//   useAdminApplications,
//   useReviewApplication,
//   useAddModerator,
//   useRemoveModerator,
//   useUpdateModeratorPermissions,
//   useAdminDeletePost,
//   useAdminPinPost,
//   useAdminBanUser,
//   useAdminUnbanUser,
//   useAdminModeratorList,
//   useAdminBannedUsers,
//   useAdminReports,
//   useAdminBoardPosts,
// } from "@/hooks/useAdminModerator";

// ============ 版主申请管理模块 ============
// ============ 版主申请管理模块 ============
function ApplicationManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>();
  const [status, setStatus] = useState<"pending" | "approved" | "rejected">("pending");
  const [selectedApp, setSelectedApp] = useState<any>(null);
  const [reviewNote, setReviewNote] = useState("");
  const [permissions, setPermissions] = useState({
    can_delete_post: true,
    can_pin_post: true,
    can_edit_any_post: true,
    can_manage_moderator: false,
    can_ban_user: true,
  });
  
  const { data: applications, isLoading, refetch } = useAdminApplications({
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
        ...(approve && permissions) // 只有通过时才设置权限
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

  if (isLoading) return <div className="flex justify-center py-8"><span className="loading loading-spinner loading-lg"></span></div>;

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <input
          type="number"
          placeholder="板块 ID (可选)"
          value={selectedBoardId || ""}
          onChange={(e) => setSelectedBoardId(e.target.value ? parseInt(e.target.value) : undefined)}
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
        <button className="btn btn-outline" onClick={() => refetch()}>刷新</button>
      </div>

      <div className="space-y-3">
        {applications?.list?.map((app: any) => (
          <div key={app.id} className="card bg-base-100 shadow-sm border border-base-200">
            <div className="card-body">
              <div className="flex justify-between items-start">
                <div className="space-y-2 flex-1">
                  <div className="flex items-center gap-3">
                    <div className="avatar placeholder">
                      <div className="bg-primary/10 rounded-full w-10">
                        <span className="text-primary">{app.user?.username?.[0]?.toUpperCase()}</span>
                      </div>
                    </div>
                    <div>
                      <p className="font-medium">{app.user?.username}</p>
                      <p className="text-sm text-gray-500">申请版块 ID: {app.board_id}</p>
                    </div>
                  </div>
                  <p className="text-sm text-gray-600">{app.reason}</p>
                  {app.review_note && (
                    <p className="text-sm text-gray-500">审批备注: {app.review_note}</p>
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
                <div className={`badge badge-lg ${
                  app.status === "pending" ? "badge-warning" : 
                  app.status === "approved" ? "badge-success" : "badge-error"
                }`}>
                  {app.status === "pending" ? "待审批" : app.status === "approved" ? "已通过" : "已拒绝"}
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
        <dialog className="modal modal-open" onClick={() => setSelectedApp(null)}>
          <div className="modal-box max-w-md" onClick={(e) => e.stopPropagation()}>
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
                    onChange={(e) => setPermissions({ ...permissions, can_delete_post: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以删除帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_pin_post}
                    onChange={(e) => setPermissions({ ...permissions, can_pin_post: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以置顶帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_edit_any_post}
                    onChange={(e) => setPermissions({ ...permissions, can_edit_any_post: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以编辑任何帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_manage_moderator}
                    onChange={(e) => setPermissions({ ...permissions, can_manage_moderator: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以管理其他版主</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_ban_user}
                    onChange={(e) => setPermissions({ ...permissions, can_ban_user: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以禁言用户</span>
                </label>
              </div>
            </div>
            <div className="modal-action">
              <button className="btn" onClick={() => setSelectedApp(null)}>取消</button>
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

// ============ 版主任命管理模块 ============
function ModeratorManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [userId, setUserId] = useState("");
  const [permissions, setPermissions] = useState({
    can_delete_post: true,
    can_pin_post: true,
    can_edit_any_post: true,
    can_manage_moderator: false,
    can_ban_user: true,
  });
  
  const { data: moderators, isLoading, refetch } = useAdminModeratorList(selectedBoardId);
  const addModerator = useAddModerator(selectedBoardId);
  const removeModerator = useRemoveModerator(selectedBoardId);

  const handleAddModerator = () => {
    if (!userId) {
      toast.error("请输入用户 ID");
      return;
    }
    addModerator.mutate({
      user_id: parseInt(userId),
      ...permissions,
    });
    setIsAddDialogOpen(false);
    setUserId("");
  };

  const handleRemoveModerator = (userId: number) => {
    if (confirm("确定要移除该版主吗？")) {
      removeModerator.mutate(userId);
    }
  };

  if (isLoading) return <div className="flex justify-center py-8"><span className="loading loading-spinner loading-lg"></span></div>;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center gap-4">
        <div className="flex gap-2">
          <input 
            type="number"
            placeholder="板块 ID"
            value={selectedBoardId}
            onChange={(e) => setSelectedBoardId(parseInt(e.target.value))}
            className="input input-bordered w-32"
          />
          <button className="btn btn-outline" onClick={() => refetch()}>刷新</button>
        </div>
        <button className="btn btn-primary" onClick={() => setIsAddDialogOpen(true)}>
          <UserPlus className="w-4 h-4 mr-1" />
          任命版主
        </button>
      </div>

      <div className="space-y-3">
        {moderators?.map((mod: any) => (
          <div key={mod.user_id} className="card bg-base-100 shadow-sm border border-base-200">
            <div className="card-body">
              <div className="flex justify-between items-center">
                <div className="flex items-center gap-3">
                  <div className="avatar placeholder">
                    <div className="bg-primary/10 rounded-full w-10">
                      <span className="text-primary">{mod.user?.username?.[0]?.toUpperCase()}</span>
                    </div>
                  </div>
                  <div>
                    <p className="font-medium">{mod.user?.username}</p>
                    <div className="flex gap-2 mt-1">
                      {mod.can_delete_post && <span className="badge badge-sm">删除帖子</span>}
                      {mod.can_pin_post && <span className="badge badge-sm">置顶</span>}
                      {mod.can_ban_user && <span className="badge badge-sm">禁言</span>}
                    </div>
                  </div>
                </div>
                <button 
                  className="btn btn-ghost btn-sm"
                  onClick={() => handleRemoveModerator(mod.user_id)}
                  disabled={removeModerator.isPending}
                >
                  <UserMinus className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ))}
        {moderators?.length === 0 && (
          <div className="text-center py-8 text-gray-500">暂无版主</div>
        )}
      </div>

      {/* 任命版主对话框 */}
      {isAddDialogOpen && (
        <dialog className="modal modal-open" onClick={() => setIsAddDialogOpen(false)}>
          <div className="modal-box" onClick={(e) => e.stopPropagation()}>
            <h3 className="font-bold text-lg">任命版主</h3>
            <div className="space-y-4 mt-4">
              <input 
                type="number"
                placeholder="用户 ID"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                className="input input-bordered w-full"
              />
              <div className="space-y-2">
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_delete_post}
                    onChange={(e) => setPermissions({ ...permissions, can_delete_post: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以删除帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_pin_post}
                    onChange={(e) => setPermissions({ ...permissions, can_pin_post: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以置顶帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_edit_any_post}
                    onChange={(e) => setPermissions({ ...permissions, can_edit_any_post: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以编辑任何帖子</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_manage_moderator}
                    onChange={(e) => setPermissions({ ...permissions, can_manage_moderator: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以管理其他版主</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={permissions.can_ban_user}
                    onChange={(e) => setPermissions({ ...permissions, can_ban_user: e.target.checked })}
                    className="checkbox checkbox-sm"
                  />
                  <span>可以禁言用户</span>
                </label>
              </div>
            </div>
            <div className="modal-action">
              <button className="btn" onClick={() => setIsAddDialogOpen(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleAddModerator} disabled={addModerator.isPending}>
                确认任命
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}

// ============ 举报管理模块 ============
function ReportManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [status, setStatus] = useState<string>("pending");
  
  const { data: reports, isLoading, refetch } = useAdminReports(selectedBoardId, {
    page: 1,
    page_size: 20,
    status,
  });

  if (isLoading) return <div className="flex justify-center py-8"><span className="loading loading-spinner loading-lg"></span></div>;

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
        <button className="btn btn-outline" onClick={() => refetch()}>刷新</button>
      </div>
      
      <div className="space-y-3">
        {reports?.list?.map((report: any) => (
          <div key={report.id} className="card bg-base-100 shadow-sm border border-base-200">
            <div className="card-body">
              <div className="space-y-2">
                <div className="flex justify-between">
                  <div className="badge badge-error gap-1">
                    <Flag className="w-3 h-3" />
                    举报
                  </div>
                  <span className="text-sm text-gray-500">{new Date(report.created_at).toLocaleString()}</span>
                </div>
                <p className="text-sm"><strong>举报人:</strong> {report.reporter?.username}</p>
                <p className="text-sm"><strong>被举报内容:</strong> {report.content_preview}</p>
                <p className="text-sm"><strong>原因:</strong> {report.reason}</p>
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

// ============ 帖子管理模块 ============
function PostManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [keyword, setKeyword] = useState("");
  const [debouncedKeyword, setDebouncedKeyword] = useState("");
  
  const { data: posts, isLoading, refetch } = useAdminBoardPosts(selectedBoardId, {
    page: 1,
    page_size: 20,
    keyword: debouncedKeyword,
  });
  
  const deletePost = useAdminDeletePost();
  const pinPost = useAdminPinPost();

  const handleSearch = () => {
    setDebouncedKeyword(keyword);
  };

  const handleDeletePost = (boardId: number, postId: number) => {
    if (confirm("确定要删除该帖子吗？此操作不可恢复。")) {
      deletePost.mutate({ boardId, postId });
    }
  };

  const handlePinPost = (boardId: number, postId: number, pinInBoard: boolean) => {
    pinPost.mutate({ boardId, postId, pinInBoard });
  };

  if (isLoading) return <div className="flex justify-center py-8"><span className="loading loading-spinner loading-lg"></span></div>;

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
        <div className="flex-1 flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input 
              placeholder="搜索帖子..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              className="input input-bordered w-full pl-9"
            />
          </div>
          <button className="btn btn-primary" onClick={handleSearch}>搜索</button>
          <button className="btn btn-outline" onClick={() => refetch()}>刷新</button>
        </div>
      </div>

      <div className="space-y-3">
        {posts?.list?.map((post: any) => (
          <div key={post.id} className="card bg-base-100 shadow-sm border border-base-200">
            <div className="card-body">
              <div className="space-y-2">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <h3 className="font-medium flex items-center gap-2">
                      {post.title}
                      {post.is_pinned && <Pin className="w-4 h-4 text-primary" />}
                    </h3>
                    <p className="text-sm text-gray-500 mt-1">
                      作者: {post.author?.username} | 回复: {post.reply_count} | 浏览: {post.view_count}
                    </p>
                  </div>
                  <div className="flex gap-1">
                    <button 
                      className="btn btn-ghost btn-sm"
                      onClick={() => handlePinPost(selectedBoardId, post.id, !post.is_pinned)}
                      disabled={pinPost.isPending}
                    >
                      <Pin className="w-4 h-4" />
                    </button>
                    <button 
                      className="btn btn-ghost btn-sm text-error"
                      onClick={() => handleDeletePost(selectedBoardId, post.id)}
                      disabled={deletePost.isPending}
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <p className="text-sm text-gray-600 line-clamp-2">{post.content}</p>
                <div className="flex gap-2">
                  {post.is_pinned && <span className="badge badge-primary">置顶</span>}
                  {post.status === "deleted" && <span className="badge badge-error">已删除</span>}
                </div>
              </div>
            </div>
          </div>
        ))}
        {posts?.list?.length === 0 && (
          <div className="text-center py-8 text-gray-500">暂无帖子</div>
        )}
      </div>
    </div>
  );
}

// ============ 禁言管理模块 ============
function BanManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [isBanDialogOpen, setIsBanDialogOpen] = useState(false);
  const [banUserId, setBanUserId] = useState("");
  const [banReason, setBanReason] = useState("");
  const [banDuration, setBanDuration] = useState("");
  
  const { data: bannedUsers, isLoading, refetch } = useAdminBannedUsers(selectedBoardId);
  const banUser = useAdminBanUser(selectedBoardId);
  const unbanUser = useAdminUnbanUser(selectedBoardId);

  const handleBanUser = () => {
    if (!banUserId || !banReason) {
      toast.error("请填写完整信息");
      return;
    }
    
    let expiresAt;
    if (banDuration) {
      const durationMap: Record<string, number> = {
        "1d": 86400000,
        "3d": 259200000,
        "7d": 604800000,
        "30d": 2592000000,
      };
      expiresAt = new Date(Date.now() + durationMap[banDuration]).toISOString();
    }
    
    banUser.mutate({
      user_id: parseInt(banUserId),
      reason: banReason,
      expires_at: expiresAt,
    }, {
      onSuccess: () => {
        setIsBanDialogOpen(false);
        setBanUserId("");
        setBanReason("");
        setBanDuration("");
      }
    });
  };

  const handleUnbanUser = (userId: number) => {
    if (confirm("确定要解除禁言吗？")) {
      unbanUser.mutate(userId);
    }
  };

  if (isLoading) return <div className="flex justify-center py-8"><span className="loading loading-spinner loading-lg"></span></div>;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center gap-4">
        <div className="flex gap-2">
          <input 
            type="number"
            placeholder="板块 ID"
            value={selectedBoardId}
            onChange={(e) => setSelectedBoardId(parseInt(e.target.value))}
            className="input input-bordered w-32"
          />
          <button className="btn btn-outline" onClick={() => refetch()}>刷新</button>
        </div>
        <button className="btn btn-primary" onClick={() => setIsBanDialogOpen(true)}>
          <Ban className="w-4 h-4 mr-1" />
          禁言用户
        </button>
      </div>

      <div className="space-y-3">
        {bannedUsers?.list?.map((ban: any) => (
          <div key={ban.user_id} className="card bg-base-100 shadow-sm border border-base-200">
            <div className="card-body">
              <div className="flex justify-between items-center">
                <div>
                  <p className="font-medium">{ban.user?.username}</p>
                  <p className="text-sm text-gray-500">禁言原因: {ban.reason}</p>
                  {ban.expires_at && (
                    <p className="text-sm text-gray-500">
                      到期时间: {new Date(ban.expires_at).toLocaleString()}
                    </p>
                  )}
                  {ban.created_at && (
                    <p className="text-sm text-gray-500">
                      禁言时间: {new Date(ban.created_at).toLocaleString()}
                    </p>
                  )}
                </div>
                <button 
                  className="btn btn-ghost btn-sm"
                  onClick={() => handleUnbanUser(ban.user_id)}
                  disabled={unbanUser.isPending}
                >
                  <Unlock className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ))}
        {bannedUsers?.list?.length === 0 && (
          <div className="text-center py-8 text-gray-500">暂无禁言用户</div>
        )}
      </div>

      {/* 禁言对话框 */}
      {isBanDialogOpen && (
        <dialog className="modal modal-open" onClick={() => setIsBanDialogOpen(false)}>
          <div className="modal-box" onClick={(e) => e.stopPropagation()}>
            <h3 className="font-bold text-lg">禁言用户</h3>
            <div className="space-y-4 mt-4">
              <input 
                type="number"
                placeholder="用户 ID"
                value={banUserId}
                onChange={(e) => setBanUserId(e.target.value)}
                className="input input-bordered w-full"
              />
              <textarea 
                placeholder="禁言原因"
                value={banReason}
                onChange={(e) => setBanReason(e.target.value)}
                className="textarea textarea-bordered w-full"
                rows={3}
              />
              <select 
                className="select select-bordered w-full"
                value={banDuration} 
                onChange={(e) => setBanDuration(e.target.value)}
              >
                <option value="">选择禁言时长</option>
                <option value="1d">1 天</option>
                <option value="3d">3 天</option>
                <option value="7d">7 天</option>
                <option value="30d">30 天</option>
                <option value="">永久禁言</option>
              </select>
            </div>
            <div className="modal-action">
              <button className="btn" onClick={() => setIsBanDialogOpen(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleBanUser} disabled={banUser.isPending}>
                确认禁言
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}

// ============ 主组件 ============
export function AdminTasks() {
  const [activeTab, setActiveTab] = useState("applications");

  const tabs = [
    { id: "applications", label: "版主申请", icon: FileText },
    { id: "moderators", label: "版主任命", icon: Shield },
    { id: "reports", label: "举报管理", icon: Flag },
    { id: "posts", label: "帖子管理", icon: FileText },
    { id: "bans", label: "禁言管理", icon: Ban },
  ];

  return (
    <div className="flex flex-col gap-6">

      <div role="tablist" className="tabs tabs-lifted">
        {tabs.map((tab) => (
          <a
            key={tab.id}
            role="tab"
            className={`tab ${activeTab === tab.id ? "tab-active" : ""}`}
            onClick={() => setActiveTab(tab.id)}
          >
            <tab.icon className="w-4 h-4 mr-2" />
            {tab.label}
          </a>
        ))}
      </div>

      <div className="mt-6">
        {activeTab === "applications" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <FileText className="w-5 h-5" />
                版主申请审批
              </h2>
              <p className="text-gray-500">审核用户提交的版主申请，通过后用户将成为对应板块的版主</p>
              <div className="mt-4">
                <ApplicationManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "moderators" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <Shield className="w-5 h-5" />
                版主任命与管理
              </h2>
              <p className="text-gray-500">任命新版主、管理现有版主及其权限</p>
              <div className="mt-4">
                <ModeratorManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "reports" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <Flag className="w-5 h-5" />
                举报处理
              </h2>
              <p className="text-gray-500">处理用户举报，维护社区秩序</p>
              <div className="mt-4">
                <ReportManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "posts" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <FileText className="w-5 h-5" />
                帖子管理
              </h2>
              <p className="text-gray-500">搜索、置顶、删除帖子，管理板块内容</p>
              <div className="mt-4">
                <PostManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "bans" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <Ban className="w-5 h-5" />
                禁言管理
              </h2>
              <p className="text-gray-500">禁言违规用户，查看禁言列表，解除禁言</p>
              <div className="mt-4">
                <BanManagement />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}