import { LayoutDashboard } from "lucide-react";

// 头部组件
export function AdminHeader({ t }: { t: (key: string) => string }) {
  return (
    <div className="flex items-center gap-3 mb-6">
      <LayoutDashboard className="w-6 h-6 text-primary" />
      <h1 className="text-2xl font-bold">{t("title")}</h1>
    </div>
  );
}
