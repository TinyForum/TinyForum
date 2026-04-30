// app/[locale]/dashboard/user/page.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Pagination } from "@/features/admin/components/Pagination";
import { FavoritesList } from "@/features/user/components/FavoritesList";
import { MyCommentsTable } from "@/features/user/components/MyCommentsTable";
import { MyPostsTable } from "@/features/user/components/MyPostsTable";
// import { NotificationsList } from "@/features/user/components/NotificationsList";
import { ViolationPanel } from "@/features/user/components/ViolationPanel";
import { StatCard } from "@/features/user/components/StatCard";
import { UserInfoCard } from "@/features/user/components/UserInfoCard";
import { NotificationsList } from "@/features/user/components/NotificationsList";
import SearchBar from "@/features/moderator/components/SearchBar";
import { useAuthStore } from "@/store";
import { useUserRole } from "@/features/user/hooks/useUserRole";

// 组件导入
export default function UserDashboardPage() {
  const t = useTranslations("User");
  const [activeTab, setActiveTab] = useState("overview");
  const [keyword, setKeyword] = useState("");
  const [page, setPage] = useState(1);

  // TODO: 从 hooks 获取用户信息
  // const { user } = useCurrentUser();
  const { user } = useAuthStore();
  const { isAdmin } = useUserRole();
  const { isModerator } = useUserRole();

  // TODO: 获取各模块数据
  // const { posts, total: postsTotal } = useMyPosts({ page, keyword });
  // const { comments, total: commentsTotal } = useMyComments({ page });
  // const { favorites, total: favoritesTotal } = useMyFavorites({ page });
  // const { notifications, unreadCount, total: notificationsTotal } = useMyNotifications({ page });

  // TODO: 统计数据
  const stats = {
    posts: 0,
    comments: 0,
    favorites: 0,
    unreadNotif: 0,
  };

  const tabs = [
    { id: "overview", label: t("overview"), icon: "📊" },
    { id: "posts", label: t("my_posts"), icon: "📝", badge: stats.posts },
    {
      id: "comments",
      label: t("my_comments"),
      icon: "💬",
      badge: stats.comments,
    },
    {
      id: "favorites",
      label: t("favorites"),
      icon: "❤️",
      badge: stats.favorites,
    },
    {
      id: "notifications",
      label: t("notifications"),
      icon: "🔔",
      badge: stats.unreadNotif,
    },
    { id: "violation_records", label: t("violation_records"), icon: "🚫" },
  ];

  const renderContent = () => {
    switch (activeTab) {
      case "overview":
        return (
          <div className="space-y-6">
            <UserInfoCard
              user={user}
              isAdmin={isAdmin}
              isModerator={isModerator}
            />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <StatCard
                title={t("total_posts")}
                value={stats.posts}
                icon="📝"
                color="text-primary"
              />
              <StatCard
                title={t("total_comments")}
                value={stats.comments}
                icon="💬"
                color="text-secondary"
              />
              <StatCard
                title={t("total_favorites")}
                value={stats.favorites}
                icon="❤️"
                color="text-error"
              />
              <StatCard
                title={t("unread_notifications")}
                value={stats.unreadNotif}
                icon="🔔"
                color="text-warning"
              />
            </div>
          </div>
        );

      case "posts":
        return (
          <div className="space-y-4">
            <SearchBar
              keyword={keyword}
              onKeywordChange={setKeyword}
              placeholder={t("search_my_posts")}
            />
            <MyPostsTable posts={[]} />
            <Pagination
              currentPage={page}
              total={stats.posts}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "comments":
        return (
          <div className="space-y-4">
            <MyCommentsTable comments={[]} />
            <Pagination
              currentPage={page}
              total={stats.comments}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "favorites":
        return (
          <div className="space-y-4">
            <FavoritesList favorites={[]} />
            <Pagination
              currentPage={page}
              total={stats.favorites}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "notifications":
        return (
          <div className="space-y-4">
            <NotificationsList notifications={[]} />
            <Pagination
              currentPage={page}
              total={stats.comments}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "violation_records":
        return <ViolationPanel />;

      default:
        return null;
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">{t("user_center")}</h1>
        <p className="text-base-content/60 mt-1">{t("manage_your_content")}</p>
      </div>

      <div className="border-b border-base-300 mb-6 overflow-x-auto">
        <div className="flex gap-2">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => {
                setActiveTab(tab.id);
                setPage(1);
              }}
              className={`px-4 py-2 font-medium transition-colors relative whitespace-nowrap ${
                activeTab === tab.id
                  ? "text-primary border-b-2 border-primary"
                  : "text-base-content/60 hover:text-base-content"
              }`}
            >
              <span className="flex items-center gap-2">
                <span>{tab.icon}</span>
                <span>{tab.label}</span>
                {tab.badge !== undefined && tab.badge > 0 && (
                  <span className="badge badge-sm badge-primary">
                    {tab.badge}
                  </span>
                )}
              </span>
            </button>
          ))}
        </div>
      </div>

      <div className="min-h-[calc(100vh-200px)]">{renderContent()}</div>
    </div>
  );
}
