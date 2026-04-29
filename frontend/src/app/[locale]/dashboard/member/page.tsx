// app/[locale]/dashboard/member/page.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Pagination } from "@/features/admin/components/Pagination";
import { MemberCommentsTable } from "@/features/member/MemberCommentsTable";
import { MemberFavorites } from "@/features/member/MemberFavorites";
import { MemberNotifications } from "@/features/member/MemberNotifications";
import { MemberPostsTable } from "@/features/member/MemberPostsTable";
import { MemberProfile } from "@/features/member/MemberProfile";
import { MemberSidebar } from "@/features/member/MemberSidebar";
import { MemberStats } from "@/features/member/MemberStats";
import SearchBar from "@/features/moderator/components/SearchBar"; // TODO: 修改为共享组件

// 组件导入

export default function MemberDashboard() {
  const t = useTranslations("Member");

  const [activeMenu, setActiveMenu] = useState("dashboard");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [page, setPage] = useState(1);

  // TODO: 获取会员信息
  // const { member, loading: memberLoading } = useCurrentMember();

  // TODO: 获取统计数据
  // const { data: stats, loading: statsLoading } = useMemberStats();

  // TODO: 获取帖子列表
  // const { data: postsData, loading: postsLoading } = useMemberPosts({
  //   page,
  //   keyword,
  //   enabled: activeMenu === "posts",
  // });

  // TODO: 获取评论列表
  // const { data: commentsData, loading: commentsLoading } = useMemberComments({
  //   page,
  //   enabled: activeMenu === "comments",
  // });

  // TODO: 获取收藏列表
  // const { data: favoritesData, loading: favoritesLoading } = useMemberFavorites({
  //   page,
  //   enabled: activeMenu === "favorites",
  // });

  // TODO: 获取通知列表
  // const { data: notificationsData, loading: notificationsLoading } = useMemberNotifications({
  //   page,
  //   enabled: activeMenu === "notifications",
  // });

  // TODO: 删除帖子
  // const { mutate: deletePost, isPending: isDeleting } = useDeleteMemberPost();

  // TODO: 删除评论
  // const { mutate: deleteComment, isPending: isDeletingComment } = useDeleteMemberComment();

  // TODO: 取消收藏
  // const { mutate: removeFavorite, isPending: isRemoving } = useRemoveFavorite();

  // TODO: 标记通知已读
  // const { mutate: markAsRead } = useMarkNotificationRead();

  // TODO: 统计数据
  const stats = {
    posts: 0,
    comments: 0,
    favorites: 0,
    unreadNotif: 0,
  };

  const menus = [
    { id: "dashboard", label: t("dashboard"), icon: "📊" },
    { id: "posts", label: t("my_posts"), icon: "📝", badge: stats.posts },
    {
      id: "comments",
      label: t("my_comments"),
      icon: "💬",
      badge: stats.comments,
    },
    {
      id: "favorites",
      label: t("my_favorites"),
      icon: "❤️",
      badge: stats.favorites,
    },
    {
      id: "notifications",
      label: t("notifications"),
      icon: "🔔",
      badge: stats.unreadNotif,
    },
    { id: "profile", label: t("profile"), icon: "👤" },
  ];

  const renderContent = () => {
    switch (activeMenu) {
      case "dashboard":
        return <MemberStats stats={stats} />;

      case "posts":
        return (
          <div className="space-y-4">
            <SearchBar
              keyword={keyword}
              onKeywordChange={(value: string) => {
                setKeyword(value);
                setPage(1);
              }}
              placeholder={t("search_my_posts")}
            />
            <MemberPostsTable
              onDelete={(id) => {
                console.log(id);
                // TODO: 删除帖子
              }}
            />
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
            <MemberCommentsTable
              onDelete={(id) => {
                console.log(id);
                // TODO: 删除评论
              }}
            />
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
            <MemberFavorites
              onRemove={(id) => {
                console.log(id);
                // TODO: 取消收藏
              }}
              favorites={[]}
            />
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
            <MemberNotifications
              onMarkRead={(id) => {
                console.log(id);
                // TODO: 标记已读
              }}
              onMarkAllRead={() => {
                // TODO: 全部标记已读
              }}
              notifications={[]}
            />
            <Pagination
              currentPage={page}
              total={stats.unreadNotif}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "profile":
        return <MemberProfile />;

      default:
        return null;
    }
  };

  return (
    <div className="flex h-screen bg-base-100">
      <MemberSidebar
        activeMenu={activeMenu}
        onMenuChange={setActiveMenu}
        collapsed={sidebarCollapsed}
        onCollapsedChange={setSidebarCollapsed}
        menus={menus}
      />

      <div className="flex-1 overflow-y-auto">
        <div className="p-6">
          <div className="mb-6">
            <h1 className="text-2xl font-bold">{t("member_center")}</h1>
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
