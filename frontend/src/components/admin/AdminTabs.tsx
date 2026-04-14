import { TabType } from "@/type/admin.types";
import { FileText, Users } from "lucide-react";

// Tab 切换组件
export function AdminTabs({ 
  activeTab, 
  onTabChange, 
  t 
}: { 
  activeTab: TabType; 
  onTabChange: (tab: TabType) => void; 
  t: (key: string) => string;
}) {
  return (
    <div className="tabs tabs-boxed bg-base-100 border border-base-300 mb-4 p-1 w-fit">
      <button
        className={`tab gap-2 ${activeTab === "users" ? "tab-active" : ""}`}
        onClick={() => onTabChange("users")}
      >
        <Users className="w-4 h-4" /> {t("user_management")}
      </button>
      <button
        className={`tab gap-2 ${activeTab === "posts" ? "tab-active" : ""}`}
        onClick={() => onTabChange("posts")}
      >
        <FileText className="w-4 h-4" /> {t("post_management")}
      </button>
    </div>
  );
}