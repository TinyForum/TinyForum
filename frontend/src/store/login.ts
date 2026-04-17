// store/login.ts
import { create } from 'zustand';
import { authApi } from '@/lib/api';
import { useAuthStore } from './auth';

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
  email: '',
  password: '',
  rememberMe: false,
  isLoading: false,
  error: null,

  // Setters
  setEmail: (email) => set({ email, error: null }),
  setPassword: (password) => set({ password, error: null }),
  setRememberMe: (rememberMe) => set({ rememberMe }),
  setError: (error) => set({ error }),
  resetForm: () => set({ email: '', password: '', rememberMe: false, error: null }),

  // 登录操作
  login: async () => {
    const { email, password, resetForm } = get();
    
    if (!email || !password) {
      set({ error: '请输入邮箱和密码' });
      return { success: false };
    }

    set({ isLoading: true, error: null });

    try {
      const response = await authApi.login({ email, password });
      
      const { user } = response.data.data;

      // 更新全局认证状态
      useAuthStore.getState().setAuth(user);
      
      // 清空表单
      resetForm();
      
      return { success: true };
    } catch (error: any) {
      const errorMessage = error?.response?.data?.message || '登录失败，请重试';
      set({ error: errorMessage });
      return { success: false };
    } finally {
      set({ isLoading: false });
    }
  },
}));