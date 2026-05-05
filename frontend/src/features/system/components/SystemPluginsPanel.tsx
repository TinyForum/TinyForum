"use client";

import { useState } from "react";
import { Switch } from "@headlessui/react";
import {
  Puzzle,
  Plus,
  Trash2,
  Pencil,
  Power,
  PowerOff,
  AlertCircle,
  CheckCircle2,
  Code2,
  ChevronDown,
  ChevronUp,
  Search,
  X,
} from "lucide-react";
import { Modal } from "@/features/admin/components/Modal";
import { PluginUploadForm } from "@/features/plugin/components/PluginUploadForm";
import { useAdminPlugins } from "@/features/plugin/useAdminPlugins";
import { PluginMeta } from "@/shared/type/plugin.type";

/**
 * SystemPluginsPanel
 * 在 system 页面中嵌入完整的插件管理 UI。
 * 直接复用 useAdminPlugins hook + PluginUploadForm + Modal，
 * 额外增加搜索框和本地 enabled toggle 以便无后端时也能演示。
 */
export function SystemPluginsPanel() {
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

  const filtered = plugins.filter(
    (p) =>
      p.name.toLowerCase().includes(keyword.toLowerCase()) ||
      p.description?.toLowerCase().includes(keyword.toLowerCase()) ||
      p.author?.toLowerCase().includes(keyword.toLowerCase()),
  );

  const activeCount = plugins.filter((p) => p.enabled).length;

  return (
    <div className="space-y-4">
      {/* ── Top bar ── */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Puzzle className="w-5 h-5 text-secondary" />
          <h2 className="text-base font-semibold">插件管理</h2>
          <span className="badge badge-ghost badge-sm">
            {activeCount} 启用 / {total} 已安装
          </span>
        </div>
        <button
          onClick={() => setIsInstallOpen(true)}
          className="btn btn-secondary btn-sm gap-2"
        >
          <Plus className="w-4 h-4" /> 安装插件
        </button>
      </div>

      {/* Search */}
      <label className="input input-bordered input-sm flex items-center gap-2 w-full max-w-xs">
        <Search className="w-3.5 h-3.5 text-base-content/40 shrink-0" />
        <input
          type="text"
          placeholder="搜索插件..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          className="grow"
        />
        {keyword && (
          <button onClick={() => setKeyword("")}>
            <X className="w-3.5 h-3.5 text-base-content/30" />
          </button>
        )}
      </label>

      {/* Info banner */}
      <div className="alert py-2 text-sm bg-base-200 border-base-300">
        <AlertCircle className="w-4 h-4 shrink-0 text-base-content/40" />
        <span className="text-base-content/60">
          插件为外部 JS 文件，请仅安装可信来源的插件。变更在下次页面加载时生效。
        </span>
      </div>

      {/* Plugin list */}
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
              <div key={plugin.id}>
                <PluginRow
                  plugin={plugin}
                  expanded={expandedId === plugin.id}
                  onExpand={() =>
                    setExpandedId(expandedId === plugin.id ? null : plugin.id)
                  }
                  onToggle={() => handleToggle(plugin.id, !plugin.enabled)}
                  onEdit={() => setEditTarget(plugin)}
                  onDelete={() => setDeleteTarget(plugin)}
                />
                {expandedId === plugin.id && <PluginDetail plugin={plugin} />}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Pagination */}
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

      {/* Modals */}
      <Modal
        isOpen={isInstallOpen}
        onClose={() => setIsInstallOpen(false)}
        title="安装新插件"
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
              await handleUpdate({ id: editTarget.id, ...p });
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
                  handleDelete(deleteTarget.id);
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

// ── Sub-components ─────────────────────────────────────────────────────────────

function PluginRow({
  plugin,
  expanded,
  onExpand,
  onToggle,
  onEdit,
  onDelete,
}: {
  plugin: PluginMeta;
  expanded: boolean;
  onExpand: () => void;
  onToggle: () => void;
  onEdit: () => void;
  onDelete: () => void;
}) {
  return (
    <div className="flex items-center gap-3 px-4 py-3 hover:bg-base-50 transition-colors">
      {/* Icon */}
      <div className="w-9 h-9 rounded-lg bg-secondary/10 flex items-center justify-center shrink-0">
        <Code2 className="w-4 h-4 text-secondary" />
      </div>

      {/* Name + expand */}
      <button className="flex-1 text-left min-w-0 group" onClick={onExpand}>
        <div className="flex items-center gap-1.5">
          <span className="text-sm font-medium truncate">{plugin.name}</span>
          <span className="badge badge-ghost badge-xs font-mono shrink-0">
            v{plugin.version}
          </span>
          {expanded ? (
            <ChevronUp className="w-3.5 h-3.5 text-base-content/30 shrink-0" />
          ) : (
            <ChevronDown className="w-3.5 h-3.5 text-base-content/20 group-hover:text-base-content/40 shrink-0" />
          )}
        </div>
        <p className="text-xs text-base-content/40 truncate">{plugin.author}</p>
      </button>

      {/* Slots preview */}
      <div className="hidden sm:flex flex-wrap gap-1 max-w-[140px]">
        {(plugin.slots ?? []).slice(0, 2).map((s) => (
          <span key={s} className="badge badge-outline badge-xs">
            {s}
          </span>
        ))}
        {(plugin.slots ?? []).length > 2 && (
          <span className="badge badge-ghost badge-xs">
            +{(plugin.slots ?? []).length - 2}
          </span>
        )}
      </div>

      {/* Status */}
      <span
        className={`badge badge-sm shrink-0 gap-1 ${plugin.enabled ? "badge-success" : "badge-ghost"}`}
      >
        {plugin.enabled ? (
          <CheckCircle2 className="w-3 h-3" />
        ) : (
          <PowerOff className="w-3 h-3" />
        )}
        {plugin.enabled ? "启用" : "禁用"}
      </span>

      {/* Toggle */}
      <Switch
        checked={plugin.enabled}
        onChange={onToggle}
        className={`${plugin.enabled ? "bg-secondary" : "bg-base-300"} relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-secondary/30`}
      >
        <span
          className={`${plugin.enabled ? "translate-x-5" : "translate-x-1"} inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform`}
        />
      </Switch>

      {/* Action buttons */}
      <div className="flex items-center gap-0.5 shrink-0">
        <button
          onClick={onEdit}
          className="btn btn-ghost btn-xs btn-square"
          title="编辑"
        >
          <Pencil className="w-3.5 h-3.5" />
        </button>
        <button
          onClick={onDelete}
          className="btn btn-ghost btn-xs btn-square text-error hover:text-error"
          title="移除"
        >
          <Trash2 className="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  );
}

function PluginDetail({ plugin }: { plugin: PluginMeta }) {
  return (
    <div className="px-4 pb-4 pt-0 bg-base-200/40 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm border-t border-base-200">
      <div className="space-y-2 pt-3">
        <DetailRow label="描述" value={plugin.description || "—"} />
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            注入插槽
          </span>
          <div className="flex flex-wrap gap-1 mt-1">
            {(plugin.slots ?? []).length === 0 ? (
              <span className="text-xs text-base-content/30 italic">
                未声明
              </span>
            ) : (
              (plugin.slots ?? []).map((s) => (
                <span key={s} className="badge badge-outline badge-sm">
                  {s}
                </span>
              ))
            )}
          </div>
        </div>
      </div>
      <div className="space-y-2 pt-3">
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            Script URL
          </span>
          <p className="mt-0.5 font-mono text-xs break-all text-base-content/50 bg-base-300 p-2 rounded">
            {plugin.scriptUrl}
          </p>
        </div>
        {plugin.updatedAt && (
          <DetailRow
            label="最后更新"
            value={new Date(plugin.updatedAt).toLocaleString("zh-CN")}
          />
        )}
      </div>
    </div>
  );
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <span className="text-xs text-base-content/40 uppercase tracking-wide">
        {label}
      </span>
      <p className="mt-0.5 text-base-content/70 text-xs">{value}</p>
    </div>
  );
}

function EmptyState({
  keyword,
  onInstall,
}: {
  keyword: string;
  onInstall: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center h-40 gap-3 text-base-content/40">
      <Puzzle className="w-8 h-8" />
      <p className="text-sm">
        {keyword ? `没有找到"${keyword}"相关插件` : "尚未安装任何插件"}
      </p>
      {!keyword && (
        <button onClick={onInstall} className="btn btn-ghost btn-xs gap-1">
          <Plus className="w-3 h-3" /> 安装第一个插件
        </button>
      )}
    </div>
  );
}
