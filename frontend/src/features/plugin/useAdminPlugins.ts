import { useState, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import {
  pluginApi,
  CreatePluginPayload,
  UpdatePluginPayload,
} from "@/shared/api/modules/plugin/plugins";

export function useAdminPlugins() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const PAGE_SIZE = 10;

  const { data, isLoading } = useQuery({
    queryKey: ["admin", "plugins", page],
    queryFn: () =>
      pluginApi.list({ page, page_size: PAGE_SIZE }).then((r) => r.data),
  });

  const createMutation = useMutation({
    mutationFn: (payload: CreatePluginPayload) => pluginApi.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "plugins"] });
      toast.success("Plugin installed successfully");
    },
    onError: () => toast.error("Failed to install plugin"),
  });

  const updateMutation = useMutation({
    mutationFn: (payload: UpdatePluginPayload) => pluginApi.update(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "plugins"] });
      toast.success("Plugin updated");
    },
    onError: () => toast.error("Failed to update plugin"),
  });

  const toggleMutation = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) =>
      pluginApi.toggle(id, enabled),
    onSuccess: (_, { enabled }) => {
      queryClient.invalidateQueries({ queryKey: ["admin", "plugins"] });
      toast.success(enabled ? "Plugin enabled" : "Plugin disabled");
    },
    onError: () => toast.error("Failed to toggle plugin"),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => pluginApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "plugins"] });
      toast.success("Plugin removed");
    },
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
    plugins: data?.data ?? [],
    total: data?.total ?? 0,
    isLoading,
    page,
    pageSize: PAGE_SIZE,
    setPage,
    handleToggle,
    handleDelete,
    handleCreate,
    handleUpdate,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}
