// app/admin/tasks/page.tsx
"use client";

import { useState } from "react";
import { Shield, Flag, FileText, Ban, Bell } from "lucide-react";
import { ApplicationManagement } from "./ApplicationManagement";
import { BanManagement } from "./BanManagement";
import { ModeratorManagement } from "./ModeratorManagement";
import { PostManagement } from "./PostManagement";
import { ReportManagement } from "./ReportManagement";
import { ReviewManagement } from "./ReviewManagement";

// ============ 主组件 ============
export function AdminTasks() {
  const [activeTab, setActiveTab] = useState("risks");

  const tabs = [
    { id: "risks", label: "风控审查", icon: Bell },
    { id: "applications", label: "版主申请", icon: FileText },
    { id: "moderators", label: "版主任命", icon: Shield },
    { id: "reports", label: "举报管理", icon: Flag },
    { id: "posts", label: "帖子管理", icon: FileText },
    { id: "bans", label: "禁言管理", icon: Ban },
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
                版主申请审批
              </h2>
              <p className="text-gray-500">
                审核用户提交的版主申请，通过后用户将成为对应板块的版主
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
                版主任命与管理
              </h2>
              <p className="text-gray-500">任命新版主、管理现有版主及其权限</p>
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
                举报处理
              </h2>
              <p className="text-gray-500">处理用户举报，维护社区秩序</p>
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
                帖子管理
              </h2>
              <p className="text-gray-500">
                搜索、置顶、删除帖子，管理板块内容
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
                禁言管理
              </h2>
              <p className="text-gray-500">
                禁言违规用户，查看禁言列表，解除禁言
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
