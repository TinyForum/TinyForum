"use client";

import { useState } from "react";
import {
  Puzzle,
  Plus,
  Trash2,
  Pencil,
  Power,
  PowerOff,
  RefreshCw,
  AlertCircle,
  CheckCircle2,
  Loader2,
  Code2,
  ChevronDown,
  ChevronUp,
} from "lucide-react";
import { Modal } from "./Modal";
import { PluginUploadForm } from "./PluginUploadForm";
import { useAdminPlugins } from "../hooks/useAdminPlugins";
import type { PluginMeta } from "@/shared/plugin/types";
import type { CreatePluginPayload } from "@/shared/api/modules/plugins";

export function PluginManager({ t }: { t: (key: string) => string }) {
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

  const [isInstallOpen, setIsInstallOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<PluginMeta | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<PluginMeta | null>(null);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const handleInstall = async (payload: CreatePluginPayload) => {
    await handleCreate(payload);
    setIsInstallOpen(false);
  };

  const handleEdit = async (payload: CreatePluginPayload) => {
    if (!editTarget) return;
    await handleUpdate({ id: editTarget.id, ...payload });
    setEditTarget(null);
  };

  const confirmDelete = () => {
    if (!deleteTarget) return;
    handleDelete(deleteTarget.id);
    setDeleteTarget(null);
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Puzzle className="w-5 h-5 text-primary" />
          <h2 className="text-lg font-semibold">Plugin Manager</h2>
          <span className="badge badge-neutral badge-sm">
            {total} installed
          </span>
        </div>
        <button
          onClick={() => setIsInstallOpen(true)}
          className="btn btn-primary btn-sm gap-2"
        >
          <Plus className="w-4 h-4" />
          Install Plugin
        </button>
      </div>

      {/* Info banner */}
      <div className="alert alert-info py-2 text-sm">
        <AlertCircle className="w-4 h-4 shrink-0" />
        <span>
          Plugins are loaded as external JavaScript bundles. Only install
          plugins from trusted sources. Changes take effect on next page load.
        </span>
      </div>

      {/* Table */}
      <div className="card bg-base-100 border border-base-300 shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center items-center h-48">
            <Loader2 className="w-6 h-6 animate-spin text-primary" />
          </div>
        ) : plugins.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-48 gap-3 text-base-content/40">
            <Puzzle className="w-10 h-10" />
            <p className="text-sm">No plugins installed yet</p>
            <button
              onClick={() => setIsInstallOpen(true)}
              className="btn btn-ghost btn-xs gap-1"
            >
              <Plus className="w-3 h-3" /> Install your first plugin
            </button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="table table-sm">
              <thead>
                <tr className="bg-base-200 text-base-content/60 text-xs uppercase tracking-wide">
                  <th>Plugin</th>
                  <th>Version</th>
                  <th>Author</th>
                  <th>Slots</th>
                  <th>Status</th>
                  <th className="text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {plugins.map((plugin) => (
                  <>
                    <tr
                      key={plugin.id}
                      className="hover:bg-base-50 transition-colors"
                    >
                      {/* Plugin name + description toggle */}
                      <td>
                        <button
                          className="flex items-center gap-2 text-left group"
                          onClick={() =>
                            setExpandedId(
                              expandedId === plugin.id ? null : plugin.id,
                            )
                          }
                        >
                          <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
                            <Code2 className="w-4 h-4 text-primary" />
                          </div>
                          <div>
                            <div className="font-medium text-sm flex items-center gap-1">
                              {plugin.name}
                              {expandedId === plugin.id ? (
                                <ChevronUp className="w-3 h-3 text-base-content/40" />
                              ) : (
                                <ChevronDown className="w-3 h-3 text-base-content/40 group-hover:text-base-content/60" />
                              )}
                            </div>
                          </div>
                        </button>
                      </td>

                      {/* Version */}
                      <td>
                        <span className="badge badge-ghost badge-sm font-mono">
                          v{plugin.version}
                        </span>
                      </td>

                      {/* Author */}
                      <td className="text-sm text-base-content/60">
                        {plugin.author}
                      </td>

                      {/* Slots */}
                      <td>
                        <div className="flex flex-wrap gap-1 max-w-[160px]">
                          {(plugin.slots ?? []).length === 0 ? (
                            <span className="text-xs text-base-content/30">
                              —
                            </span>
                          ) : (
                            (plugin.slots ?? []).slice(0, 2).map((s) => (
                              <span
                                key={s}
                                className="badge badge-outline badge-xs"
                              >
                                {s}
                              </span>
                            ))
                          )}
                          {(plugin.slots ?? []).length > 2 && (
                            <span className="badge badge-ghost badge-xs">
                              +{(plugin.slots ?? []).length - 2}
                            </span>
                          )}
                        </div>
                      </td>

                      {/* Status */}
                      <td>
                        <StatusBadge enabled={plugin.enabled} />
                      </td>

                      {/* Actions */}
                      <td className="text-right">
                        <div className="flex items-center justify-end gap-1">
                          {/* Toggle */}
                          <button
                            onClick={() =>
                              handleToggle(plugin.id, !plugin.enabled)
                            }
                            className={`btn btn-xs btn-ghost gap-1 ${
                              plugin.enabled
                                ? "text-warning hover:text-warning"
                                : "text-success hover:text-success"
                            }`}
                            title={plugin.enabled ? "Disable" : "Enable"}
                          >
                            {plugin.enabled ? (
                              <PowerOff className="w-3.5 h-3.5" />
                            ) : (
                              <Power className="w-3.5 h-3.5" />
                            )}
                          </button>

                          {/* Edit */}
                          <button
                            onClick={() => setEditTarget(plugin)}
                            className="btn btn-xs btn-ghost"
                            title="Edit"
                          >
                            <Pencil className="w-3.5 h-3.5" />
                          </button>

                          {/* Delete */}
                          <button
                            onClick={() => setDeleteTarget(plugin)}
                            className="btn btn-xs btn-ghost text-error hover:text-error"
                            title="Remove"
                          >
                            <Trash2 className="w-3.5 h-3.5" />
                          </button>
                        </div>
                      </td>
                    </tr>

                    {/* Expanded detail row */}
                    {expandedId === plugin.id && (
                      <tr
                        key={`${plugin.id}-detail`}
                        className="bg-base-200/50"
                      >
                        <td colSpan={6} className="py-3 px-4">
                          <PluginDetailPanel plugin={plugin} />
                        </td>
                      </tr>
                    )}
                  </>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-1">
          <button
            className="btn btn-xs btn-ghost"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            «
          </button>
          {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
            <button
              key={p}
              className={`btn btn-xs ${p === page ? "btn-primary" : "btn-ghost"}`}
              onClick={() => setPage(p)}
            >
              {p}
            </button>
          ))}
          <button
            className="btn btn-xs btn-ghost"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            »
          </button>
        </div>
      )}

      {/* ── Install Modal ── */}
      <Modal
        isOpen={isInstallOpen}
        onClose={() => setIsInstallOpen(false)}
        title="Install New Plugin"
      >
        <PluginUploadForm
          onSubmit={handleInstall}
          onCancel={() => setIsInstallOpen(false)}
          isLoading={isCreating}
        />
      </Modal>

      {/* ── Edit Modal ── */}
      <Modal
        isOpen={!!editTarget}
        onClose={() => setEditTarget(null)}
        title={`Edit Plugin: ${editTarget?.name}`}
      >
        {editTarget && (
          <PluginUploadForm
            initial={editTarget}
            onSubmit={handleEdit}
            onCancel={() => setEditTarget(null)}
            isLoading={isUpdating}
          />
        )}
      </Modal>

      {/* ── Delete Confirm Modal ── */}
      <Modal
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        title="Remove Plugin"
      >
        <div className="space-y-4">
          <div className="flex items-start gap-3 p-3 bg-error/10 rounded-lg border border-error/20">
            <AlertCircle className="w-5 h-5 text-error shrink-0 mt-0.5" />
            <div>
              <p className="font-medium text-sm">
                Are you sure you want to remove{" "}
                <span className="text-error">{deleteTarget?.name}</span>?
              </p>
              <p className="text-xs text-base-content/60 mt-1">
                This will permanently delete the plugin configuration. The
                script file itself will not be affected.
              </p>
            </div>
          </div>
          <div className="flex justify-end gap-2">
            <button
              onClick={() => setDeleteTarget(null)}
              className="btn btn-ghost btn-sm"
            >
              Cancel
            </button>
            <button
              onClick={confirmDelete}
              className="btn btn-error btn-sm gap-2"
            >
              <Trash2 className="w-4 h-4" />
              Remove Plugin
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

