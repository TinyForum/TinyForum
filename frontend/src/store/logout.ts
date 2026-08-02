// store/logout.ts
import { create } from "zustand";
import { useAuthStore } from "./auth";
import { authApi } from "@/shared/api/modules/auth";

interface LogoutState {
  isLoading: boolean;
  logout: () => Promise<void>;
  forceLogout: () => void;
}

export const useLogoutStore = create<LogoutState>()((set) => ({
  isLoading: false,

  logout: async () => {
    set({ isLoading: true });
    try {
      // 只清除本地状态，不调用 API（同步方法）
      useAuthStore.getState().logout();
    } catch (error) {
      console.error("登出失败:", error);
    } finally {
      set({ isLoading: false });
    }
  },

  forceLogout: () => {
    useAuthStore.getState().logout();
    set({ isLoading: false });
  },
}));
