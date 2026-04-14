"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useAdminAuth } from "@/hooks/admin/useAdminAuth";
import { useUsersData } from "@/hooks/admin/useUsersData";
import { usePostsData } from "@/hooks/admin/usePostsData";
import { useAnnouncementsData } from "@/hooks/admin/useAnnouncementsData";
import { useQAData } from "@/hooks/admin/useQAData";
import { usePointsData } from "@/hooks/admin/usePointsData";
import { useStatsData } from "@/hooks/admin/useStatsData";
import { AnnouncementsManager } from "@/components/admin/AnnouncementsManager";
import { AdminSearchBar } from "@/components/admin/AdminSearchBar";
import { Pagination } from "@/components/admin/Pagination";
import { PointsManager } from "@/components/admin/PointsManager";
import { PostsTable } from "@/components/admin/PostsTable";
import { SidebarMenu } from "@/components/admin/SidebarMenu";
import { Dashboard } from "@/components/admin/Dashboard";
import { UsersTable } from "@/hooks/admin/UsersTable";
import { Statistics } from "@/components/admin/Statistics";

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
  const pointsData = usePointsData(activeMenu === "points" && isAdmin);
  const statsData = useStatsData(activeMenu === "statistics" && isAdmin);

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
              onToggleActive={usersData.toggleActive}
              isToggling={usersData.isToggling}
              t={t}
            />
            <Pagination
              currentPage={page}
              total={usersData.total}
              onPageChange={setPage}
            />
          </div>
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
                  {t("qa_management_coming")}
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
