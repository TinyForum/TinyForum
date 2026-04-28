// store/register.ts
import { create } from "zustand";
import { authApi } from "@/lib/api";
import { useAuthStore } from "./auth";
import { ApiError } from "@/shared/api/types/basic.type";

interface RegisterState {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  agreeToTerms: boolean;
  errors: Record<string, string>;
  isValid: boolean;
  isLoading: boolean;
  serverError: string | null;

  setUsername: (username: string) => void;
  setEmail: (email: string) => void;
  setPassword: (password: string) => void;
  setConfirmPassword: (confirmPassword: string) => void;
  setAgreeToTerms: (agree: boolean) => void;
  setServerError: (error: string | null) => void;
  clearErrors: () => void;
  resetForm: () => void;
  validateForm: () => boolean;
  register: () => Promise<{ success: boolean; message?: string }>;
}

export const useRegisterStore = create<RegisterState>()((set, get) => ({
  username: "",
  email: "",
  password: "",
  confirmPassword: "",
  agreeToTerms: false,
  errors: {},
  isValid: false,
  isLoading: false,
  serverError: null,

  setUsername: (username) => {
    set({ username });
    get().validateForm();
  },

  setEmail: (email) => {
    set({ email });
    get().validateForm();
  },

  setPassword: (password) => {
    set({ password });
    get().validateForm();
  },

  setConfirmPassword: (confirmPassword) => {
    set({ confirmPassword });
    get().validateForm();
  },

  setAgreeToTerms: (agreeToTerms) => {
    set({ agreeToTerms });
    get().validateForm();
  },

  setServerError: (serverError) => set({ serverError }),

  clearErrors: () => set({ errors: {}, serverError: null }),

  resetForm: () =>
    set({
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      agreeToTerms: false,
      errors: {},
      isValid: false,
      serverError: null,
    }),

  validateForm: () => {
    const { username, email, password, confirmPassword } = get();
    const errors: Record<string, string> = {};

    if (!username) {
      errors.username = "请输入用户名";
    } else if (username.length < 2) {
      errors.username = "用户名至少 2 个字符";
    } else if (username.length > 50) {
      errors.username = "用户名最多 50 个字符";
    } else if (!/^[a-zA-Z0-9_\u4e00-\u9fa5]+$/.test(username)) {
      errors.username = "用户名只能包含字母、数字、下划线和中文";
    }

    if (!email) {
      errors.email = "请输入邮箱";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "请输入有效的邮箱地址";
    }

    if (!password) {
      errors.password = "请输入密码";
    } else if (password.length < 6) {
      errors.password = "密码至少 6 个字符";
    } else if (password.length > 32) {
      errors.password = "密码最多 32 个字符";
    }

    if (password !== confirmPassword) {
      errors.confirmPassword = "两次输入的密码不一致";
    }

    const isValid = Object.keys(errors).length === 0;
    set({ errors, isValid });
    return isValid;
  },

  register: async () => {
    const { username, email, password, validateForm, resetForm } = get();

    if (!validateForm()) {
      return { success: false, message: "请检查表单填写" };
    }

    set({ isLoading: true, serverError: null });

    try {
      const res = await authApi.register({ username, email, password });

      if (res.data.data?.user) {
        const { user } = res.data.data;
        useAuthStore.getState().setAuth(user);
        resetForm();
        return { success: true, message: "注册成功" };
      }

      return { success: false, message: "响应数据格式错误" };
    } catch (err: unknown) {
      let errorMessage = "注册失败，请重试";
      const error = err as ApiError; // ← 使用 ApiError 类型

      if (error.response?.data?.errors) {
        const fieldErrors: Record<string, string> = {};
        error.response.data.errors.forEach((errItem) => {
          fieldErrors[errItem.field] = errItem.message;
        });
        set({ errors: fieldErrors });
        errorMessage = "请检查表单填写";
      } else if (error.response?.data?.message) {
        errorMessage = error.response.data.message;
      } else if (error.message) {
        errorMessage = error.message;
      }

      set({ serverError: errorMessage });
      return { success: false, message: errorMessage };
    } finally {
      set({ isLoading: false });
    }
  },
}));
