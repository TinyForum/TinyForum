import { useAdminPlugins } from "../useAdminPlugins";

// ==================== 插件配置（展示各插件独立设置） ====================
export function PluginConfigTab() {
  const { plugins } = useAdminPlugins();
  // 实际需要为每个插件渲染其配置表单，此处仅为框架
  return (
    <div className="space-y-4">
      <p className="text-sm text-base-content/60">
        每个插件的独立配置项（需插件支持）
      </p>
      {plugins.map((plugin) => (
        <div key={plugin.slug} className="card card-sm bg-base-200 p-4">
          <h4 className="font-semibold">{plugin.name}</h4>
          <div className="text-xs text-base-content/50">
            该插件暂无配置项或未暴露设置界面。
          </div>
        </div>
      ))}
    </div>
  );
}
