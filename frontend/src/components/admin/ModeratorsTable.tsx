import { useState, useEffect } from "react";
import {
  Shield,
  UserPlus,
  UserMinus,
  Edit,
  Search,
  ShieldCheck,
  ShieldAlert,
  ShieldX,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import {
  moderatorApi,
  Moderator,
  AddModeratorRequest,
} from "@/lib/api/modules/moderator";
import { useBoard } from "@/hooks/useBoard";

// 扩展版主类型，包含板块信息
interface ExtendedModerator extends Moderator {
  board_name?: string;
  board_slug?: string;
}

// ==================== 管理员管理版主组件 ====================
export function ModeratorsTable() {
  const { boards, loading: boardsLoading } = useBoard({
    autoLoad: true,
    pageSize: 500,
  });
  const [moderators, setModerators] = useState<ExtendedModerator[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedModerator, setSelectedModerator] =
    useState<ExtendedModerator | null>(null);
  const [searchKeyword, setSearchKeyword] = useState("");
  const [selectedBoardFilter, setSelectedBoardFilter] = useState<
    number | "all"
  >("all");

  // 分页状态
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);

  // 表单状态
  const [userId, setUserId] = useState("");
  const [selectedBoardId, setSelectedBoardId] = useState<number>(
    boards[0]?.id || 0,
  );
  const [permissions, setPermissions] = useState({
    can_delete_post: false,
    can_pin_post: false,
    can_edit_any_post: false,
    can_manage_moderator: false,
    can_ban_user: false,
  });

  const [submitting, setSubmitting] = useState(false);

  // 加载所有版主（从所有板块）
  const loadAllModerators = async () => {
    setLoading(true);
    try {
      // 获取所有板块的版主
      const allModerators: ExtendedModerator[] = [];

      for (const board of boards) {
        try {
          const res = await moderatorApi.getModerators(board.id);
          const boardModerators = (res.data.data || []).map((m: Moderator) => ({
            ...m,
            board_name: board.name,
            board_slug: board.slug,
          }));
          allModerators.push(...boardModerators);
        } catch (error) {
          console.error(`加载板块 ${board.id} 版主失败:`, error);
        }
      }

      setModerators(allModerators);
    } catch (error) {
      console.error("加载版主列表失败:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (boards.length > 0) {
      loadAllModerators();
      setSelectedBoardId(boards[0]?.id || 0);
    }
  }, [boards]);

  // 筛选版主
  const filteredModerators = moderators.filter((moderator) => {
    // 板块筛选
    if (
      selectedBoardFilter !== "all" &&
      moderator.board_id !== selectedBoardFilter
    ) {
      return false;
    }
    // 关键词筛选
    if (
      searchKeyword &&
      !moderator.user?.username
        ?.toLowerCase()
        .includes(searchKeyword.toLowerCase()) &&
      !moderator.user_id?.toString().includes(searchKeyword) &&
      !moderator.board_name?.toLowerCase().includes(searchKeyword.toLowerCase())
    ) {
      return false;
    }
    return true;
  });

  // 分页数据
  const paginatedModerators = filteredModerators.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize,
  );
  const totalPages = Math.ceil(filteredModerators.length / pageSize);

  // 统计信息
  const stats = {
    totalModerators: moderators.length,
    totalBoards: new Set(moderators.map((m) => m.board_id)).size,
    permissionsStats: {
      canDeletePost: moderators.filter((m) => m.permissions?.can_delete_post)
        .length,
      canPinPost: moderators.filter((m) => m.permissions?.can_pin_post).length,
      canBanUser: moderators.filter((m) => m.permissions?.can_ban_user).length,
    },
  };

  // 添加版主
  const handleAddModerator = async () => {
    if (!userId) {
      alert("请输入用户ID");
      return;
    }
    if (!selectedBoardId) {
      alert("请选择板块");
      return;
    }

    setSubmitting(true);
    try {
      const data: AddModeratorRequest = {
        user_id: Number(userId),
        ...permissions,
      };
      await moderatorApi.addModerator(selectedBoardId, data);
      alert("添加版主成功");
      setShowAddModal(false);
      resetForm();
      loadAllModerators();
    } catch (error: any) {
      alert(error.response?.data?.message || "操作失败");
    } finally {
      setSubmitting(false);
    }
  };

  // 移除版主
  const handleRemoveModerator = async (moderator: ExtendedModerator) => {
    if (
      !confirm(
        `确定要移除 ${moderator.user?.username || moderator.user_id} 的版主权限吗？`,
      )
    )
      return;

    try {
      await moderatorApi.removeModerator(moderator.board_id, moderator.user_id);
      alert("移除版主成功");
      loadAllModerators();
    } catch (error: any) {
      alert(error.response?.data?.message || "操作失败");
    }
  };

  // 更新权限
  const handleUpdatePermissions = async () => {
    if (!selectedModerator) return;

    setSubmitting(true);
    try {
      await moderatorApi.updateModeratorPermissions(
        selectedModerator.board_id,
        selectedModerator.user_id,
        permissions,
      );
      alert("权限更新成功");
      setShowEditModal(false);
      resetForm();
      loadAllModerators();
    } catch (error: any) {
      alert(error.response?.data?.message || "操作失败");
    } finally {
      setSubmitting(false);
    }
  };

  const resetForm = () => {
    setUserId("");
    setPermissions({
      can_delete_post: false,
      can_pin_post: false,
      can_edit_any_post: false,
      can_manage_moderator: false,
      can_ban_user: false,
    });
    setSelectedModerator(null);
  };

  const openEditModal = (moderator: ExtendedModerator) => {
    setSelectedModerator(moderator);
    setPermissions(
      moderator.permissions || {
        can_delete_post: false,
        can_pin_post: false,
        can_edit_any_post: false,
        can_manage_moderator: false,
        can_ban_user: false,
      },
    );
    setShowEditModal(true);
  };

  // 权限徽章组件
  const PermissionBadge = ({
    label,
    hasPermission,
  }: {
    label: string;
    hasPermission: boolean;
  }) => (
    <span
      className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium ${
        hasPermission
          ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
          : "bg-gray-100 text-gray-500 dark:bg-gray-700/30 dark:text-gray-400"
      }`}
    >
      {hasPermission ? (
        <ShieldCheck className="w-3 h-3" />
      ) : (
        <ShieldX className="w-3 h-3" />
      )}
      {label}
    </span>
  );

  const isLoading = loading || boardsLoading;

  return (
    <div className="space-y-6">
      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-primary">
            <Shield className="w-6 h-6" />
          </div>
          <div className="stat-title">版主总数</div>
          <div className="stat-value text-primary">{stats.totalModerators}</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-secondary">
            <ShieldCheck className="w-6 h-6" />
          </div>
          <div className="stat-title">可删除帖子</div>
          <div className="stat-value text-secondary">
            {stats.permissionsStats.canDeletePost}
          </div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-accent">
            <ShieldAlert className="w-6 h-6" />
          </div>
          <div className="stat-title">可禁言用户</div>
          <div className="stat-value text-accent">
            {stats.permissionsStats.canBanUser}
          </div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-info">
            <Shield className="w-6 h-6" />
          </div>
          <div className="stat-title">管理板块数</div>
          <div className="stat-value text-info">{stats.totalBoards}</div>
        </div>
      </div>

      {/* 操作栏 */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body">
          <div className="flex flex-wrap justify-between items-center gap-4 mb-4">
            <h3 className="font-semibold">版主管理</h3>
            <div className="flex flex-wrap gap-2">
              {/* 板块筛选 */}
              <select
                className="select select-bordered select-sm"
                value={selectedBoardFilter}
                onChange={(e) => {
                  setSelectedBoardFilter(
                    e.target.value === "all" ? "all" : Number(e.target.value),
                  );
                  setCurrentPage(1);
                }}
              >
                <option value="all">所有板块</option>
                {boards.map((board) => (
                  <option key={board.id} value={board.id}>
                    {board.name}
                  </option>
                ))}
              </select>

              {/* 搜索框 */}
              <div className="form-control">
                <div className="flex gap-2">
                  <input
                    type="text"
                    className="input input-bordered input-sm"
                    placeholder="搜索用户/板块"
                    value={searchKeyword}
                    onChange={(e) => {
                      setSearchKeyword(e.target.value);
                      setCurrentPage(1);
                    }}
                  />
                  <button className="btn btn-ghost btn-sm">
                    <Search className="w-4 h-4" />
                  </button>
                </div>
              </div>

              {/* 添加按钮 */}
              <button
                className="btn btn-primary btn-sm"
                onClick={() => setShowAddModal(true)}
              >
                <UserPlus className="w-4 h-4" /> 添加版主
              </button>
            </div>
          </div>

          {/* 版主列表表格 */}
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>用户</th>
                  <th>板块</th>
                  <th>权限</th>
                  <th>添加时间</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {isLoading ? (
                  <tr>
                    <td colSpan={5} className="text-center py-8">
                      <span className="loading loading-spinner loading-md" />
                    </td>
                  </tr>
                ) : paginatedModerators.length > 0 ? (
                  paginatedModerators.map((moderator) => (
                    <tr key={`${moderator.board_id}-${moderator.user_id}`}>
                      <td>
                        <div className="flex items-center gap-2">
                          {moderator.user?.avatar && (
                            <img
                              src={moderator.user.avatar}
                              alt=""
                              className="w-8 h-8 rounded-full"
                            />
                          )}
                          <div>
                            <div className="font-medium">
                              {moderator.user?.username ||
                                `用户 #${moderator.user_id}`}
                            </div>
                            <div className="text-xs text-base-content/50">
                              ID: {moderator.user_id}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td>
                        <span className="badge badge-ghost">
                          {moderator.board_name ||
                            `板块 #${moderator.board_id}`}
                        </span>
                      </td>
                      <td>
                        <div className="flex flex-wrap gap-1">
                          <PermissionBadge
                            label="删帖"
                            hasPermission={
                              moderator.permissions?.can_delete_post
                            }
                          />
                          <PermissionBadge
                            label="置顶"
                            hasPermission={moderator.permissions?.can_pin_post}
                          />
                          <PermissionBadge
                            label="编辑"
                            hasPermission={
                              moderator.permissions?.can_edit_any_post
                            }
                          />
                          <PermissionBadge
                            label="管理版主"
                            hasPermission={
                              moderator.permissions?.can_manage_moderator
                            }
                          />
                          <PermissionBadge
                            label="禁言"
                            hasPermission={moderator.permissions?.can_ban_user}
                          />
                        </div>
                      </td>
                      <td className="text-sm text-base-content/60">
                        {new Date(moderator.created_at).toLocaleDateString()}
                      </td>
                      <td>
                        <div className="flex gap-2">
                          <button
                            className="btn btn-ghost btn-xs"
                            onClick={() => openEditModal(moderator)}
                          >
                            <Edit className="w-4 h-4" />
                          </button>
                          <button
                            className="btn btn-ghost btn-xs text-error"
                            onClick={() => handleRemoveModerator(moderator)}
                          >
                            <UserMinus className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td
                      colSpan={5}
                      className="text-center py-8 text-base-content/50"
                    >
                      暂无版主数据
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* 分页 */}
          {totalPages > 1 && (
            <div className="flex justify-center items-center gap-2 mt-4">
              <button
                className="btn btn-ghost btn-sm"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(currentPage - 1)}
              >
                <ChevronLeft className="w-4 h-4" />
              </button>
              <span className="text-sm">
                第 {currentPage} / {totalPages} 页
              </span>
              <button
                className="btn btn-ghost btn-sm"
                disabled={currentPage === totalPages}
                onClick={() => setCurrentPage(currentPage + 1)}
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          )}
        </div>
      </div>

      {/* 添加版主模态框 */}
      {showAddModal && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">添加版主</h3>

            <div className="form-control mt-4">
              <label className="label">
                <span className="label-text">选择板块</span>
              </label>
              <select
                className="select select-bordered"
                value={selectedBoardId}
                onChange={(e) => setSelectedBoardId(Number(e.target.value))}
              >
                {boards.map((board) => (
                  <option key={board.id} value={board.id}>
                    {board.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-control mt-4">
              <label className="label">
                <span className="label-text">用户ID</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder="输入用户ID"
              />
            </div>

            <div className="form-control mt-4">
              <label className="label">
                <span className="label-text">权限设置</span>
              </label>
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
                  删除帖子
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
                  置顶帖子
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
                  编辑任意帖子
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
                  管理版主
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
                  禁言用户
                </label>
              </div>
            </div>

            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => {
                  setShowAddModal(false);
                  resetForm();
                }}
              >
                取消
              </button>
              <button
                className="btn btn-primary"
                onClick={handleAddModerator}
                disabled={submitting}
              >
                {submitting ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  "确认"
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 编辑权限模态框 */}
      {showEditModal && selectedModerator && (
        <div className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg">
              编辑权限 -{" "}
              {selectedModerator.user?.username ||
                `用户 #${selectedModerator.user_id}`}
            </h3>
            <p className="text-sm text-base-content/60 mt-1">
              板块：
              {selectedModerator.board_name ||
                `板块 #${selectedModerator.board_id}`}
            </p>

            <div className="form-control mt-4">
              <label className="label">
                <span className="label-text">权限设置</span>
              </label>
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
                  删除帖子
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
                  置顶帖子
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
                  编辑任意帖子
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
                  管理版主
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
                  禁言用户
                </label>
              </div>
            </div>

            <div className="modal-action">
              <button
                className="btn btn-ghost"
                onClick={() => {
                  setShowEditModal(false);
                  resetForm();
                }}
              >
                取消
              </button>
              <button
                className="btn btn-primary"
                onClick={handleUpdatePermissions}
                disabled={submitting}
              >
                {submitting ? (
                  <span className="loading loading-spinner loading-sm" />
                ) : (
                  "保存"
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
