import {
  useAdminModeratorList,
  useAddModerator,
  useRemoveModerator,
} from "@/hooks/admin/useAdminModerator";
import { UserPlus, UserMinus } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";

export function ModeratorManagement() {
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

  const {
    data: moderators,
    isLoading,
    refetch,
  } = useAdminModeratorList(selectedBoardId);
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

  if (isLoading)
    return (
      <div className="flex justify-center py-8">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );

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
          <button className="btn btn-outline" onClick={() => refetch()}>
            刷新
          </button>
        </div>
        <button
          className="btn btn-primary"
          onClick={() => setIsAddDialogOpen(true)}
        >
          <UserPlus className="w-4 h-4 mr-1" />
          任命版主
        </button>
      </div>

      <div className="space-y-3">
        {moderators?.map((mod: any) => (
          <div
            key={mod.user_id}
            className="card bg-base-100 shadow-sm border border-base-200"
          >
            <div className="card-body">
              <div className="flex justify-between items-center">
                <div className="flex items-center gap-3">
                  <div className="avatar placeholder">
                    <div className="bg-primary/10 rounded-full w-10">
                      <span className="text-primary">
                        {mod.user?.username?.[0]?.toUpperCase()}
                      </span>
                    </div>
                  </div>
                  <div>
                    <p className="font-medium">{mod.user?.username}</p>
                    <div className="flex gap-2 mt-1">
                      {mod.can_delete_post && (
                        <span className="badge badge-sm">删除帖子</span>
                      )}
                      {mod.can_pin_post && (
                        <span className="badge badge-sm">置顶</span>
                      )}
                      {mod.can_ban_user && (
                        <span className="badge badge-sm">禁言</span>
                      )}
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
        <dialog
          className="modal modal-open"
          onClick={() => setIsAddDialogOpen(false)}
        >
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
              <button className="btn" onClick={() => setIsAddDialogOpen(false)}>
                取消
              </button>
              <button
                className="btn btn-primary"
                onClick={handleAddModerator}
                disabled={addModerator.isPending}
              >
                确认任命
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}
