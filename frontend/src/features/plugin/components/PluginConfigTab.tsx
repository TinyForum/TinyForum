import { useTranslations } from "next-intl";
import { useAdminPlugins } from "../useAdminPlugins";

// ==================== 插件配置（展示各插件独立设置） ====================
export function PluginConfigTab() {
  const t = useTranslations("Plugin");
  const { plugins } = useAdminPlugins();
  // 实际需要为每个插件渲染其配置表单，此处仅为框架
  return (
    <div className="space-y-4">
      <p className="text-sm text-base-content/60">
        {t("config_description")}
      </p>
      {plugins.map((plugin) => (
        <div key={plugin.slug} className="card card-sm bg-base-200 p-4">
          <h4 className="font-semibold">{plugin.name}</h4>
          <div className="text-xs text-base-content/50">
            {t("no_config")}
          </div>
        </div>
      ))}
    </div>
  );
}
