// app/[locale]/dashboard/moderator/page.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  useModeratorPermissions,
  useBanUser,
  useUnbanUser,
  useDeletePost,
  usePinPost,
  useMyModeratorBoards,
} from "@/hooks/moderators/useModerator";
import { useModeratorAuth } from "@/hooks/moderators/useModeratorAuth";
import { useModeratorBannedUsers } from "@/hooks/moderators/useModeratorBannedUsers";
import { useModeratorPosts } from "@/hooks/moderators/useModeratorPosts";
import { useModeratorReports } from "@/hooks/moderators/useModeratorReports";
import { Pagination } from "@/components/admin/Pagination";
import { BannedUsersTable } from "@/components/moderator/BannedUsersTable";
import { BanUserModal } from "@/components/moderator/BanUserModal";
import { ReportedContentTable } from "@/components/moderator/ReportedContentTable";
import { ModeratorDashboard } from "@/components/moderator/ModeratorDashboard";
import { ModeratorSidebar } from "@/components/moderator/ModeratorSidebar";
import { PendingPostsTable } from "@/components/moderator/PendingPostsTable";
import SearchBar from "@/components/nav/SearchBar";
import { ModeratorBoard } from "@/lib/api/modules/moderator";

export default function ModeratorPage() {
  const t = useTranslations("Moderator");
  const { isCheckingAuth, isModerator, user } = useModeratorAuth();

  const [activeMenu, setActiveMenu] = useState("dashboard");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [page, setPage] = useState(1);
  const [selectedBoardId, setSelectedBoardId] = useState<number | null>(null);

  // 获取当前用户管理的板块
  const { data: boardsData, isLoading: boardsLoading } = useMyModeratorBoards();
  const boards: ModeratorBoard[] = boardsData || []; // 明确类型为数组

  // 选择第一个板块作为默认板块
  const currentBoardId =
    selectedBoardId || (boards.length > 0 ? boards[0]?.id : null);
  const currentBoard = boards.find(
    (b: ModeratorBoard) => b.id === currentBoardId,
  );

  // 当前板块的版主权限
  const permissions = useModeratorPermissions(currentBoardId || 0);

  // 数据获取（只在对应菜单激活时获取）
  const postsData = useModeratorPosts(
    currentBoardId || 0,
    page,
    keyword,
    activeMenu === "posts" && !!currentBoardId,
  );

  const reportsData = useModeratorReports(
    currentBoardId || 0,
    page,
    activeMenu === "reports" && !!currentBoardId,
  );

  const bannedUsersData = useModeratorBannedUsers(
    currentBoardId || 0,
    page,
    activeMenu === "bans" && !!currentBoardId,
  );

  // Mutations
  const { mutate: banUser, isPending: isBanning } = useBanUser(
    currentBoardId || 0,
  );
  const { mutate: unbanUser, isPending: isUnbanning } = useUnbanUser(
    currentBoardId || 0,
  );
  const { mutate: deletePost, isPending: isDeleting } = useDeletePost(
    currentBoardId || 0,
  );
  const { mutate: pinPost, isPending: isPinning } = usePinPost(
    currentBoardId || 0,
  );

  // 处理禁言
  const handleBanUser = (data: {
    userId: number;
    reason: string;
    expiresAt?: string;
  }) => {
    banUser({
      user_id: data.userId,
      reason: data.reason,
      expires_at: data.expiresAt,
    });
  };

  // 处理解禁
  const handleUnbanUser = (userId: number) => {
    if (confirm(t("confirm_unban"))) {
      unbanUser(userId);
    }
  };

  // 处理删除帖子
  const handleDeletePost = (postId: number) => {
    if (confirm(t("confirm_delete_post"))) {
      deletePost(postId);
    }
  };

  // 处理置顶帖子
  const handlePinPost = (postId: number, pin: boolean) => {
    pinPost({ postId, pinInBoard: pin });
  };

  // 加载状态
  if (isCheckingAuth || boardsLoading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  if (!isModerator) {
    return null;
  }

  // 如果没有管理的板块
  if (boards.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <div className="card bg-base-100 border border-base-300 p-8 text-center">
          <h2 className="text-xl font-bold mb-2">{t("no_managed_boards")}</h2>
          <p className="text-base-content/60 mb-4">
            {t("no_managed_boards_desc")}
          </p>
          <a href="/boards/apply" className="btn btn-primary">
            {t("apply_for_moderator")}
          </a>
        </div>
      </div>
    );
  }

  // 渲染右侧内容
  const renderContent = () => {
    if (!currentBoardId) {
      return (
        <div className="card bg-base-100 border border-base-300">
          <div className="card-body items-center text-center py-12">
            <p className="text-base-content/50">{t("no_board_selected")}</p>
          </div>
        </div>
      );
    }

    switch (activeMenu) {
      case "dashboard":
        return (
          <ModeratorDashboard
            board={currentBoard}
            permissions={permissions}
            stats={{
              postCount: postsData.total || 0,
              reportCount: reportsData.total || 0,
              bannedCount: bannedUsersData.total || 0,
            }}
            t={t}
          />
        );

      case "posts":
        return (
          <div className="space-y-4">
            <SearchBar
              keyword={keyword}
              onKeywordChange={(value: string) => {
                setKeyword(value);
                setPage(1);
              }}
              placeholder={t("search_posts")}
            />
            <PendingPostsTable
              posts={postsData.posts || []}
              onDelete={
                permissions.canDeletePost ? handleDeletePost : undefined
              }
              onPin={permissions.canPinPost ? handlePinPost : undefined}
              isDeleting={isDeleting}
              isPinning={isPinning}
              permissions={permissions}
              t={t}
            />
            <Pagination
              currentPage={page}
              total={postsData.total || 0}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "reports":
        return (
          <div className="space-y-4">
            <ReportedContentTable
              reports={reportsData.reports || []}
              onDeletePost={
                permissions.canDeletePost ? handleDeletePost : undefined
              }
              onBanUser={permissions.canBanUser ? handleBanUser : undefined}
              isDeleting={isDeleting}
              isBanning={isBanning}
              permissions={permissions}
              t={t}
            />
            <Pagination
              currentPage={page}
              total={reportsData.total || 0}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "bans":
        return (
          <div className="space-y-4">
            {permissions.canBanUser && (
              <div className="flex justify-end">
                <BanUserModal
                  boardId={currentBoardId}
                  onBan={handleBanUser}
                  isBanning={isBanning}
                  t={t}
                />
              </div>
            )}
            <BannedUsersTable
              users={bannedUsersData.users || []}
              onUnban={permissions.canBanUser ? handleUnbanUser : undefined}
              isUnbanning={isUnbanning}
              t={t}
            />
            <Pagination
              currentPage={page}
              total={bannedUsersData.total || 0}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="flex h-screen bg-base-100">
      {/* 左侧菜单 */}
      <ModeratorSidebar
        activeMenu={activeMenu}
        onMenuChange={setActiveMenu}
        collapsed={sidebarCollapsed}
        onCollapsedChange={setSidebarCollapsed}
        boards={boards} // 传递板块列表
        currentBoardId={currentBoardId} // 当前选中的板块ID
        onBoardChange={setSelectedBoardId} // 切换板块的回调
        permissions={permissions}
        t={t}
      />

      {/* 右侧内容区域 */}
      <div className="flex-1 overflow-y-auto">
        <div className="p-6">
          {/* 页面标题 */}
          <div className="mb-6">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-2xl font-bold">
                  {currentBoard?.name || t("moderator_panel")} - {t(activeMenu)}
                </h1>
                <p className="text-sm text-base-content/60 mt-1">
                  {t(`${activeMenu}_description`)}
                </p>
              </div>
            </div>

            {/* 权限标签 */}
            {permissions.isModerator && (
              <div className="flex gap-2 mt-2">
                {permissions.canDeletePost && (
                  <span className="badge badge-sm badge-success">
                    {t("can_delete_post")}
                  </span>
                )}
                {permissions.canPinPost && (
                  <span className="badge badge-sm badge-info">
                    {t("can_pin_post")}
                  </span>
                )}
                {permissions.canBanUser && (
                  <span className="badge badge-sm badge-warning">
                    {t("can_ban_user")}
                  </span>
                )}
                {permissions.canManageModerator && (
                  <span className="badge badge-sm badge-error">
                    {t("can_manage_moderator")}
                  </span>
                )}
              </div>
            )}
          </div>

          {/* 内容 */}
          <div className="min-h-[calc(100vh-120px)]">{renderContent()}</div>
        </div>
      </div>
    </div>
  );
}
