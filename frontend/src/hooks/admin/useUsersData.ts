// hooks/admin/useUsersData.ts
import { useEffect, useState } from "react";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { adminApi } from "@/lib/api";
import toast from "react-hot-toast";
// import { ResetPasswordRequest } from "@/lib/api/modules/admin";

// 用户数据管理 Hook（扩展版）
export function useUsersData(page: number, keyword: string, enabled: boolean) {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");

  // 获取用户列表
  const { data, isLoading } = useQuery({
    queryKey: ["admin-users", page, keyword],
    queryFn: () =>
      adminApi
        .listUsers({ page, page_size: 20, keyword })
        .then((r) => r.data.data),
    enabled,
  });

  // 切换激活状态
  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, active }: { id: number; active: boolean }) =>
      adminApi.setUserActive(id, active),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      toast.success(t("operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 切换封禁状态
  const toggleBlockMutation = useMutation({
    mutationFn: ({ id, blocked }: { id: number; blocked: boolean }) =>
      adminApi.setUserBlocked(id, blocked),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      toast.success(t("operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 切换角色
  const toggleRoleMutation = useMutation({
    mutationFn: ({ id, role }: { id: number; role: string }) =>
      adminApi.setUserRole(id, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      toast.success(t("operation_successful"));
    },
    onError: () => toast.error(t("operation_failed")),
  });

  // 删除用户
  const deleteUserMutation = useMutation({
    mutationFn: ({ id }: { id: number }) => adminApi.deleteUser(id),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      toast.success(t("user_deleted_successfully"));
    },
    onError: () => toast.error(t("delete_failed")),
  });

  // 重置密码（只需要传用户 ID）
  const resetPasswordMutation = useMutation({
    mutationFn: (id: number) => adminApi.resetUserPassword(id),
    onSuccess: (_, userId) => {
      toast.success(t("password_reset_and_notified"), {
        duration: 5000,
      });
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
    },
    onError: (error: any) => {
      // 根据错误码显示不同的错误信息
      const errorCode = error?.response?.data?.code;
      if (errorCode === 20011) {
        toast.error(t("cannot_modify_self"));
      } else if (errorCode === 20012) {
        toast.error(t("cannot_modify_super_admin"));
      } else if (errorCode === 40001) {
        toast.error(t("insufficient_permission"));
      } else {
        toast.error(t("operation_failed"));
      }
    },
  });

  // 包装函数
  const handleToggleActive = (id: number, active: boolean) => {
    toggleActiveMutation.mutate({ id, active });
  };

  const handleToggleBlock = (id: number, blocked: boolean) => {
    toggleBlockMutation.mutate({ id, blocked });
  };

  const handleToggleRole = (id: number, role: string) => {
    toggleRoleMutation.mutate({ id, role });
  };

  const handleDeleteUser = (id: number, username: string) => {
    deleteUserMutation.mutate({ id });
  };

  const handleResetPassword = (id: number) => {
    resetPasswordMutation.mutate(id);
  };

  return {
    users: data?.list ?? [],
    total: data?.total ?? 0,
    isLoading,
    // 激活/停用
    toggleActive: handleToggleActive,
    isTogglingActive: toggleActiveMutation.isPending,
    // 封禁/解封
    toggleBlock: handleToggleBlock,
    isTogglingBlock: toggleBlockMutation.isPending,
    // 角色管理
    toggleRole: handleToggleRole,
    isUpdatingRole: toggleRoleMutation.isPending,
    // 删除用户
    deleteUser: handleDeleteUser,
    isDeleting: deleteUserMutation.isPending,
    // 重置密码
    resetPassword: handleResetPassword,
    isResettingPassword: resetPasswordMutation.isPending,
  };
}
