// store/register.ts
import { create } from 'zustand';
import { authApi } from '@/lib/api';

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
  
  // Actions
  setUsername: (username: string) => void;
  setEmail: (email: string) => void;
  setPassword: (password: string) => void;
  setConfirmPassword: (confirmPassword: string) => void;
  setAgreeToTerms: (agree: boolean) => void;
  setErrors: (errors: Record<string, string>) => void;
  resetForm: () => void;
  
  // 验证
  validateForm: () => boolean;
  
  // 注册操作
  register: () => Promise<{ success: boolean; message?: string }>;
}

export const useRegisterStore = create<RegisterState>()((set, get) => ({
  // 初始状态
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeToTerms: false,
  errors: {},
  isValid: false,
  isLoading: false,

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
  
  setErrors: (errors) => set({ errors }),
  
  resetForm: () => set({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    agreeToTerms: false,
    errors: {},
    isValid: false,
  }),

  // 表单验证
  validateForm: () => {
    const { username, email, password, confirmPassword, agreeToTerms } = get();
    const errors: Record<string, string> = {};

    // 用户名验证
    if (!username) {
      errors.username = '请输入用户名';
    } else if (username.length < 3) {
      errors.username = '用户名至少 3 个字符';
    } else if (username.length > 20) {
      errors.username = '用户名最多 20 个字符';
    } else if (!/^[a-zA-Z0-9_]+$/.test(username)) {
      errors.username = '用户名只能包含字母、数字和下划线';
    }

    // 邮箱验证
    if (!email) {
      errors.email = '请输入邮箱';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = '请输入有效的邮箱地址';
    }

    // 密码验证
    if (!password) {
      errors.password = '请输入密码';
    } else if (password.length < 6) {
      errors.password = '密码至少 6 个字符';
    } else if (password.length > 32) {
      errors.password = '密码最多 32 个字符';
    }

    // 确认密码验证
    if (password !== confirmPassword) {
      errors.confirmPassword = '两次输入的密码不一致';
    }

    // 协议验证
    if (!agreeToTerms) {
      errors.agreeToTerms = '请阅读并同意用户协议';
    }

    const isValid = Object.keys(errors).length === 0;
    set({ errors, isValid });
    
    return isValid;
  },

  // 注册操作
  register: async () => {
    const { username, email, password, validateForm, resetForm } = get();
    
    if (!validateForm()) {
      return { success: false, message: '请检查表单填写' };
    }

    set({ isLoading: true });

    try {
      await authApi.register({ username, email, password });
      
      resetForm();
      
      return { 
        success: true, 
        message: '注册成功，请登录' 
      };
    } catch (error: any) {
      const errorMessage = error?.response?.data?.message || '注册失败，请重试';
      
      // 处理特定字段错误
      if (error?.response?.data?.errors) {
        const fieldErrors: Record<string, string> = {};
        error.response.data.errors.forEach((err: any) => {
          fieldErrors[err.field] = err.message;
        });
        set({ errors: fieldErrors });
      }
      
      return { success: false, message: errorMessage };
    } finally {
      set({ isLoading: false });
    }
  },
}));