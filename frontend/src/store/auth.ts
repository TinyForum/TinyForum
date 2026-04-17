// store/auth.ts
import { User, userApi, authApi } from '@/lib/api';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isHydrated: boolean;
  
  setAuth: (user: User) => void;
  logout: () => Promise<void>;
  updateUser: (user: Partial<User>) => void;
  setHydrated: (state: boolean) => void;
  refreshUser: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isHydrated: false,

      setAuth: (user) => {
        set({ user, isAuthenticated: true });
      },

      logout: async () => {
        await authApi.logout();
        // 清除所有存储
        localStorage.removeItem('tiny-auth');
        sessionStorage.clear();
        set({ user: null, isAuthenticated: false });
      },

      updateUser: (partial) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...partial } : null,
        })),

      setHydrated: (state) => set({ isHydrated: state }),
      
      // 刷新用户信息
    // store/auth.ts
refreshUser: async () => {
  try {
    const currentUser = get().user;
    if (!currentUser) return;
    
    // 获取最新的角色信息
    const roleResponse = await userApi.getMeRole();
    if (roleResponse.data?.data) {
      const { role } = roleResponse.data.data;
      
      // 更新用户的角色，保持其他信息不变
      const updatedUser = {
        ...currentUser,
        role: role,
      };
      
      // ✅ 只调用 set，让 Zustand persist 中间件自动处理 localStorage
      set({ user: updatedUser, isAuthenticated: true });
      
      console.log('用户角色已刷新:', role);
    }
  } catch (error) {
    console.error('刷新用户信息失败:', error);
    // 如果刷新失败，可能是 token 过期，自动登出
    await get().logout();
  }
},
    }),
    {
      name: 'tiny-auth',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        state?.setHydrated(true);
        // 恢复后验证用户信息
        if (state?.user && state.isAuthenticated) {
          state.refreshUser().catch(console.error);
        }
      },
    }
  )
);