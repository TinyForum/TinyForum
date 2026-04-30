// hooks/admin/useAdminUsers.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { User, PageData } from "@/shared/api/types";
import { toast } from "react-hot-toast";
import {
  adminUsersApi,
  ResetPasswordResponse,
} from "@/shared/api/modules/admin/user";

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
  return useQuery({
    queryKey: adminUsersKeys.list(params || {}),
    queryFn: async () => {
      const res = await adminUsersApi.listUsers(params);
      if (res.data.code !== 0)
        throw new Error(res.data.message || "获取用户列表失败");
      return res.data.data as PageData<User>;
    },
  });
}

// ========== 设置用户激活状态 ==========
export function useAdminSetUserActive() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, isActive }: { id: number; isActive: boolean }) =>
      adminUsersApi.setUserActive(id, isActive),
    onSuccess: (_, { isActive }) => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success(isActive ? "用户已激活" : "用户已停用");
    },
    onError: (error: Error) => {
      toast.error(error.message || "操作失败");
    },
  });
}

// ========== 设置用户封禁状态 ==========
export function useAdminSetUserBlocked() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, isBlocked }: { id: number; isBlocked: boolean }) =>
      adminUsersApi.setUserBlocked(id, isBlocked),
    onSuccess: (_, { isBlocked }) => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success(isBlocked ? "用户已封禁" : "用户已解封");
    },
    onError: (error: Error) => {
      toast.error(error.message || "操作失败");
    },
  });
}

// ========== 设置用户角色 ==========
export function useAdminSetUserRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, role }: { id: number; role: string }) =>
      adminUsersApi.setUserRole(id, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success("角色已更新");
    },
    onError: (error: Error) => {
      toast.error(error.message || "更新角色失败");
    },
  });
}

// ========== 删除用户 ==========
export function useAdminDeleteUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => adminUsersApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.all });
      toast.success("用户已删除");
    },
    onError: (error: Error) => {
      toast.error(error.message || "删除失败");
    },
  });
}

// ========== 重置用户密码 ==========
export function useAdminResetUserPassword() {
  return useMutation({
    mutationFn: (id: number) => adminUsersApi.resetUserPassword(id),
    onSuccess: (res) => {
      const data = res.data.data as ResetPasswordResponse;
      toast.success(data?.message || "密码已重置");
      // 可在此处展示新密码（如果后端返回）
      if (data?.message?.includes("新密码")) {
        // 例如通过 alert 或额外弹窗展示
        console.log("重置结果:", data);
      }
    },
    onError: (error: Error) => {
      toast.error(error.message || "重置密码失败");
    },
  });
}
