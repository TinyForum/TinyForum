import { useState } from "react";
import { useAdminPlugins } from "../useAdminPlugins";
import { CreatePluginPayload } from "@/shared/api/modules/plugin/plugins";
import Image from "next/image";

export function PluginMarketTab() {
  const {
    plugins,
    total,
    isLoading,
    page,
    pageSize,
    setPage,
    upload,
    handleToggle,
    handleDelete,
    handleCreate,
    isCreating,
    isUpdating,
    isDeleting,
  } = useAdminPlugins();

  const [searchTerm, setSearchTerm] = useState("");
  const [uploading, setUploading] = useState(false);

  // 前端过滤当前页插件（根据名称或描述）
  const filteredPlugins = plugins.filter((plugin) => {
    const term = searchTerm.toLowerCase();
    return (
      plugin.name?.toLowerCase().includes(term) ||
      plugin.description?.toLowerCase().includes(term) ||
      plugin.version?.toLowerCase().includes(term)
    );
  });

  // 处理文件上传安装
  const handleFileUpload = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setUploading(true);
    try {
      // 1. 上传插件包，后端解析返回插件元数据
      const pluginInfo = await upload(file);
      // 确保 pluginInfo 是一个有效对象
      if (pluginInfo && typeof pluginInfo === "object") {
        // 2. 调用创建接口完成安装（默认启用）
        await handleCreate({
          ...(pluginInfo as CreatePluginPayload),
          enabled: true,
        });
      } else {
        throw new Error("Invalid plugin metadata");
      }
      // 清空 input 以便重复上传同一文件仍可触发
      event.target.value = "";
    } catch (error) {
      console.error("Install failed:", error);
    } finally {
      setUploading(false);
    }
  };

  // 分页控制
  const totalPages = Math.ceil(total / pageSize);
  const canPrev = page > 1;
  const canNext = page < totalPages;

  return (
    <div className="space-y-4">
      {/* 头部：标题 + 搜索 + 安装按钮 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h3 className="text-base font-medium">官方插件市场</h3>
        <div className="flex items-center gap-2">
          <input
            type="text"
            placeholder="搜索插件（当前页）..."
            className="input input-bordered input-sm w-64"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
          <label className="btn btn-sm btn-primary">
            {uploading ? (
              <span className="loading loading-spinner loading-xs" />
            ) : (
              "安装插件"
            )}
            <input
              type="file"
              accept=".zip"
              className="hidden"
              onChange={handleFileUpload}
              disabled={uploading || isCreating}
            />
          </label>
        </div>
      </div>

      {/* 插件列表 */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <span className="loading loading-spinner loading-md" />
        </div>
      ) : filteredPlugins.length === 0 ? (
        <div className="py-12 text-center text-base-content/70">
          {plugins.length === 0
            ? "暂无已安装插件，点击上方按钮安装第一个插件"
            : "没有找到匹配的插件"}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filteredPlugins.map((plugin) => (
            <div
              key={plugin.slug}
              className="card card-bordered bg-base-100 shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="card-body p-4">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      {plugin.iconUrl && (
                        <Image
                          src={plugin.iconUrl}
                          alt=""
                          className="w-6 h-6"
                        />
                      )}
                      <h4 className="card-title text-base">{plugin.name}</h4>
                      <span className="badge badge-sm badge-ghost">
                        v{plugin.version}
                      </span>
                    </div>
                    <p className="text-sm text-base-content/70 mt-1 line-clamp-2">
                      {plugin.description || "暂无描述"}
                    </p>
                    <div className="mt-2 text-xs text-base-content/50">
                      作者：{plugin.author || "未知"}
                    </div>
                  </div>
                  <div className="flex items-center gap-1">
                    {/* 启用/禁用开关 */}
                    <button
                      className={`btn btn-xs btn-square ${
                        plugin.enabled ? "btn-success" : "btn-outline"
                      }`}
                      onClick={() => handleToggle(plugin.slug, !plugin.enabled)}
                      disabled={isUpdating || isDeleting}
                      title={plugin.enabled ? "禁用" : "启用"}
                    >
                      {plugin.enabled ? (
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M5 13l4 4L19 7"
                          />
                        </svg>
                      ) : (
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M6 18L18 6M6 6l12 12"
                          />
                        </svg>
                      )}
                    </button>
                    {/* 删除按钮 */}
                    <button
                      className="btn btn-xs btn-square btn-outline btn-error"
                      onClick={() => handleDelete(plugin.slug)}
                      disabled={isDeleting}
                      title="卸载"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 分页组件 */}
      {totalPages > 1 && (
        <div className="flex justify-center pt-2">
          <div className="join">
            <button
              className="join-item btn btn-sm"
              onClick={() => setPage(page - 1)}
              disabled={!canPrev || isLoading}
            >
              «
            </button>
            <span className="join-item btn btn-sm btn-disabled">
              第 {page} / {totalPages} 页
            </span>
            <button
              className="join-item btn btn-sm"
              onClick={() => setPage(page + 1)}
              disabled={!canNext || isLoading}
            >
              »
            </button>
          </div>
        </div>
      )}

      {/* 底部信息 */}
      <div className="text-xs text-center text-base-content/50">
        共 {total} 个插件，当前页显示 {filteredPlugins.length} 个
      </div>
    </div>
  );
}
