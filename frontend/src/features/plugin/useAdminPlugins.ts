import { useState, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import {
  pluginApi,
  CreatePluginPayload,
  UpdatePluginPayload,
} from "@/shared/api/modules/plugin/plugins";
import { PluginMeta } from "@/shared/api/types/plugin.model";
import { pluginKeys } from "./hooks/usePluginKeys";

const PAGE_SIZE = 10;

export function useAdminPlugins() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [isUploading, setIsUploading] = useState(false);

  // ── 分页列表（管理端） ───────────────────────────────────────────────────────
  const { data: pageRaw, isLoading } = useQuery({
    queryKey: pluginKeys.list(page),
    queryFn: async () => {
      const res = await pluginApi.list({ page, page_size: PAGE_SIZE });
      return res.data?.data ?? null;
    },
  });

  const plugins: PluginMeta[] = pageRaw?.list ?? [];
  const total: number = pageRaw?.total ?? 0;

  // ── 已启用插件列表（仅展示"已启用数量"等，加载由 PluginContext 完成） ────────
  const { data: enabledPluginsRaw } = useQuery({
    queryKey: pluginKeys.enabled(),
    queryFn: async () => {
      const res = await pluginApi.listEnabled();
      return res.data?.data ?? null;
    },
  });
  const enabledPlugins: PluginMeta[] = enabledPluginsRaw?.list ?? [];
  const pageTotal = Math.ceil(total / PAGE_SIZE);

  // ── 通用成功回调 ──────────────────────────────────────────────────────────
  const invalidateAndToast = (msg: string) => {
    queryClient.invalidateQueries({ queryKey: pluginKeys.all });
    toast.success(msg);
  };

  // ── 文件上传 ──────────────────────────────────────────────────────────────
  const upload = async (file: File) => {
    setIsUploading(true);
    try {
      const res = await pluginApi.upload(file);
      return res.data.data;
    } finally {
      setIsUploading(false);
    }
  };

  // ── Mutations ─────────────────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: (payload: CreatePluginPayload) => pluginApi.create(payload),
    onSuccess: () => invalidateAndToast("Plugin installed successfully"),
    onError: () => toast.error("Failed to install plugin"),
  });

  const updateMutation = useMutation({
    mutationFn: (payload: UpdatePluginPayload) => pluginApi.update(payload),
    onSuccess: () => invalidateAndToast("Plugin updated"),
    onError: () => toast.error("Failed to update plugin"),
  });

  const toggleMutation = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) =>
      pluginApi.toggle(id, enabled),
    onSuccess: (_, { enabled }) =>
      invalidateAndToast(enabled ? "Plugin enabled" : "Plugin disabled"),
    onError: () => toast.error("Failed to toggle plugin"),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => pluginApi.delete(id),
    onSuccess: () => invalidateAndToast("Plugin removed"),
    onError: () => toast.error("Failed to remove plugin"),
  });

  // ── 稳定引用的处理函数 ────────────────────────────────────────────────────
  const handleToggle = useCallback(
    (id: string, enabled: boolean) => toggleMutation.mutate({ id, enabled }),
    [toggleMutation],
  );
  const handleDelete = useCallback(
    (id: string) => deleteMutation.mutate(id),
    [deleteMutation],
  );
  const handleCreate = useCallback(
    (payload: CreatePluginPayload) => createMutation.mutateAsync(payload),
    [createMutation],
  );
  const handleUpdate = useCallback(
    (payload: UpdatePluginPayload) => updateMutation.mutateAsync(payload),
    [updateMutation],
  );

  return {
    plugins,
    total,
    isLoading,
    page,
    pageSize: PAGE_SIZE,
    pageTotal,
    enabledPlugins, // PluginMeta[]，正确的数组类型
    setPage,
    upload,
    isUploading,
    handleToggle,
    handleDelete,
    handleCreate,
    handleUpdate,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}
