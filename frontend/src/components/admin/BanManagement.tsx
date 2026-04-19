import { useAdminBannedUsers, useAdminBanUser, useAdminUnbanUser } from "@/hooks/admin/useAdminModerator";
import { Ban, Unlock } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";

export function BanManagement() {
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