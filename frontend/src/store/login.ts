// store/login.ts
import { create } from "zustand";
import { useAuthStore } from "./auth";
import { ApiError } from "@/shared/api/types/basic.model";
import { authApi } from "@/shared/api/modules/auth";

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
      set({ error: "enter_email_and_password" });
      return { success: false };
    }

    set({ isLoading: true, error: null });

    try {
      const response = await authApi.login({ email, password });

      if (response.data.data) {
        const user = response.data.data;
        useAuthStore.getState().setAuth(user);
        resetForm();
        return { success: true };
      }

      // 响应格式错误
      set({ error: "invalid_response_format" });
      return { success: false };
    } catch (err: unknown) {
      const error = err as ApiError;
      const errorMessage = error.response?.data?.message || "login_failed";
      set({ error: errorMessage });
      return { success: false };
    } finally {
      set({ isLoading: false });
    }
  },
}));
