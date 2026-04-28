import {
  useAdminModeratorList,
  useAddModerator,
  useRemoveModerator,
} from "@/hooks/admin/useAdminModerator";
import { UserPlus, UserMinus } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";

// 类型定义
interface User {
  id: number;
  username: string;
  avatar?: string;
  email?: string;
}

interface ModeratorPermissions {
  can_delete_post: boolean;
  can_pin_post: boolean;
  can_edit_any_post: boolean;
  can_manage_moderator: boolean;
  can_ban_user: boolean;
}

interface Moderator {
  id: number;
  user_id: number;
  board_id: number;
  user?: User;
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

export function ModeratorManagement() {
  const [selectedBoardId, setSelectedBoardId] = useState<number>(1);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState<boolean>(false);
  const [userId, setUserId] = useState<string>("");
  const [permissions, setPermissions] = useState<ModeratorPermissions>(DEFAULT_PERMISSIONS);

  const {
    data: moderators,
    isLoading,
    refetch,
  } = useAdminModeratorList(selectedBoardId);
  const addModerator = useAddModerator(selectedBoardId);
  const removeModerator = useRemoveModerator(selectedBoardId);

  // 获取版主列表
  const moderatorList = (moderators as unknown as Moderator[]) || [];

  const handleAddModerator = async () => {
    if (!userId) {
      toast.error("请输入用户 ID");
      return;
    }
    
    try {
      await addModerator.mutateAsync({
        user_id: parseInt(userId),
        ...permissions,
      });
      toast.success("版主任命成功");
      setIsAddDialogOpen(false);
      setUserId("");
      setPermissions(DEFAULT_PERMISSIONS);
      refetch();
    } catch {
      toast.error("操作失败");
    }
  };

  const handleRemoveModerator = async (userId: number, username?: string) => {
    if (confirm(`确定要移除版主「${username || userId}」吗？`)) {
      try {
        await removeModerator.mutateAsync(userId);
        toast.success("版主已移除");
        refetch();
      } catch {
        toast.error("操作失败");
      }
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
            onClick={() => setIsAddDialogOpen(true)}
          >
            <UserPlus className="w-4 h-4 mr-1" />
            任命版主
          </button>
        </div>
      </div>

      <div className="space-y-3">
        {moderatorList.map((mod: Moderator) => (
          <div
            key={mod.user_id}
            className="card bg-base-100 shadow-sm border border-base-200 hover:shadow-md transition-shadow"
          >
            <div className="card-body p-4">
              <div className="flex justify-between items-center flex-wrap gap-2">
                <div className="flex items-center gap-3">
                  <div className="avatar placeholder">
                    <div className="bg-primary/10 rounded-full w-10 h-10 flex items-center justify-center">
                      <span className="text-primary font-medium">
                        {mod.user?.username?.[0]?.toUpperCase() || "U"}
                      </span>
                    </div>
                  </div>
                  <div>
                    <p className="font-medium">{mod.user?.username || `用户${mod.user_id}`}</p>
                    <div className="flex gap-2 mt-1 flex-wrap">
                      {mod.can_delete_post && (
                        <span className="badge badge-sm badge-ghost">删除帖子</span>
                      )}
                      {mod.can_pin_post && (
                        <span className="badge badge-sm badge-ghost">置顶</span>
                      )}
                      {mod.can_ban_user && (
                        <span className="badge badge-sm badge-ghost">禁言</span>
                      )}
                      {mod.can_edit_any_post && (
                        <span className="badge badge-sm badge-ghost">编辑</span>
                      )}
                      {mod.can_manage_moderator && (
                        <span className="badge badge-sm badge-warning">管理版主</span>
                      )}
                    </div>
                  </div>
                </div>
                <button
                  className="btn btn-ghost btn-sm text-error hover:text-error"
                  onClick={() => handleRemoveModerator(mod.user_id, mod.user?.username)}
                  disabled={removeModerator.isPending}
                  title="移除版主"
                >
                  <UserMinus className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ))}
        
        {moderatorList.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            <Shield className="w-12 h-12 mx-auto mb-3 opacity-30" />
            <p>暂无版主</p>
          </div>
        )}
      </div>

      {/* 任命版主对话框 */}
      {isAddDialogOpen && (
        <dialog
          className="modal modal-open"
          onClick={(e) => {
            if (e.target === e.currentTarget) {
              setIsAddDialogOpen(false);
            }
          }}
        >
          <div className="modal-box" onClick={(e) => e.stopPropagation()}>
            <h3 className="font-bold text-lg">任命版主</h3>
            <div className="space-y-4 mt-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">用户 ID</span>
                </label>
                <input
                  type="number"
                  placeholder="用户 ID"
                  value={userId}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => 
                    setUserId(e.target.value)
                  }
                  className="input input-bordered w-full"
                  autoFocus
                />
              </div>
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
                  <span>可以删除帖子</span>
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
                  <span>可以置顶帖子</span>
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
                  <span>可以编辑任何帖子</span>
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
                  <span>可以管理其他版主</span>
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
                  <span>可以禁言用户</span>
                </label>
              </div>
            </div>
            <div className="modal-action">
              <button 
                className="btn" 
                onClick={() => setIsAddDialogOpen(false)}
              >
                取消
              </button>
              <button
                className="btn btn-primary"
                onClick={handleAddModerator}
                disabled={addModerator.isPending}
              >
                {addModerator.isPending ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : null}
                确认任命
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}

// 需要导入 Shield 图标
import { Shield } from "lucide-react";