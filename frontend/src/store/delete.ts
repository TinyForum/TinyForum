// store/deleteAccount.ts
import { create } from "zustand";
import { authApi } from "@/lib/api";
import { useAuthStore } from "./auth";
import { ApiError } from "@/shared/api/types/basic.type";

interface DeletionStatus {
  is_deleted: boolean;
  deleted_at?: string;
  can_restore: boolean;
  remaining_days?: number;
}

interface DeleteAccountState {
  // 表单数据
  confirmText: string;
  password: string;

  // 状态
  isModalOpen: boolean;
  isLoading: boolean;
  error: string | null;

  // 删除状态信息
  deletionStatus: DeletionStatus | null;

  // Actions
  setConfirmText: (text: string) => void;
  setPassword: (password: string) => void;
  setModalOpen: (open: boolean) => void;
  setError: (error: string | null) => void;
  resetForm: () => void;
  setDeletionStatus: (status: DeletionStatus | null) => void;

  // API 操作
  deleteAccount: () => Promise<{ success: boolean }>;
  cancelDeletion: () => Promise<{ success: boolean }>;
  confirmDeletion: () => Promise<{ success: boolean }>;
  fetchDeletionStatus: () => Promise<void>;
}

export const useDeleteAccountStore = create<DeleteAccountState>()(
  (set, get) => ({
    // 初始状态
    confirmText: "",
    password: "",
    isModalOpen: false,
    isLoading: false,
    error: null,
    deletionStatus: null,

    // Setters
    setConfirmText: (confirmText) => set({ confirmText, error: null }),
    setPassword: (password) => set({ password, error: null }),
    setModalOpen: (isModalOpen) => set({ isModalOpen, error: null }),
    setError: (error) => set({ error }),
    setDeletionStatus: (deletionStatus) => set({ deletionStatus }),
    resetForm: () => set({ confirmText: "", password: "", error: null }),

    // 获取删除状态
    fetchDeletionStatus: async () => {
      try {
        const response = await authApi.getDeletionStatus();
        if (response.data.data) {
          set({ deletionStatus: response.data.data as DeletionStatus });
        }
      } catch (err: unknown) {
        const error = err as ApiError;
        console.error(
          "获取删除状态失败:",
          error.response?.data?.message || error.message,
        );
      }
    },

    // 请求注销（软删除）
    deleteAccount: async () => {
      const { confirmText, password } = get();

      if (confirmText !== "DELETE") {
        set({ error: "请输入 DELETE 确认删除" });
        return { success: false };
      }

      set({ isLoading: true, error: null });

      try {
        await authApi.deleteAccount({ confirm: confirmText, password });

        // 刷新删除状态
        await get().fetchDeletionStatus();

        // 关闭模态框并重置表单
        set({ isModalOpen: false });
        get().resetForm();

        return { success: true };
      } catch (err: unknown) {
        const error = err as ApiError;
        const errorMessage =
          error.response?.data?.message || "注销请求失败，请重试";
        set({ error: errorMessage });
        return { success: false };
      } finally {
        set({ isLoading: false });
      }
    },

    // 取消注销（恢复账户）
    cancelDeletion: async () => {
      set({ isLoading: true, error: null });

      try {
        await authApi.cancelDeletion();

        // 刷新删除状态
        await get().fetchDeletionStatus();

        // 重置表单
        get().resetForm();

        return { success: true };
      } catch (err: unknown) {
        const error = err as ApiError;
        const errorMessage =
          error.response?.data?.message || "取消注销失败，请重试";
        set({ error: errorMessage });
        return { success: false };
      } finally {
        set({ isLoading: false });
      }
    },

    // 确认永久删除（硬删除）
    confirmDeletion: async () => {
      const { confirmText, password } = get();

      if (confirmText !== "PERMANENT_DELETE") {
        set({ error: "请输入 PERMANENT_DELETE 确认永久删除" });
        return { success: false };
      }

      set({ isLoading: true, error: null });

      try {
        await authApi.confirmDeletion({ confirm: confirmText, password });

        // 调用 auth store 的 logout 方法清除认证状态
        await useAuthStore.getState().logout();

        // 关闭模态框并重置表单
        set({ isModalOpen: false });
        get().resetForm();

        return { success: true };
      } catch (err: unknown) {
        const error = err as ApiError;
        const errorMessage =
          error.response?.data?.message || "永久删除失败，请重试";
        set({ error: errorMessage });
        return { success: false };
      } finally {
        set({ isLoading: false });
      }
    },
  }),
);
