import { Modal } from "@/features/admin/components/Modal";
import { EmptyState } from "@/features/system/components/DetailRow";
import { PluginDetail } from "@/features/system/components/PluginDetail";
import { PluginRow } from "@/features/system/components/PluginRow";
import { PluginMeta } from "@/shared/type/plugin.type";
import { Puzzle, Plus, Search, X, AlertCircle, Trash2 } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";
import { useAdminPlugins } from "../useAdminPlugins";
import { PluginUploadForm } from "./PluginUploadForm";

// ==================== 子组件：插件管理（含自定义安装入口） ====================
export function PluginManagementTab() {
  const {
    plugins,
    total,
    isLoading,
    page,
    pageSize,
    setPage,
    handleToggle,
    handleDelete,
    handleCreate,
    handleUpdate,
    isCreating,
    isUpdating,
  } = useAdminPlugins();

  const [keyword, setKeyword] = useState("");
  const [isInstallOpen, setIsInstallOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<PluginMeta | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<PluginMeta | null>(null);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const t = useTranslations("Plugin");

  const filtered = plugins.filter((plugin) => {
    const matchName = plugin.name
      ?.toLowerCase()
      .includes(keyword.toLowerCase());
    const matchDesc = plugin.description
      ?.toLowerCase()
      .includes(keyword.toLowerCase());
    const matchAuthor = plugin.author
      ?.toLowerCase()
      .includes(keyword.toLowerCase());
    return matchName || matchDesc || matchAuthor;
  });

  const activeCount = plugins.filter((p) => p.enabled).length;

  return (
    <div className="space-y-4">
      {/* 头部 + 自定义安装按钮 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Puzzle className="w-5 h-5 text-secondary" />
          <h2 className="text-base font-semibold">{t("installed_plugins")}</h2>
          <span className="badge badge-ghost badge-sm">
            {activeCount} 启用 / {total} 已安装
          </span>
        </div>
        <button
          onClick={() => setIsInstallOpen(true)}
          className="btn btn-secondary btn-sm gap-2"
        >
          <Plus className="w-4 h-4" /> 自定义安装
        </button>
      </div>

      {/* 搜索 */}
      <label className="input input-bordered input-sm flex items-center gap-2 w-full max-w-xs">
        <Search className="w-3.5 h-3.5" />
        <input
          type="text"
          placeholder="搜索插件..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          className="grow"
        />
        {keyword && (
          <button onClick={() => setKeyword("")}>
            <X className="w-3.5 h-3.5" />
          </button>
        )}
      </label>

      <div className="alert py-2 text-sm bg-base-200 border-base-300">
        <AlertCircle className="w-4 h-4 shrink-0" />
        <span className="text-base-content/60">
          插件为外部 JS 文件，请仅安装可信来源的插件。
        </span>
      </div>

      <div className="card bg-base-100 border border-base-300 shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center items-center h-40">
            <span className="loading loading-spinner loading-md text-secondary" />
          </div>
        ) : filtered.length === 0 ? (
          <EmptyState
            keyword={keyword}
            onInstall={() => setIsInstallOpen(true)}
          />
        ) : (
          <div className="divide-y divide-base-200">
            {filtered.map((plugin) => (
              <div key={plugin.slug}>
                <PluginRow
                  plugin={plugin}
                  expanded={expandedId === plugin.slug}
                  onExpand={() =>
                    setExpandedId(
                      expandedId === plugin.slug ? null : plugin.slug,
                    )
                  }
                  onToggle={() => handleToggle(plugin.slug, !plugin.enabled)}
                  onEdit={() => setEditTarget(plugin)}
                  onDelete={() => setDeleteTarget(plugin)}
                />
                {expandedId === plugin.slug && <PluginDetail plugin={plugin} />}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 分页 */}
      {Math.ceil(total / pageSize) > 1 && (
        <div className="flex justify-center gap-1">
          <button
            className="btn btn-xs btn-ghost"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            «
          </button>
          {Array.from(
            { length: Math.ceil(total / pageSize) },
            (_, i) => i + 1,
          ).map((p) => (
            <button
              key={p}
              className={`btn btn-xs ${p === page ? "btn-secondary" : "btn-ghost"}`}
              onClick={() => setPage(p)}
            >
              {p}
            </button>
          ))}
          <button
            className="btn btn-xs btn-ghost"
            disabled={page === Math.ceil(total / pageSize)}
            onClick={() => setPage((p) => p + 1)}
          >
            »
          </button>
        </div>
      )}

      {/* 安装/编辑/删除 Modal 保持不变 */}
      <Modal
        isOpen={isInstallOpen}
        onClose={() => setIsInstallOpen(false)}
        title="自定义安装插件"
      >
        <PluginUploadForm
          onSubmit={async (p) => {
            await handleCreate(p);
            setIsInstallOpen(false);
          }}
          onCancel={() => setIsInstallOpen(false)}
          isLoading={isCreating}
        />
      </Modal>
      <Modal
        isOpen={!!editTarget}
        onClose={() => setEditTarget(null)}
        title={`编辑插件：${editTarget?.name}`}
      >
        {editTarget && (
          <PluginUploadForm
            initial={editTarget}
            onSubmit={async (p) => {
              await handleUpdate({ id: editTarget.slug, ...p });
              setEditTarget(null);
            }}
            onCancel={() => setEditTarget(null)}
            isLoading={isUpdating}
          />
        )}
      </Modal>
      <Modal
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        title="确认移除插件"
      >
        <div className="space-y-4">
          <div className="flex gap-3 p-3 rounded-lg bg-error/10 border border-error/20">
            <AlertCircle className="w-5 h-5 text-error shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-medium">
                确定要移除{" "}
                <span className="text-error">{deleteTarget?.name}</span> 吗？
              </p>
              <p className="text-xs text-base-content/50 mt-1">
                此操作将删除插件配置，脚本文件不受影响。
              </p>
            </div>
          </div>
          <div className="flex justify-end gap-2">
            <button
              onClick={() => setDeleteTarget(null)}
              className="btn btn-ghost btn-sm"
            >
              取消
            </button>
            <button
              onClick={() => {
                if (deleteTarget) {
                  handleDelete(deleteTarget.slug);
                  setDeleteTarget(null);
                }
              }}
              className="btn btn-error btn-sm gap-2"
            >
              <Trash2 className="w-4 h-4" /> 确认移除
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
