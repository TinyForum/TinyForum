"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useAdminAuth } from "@/hooks/admin/useAdminAuth";
import { useUsersData } from "@/hooks/admin/useUsersData";
import { usePostsData } from "@/hooks/admin/usePostsData";
import { useAnnouncementsData } from "@/hooks/admin/useAnnouncementsData";
import { useQAData } from "@/hooks/admin/useQAData";
import { useScoreData } from "@/hooks/admin/useScoreData";
import { useStatsData } from "@/hooks/admin/useStatsData";
import { AnnouncementsManager } from "@/components/admin/AnnouncementsManager";
import { AdminSearchBar } from "@/components/admin/AdminSearchBar";
import { Pagination } from "@/components/admin/Pagination";
import { PointsManager } from "@/components/admin/PointsManager";
import { PostsTable } from "@/components/admin/PostsTable";
import { SidebarMenu } from "@/components/admin/SidebarMenu";
import { Dashboard } from "@/components/admin/Dashboard";
import { UsersTable } from "@/components/admin/UsersTable"; 
import { Statistics } from "@/components/admin/Statistics";
import { ModeratorsTable } from "@/components/admin/ModeratorsTable";
import { AdminTasks } from "@/components/admin/AdminTasks";

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
  const announcementsData = useAnnouncementsData(
    activeMenu === "announcements" && isAdmin,
  );
  const qaData = useQAData(page, keyword, activeMenu === "qa" && isAdmin);
  
  // 修复：useScoreData 期望可选的 userId 参数，不是布尔值
  const pointsData = useScoreData(activeMenu === "points" && isAdmin ? undefined : undefined);
  
  const statsData = useStatsData(activeMenu === "statistics" && isAdmin);

  // 处理用户激活/停用
  const handleToggleActive = async (userId: number, active: boolean) => {
    await usersData.toggleActive(userId, active);
  };

  // 处理用户封禁/解封
  const handleToggleBlock = async (userId: number, blocked: boolean) => {
    // 需要确认 usersData 中是否有 toggleBlock 方法
    // 如果没有，需要添加这个功能到 useUsersData hook 中
    await usersData.toggleBlock?.(userId, blocked);
  };

  // 处理用户角色变更
  const handleToggleRole = async (userId: number, role: string) => {
    // 需要确认 usersData 中是否有 toggleRole 方法
    await usersData.toggleRole?.(userId, role);
  };

  // 处理删除用户
  const handleDeleteUser = async (userId: number, username: string) => {
    // 需要确认 usersData 中是否有 deleteUser 方法
    await usersData.deleteUser?.(userId, username);
  };

  // 处理重置密码
  const handleResetPassword = async (userId: number, username: string) => {
    // 需要确认 usersData 中是否有 resetPassword 方法
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
      // MARK: 统计
      case "dashboard":
        return <Dashboard t={t} />;
      case "tasks":
        return <AdminTasks/>
        
      // MARK: 公告
      case "announcements":
        return <AnnouncementsManager t={t} />;

      // MARK: 用户
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
              users={usersData.users}
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
              t={t}
            />
            <Pagination
              currentPage={page}
              total={usersData.total}
              onPageChange={setPage}
            />
          </div>
        );
      
  
        
      case "moderators_management":
        return (
          <ModeratorsTable />
        );
      
      // MARK: 帖子
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

      // MARK: QA
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
            {/* QA 表格组件 */}
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body">
                <p className="text-center text-base-content/50">
                  {/* {t("qa_management_coming")} */}
                  TODO
                </p>
              </div>
            </div>
            <Pagination
              currentPage={page}
              total={qaData.total}
              onPageChange={setPage}
            />
          </div>
        );
      
      // MARK: 积分
      case "points":
        return <PointsManager t={t} />;
      
      // MARK: 统计
      case "statistics":
        return <Statistics t={t} />;
      
      // MARK: 设置
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
      {/* 左侧业务面板 */}
      <SidebarMenu
        activeMenu={activeMenu}
        onMenuChange={setActiveMenu}
        collapsed={sidebarCollapsed}
        onCollapsedChange={setSidebarCollapsed}
        t={t}
      />

      {/* 右侧内容区域 */}
      <div className="flex-1 overflow-y-auto">
        <div className="p-6">
          {/* 页面标题 */}
          <div className="mb-6">
            <h1 className="text-2xl font-bold">{t(activeMenu)}</h1>
            <p className="text-sm text-base-content/60 mt-1">
              {t(`${activeMenu}_description`)}
            </p>
          </div>

          {/* 内容 */}
          <div className="min-h-[calc(100vh-120px)]">{renderContent()}</div>
        </div>
      </div>
    </div>
  );
}