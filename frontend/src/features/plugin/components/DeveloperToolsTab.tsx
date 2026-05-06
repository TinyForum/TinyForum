import { Wrench } from "lucide-react";

// ==================== 开发者工具 ====================
export function DeveloperToolsTab() {
  return (
    <div className="space-y-4">
      <div className="alert alert-info shadow-lg">
        <div>
          <Wrench className="w-5 h-5" />
          <span>开发者工具：钩子调试、日志查看、插件打包上传</span>
        </div>
      </div>
      <details className="collapse collapse-arrow bg-base-200">
        <summary className="collapse-title text-sm font-medium">
          钩子列表
        </summary>
        <div className="collapse-content">
          <pre className="text-xs">post:create, user:login, ...</pre>
        </div>
      </details>
      <button className="btn btn-outline btn-sm">导出调试日志</button>
    </div>
  );
}
