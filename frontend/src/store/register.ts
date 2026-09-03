// store/register.ts
import { create } from "zustand";
import { useAuthStore } from "./auth";
import { ApiError } from "@/shared/api/types/basic.model";
import { authApi } from "@/shared/api/modules/auth";

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
      errors.username = "enter_username";
    } else if (username.length < 2) {
      errors.username = "username_min_length";
    } else if (username.length > 50) {
      errors.username = "username_max_length";
    } else if (!/^[a-zA-Z0-9_\u4e00-\u9fa5]+$/.test(username)) {
      errors.username = "username_invalid_chars";
    }

    if (!email) {
      errors.email = "enter_email";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "invalid_email_address";
    }

    if (!password) {
      errors.password = "enter_password";
    } else if (password.length < 6) {
      errors.password = "password_min_length";
    } else if (password.length > 32) {
      errors.password = "password_max_length";
    }

    if (password !== confirmPassword) {
      errors.confirmPassword = "passwords_not_match";
    }

    const isValid = Object.keys(errors).length === 0;
    set({ errors, isValid });
    return isValid;
  },

  register: async () => {
    const { username, email, password, validateForm, resetForm } = get();

    if (!validateForm()) {
      return { success: false, message: "check_form_fields" };
    }

    set({ isLoading: true, serverError: null });

    try {
      const res = await authApi.register({ username, email, password });

      if (res.data.data?.user) {
        const { user } = res.data.data;
        useAuthStore.getState().setAuth(user);
        resetForm();
        return { success: true, message: "registration_successful" };
      }

      return { success: false, message: "invalid_response_format" };
    } catch (err: unknown) {
      let errorMessage = "registration_failed";
      const error = err as ApiError; // ← 使用 ApiError 类型

      if (error.response?.data?.errors) {
        const fieldErrors: Record<string, string> = {};
        error.response.data.errors.forEach((errItem) => {
          fieldErrors[errItem.field] = errItem.message;
        });
        set({ errors: fieldErrors });
        errorMessage = "check_form_fields";
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
