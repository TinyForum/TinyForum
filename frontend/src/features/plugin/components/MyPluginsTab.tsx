"use client";

import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  CheckCircle2,
  Trash2,
  Power,
  PowerOff,
  Package,
  Upload,
} from "lucide-react";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { useAdminPlugins } from "../useAdminPlugins";
import { Pagination } from "@/shared/ui/common/Pagination";
import { pluginKeys } from "../hooks/usePluginKeys";

// 假设 getUserPluginsList 现在返回 { list: FileInfo[], total: number }
// 如果实际 API 不支持，请先修改后端或增加一个获取总数的接口

export function MyPluginsTab() {
  const t = useTranslations("Plugin");
  const queryClient = useQueryClient();
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [uploadedFlag, setUploadedFlag] = useState(false);

  const {
    isUploading,
    upload,
    plugins,
    page,
    pageSize,
    total,
    pageTotal,
    isLoading,
    setPage,
    handleDelete,
    handleToggle,
  } = useAdminPlugins();

  // 删除插件
  const handleDeletePlugin = (pluginSlug: string) => {
    handleDelete(pluginSlug);
    queryClient.invalidateQueries({ queryKey: pluginKeys.all });
    toast.success(t("plugin_deleted"));
  };

  // 上传文件处理
  const [file, setFile] = useState<File | null>(null);
  const handleUploadSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) {
      toast.error(t("select_zip"));
      return;
    }
    const result = await upload(file);
    if (result) {
      setUploadedFlag(true);
      setFile(null);
      setShowUploadForm(false);
      queryClient.invalidateQueries({ queryKey: pluginKeys.all });
      setTimeout(() => setUploadedFlag(false), 3000);
    } else {
      toast.error(t("upload_failed"));
    }
  };

  // 成功上传后的提示
  if (uploadedFlag) {
    return (
      <div className="card bg-base-200 p-6 text-center">
        <CheckCircle2 className="w-12 h-12 text-success mx-auto mb-3" />
        <h3 className="text-lg font-semibold">{t("upload_success")}</h3>
        <p className="text-sm text-base-content/60">
          {t("upload_success_desc")}
        </p>
        <button
          className="btn btn-sm btn-outline mt-4"
          onClick={() => setUploadedFlag(false)}
        >
          {t("back_to_my_plugins")}
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 头部 */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Package className="w-5 h-5 text-secondary" />
          <h2 className="text-xl font-semibold">{t("my_plugins")}</h2>
          <span className="badge badge-ghost badge-sm">{t("plugin_count", { count: total })}</span>
        </div>
        <button
          onClick={() => setShowUploadForm(!showUploadForm)}
          className="btn btn-secondary btn-sm gap-2"
        >
          <Upload className="w-4 h-4" />
          {showUploadForm ? t("cancel_upload") : t("upload_plugin")}
        </button>
      </div>

      {/* 上传表单 */}
      {showUploadForm && (
        <div className="card bg-base-200 border border-base-300 p-4">
          <form onSubmit={handleUploadSubmit} className="space-y-4">
            <div>
              <label className="label-text">{t("plugin_zip_label")}</label>
              <input
                type="file"
                accept=".zip"
                required
                onChange={(e) => setFile(e.target.files?.[0] || null)}
                className="file-input file-input-bordered w-full"
              />
              <p className="text-xs text-base-content/50 mt-1">
                {t("zip_instruction")}
              </p>
            </div>
            <div className="flex justify-end gap-2">
              <button
                type="submit"
                className="btn btn-secondary btn-sm"
                disabled={isUploading}
              >
                {isUploading && (
                  <span className="loading loading-spinner loading-xs" />
                )}
                {t("submit_review")}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* 插件列表 */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <span className="loading loading-spinner loading-md" />
        </div>
      ) : plugins.length === 0 ? (
        <div className="card bg-base-200 p-8 text-center">
          <Package className="w-12 h-12 text-base-content/30 mx-auto mb-3" />
          <p className="text-base-content/60">
            {t("no_plugins")}
          </p>
        </div>
      ) : (
        <>
          <div className="grid gap-4">
            {plugins.map((plugin) => (
              <div
                key={plugin.slug}
                className="card card-side bg-base-100 border border-base-300 shadow-sm"
              >
                <div className="card-body flex-row flex-wrap items-center justify-between p-4">
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center gap-2">
                      <h3 className="font-semibold text-base">{plugin.name}</h3>
                      <span className="badge badge-ghost badge-sm">
                        v{plugin.version}
                      </span>
                      {plugin.enabled ? (
                        <span className="badge badge-success badge-sm">
                          {t("enabled")}
                        </span>
                      ) : (
                        <span className="badge badge-warning badge-sm">
                          {t("disabled")}
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-base-content/70">
                      {plugin.description}
                    </p>
                    <div className="flex gap-3 text-xs text-base-content/50">
                      <span>{t("author_label", { author: plugin.author ?? "" })}</span>
                      <span>{t("upload_time_label", { time: plugin.createdAt ?? "" })}</span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      className="btn btn-sm btn-ghost gap-1"
                      onClick={() => handleToggle(plugin.slug, !plugin.enabled)}
                    >
                      {plugin.enabled ? (
                        <PowerOff className="w-4 h-4" />
                      ) : (
                        <Power className="w-4 h-4" />
                      )}
                    </button>
                    <button
                      className="btn btn-sm btn-ghost text-error gap-1"
                      onClick={() => handleDeletePlugin(plugin.slug)}
                    >
                      <Trash2 className="w-4 h-4" />
                      {t("delete")}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* 使用通用分页组件 */}
          <Pagination
            totalPages={pageTotal}
            currentPage={page}
            onPageChange={setPage}
            totalItems={total}
            pageSize={pageSize}
          />
        </>
      )}
    </div>
  );
}
