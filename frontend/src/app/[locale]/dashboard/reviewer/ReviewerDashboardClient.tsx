"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { ReviewerSidebar } from "@/layout/reviewer/ReviewerSidebar";
import { ReviewerStats } from "@/layout/reviewer/ReviewerStats";
import { ReviewSettings } from "@/layout/reviewer/ReviewSettings";
import { ReviewerTaskList } from "@/features/reviewer/components/ReviewerTaskList";
import { ReviewerHistoryList } from "@/features/reviewer/components/ReviewerHistoryList";
import { useReviewerTasks } from "@/features/reviewer/hooks/useReviewerTasks";

export default function ReviewerDashboardClient() {
  const t = useTranslations("Reviewer");

  const [activeMenu, setActiveMenu] = useState("pending");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  // 待审核任务作为统计与列表的统一数据源
  const { list, total } = useReviewerTasks();

  const pendingCount = total;
  const reportedCount = list.filter(
    (task) => task.trigger_type === "report_aggregate",
  ).length;
  const stats = { pending: pendingCount, reported: reportedCount, reviewedToday: 0 };

  const labels: Record<string, string> = {
    pending: t("pending_review"),
    reports: t("manage_reports"),
    history: t("review_history"),
    settings: t("settings"),
  };

  const descriptions: Record<string, string> = {
    pending: t("pending_description"),
    reports: t("reports_description"),
    history: t("history_description"),
    settings: t("settings_description"),
  };

  const menus = [
    { id: "pending", label: labels.pending, icon: "📋" },
    { id: "reports", label: labels.reports, icon: "🚫" },
    { id: "history", label: labels.history, icon: "📜" },
    { id: "settings", label: labels.settings, icon: "⚙️" },
  ];

  const renderContent = () => {
    switch (activeMenu) {
      case "pending":
        return <ReviewerTaskList />;
      case "reports":
        return <ReviewerTaskList triggerTypeFilter="report_aggregate" />;
      case "history":
        return <ReviewerHistoryList />;
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
              {t("reviewer_panel")} - {labels[activeMenu]}
            </h1>
            <p className="text-sm text-base-content/60 mt-1">
              {descriptions[activeMenu]}
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
