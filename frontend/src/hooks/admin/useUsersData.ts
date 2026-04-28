// hooks/admin/useUsersData.ts
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { adminApi } from "@/shared/api";
import toast from "react-hot-toast";
import { ApiResponse } from "@/shared/api/types";

// 类型定义
interface UserListResponse {
  list: User[];
  total: number;
  page: number;
  page_size: number;
}

interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  is_active: boolean;
  is_blocked: boolean;
  // 其他字段
}

interface ErrorResponse {
  response?: {
    data?: {
      code?: number;
      message?: string;
    };
  };
  message?: string;
}

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
        .then((r: { data: ApiResponse<UserListResponse> }) => r.data.data),
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

  // 删除用户 - 修复未使用的参数
  const deleteUserMutation = useMutation({
    mutationFn: ({ id }: { id: number }) => adminApi.deleteUser(id),
    onSuccess: () => {
      // 移除未使用的 _ 和 variables 参数
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
      toast.success(t("user_deleted_successfully"));
    },
    onError: () => toast.error(t("delete_failed")),
  });

  // 重置密码
  const resetPasswordMutation = useMutation({
    mutationFn: (id: number) => adminApi.resetUserPassword(id),
    onSuccess: () => {
      // 移除未使用的 userId 参数
      toast.success(t("password_reset_and_notified"), {
        duration: 5000,
      });
      queryClient.invalidateQueries({ queryKey: ["admin-users"] });
    },
    onError: (error: unknown) => {
      // 修复 any 类型
      const err = error as ErrorResponse;
      const errorCode = err?.response?.data?.code;

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

  // 包装函数 - 修复未使用的 username 参数
  const handleDeleteUser = (id: number) => {
    // 添加下划线前缀表示未使用
    deleteUserMutation.mutate({ id });
  };

  const handleToggleActive = (id: number, active: boolean) => {
    toggleActiveMutation.mutate({ id, active });
  };

  const handleToggleBlock = (id: number, blocked: boolean) => {
    toggleBlockMutation.mutate({ id, blocked });
  };

  const handleToggleRole = (id: number, role: string) => {
    toggleRoleMutation.mutate({ id, role });
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
