// store/password-reset.ts
import { create } from "zustand";
import { authApi } from "@/shared/api";
import { ApiError } from "@/shared/api/types/basic.type";

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
      set({ isValid: false, error: "token is required" });
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
        error: isValid ? null : "Invalid or expired token",
      });

      return isValid;
    } catch (err: unknown) {
      const error = err as ApiError;
      console.error("Token validation error:", error);
      set({
        isValid: false,
        isLoading: false,
        error: error.response?.data?.message || "Token validation failed",
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
        set({ error: "请输入邮箱地址" });
        return { success: false };
      }

      // 简单的邮箱格式验证
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(email)) {
        set({ error: "请输入有效的邮箱地址" });
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

      console.log("=== ResetPassword Debug ===");
      console.log("1. Token:", token);
      console.log("2. Password length:", password?.length);
      console.log("3. Confirm password length:", confirmPassword?.length);
      console.log("4. Passwords match:", password === confirmPassword);

      if (!token) {
        console.log("❌ No token");
        set({ error: "无效的重置链接" });
        return { success: false };
      }

      if (!password) {
        console.log("❌ No password");
        set({ error: "请输入新密码" });
        return { success: false };
      }

      if (password.length < 6) {
        console.log("❌ Password too short");
        set({ error: "密码长度至少为6个字符" });
        return { success: false };
      }

      if (password !== confirmPassword) {
        console.log("❌ Passwords don't match");
        set({ error: "两次输入的密码不一致" });
        return { success: false };
      }

      console.log("✅ Validation passed, calling API...");
      set({ isLoading: true, error: null });

      try {
        console.log("Calling authApi.resetPassword with:", {
          token,
          password: "***",
        });
        const response = await authApi.resetPasswordWithToken({
          token,
          password,
        });

        console.log("5. Response status:", response.status);
        console.log("6. Response data:", response.data);
        console.log("7. Response code:", response.data?.code);
        console.log("8. Response success:", response.data?.data?.success);

        if (response.data.code === 0 && response.data.data?.success) {
          console.log("✅ Reset successful!");
          set({
            isSuccess: true,
            isLoading: false,
            error: null,
          });
          return { success: true };
        }

        console.log("❌ Reset failed - unexpected response");
        set({
          error: response.data.message || "密码重置失败，请重试",
          isLoading: false,
        });
        return { success: false };
      } catch (err: unknown) {
        const error = err as ApiError;
        console.error("❌ Reset password error:", error);
        console.error("Error response:", error.response);
        console.error("Error message:", error.message);
        console.error("Error code:", error.response?.data?.code);
        console.error("Error data:", error.response?.data);

        set({
          error: error.response?.data?.message || "密码重置失败，请重试",
          isLoading: false,
        });
        return { success: false };
      }
    },
  }),
);
