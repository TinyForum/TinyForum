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

      // 登出：服务端调用尽力而为，无论成功与否都清除本地会话
      logout: async () => {
        try {
          await authApi.logout();
        } catch (error) {
          console.error("登出接口调用失败，仍清除本地会话:", error);
        } finally {
          localStorage.removeItem("tiny-auth");
          sessionStorage.clear();
          set({ user: null, isAuthenticated: false });
        }
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
          }
        } catch (error) {
          // 仅在认证失效（401/403）时登出；网络或服务端瞬时错误保留本地会话
          const status = (error as { response?: { status?: number } })?.response
            ?.status;
          if (status === 401 || status === 403) {
            await get().logout();
          } else {
            console.error("刷新用户信息失败:", error);
          }
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
