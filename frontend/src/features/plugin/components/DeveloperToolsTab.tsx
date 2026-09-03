import { Wrench } from "lucide-react";
import { useTranslations } from "next-intl";

// ==================== 开发者工具 ====================
export function DeveloperToolsTab() {
  const t = useTranslations("Plugin");

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 text-sm font-medium">
        <div>
          <Wrench className="w-5 h-5" />
          <span>{t("dev_tools_title")}</span>
        </div>
      </div>
      <details className="collapse collapse-arrow bg-base-200">
        <summary className="collapse-title text-sm font-medium">
          {t("hooks_list")}
        </summary>
        <div className="collapse-content">
          <pre className="text-xs">post:create, user:login, ...</pre>
        </div>
      </details>
      <button className="btn btn-outline btn-sm">{t("export_debug_logs")}</button>
    </div>
  );
}
