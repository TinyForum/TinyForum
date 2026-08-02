// dashboard/member/page.tsx 客户端组件
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { MemberProfile } from "@/features/member/MemberProfile";
import { MemberSidebar } from "@/features/member/MemberSidebar";
import { MemberStats } from "@/features/member/MemberStats";
import { useUserStats } from "@/features/user/hooks/useUserStats";

export default function MemberDashboardClient() {
  const t = useTranslations("Member");

  const [activeMenu, setActiveMenu] = useState("dashboard");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const { total_post, total_comment, total_favorite, isLoading } =
    useUserStats();

  const stats = isLoading
    ? { posts: 0, comments: 0, favorites: 0, unreadNotif: 0 }
    : {
        posts: total_post ?? 0,
        comments: total_comment ?? 0,
        favorites: total_favorite ?? 0,
        unreadNotif: 0,
      };

  const menus = [
    { id: "dashboard", label: t("dashboard"), icon: "📊" },
    { id: "posts", label: t("my_posts"), icon: "📝", badge: stats.posts },
    { id: "comments", label: t("my_comments"), icon: "💬", badge: stats.comments },
    { id: "favorites", label: t("my_favorites"), icon: "❤️", badge: stats.favorites },
    { id: "notifications", label: t("notifications"), icon: "🔔", badge: stats.unreadNotif },
    { id: "profile", label: t("profile"), icon: "👤" },
  ];

  const renderPlaceholder = (feature: string) => (
    <div className="flex flex-col items-center justify-center py-20 text-center">
      <span className="text-4xl mb-4">🚧</span>
      <h3 className="text-lg font-semibold">{t("feature_under_development")}</h3>
      <p className="text-sm text-base-content/60 mt-2">
        {t("feature_coming_soon", { feature })}
      </p>
    </div>
  );

  const renderContent = () => {
    switch (activeMenu) {
      case "dashboard":
        return <MemberStats stats={stats} />;
      case "posts":
        return renderPlaceholder(t("my_posts"));
      case "comments":
        return renderPlaceholder(t("my_comments"));
      case "favorites":
        return renderPlaceholder(t("my_favorites"));
      case "notifications":
        return renderPlaceholder(t("notifications"));
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
