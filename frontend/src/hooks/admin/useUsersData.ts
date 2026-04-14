

import { useEffect, useState } from "react";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { adminApi } from "@/lib/api";
import toast from "react-hot-toast";
// 用户数据管理 Hook
export function useUsersData(page: number, keyword: string, enabled: boolean) {
  const queryClient = useQueryClient();
  const t = useTranslations("admin");

  const { data, isLoading } = useQuery({
    queryKey: ["admin-users", page, keyword],
    queryFn: () =>
      adminApi
        .listUsers({ page, page_size: 20, keyword })
        .then((r) => r.data.data),
    enabled,
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, active }: { id: number; active: boolean }) =>
      adminApi.setUserActive(id, active),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      toast.success(t("operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // ✅ 修复：创建一个包装函数，接收两个单独的参数
  const handleToggleActive = (id: number, active: boolean) => {
    toggleActiveMutation.mutate({ id, active });
  };

  return {
    users: data?.list ?? [],
    total: data?.total ?? 0,
    isLoading,
    toggleActive: handleToggleActive,
    isToggling: toggleActiveMutation.isPending,
  };
}
