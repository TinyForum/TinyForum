"use client";

import { useEffect, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  CheckCircle2,
  Trash2,
  Power,
  PowerOff,
  Package,
  Upload,
} from "lucide-react";
import toast from "react-hot-toast";
import { useUpload } from "@/features/plugin/hooks/useUpload";
import { useAdminPlugins } from "../useAdminPlugins";
import { Pagination } from "@/shared/ui/common/Pagination";

// 假设 getUserPluginsList 现在返回 { list: FileInfo[], total: number }
// 如果实际 API 不支持，请先修改后端或增加一个获取总数的接口

export function MyPluginsTab() {
  const queryClient = useQueryClient();
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [uploadedFlag, setUploadedFlag] = useState(false);
  // 分页状态
  const [currentPage, setCurrentPage] = useState(1);

  const {
    isUploading,

    error: uploadError,
    resetError,
    getUserPluginsList, // 应返回 { list: FileInfo[], total: number }
    deleteFile,
  } = useUpload();
  const { upload, plugins, page, pageSize, total, pageTotal } =
    useAdminPlugins();
  // setCurrentPage(page);

  useEffect(() => {
    setCurrentPage(page);
  }, [page]); // 仅在 page 变化时更新
  // 查询我的插件列表（假设返回分页对象）
  const { isLoading } = useQuery({
    queryKey: ["myPlugins", page, pageSize],
    queryFn: () => getUserPluginsList({ page, page_size: pageSize }),
    // 如果 getUserPluginsList 仍返回数组，请使用下方 adapter
    // select: (data) => ({ list: data, total: data.length }) // 临时兼容
  });

  // 兼容两种情况：后端返回 { list, total } 或仅返回数组
  // const plugins = (data as any)?.list || data || [];
  // const total = plugins.length;
  // const totalPages = Math.ceil(total / pageSize);

  // 删除 mutation
  const deleteMutation = useMutation({
    mutationFn: (pluginId: string) => deleteFile("plugin", pluginId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["myPlugins"] });
      toast.success("插件已删除");
    },
    onError: () => toast.error("删除失败，请重试"),
  });

  // 上传文件处理
  const [file, setFile] = useState<File | null>(null);
  const handleUploadSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) {
      toast.error("请选择 ZIP 文件");
      return;
    }
    resetError();
    const result = await upload(file);
    if (result) {
      setUploadedFlag(true);
      setFile(null);
      setShowUploadForm(false);
      queryClient.invalidateQueries({ queryKey: ["myPlugins"] });
      setTimeout(() => setUploadedFlag(false), 3000);
    } else {
      toast.error(uploadError || "上传失败");
    }
  };

  // 成功上传后的提示
  if (uploadedFlag) {
    return (
      <div className="card bg-base-200 p-6 text-center">
        <CheckCircle2 className="w-12 h-12 text-success mx-auto mb-3" />
        <h3 className="text-lg font-semibold">上传成功！</h3>
        <p className="text-sm text-base-content/60">
          插件已提交审核，稍后可在列表中查看。
        </p>
        <button
          className="btn btn-sm btn-outline mt-4"
          onClick={() => setUploadedFlag(false)}
        >
          返回我的插件
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
          <h2 className="text-xl font-semibold">我的插件</h2>
          <span className="badge badge-ghost badge-sm">{total} 个插件</span>
        </div>
        <button
          onClick={() => setShowUploadForm(!showUploadForm)}
          className="btn btn-secondary btn-sm gap-2"
        >
          <Upload className="w-4 h-4" />
          {showUploadForm ? "取消上传" : "上传插件"}
        </button>
      </div>

      {/* 上传表单 */}
      {showUploadForm && (
        <div className="card bg-base-200 border border-base-300 p-4">
          <form onSubmit={handleUploadSubmit} className="space-y-4">
            <div>
              <label className="label-text">插件 ZIP 包 *</label>
              <input
                type="file"
                accept=".zip"
                required
                onChange={(e) => setFile(e.target.files?.[0] || null)}
                className="file-input file-input-bordered w-full"
              />
              <p className="text-xs text-base-content/50 mt-1">
                请打包符合规范的插件目录为 ZIP 文件，系统将自动读取插件信息。
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
                提交审核
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
            暂无插件，点击上方按钮上传你的第一个插件吧~
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
                          已启用
                        </span>
                      ) : (
                        <span className="badge badge-warning badge-sm">
                          未启用
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-base-content/70">
                      {plugin.description}
                    </p>
                    <div className="flex gap-3 text-xs text-base-content/50">
                      <span>作者: {plugin.author}</span>
                      <span>上传时间: {plugin.createdAt}</span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      className="btn btn-sm btn-ghost gap-1"
                      onClick={() => {
                        toast("启用/禁用功能开发中");
                      }}
                    >
                      {plugin.enabled ? (
                        <PowerOff className="w-4 h-4" />
                      ) : (
                        <Power className="w-4 h-4" />
                      )}
                    </button>
                    <button
                      className="btn btn-sm btn-ghost text-error gap-1"
                      onClick={() => deleteMutation.mutate(plugin.slug)}
                      disabled={deleteMutation.isPending}
                    >
                      <Trash2 className="w-4 h-4" />
                      删除
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* 使用通用分页组件 */}
          <Pagination
            totalPages={pageTotal}
            currentPage={currentPage}
            onPageChange={(newPage) => setCurrentPage(newPage)}
            totalItems={total}
            pageSize={pageSize}
          />
        </>
      )}
    </div>
  );
}
