import { useState, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import {
  pluginApi,
  CreatePluginPayload,
  UpdatePluginPayload,
} from "@/shared/api/modules/plugin/plugins";

const PAGE_SIZE = 10;

export function useAdminPlugins() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);

  // 直接拿到分页数据
  const { data: pageData, isLoading } = useQuery({
    queryKey: ["admin", "plugins", page],
    queryFn: async () => {
      const res = await pluginApi.list({ page, page_size: PAGE_SIZE });
      console.log("请求插件： ", res);
      return (
        res.data?.data ?? { list: [], total: 0, page: 1, page_size: PAGE_SIZE }
      );
    },
  });

  // 提取插件列表和总数
  const plugins = pageData?.list ?? [];
  const total = pageData?.total ?? 0;

  // 通用成功回调
  const invalidateAndToast = (successMsg: string) => {
    queryClient.invalidateQueries({ queryKey: ["admin", "plugins"] });
    toast.success(successMsg);
  };

  const upload = async (file: File) => {
    const res = await pluginApi.upload(file);
    return res.data.data;
  };
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
    plugins, // 插件数组（PluginVO[]）
    total, // 总数（number）
    isLoading,
    page,
    pageSize: PAGE_SIZE,
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
