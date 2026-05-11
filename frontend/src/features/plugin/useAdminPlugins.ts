import { useState, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import {
  pluginApi,
  CreatePluginPayload,
  UpdatePluginPayload,
} from "@/shared/api/modules/plugin/plugins";
import { PluginMeta } from "@/shared/type/plugin.type";

const PAGE_SIZE = 10;

// ── 响应数据解析辅助函数 ──────────────────────────────────────────────────────
// 后端可能返回两种结构，统一处理：
// 1. { data: { list: PluginMeta[], total: number } }  分页结构
// 2. { data: PluginMeta[] }                           数组结构
function extractPluginList(responseData: unknown): PluginMeta[] {
  if (!responseData) return [];
  if (Array.isArray(responseData)) return responseData;
  const d = responseData as Record<string, unknown>;
  if (Array.isArray(d.list)) return d.list as PluginMeta[];
  if (Array.isArray(d.data)) return d.data as PluginMeta[];
  if (d.data && Array.isArray((d.data as Record<string, unknown>).list)) {
    return (d.data as Record<string, unknown>).list as PluginMeta[];
  }
  return [];
}

function extractTotal(responseData: unknown): number {
  if (!responseData) return 0;
  const d = responseData as Record<string, unknown>;
  if (typeof d.total === "number") return d.total;
  if (d.data && typeof (d.data as Record<string, unknown>).total === "number") {
    return (d.data as Record<string, unknown>).total as number;
  }
  if (Array.isArray(responseData)) return responseData.length;
  return 0;
}

export function useAdminPlugins() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);

  // ── 分页列表（管理端） ───────────────────────────────────────────────────────
  const { data: pageRaw, isLoading } = useQuery({
    queryKey: ["admin", "plugins", page],
    queryFn: async () => {
      const res = await pluginApi.list({ page, page_size: PAGE_SIZE });
      return res.data?.data ?? null;
    },
  });

  const plugins: PluginMeta[] = extractPluginList(pageRaw);
  const total: number = extractTotal(pageRaw);

  // ── 已启用插件列表（仅供展示，加载由 PluginContext 用 fetch 完成） ─────────
  // 注意：这里的 enabledPlugins 仅用于管理 UI 展示"已启用数量"等信息
  // PluginContext 不能调用这个 hook（违反 Hook 规则），它直接用 fetch
  const { data: enabledPluginsRaw } = useQuery({
    queryKey: ["admin", "plugins", "enabled"],
    queryFn: async () => {
      const res = await pluginApi.listEnabled();
      return res.data?.data ?? null;
    },
  });
  const enabledPlugins: PluginMeta[] = extractPluginList(enabledPluginsRaw);
  const pageTotal = Math.ceil(total / PAGE_SIZE);

  // ── 通用成功回调 ──────────────────────────────────────────────────────────
  const invalidateAndToast = (msg: string) => {
    queryClient.invalidateQueries({ queryKey: ["admin", "plugins"] });
    toast.success(msg);
  };

  // ── 文件上传 ──────────────────────────────────────────────────────────────
  const upload = async (file: File) => {
    const res = await pluginApi.upload(file);
    return res.data.data;
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
    handleToggle,
    handleDelete,
    handleCreate,
    handleUpdate,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}
