// store/auth.ts
import { User } from '@/lib/api';
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authApi } from '@/lib/api';


interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isHydrated: boolean; // 添加 hydration 状态
  setAuth: (user: User) => void;
  logout: () => Promise<void>;
  updateUser: (user: Partial<User>) => void;
  setHydrated: (state: boolean) => void; // 添加 setter
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      isHydrated: false, // 初始为 false

      setAuth: (user) => {
        set({ user, isAuthenticated: true });
      },

      logout: async () => {
        await authApi.logout();
        set({ user: null, isAuthenticated: false });
      },

      updateUser: (partial) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...partial } : null,
        })),
        
      setHydrated: (state) => set({ isHydrated: state }),
    }),
    {
      name: 'tiny-auth',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // hydration 完成后调用
        state?.setHydrated(true);
      },
    }
  )
);