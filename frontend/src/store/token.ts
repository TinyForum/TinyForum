// store/password-reset.ts
import { create } from "zustand";
import { ApiError } from "@/shared/api/types/basic.model";
import { authApi } from "@/shared/api/modules/auth";

interface ForgotPasswordState {
  // 表单数据
  email: string;

  // 状态
  isLoading: boolean;
  isEmailSent: boolean;
  error: string | null;

  // Actions
  setEmail: (email: string) => void;
  setError: (error: string | null) => void;
  resetForm: () => void;
  resetState: () => void;

  // 发送重置邮件
  sendResetEmail: () => Promise<{ success: boolean }>;
}

interface ResetPasswordState {
  // 表单数据
  token: string;
  password: string;
  confirmPassword: string;

  // 状态
  isLoading: boolean;
  isSuccess: boolean;
  error: string | null;

  // Actions
  setToken: (token: string) => void;
  setPassword: (password: string) => void;
  setConfirmPassword: (password: string) => void;
  setError: (error: string | null) => void;
  resetForm: () => void;

  // 重置密码
  resetPassword: () => Promise<{ success: boolean }>;
}

interface ValidateTokenState {
  isValid: boolean | null;
  isLoading: boolean;
  error: string | null;
  validateToken: (token: string) => Promise<boolean>;
  resetValidation: () => void;
}

export const useValidateTokenStore = create<ValidateTokenState>()((set) => ({
  isValid: null,
  isLoading: false,
  error: null,

  validateToken: async (token: string) => {
    if (!token) {
      set({ isValid: false, error: "invalid_reset_link_expired" });
      return false;
    }

    set({ isLoading: true, error: null });

    try {
      // ✅ 使用修复后的 API 方法
      const response = await authApi.validateToken({ token: token });

      // 根据你的后端响应结构调整
      const isValid = response.data.data?.valid || false;

      set({
        isValid,
        isLoading: false,
        error: isValid ? null : "invalid_reset_link_expired",
      });

      return isValid;
    } catch (err: unknown) {
      const error = err as ApiError;
      set({
        isValid: false,
        isLoading: false,
        error: error.response?.data?.message || "invalid_reset_link_expired",
      });
      return false;
    }
  },

  resetValidation: () => set({ isValid: null, isLoading: false, error: null }),
}));

// 忘记密码 Store
export const useForgotPasswordStore = create<ForgotPasswordState>()(
  (set, get) => ({
    email: "",
    isLoading: false,
    isEmailSent: false,
    error: null,

    setEmail: (email) => set({ email, error: null, isEmailSent: false }),
    setError: (error) => set({ error }),

    resetForm: () => set({ email: "", error: null }),
    resetState: () =>
      set({ isLoading: false, isEmailSent: false, error: null }),

    sendResetEmail: async () => {
      const { email } = get();

      if (!email) {
        set({ error: "enter_email_address" });
        return { success: false };
      }

      // 简单的邮箱格式验证
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(email)) {
        set({ error: "invalid_email_address" });
        return { success: false };
      }

      set({ isLoading: true, error: null });

      try {
        await authApi.forgotPassword({ email });

        set({
          isEmailSent: true,
          isLoading: false,
          error: null,
        });

        return { success: true };
      } catch {
        // const error = err as ApiError;
        // 为了安全，即使邮箱不存在也返回成功（防止邮箱枚举攻击）
        set({
          isEmailSent: true, // 总是显示成功，但实际可能没发送
          isLoading: false,
          error: null,
        });
        return { success: true };
      }
    },
  }),
);

// 重置密码 Store
export const useResetPasswordStore = create<ResetPasswordState>()(
  (set, get) => ({
    token: "",
    password: "",
    confirmPassword: "",
    isLoading: false,
    isSuccess: false,
    error: null,

    setToken: (token) => set({ token, error: null }),
    setPassword: (password) => set({ password, error: null }),
    setConfirmPassword: (confirmPassword) =>
      set({ confirmPassword, error: null }),
    setError: (error) => set({ error }),

    resetForm: () =>
      set({
        password: "",
        confirmPassword: "",
        error: null,
        isSuccess: false,
      }),

    // store/password-reset.ts - 修改 resetPassword 方法

    resetPassword: async () => {
      const { token, password, confirmPassword } = get();

      if (!token) {
        set({ error: "invalid_reset_link" });
        return { success: false };
      }

      if (!password) {
        set({ error: "enter_new_password" });
        return { success: false };
      }

      if (password.length < 6) {
        set({ error: "password_too_short" });
        return { success: false };
      }

      if (password !== confirmPassword) {
        set({ error: "passwords_not_match" });
        return { success: false };
      }

      set({ isLoading: true, error: null });

      try {
        const response = await authApi.resetPasswordWithToken({
          token,
          password,
        });

        if (response.data.code === 0 && response.data.data?.success) {
          set({
            isSuccess: true,
            isLoading: false,
            error: null,
          });
          return { success: true };
        }

        set({
          error: response.data.message || "reset_password_failed",
          isLoading: false,
        });
        return { success: false };
      } catch (err: unknown) {
        const error = err as ApiError;

        set({
          error: error.response?.data?.message || "reset_password_failed",
          isLoading: false,
        });
        return { success: false };
      }
    },
  }),
);
