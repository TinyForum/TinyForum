// app/admin/tasks/page.tsx
"use client";

import { useState } from "react";
import { Shield, Flag, FileText, Ban, Bell } from "lucide-react";
import { useTranslations } from "next-intl";
import { ApplicationManagement } from "./ApplicationManagement";
import { BanManagement } from "./BanManagement";
import { ModeratorManagement } from "./ModeratorManagement";
import { PostManagement } from "./PostManagement";
import { ReportManagement } from "./ReportManagement";
import { ReviewManagement } from "./ReviewManagement";

// ============ 主组件 ============
export function AdminTasks() {
  const t = useTranslations("Admin");
  const [activeTab, setActiveTab] = useState("risks");

  const tabs = [
    { id: "risks", label: t("risk_review_tab"), icon: Bell },
    { id: "applications", label: t("applications_tab"), icon: FileText },
    { id: "moderators", label: t("moderators_tab"), icon: Shield },
    { id: "reports", label: t("reports_tab"), icon: Flag },
    { id: "posts", label: t("posts_tab"), icon: FileText },
    { id: "bans", label: t("bans_tab"), icon: Ban },
  ];

  return (
    <div className="flex flex-col gap-6">
      <div role="tablist" className="tabs tabs-lifted">
        {tabs.map((tab) => (
          <a
            key={tab.id}
            role="tab"
            className={`tab ${activeTab === tab.id ? "tab-active" : ""}`}
            onClick={() => setActiveTab(tab.id)}
          >
            <tab.icon className="w-4 h-4 mr-2" />
            {tab.label}
          </a>
        ))}
      </div>

      <div className="mt-6">
        {activeTab === "applications" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <FileText className="w-5 h-5" />
                {t("application_approval")}
              </h2>
              <p className="text-gray-500">
                {t("application_approval_desc")}
              </p>
              <div className="mt-4">
                <ApplicationManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "risks" && (
          <>
            <ReviewManagement />
          </>
        )}
        {activeTab === "moderators" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <Shield className="w-5 h-5" />
                {t("moderator_appointment")}
              </h2>
              <p className="text-gray-500">{t("moderator_appointment_desc")}</p>
              <div className="mt-4">
                <ModeratorManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "reports" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <Flag className="w-5 h-5" />
                {t("report_handling")}
              </h2>
              <p className="text-gray-500">{t("report_handling_desc")}</p>
              <div className="mt-4">
                <ReportManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "posts" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <FileText className="w-5 h-5" />
                {t("post_management")}
              </h2>
              <p className="text-gray-500">
                {t("post_management_desc")}
              </p>
              <div className="mt-4">
                <PostManagement />
              </div>
            </div>
          </div>
        )}

        {activeTab === "bans" && (
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">
                <Ban className="w-5 h-5" />
                {t("ban_management")}
              </h2>
              <p className="text-gray-500">
                {t("ban_management_desc")}
              </p>
              <div className="mt-4">
                <BanManagement />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
