// store/register.ts
import { create } from "zustand";
import { authApi } from "@/lib/api";
import { useAuthStore } from "./auth";

interface RegisterState {
  // 表单数据
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  agreeToTerms: boolean;

  // 验证状态
  errors: Record<string, string>;
  isValid: boolean;

  // 状态
  isLoading: boolean;
  serverError: string | null;

  // Actions
  setUsername: (username: string) => void;
  setEmail: (email: string) => void;
  setPassword: (password: string) => void;
  setConfirmPassword: (confirmPassword: string) => void;
  setAgreeToTerms: (agree: boolean) => void;
  setServerError: (error: string | null) => void;
  clearErrors: () => void;
  resetForm: () => void;

  // 验证
  validateForm: () => boolean;

  // 注册操作
  register: () => Promise<{ success: boolean; message?: string }>;
}

export const useRegisterStore = create<RegisterState>()((set, get) => ({
  // 初始状态
  username: "",
  email: "",
  password: "",
  confirmPassword: "",
  agreeToTerms: false,
  errors: {},
  isValid: false,
  isLoading: false,
  serverError: null,

  // Setters
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

  // 表单验证
  validateForm: () => {
    const { username, email, password, confirmPassword, agreeToTerms } = get();
    const errors: Record<string, string> = {};

    // 用户名验证
    if (!username) {
      errors.username = "请输入用户名";
    } else if (username.length < 2) {
      errors.username = "用户名至少 2 个字符";
    } else if (username.length > 50) {
      errors.username = "用户名最多 50 个字符";
    } else if (!/^[a-zA-Z0-9_\u4e00-\u9fa5]+$/.test(username)) {
      errors.username = "用户名只能包含字母、数字、下划线和中文";
    }

    // 邮箱验证
    if (!email) {
      errors.email = "请输入邮箱";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "请输入有效的邮箱地址";
    }

    // 密码验证
    if (!password) {
      errors.password = "请输入密码";
    } else if (password.length < 6) {
      errors.password = "密码至少 6 个字符";
    } else if (password.length > 32) {
      errors.password = "密码最多 32 个字符";
    }

    // 确认密码验证
    if (password !== confirmPassword) {
      errors.confirmPassword = "两次输入的密码不一致";
    }

    // 协议验证（可选）
    // if (!agreeToTerms) {
    //   errors.agreeToTerms = '请阅读并同意用户协议';
    // }

    const isValid = Object.keys(errors).length === 0;
    set({ errors, isValid });

    return isValid;
  },

  // 注册操作
  register: async () => {
    const { username, email, password, validateForm, resetForm } = get();

    if (!validateForm()) {
      return { success: false, message: "请检查表单填写" };
    }

    set({ isLoading: true, serverError: null });

    try {
      const res = await authApi.register({ username, email, password });
      const { user } = res.data.data;

      // 更新全局认证状态
      useAuthStore.getState().setAuth(user);

      // 清空表单
      resetForm();

      return {
        success: true,
        message: "注册成功",
      };
    } catch (error: any) {
      let errorMessage = "注册失败，请重试";

      // 处理特定字段错误
      if (error?.response?.data?.errors) {
        const fieldErrors: Record<string, string> = {};
        error.response.data.errors.forEach((err: any) => {
          fieldErrors[err.field] = err.message;
        });
        set({ errors: fieldErrors });
        errorMessage = "请检查表单填写";
      } else if (error?.response?.data?.message) {
        errorMessage = error.response.data.message;
      }

      set({ serverError: errorMessage });
      return { success: false, message: errorMessage };
    } finally {
      set({ isLoading: false });
    }
  },
}));