// ── Sub-components ────────────────────────────────────────────────────────────

function StatusBadge({ enabled }: { enabled: boolean }) {
  if (enabled) {
    return (
      <span className="badge badge-success badge-sm gap-1">
        <CheckCircle2 className="w-3 h-3" /> Active
      </span>
    );
  }
  return (
    <span className="badge badge-ghost badge-sm gap-1 text-base-content/40">
      <PowerOff className="w-3 h-3" /> Disabled
    </span>
  );
}

function PluginDetailPanel({ plugin }: { plugin: PluginMeta }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
      <div className="space-y-2">
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            Description
          </span>
          <p className="mt-0.5 text-base-content/80">
            {plugin.description || (
              <span className="italic text-base-content/30">
                No description
              </span>
            )}
          </p>
        </div>
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            All Slots
          </span>
          <div className="flex flex-wrap gap-1 mt-1">
            {(plugin.slots ?? []).length === 0 ? (
              <span className="text-xs text-base-content/30 italic">
                None declared
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
      <div className="space-y-2">
        <div>
          <span className="text-xs text-base-content/40 uppercase tracking-wide">
            Script URL
          </span>
          <p className="mt-0.5 font-mono text-xs break-all text-base-content/60 bg-base-300 p-2 rounded">
            {plugin.scriptUrl}
          </p>
        </div>
        {plugin.updatedAt && (
          <div>
            <span className="text-xs text-base-content/40 uppercase tracking-wide">
              Last Updated
            </span>
            <p className="mt-0.5 text-base-content/60">
              {new Date(plugin.updatedAt).toLocaleString()}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
