// store/auth.ts
import { authApi } from "@/shared/api/modules/auth";
import { userApi } from "@/shared/api/modules/user";
import { UserDO } from "@/shared/api/types/user.model.do";

import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthState {
  user: UserDO | null;
  isAuthenticated: boolean;
  isHydrated: boolean;

  setAuth: (user: UserDO) => void;
  logout: () => Promise<void>;
  updateUser: (user: Partial<UserDO>) => void;
  setHydrated: (state: boolean) => void;
  refreshUser: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isHydrated: false,

      // 设置验证信息
      setAuth: (user) => {
        set({ user, isAuthenticated: true });
      },

      // 登出
      logout: async () => {
        await authApi.logout();
        localStorage.removeItem("tiny-auth");
        sessionStorage.clear();
        set({ user: null, isAuthenticated: false });
      },

      // 更新用户信息
      updateUser: (partial) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...partial } : null,
        })),

      // 设置是否已经从持久化存储中恢复
      setHydrated: (state) => set({ isHydrated: state }),

      // 刷新用户信息
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

            set({ user: updatedUser, isAuthenticated: true });

            console.log("用户角色已刷新:", role);
          }
        } catch (error) {
          console.error("刷新用户信息失败:", error);
          await get().logout();
        }
      },
    }),
    {
      name: "tiny-auth",
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
    },
  ),
);
