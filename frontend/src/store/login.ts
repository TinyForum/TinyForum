// store/login.ts
import { create } from "zustand";
import { authApi } from "@/shared/api";
import { useAuthStore } from "./auth";
import { ApiError } from "@/shared/api/types/basic.model";

interface LoginState {
  // 表单数据
  email: string;
  password: string;
  rememberMe: boolean;

  // 状态
  isLoading: boolean;
  error: string | null;

  // Actions
  setEmail: (email: string) => void;
  setPassword: (password: string) => void;
  setRememberMe: (remember: boolean) => void;
  setError: (error: string | null) => void;
  resetForm: () => void;

  // 登录操作
  login: () => Promise<{ success: boolean }>;
}

export const useLoginStore = create<LoginState>()((set, get) => ({
  // 初始状态
  email: "",
  password: "",
  rememberMe: false,
  isLoading: false,
  error: null,

  // Setters
  setEmail: (email) => set({ email, error: null }),
  setPassword: (password) => set({ password, error: null }),
  setRememberMe: (rememberMe) => set({ rememberMe }),
  setError: (error) => set({ error }),
  resetForm: () =>
    set({ email: "", password: "", rememberMe: false, error: null }),

  // 登录操作
  login: async () => {
    const { email, password, resetForm } = get();

    if (!email || !password) {
      set({ error: "请输入邮箱和密码" });
      return { success: false };
    }

    set({ isLoading: true, error: null });

    try {
      const response = await authApi.login({ email, password });

      // 检查响应数据是否存在
      if (response.data.data?.user) {
        const { user } = response.data.data;
        // 更新全局认证状态
        useAuthStore.getState().setAuth(user);
        // 清空表单
        resetForm();
        return { success: true };
      }

      // 响应格式错误
      set({ error: "响应数据格式错误" });
      return { success: false };
    } catch (err: unknown) {
      const error = err as ApiError;
      const errorMessage = error.response?.data?.message || "登录失败，请重试";
      set({ error: errorMessage });
      return { success: false };
    } finally {
      set({ isLoading: false });
    }
  },
}));
