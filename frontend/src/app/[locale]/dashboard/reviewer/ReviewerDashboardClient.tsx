"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { ReviewerSidebar } from "@/layout/reviewer/ReviewerSidebar";
import { ReviewerStats } from "@/layout/reviewer/ReviewerStats";
import { ReviewSettings } from "@/layout/reviewer/ReviewSettings";

export default function ReviewerDashboardClient() {
  const t = useTranslations("Reviewer");

  const [activeMenu, setActiveMenu] = useState("pending");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const stats = { pending: 0, reported: 0, reviewedToday: 0 };

  const menus = [
    { id: "pending", label: t("pending_review"), icon: "📋" },
    { id: "reports", label: t("reported_content"), icon: "🚫" },
    { id: "history", label: t("review_history"), icon: "📜" },
    { id: "settings", label: t("settings"), icon: "⚙️" },
  ];

  const renderPlaceholder = (feature: string) => (
    <div className="flex flex-col items-center justify-center py-20 text-center">
      <span className="text-4xl mb-4">🚧</span>
      <h3 className="text-lg font-semibold">
        {t("feature_under_development")}
      </h3>
      <p className="text-sm text-base-content/60 mt-2">
        {t("feature_coming_soon", { feature })}
      </p>
    </div>
  );

  const renderContent = () => {
    switch (activeMenu) {
      case "pending":
        return renderPlaceholder(t("pending_review"));
      case "reports":
        return renderPlaceholder(t("reported_content"));
      case "history":
        return renderPlaceholder(t("review_history"));
      case "settings":
        return <ReviewSettings />;
      default:
        return null;
    }
  };

  return (
    <div className="flex h-screen bg-base-100">
      <ReviewerSidebar
        activeMenu={activeMenu}
        onMenuChange={setActiveMenu}
        collapsed={sidebarCollapsed}
        onCollapsedChange={setSidebarCollapsed}
        menus={menus}
        stats={stats}
      />

      <div className="flex-1 overflow-y-auto">
        <div className="p-6">
          <div className="mb-6">
            <h1 className="text-2xl font-bold">
              {t("reviewer_panel")} - {t(activeMenu)}
            </h1>
            <p className="text-sm text-base-content/60 mt-1">
              {t(`${activeMenu}_description`)}
            </p>
          </div>

          {activeMenu === "pending" && (
            <div className="mt-4">
              <ReviewerStats stats={stats} />
            </div>
          )}

          <div className="min-h-[calc(100vh-120px)]">{renderContent()}</div>
        </div>
      </div>
    </div>
  );
}
