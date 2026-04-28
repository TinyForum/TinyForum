// store/logout.ts
import { create } from "zustand";
import { authApi } from "@/lib/api";
import { useAuthStore } from "./auth";

interface LogoutState {
  isLoading: boolean;

  // 登出操作
  logout: () => Promise<void>;

  // 强制登出（即使 API 失败也清除本地状态）
  forceLogout: () => void;
}

export const useLogoutStore = create<LogoutState>()((set ) => ({
  isLoading: false,

  logout: async () => {
    set({ isLoading: true });

    try {
      // 调用后端登出 API
      await authApi.logout();
    } catch (error) {
      console.error("登出 API 失败:", error);
      // 即使 API 失败，也继续清除本地状态
    } finally {
      // 清除认证状态
      useAuthStore.getState().logout().catch(console.error);

      // 清除其他 Store 状态（如果需要）
      //   useUserStore.getState().reset();
      // useCartStore.getState().clear();

      set({ isLoading: false });
    }
  },

  forceLogout: () => {
    useAuthStore.getState().logout().catch(console.error);
    set({ isLoading: false });
  },
}));
