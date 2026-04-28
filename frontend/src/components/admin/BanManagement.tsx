import {
  useAdminBannedUsers,
  useAdminBanUser,
  useAdminUnbanUser,
} from "@/hooks/admin/useAdminModerator";
import { Ban, Unlock, Shield } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";

// 类型定义
interface User {
  id: number;
  username: string;
  avatar?: string;
  email?: string;
}

interface BanRecord {
  id: number;
  user_id: number;
  board_id: number;
  user?: User;
  reason: string;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
}

interface BannedUsersResponse {
  list: BanRecord[];
  total: number;
  page: number;
  page_size: number;
}

// 禁言时长映射
const DURATION_MAP: Record<string, number> = {
  "1d": 86400000,      // 1 天
  "3d": 259200000,     // 3 天
  "7d": 604800000,     // 7 天
  "30d": 2592000000,   // 30 天
};

export function BanManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [isBanDialogOpen, setIsBanDialogOpen] = useState<boolean>(false);
  const [banUserId, setBanUserId] = useState<string>("");
  const [banReason, setBanReason] = useState<string>("");
  const [banDuration, setBanDuration] = useState<string>("");

  const {
    data: bannedUsers,
    isLoading,
    refetch,
  } = useAdminBannedUsers(selectedBoardId);
  const banUser = useAdminBanUser(selectedBoardId);
  const unbanUser = useAdminUnbanUser(selectedBoardId);

  // 获取禁言列表
  const bannedList = (bannedUsers as unknown as BannedUsersResponse)?.list || [];

  const handleBanUser = () => {
    if (!banUserId || !banReason) {
      toast.error("请填写完整信息");
      return;
    }

    let expiresAt: string | undefined;
    if (banDuration && DURATION_MAP[banDuration]) {
      expiresAt = new Date(Date.now() + DURATION_MAP[banDuration]).toISOString();
    }

    banUser.mutate(
      {
        user_id: parseInt(banUserId),
        reason: banReason,
        expires_at: expiresAt,
      },
      {
        onSuccess: () => {
          toast.success("禁言成功");
          setIsBanDialogOpen(false);
          setBanUserId("");
          setBanReason("");
          setBanDuration("");
          refetch();
        },
        onError: () => {
          toast.error("操作失败");
        },
      },
    );
  };

  const handleUnbanUser = (userId: number, username?: string) => {
    if (confirm(`确定要解除用户「${username || userId}」的禁言吗？`)) {
      unbanUser.mutate(userId, {
        onSuccess: () => {
          toast.success("解除禁言成功");
          refetch();
        },
        onError: () => {
          toast.error("操作失败");
        },
      });
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
      <div className="flex justify-between items-center gap-4 flex-wrap">
        <div className="flex gap-2">
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
              <span className="label-text">&nbsp;</span>
            </label>
            <button className="btn btn-outline" onClick={() => refetch()}>
              刷新
            </button>
          </div>
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text">&nbsp;</span>
          </label>
          <button
            className="btn btn-primary"
            onClick={() => setIsBanDialogOpen(true)}
          >
            <Ban className="w-4 h-4 mr-1" />
            禁言用户
          </button>
        </div>
      </div>

      <div className="space-y-3">
        {bannedList.map((ban: BanRecord) => (
          <div
            key={ban.id}
            className="card bg-base-100 shadow-sm border border-base-200 hover:shadow-md transition-shadow"
          >
            <div className="card-body p-4">
              <div className="flex justify-between items-center flex-wrap gap-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <div className="avatar placeholder">
                      <div className="bg-error/10 rounded-full w-8 h-8 flex items-center justify-center">
                        <span className="text-error font-medium text-xs">
                          {ban.user?.username?.[0]?.toUpperCase() || "U"}
                        </span>
                      </div>
                    </div>
                    <p className="font-medium">
                      {ban.user?.username || `用户${ban.user_id}`}
                    </p>
                    <span className="badge badge-error badge-sm">已禁言</span>
                  </div>
                  <div className="mt-2 space-y-1">
                    <p className="text-sm text-base-content/70">
                      <strong className="font-medium">禁言原因:</strong> {ban.reason}
                    </p>
                    {ban.expires_at && (
                      <p className="text-sm text-base-content/50">
                        <strong>到期时间:</strong> {new Date(ban.expires_at).toLocaleString()}
                      </p>
                    )}
                    {ban.expires_at === null && (
                      <p className="text-sm text-base-content/50">
                        <strong>类型:</strong> 永久禁言
                      </p>
                    )}
                    <p className="text-sm text-base-content/50">
                      <strong>禁言时间:</strong> {new Date(ban.created_at).toLocaleString()}
                    </p>
                  </div>
                </div>
                <button
                  className="btn btn-ghost btn-sm text-success hover:text-success"
                  onClick={() => handleUnbanUser(ban.user_id, ban.user?.username)}
                  disabled={unbanUser.isPending}
                  title="解除禁言"
                >
                  <Unlock className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ))}
        
        {bannedList.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            <Shield className="w-12 h-12 mx-auto mb-3 opacity-30" />
            <p>暂无禁言用户</p>
          </div>
        )}
      </div>

      {/* 禁言对话框 */}
      {isBanDialogOpen && (
        <dialog
          className="modal modal-open"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setIsBanDialogOpen(false);
            }
          }}
        >
          <div className="modal-box" onClick={(e) => e.stopPropagation()}>
            <h3 className="font-bold text-lg">禁言用户</h3>
            <div className="space-y-4 mt-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">用户 ID</span>
                </label>
                <input
                  type="number"
                  placeholder="用户 ID"
                  value={banUserId}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
                    setBanUserId(e.target.value)
                  }
                  className="input input-bordered w-full"
                  autoFocus
                />
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">禁言原因</span>
                </label>
                <textarea
                  placeholder="请输入禁言原因..."
                  value={banReason}
                  onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => 
                    setBanReason(e.target.value)
                  }
                  className="textarea textarea-bordered w-full"
                  rows={3}
                />
              </div>
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">禁言时长</span>
                </label>
                <select
                  className="select select-bordered w-full"
                  value={banDuration}
                  onChange={(e: React.ChangeEvent<HTMLSelectElement>) => 
                    setBanDuration(e.target.value)
                  }
                >
                  <option value="">选择禁言时长</option>
                  <option value="1d">1 天</option>
                  <option value="3d">3 天</option>
                  <option value="7d">7 天</option>
                  <option value="30d">30 天</option>
                  <option value="">永久禁言</option>
                </select>
                <label className="label">
                  <span className="label-text-alt text-base-content/50">
                    永久禁言将无限期限制用户发言
                  </span>
                </label>
              </div>
            </div>
            <div className="modal-action">
              <button 
                className="btn" 
                onClick={() => setIsBanDialogOpen(false)}
              >
                取消
              </button>
              <button
                className="btn btn-primary"
                onClick={handleBanUser}
                disabled={banUser.isPending}
              >
                {banUser.isPending ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : null}
                确认禁言
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}