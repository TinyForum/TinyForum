// store/auth.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
// import type { User } from '@/types';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean; // 添加加载状态
  setAuth: (user: User, token: string) => void;
  logout: () => void;
  updateUser: (user: Partial<User>) => void;
  checkAuth: () => void; // 添加检查认证的方法
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: true,

      setAuth: (user, token) => {
        if (typeof window !== 'undefined') {
          localStorage.setItem('bbs_token', token);
        }
        set({ user, token, isAuthenticated: true, isLoading: false });
      },

      logout: () => {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('bbs_token');
        }
        set({ user: null, token: null, isAuthenticated: false, isLoading: false });
      },

      updateUser: (partial) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...partial } : null,
        })),

      checkAuth: () => {
        const token = get().token;
        const user = get().user;
        
        if (token && user) {
          // 可选：验证 token 是否过期
          set({ isAuthenticated: true, isLoading: false });
        } else {
          set({ isAuthenticated: false, isLoading: false });
        }
      },
    }),
    {
      name: 'tiny-auth',
      partialize: (state) => ({ 
        user: state.user, 
        token: state.token, 
        isAuthenticated: state.isAuthenticated 
      }),
    }
  )
);

// 在应用启动时恢复状态
if (typeof window !== 'undefined') {
  // 从 localStorage 恢复 token
  const storedToken = localStorage.getItem('bbs_token');
  const storedAuth = localStorage.getItem('tiny-auth');
  
  if (storedToken && storedAuth) {
    try {
      const authState = JSON.parse(storedAuth);
      useAuthStore.setState({
        user: authState.state.user,
        token: authState.state.token,
        isAuthenticated: !!authState.state.token,
        isLoading: false,
      });
    } catch (error) {
      console.error('Failed to restore auth state:', error);
    }
  } else {
    useAuthStore.setState({ isLoading: false });
  }
}