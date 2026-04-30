// app/[locale]/dashboard/reviewer/page.tsx
"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Pagination } from "@/features/admin/components/Pagination";
import { PendingContentTable } from "@/layout/reviewer/PendingContentTable";
import { ReviewerSidebar } from "@/layout/reviewer/ReviewerSidebar";
import { ReviewerStats } from "@/layout/reviewer/ReviewerStats";
import { ReviewHistoryTable } from "@/layout/reviewer/ReviewHistoryTable";
import { ReviewSettings } from "@/layout/reviewer/ReviewSettings";
import { ReportedContentTable } from "@/layout/reviewer/ReportedContentTable"; // TODO：应该合并
import SearchBar from "@/shared/ui/nav/SearchBar";
// import { ReportedContentTable } from "@/features/moderator/components/ReportedContentTable";
// 组件导入

export default function ReviewerPage() {
  const t = useTranslations("Reviewer");

  const [activeMenu, setActiveMenu] = useState("pending");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [page, setPage] = useState(1);

  // TODO: 获取审核员权限
  // const { permissions } = useReviewerPermissions();

  // TODO: 获取待审核内容
  // const { data: pendingData, isLoading: pendingLoading } = usePendingContent({
  //   page,
  //   keyword,
  //   enabled: activeMenu === "pending",
  // });

  // TODO: 获取举报内容
  // const { data: reportsData, isLoading: reportsLoading } = useReviewedReports({
  //   page,
  //   enabled: activeMenu === "reports",
  // });

  // TODO: 获取审核历史
  // const { data: historyData, isLoading: historyLoading } = useReviewHistory({
  //   page,
  //   enabled: activeMenu === "history",
  // });

  // TODO: 审核操作
  // const { mutate: approveContent, isPending: isApproving } = useApproveContent();
  // const { mutate: rejectContent, isPending: isRejecting } = useRejectContent();
  // const { mutate: deleteContent, isPending: isDeleting } = useDeleteContent();

  // TODO: 处理审核通过
  const handleApprove = (id: number, type: string) => {
    console.log("approve", id, type);
    // approveContent({ id, type });
  };

  // TODO: 处理审核拒绝
  const handleReject = (id: number, type: string, reason?: string) => {
    console.log("reject", id, type, reason);
    // rejectContent({ id, type, reason });
  };

  // TODO: 处理删除
  const handleDelete = (id: number, type: string) => {
    console.log("delete", id, type);
    // if (confirm(t("confirm_delete"))) {
    //   deleteContent({ id, type });
    // }
  };

  const menus = [
    { id: "pending", label: t("pending_review"), icon: "📋" },
    { id: "reports", label: t("reported_content"), icon: "🚫" },
    { id: "history", label: t("review_history"), icon: "📜" },
    { id: "settings", label: t("settings"), icon: "⚙️" },
  ];

  // TODO: 统计数据
  const stats = {
    pending: 0,
    reported: 0,
    reviewedToday: 0,
  };

  const renderContent = () => {
    switch (activeMenu) {
      case "pending":
        return (
          <div className="space-y-4">
            <SearchBar
              keyword={keyword}
              onKeywordChange={(value: string) => {
                setKeyword(value);
                setPage(1);
              }}
              placeholder={t("search_pending")}
            />
            <PendingContentTable
              onApprove={handleApprove}
              onReject={handleReject}
              onDelete={handleDelete}
              reviews={[]}
            />
            <Pagination
              currentPage={page}
              total={0}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "reports":
        return (
          <div className="space-y-4">
            <SearchBar
              keyword={keyword}
              onKeywordChange={(value: string) => {
                setKeyword(value);
                setPage(1);
              }}
              placeholder={t("search_reports")}
            />
            <ReportedContentTable
              onApprove={handleApprove}
              onReject={handleReject}
              onDelete={handleDelete}
              reports={[]}
            />
            <Pagination
              currentPage={page}
              total={0}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

      case "history":
        return (
          <div className="space-y-4">
            <ReviewHistoryTable historys={[]} />
            <Pagination
              currentPage={page}
              total={0}
              pageSize={20}
              onPageChange={setPage}
            />
          </div>
        );

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
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-2xl font-bold">
                  {t("reviewer_panel")} - {t(activeMenu)}
                </h1>
                <p className="text-sm text-base-content/60 mt-1">
                  {t(`${activeMenu}_description`)}
                </p>
              </div>
            </div>

            {activeMenu === "pending" && (
              <div className="mt-4">
                <ReviewerStats stats={stats} />
              </div>
            )}
          </div>

          <div className="min-h-[calc(100vh-120px)]">{renderContent()}</div>
        </div>
      </div>
    </div>
  );
}
