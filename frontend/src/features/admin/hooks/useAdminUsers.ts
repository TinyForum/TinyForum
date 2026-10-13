// hooks/admin/useAdminUsers.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";
import {
  adminUsersApi,
  ResetPasswordResponse,
} from "@/shared/api/modules/admin/user";
import { UserDO } from "@/shared/api/types/user.model.do";
import { PageData } from "@/shared/api/types/basic.model";

// 查询键
const adminUsersKeys = {
  all: ["admin", "users"] as const,
  lists: () => [...adminUsersKeys.all, "list"] as const,
  list: (params: object) => [...adminUsersKeys.lists(), params] as const,
};

// ========== 获取用户列表 ==========
export function useAdminUsers(params?: {
  page?: number;
  page_size?: number;
  keyword?: string;
}) {
  const t = useTranslations("Admin");
  return useQuery({
    queryKey: adminUsersKeys.list(params || {}),
    queryFn: async () => {
      const res = await adminUsersApi.listUsers(params);
      if (res.data.code !== 0)
        throw new Error(res.data.message || t("get_user_list_failed"));
      return res.data.data as PageData<UserDO>;
    },
  });
}

// ========== 设置用户激活状态 ==========
export function useAdminSetUserActive() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: ({ id, isActive }: { id: number; isActive: boolean }) =>
      adminUsersApi.setUserActive(id, isActive),
    onSuccess: (_, { isActive }) => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success(isActive ? t("user_activated") : t("user_deactivated"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("operation_failed"));
    },
  });
}

// ========== 设置用户封禁状态 ==========
export function useAdminSetUserBlocked() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: ({ id, isBlocked }: { id: number; isBlocked: boolean }) =>
      adminUsersApi.setUserBlocked(id, isBlocked),
    onSuccess: (_, { isBlocked }) => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success(isBlocked ? t("user_blocked") : t("user_unblocked"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("operation_failed"));
    },
  });
}

// ========== 设置用户角色 ==========
export function useAdminSetUserRole() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: ({ id, role }: { id: number; role: string }) =>
      adminUsersApi.setUserRole(id, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success(t("role_updated"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("update_role_failed"));
    },
  });
}

// ========== 删除用户 ==========
export function useAdminDeleteUser() {
  const queryClient = useQueryClient();
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: (id: number) => adminUsersApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success(t("user_deleted"));
    },
    onError: (error: Error) => {
      toast.error(error.message || t("delete_failed"));
    },
  });
}

// ========== 重置用户密码 ==========
export function useAdminResetUserPassword() {
  const t = useTranslations("Admin");
  return useMutation({
    mutationFn: (id: number) => adminUsersApi.resetUserPassword(id),
    onSuccess: (res) => {
      const data = res.data.data as ResetPasswordResponse;
      toast.success(data?.message || t("password_reset"));
      // 可在此处展示新密码（如果后端返回）
      if (data?.message?.includes("新密码")) {
        // 例如通过 alert 或额外弹窗展示
      }
    },
    onError: (error: Error) => {
      toast.error(error.message || t("reset_password_failed"));
    },
  });
}
