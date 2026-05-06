"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  CheckCircle2,
  Trash2,
  Power,
  PowerOff,
  AlertCircle,
  Package,
  Upload,
} from "lucide-react";
import { apiClient } from "@/shared/api";
import toast from "react-hot-toast";
import { useUpload } from "@/features/plugin/hooks/useUpload";

// 插件元信息（后端返回的结构，根据实际情况调整）
interface UserPlugin {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  fileId: string;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

// 获取我的插件列表
const fetchMyPlugins = async (): Promise<UserPlugin[]> => {
  const res = await apiClient.get<{ data: UserPlugin[] }>(
    "/attachments/plugin/user/me",
  );
  return res.data.data;
};

// 删除插件
const deletePlugin = async (pluginId: string): Promise<void> => {
  await apiClient.delete(`/user/plugins/${pluginId}`);
};

export function MyPluginsTab() {
  const queryClient = useQueryClient();
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [uploadedFlag, setUploadedFlag] = useState(false);
  const { uploadPluginFile, isUploading, error, resetError } = useUpload();

  // 查询我的插件列表
  const {
    data: plugins = [],
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ["myPlugins"],
    queryFn: fetchMyPlugins,
  });

  // 删除 mutation
  const deleteMutation = useMutation({
    mutationFn: deletePlugin,
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
    const result = await uploadPluginFile(file);
    if (result) {
      setUploadedFlag(true);
      setFile(null);
      setShowUploadForm(false);
      queryClient.invalidateQueries({ queryKey: ["myPlugins"] });
      setTimeout(() => setUploadedFlag(false), 3000);
    } else {
      toast.error(error || "上传失败");
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
          <span className="badge badge-ghost badge-sm">
            {plugins.length} 个插件
          </span>
        </div>
        <button
          onClick={() => setShowUploadForm(!showUploadForm)}
          className="btn btn-secondary btn-sm gap-2"
        >
          <Upload className="w-4 h-4" />
          {showUploadForm ? "取消上传" : "上传插件"}
        </button>
      </div>

      {/* 上传表单（可折叠） */}
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
        <div className="grid gap-4">
          {plugins.map((plugin) => (
            <div
              key={plugin.id}
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
                    <span>
                      上传时间:{" "}
                      {new Date(plugin.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                </div>
                <div className="flex gap-2">
                  {/* 启用/禁用按钮（假设API支持，可根据实际调整） */}
                  <button
                    className="btn btn-sm btn-ghost gap-1"
                    onClick={() => {
                      // 这里调用启用/禁用 API，示例略
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
                    onClick={() => deleteMutation.mutate(plugin.id)}
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
      )}
    </div>
  );
}
