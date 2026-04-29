"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";

import { User } from "@/shared/api/types";
import { AdminSearchBar } from "@/features/admin/components/AdminSearchBar";
import { AdminTasks } from "@/features/admin/components/AdminTasks";
import { AnnouncementsManager } from "@/features/admin/components/AnnouncementsManager";
import { Dashboard } from "@/features/admin/components/Dashboard";
import { ModeratorsTable } from "@/features/admin/components/ModeratorsTable";
import { Pagination } from "@/features/admin/components/Pagination";
import { PointsManager } from "@/features/admin/components/PointsManager";
import { PostsTable } from "@/features/admin/components/PostsTable";
import { SidebarMenu } from "@/features/admin/components/SidebarMenu";
import { Statistics } from "@/features/admin/components/Statistics";
import { UsersTable } from "@/features/admin/components/UsersTable";
import { useAdminAuth } from "@/features/admin/hooks/useAdminAuth";
import { usePostsData } from "@/features/admin/hooks/usePostsData";
import { useQAData } from "@/features/admin/hooks/useQAData";
import { useUsersData } from "@/features/admin/hooks/useUsersData";

// 导入类型

// ==================== 主组件 ====================
export default function AdminPage() {
  const t = useTranslations("Admin");
  const { isCheckingAuth, isAdmin, user } = useAdminAuth();

  const [activeMenu, setActiveMenu] = useState("dashboard");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [page, setPage] = useState(1);

  // 数据获取（只在对应菜单激活时获取）
  const usersData = useUsersData(
    page,
    keyword,
    activeMenu === "users" && isAdmin,
  );
  const postsData = usePostsData(
    page,
    keyword,
    activeMenu === "posts" && isAdmin,
  );

  const qaData = useQAData();

  // 处理用户激活/停用
  const handleToggleActive = async (userId: number, active: boolean) => {
    await usersData.toggleActive(userId, active);
  };

  // 处理用户封禁/解封
  const handleToggleBlock = async (userId: number, blocked: boolean) => {
    await usersData.toggleBlock?.(userId, blocked);
  };

  // 处理用户角色变更
  const handleToggleRole = async (userId: number, role: string) => {
    await usersData.toggleRole?.(userId, role);
  };

  // 处理删除用户
  const handleDeleteUser = async (userId: number) => {
    await usersData.deleteUser?.(userId);
  };

  // 处理重置密码
  const handleResetPassword = async (userId: number) => {
    await usersData.resetPassword?.(userId);
  };

  // 加载状态
  if (isCheckingAuth) {
    return (
      <div className="flex justify-center items-center h-screen">
        <span className="loading loading-spinner loading-lg text-primary" />
      </div>
    );
  }

  if (!isAdmin) {
    return null;
  }

  // 渲染右侧内容
  const renderContent = () => {
    switch (activeMenu) {
      case "dashboard":
        return <Dashboard t={t} />;

      case "tasks":
        return <AdminTasks />;

      case "announcements":
        return <AnnouncementsManager t={t} />;

      case "users":
        return (
          <div className="space-y-4">
            <AdminSearchBar
              tab="users"
              keyword={keyword}
              onKeywordChange={(k) => {
                setKeyword(k);
                setPage(1);
              }}
              onPageReset={() => setPage(1)}
              t={t}
            />
            <UsersTable
              users={usersData.users as User[]}
              currentUserId={user?.id}
              onToggleActive={handleToggleActive}
              onToggleBlock={handleToggleBlock}
              onToggleRole={handleToggleRole}
              onDeleteUser={handleDeleteUser}
              onResetPassword={handleResetPassword}
              isTogglingActive={usersData.isTogglingActive || false}
              isTogglingBlock={usersData.isTogglingBlock || false}
              isDeleting={usersData.isDeleting || false}
              isUpdatingRole={usersData.isUpdatingRole || false}
            />
            <Pagination
              currentPage={page}
              total={usersData.total}
              onPageChange={setPage}
            />
          </div>
        );

      case "moderators_management":
        return <ModeratorsTable boardId={0} />;

      case "posts":
        return (
          <div className="space-y-4">
            <AdminSearchBar
              tab="posts"
              keyword={keyword}
              onKeywordChange={(k) => {
                setKeyword(k);
                setPage(1);
              }}
              onPageReset={() => setPage(1)}
              t={t}
            />
            <PostsTable
              posts={postsData.posts}
              onTogglePin={postsData.togglePin}
              isToggling={postsData.isToggling}
              t={t}
            />
            <Pagination
              currentPage={page}
              total={postsData.total}
              onPageChange={setPage}
            />
          </div>
        );

      case "qa":
        return (
          <div className="space-y-4">
            <AdminSearchBar
              tab="qa"
              keyword={keyword}
              onKeywordChange={(k) => {
                setKeyword(k);
                setPage(1);
              }}
              onPageReset={() => setPage(1)}
              t={t}
            />
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body">
                <p className="text-center text-base-content/50">TODO</p>
              </div>
            </div>
            <Pagination
              currentPage={page}
              total={qaData.total}
              onPageChange={setPage}
            />
          </div>
        );

      case "points":
        return <PointsManager />;

      case "statistics":
        return <Statistics />;

      case "settings":
        return (
          <div className="card bg-base-100 border border-base-300">
            <div className="card-body">
              <h3 className="font-semibold mb-4">{t("system_settings")}</h3>
              <p className="text-center text-base-content/50">
                {t("settings_coming")}
              </p>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="flex h-screen bg-base-100">
      <SidebarMenu
        activeMenu={activeMenu}
        onMenuChange={setActiveMenu}
        collapsed={sidebarCollapsed}
        onCollapsedChange={setSidebarCollapsed}
        t={t}
      />

      <div className="flex-1 overflow-y-auto">
        <div className="p-6">
          <div className="mb-6">
            <h1 className="text-2xl font-bold">{t(activeMenu)}</h1>
            <p className="text-sm text-base-content/60 mt-1">
              {t(`${activeMenu}_description`)}
            </p>
          </div>

          <div className="min-h-[calc(100vh-120px)]">{renderContent()}</div>
        </div>
      </div>
    </div>
  );
}
