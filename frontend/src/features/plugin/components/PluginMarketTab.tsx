import { useState } from "react";

// ==================== 插件市场 ====================
export function PluginMarketTab() {
  // mock 数据，实际对接 API
  const [loading] = useState(false);
  // useEffect 等省略，仅示例
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-base font-medium">官方插件市场</h3>
        <input
          type="text"
          placeholder="搜索插件..."
          className="input input-bordered input-sm w-64"
        />
      </div>
      {loading ? (
        <div className="flex justify-center py-12">
          <span className="loading loading-spinner loading-md" />
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* 插件卡片 */}
        </div>
      )}
    </div>
  );
}
