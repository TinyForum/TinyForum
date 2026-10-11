// store/logout.ts
import { create } from "zustand";
import { useAuthStore } from "./auth";

interface LogoutState {
  isLoading: boolean;

  // 登出操作
  logout: () => Promise<void>;

  // 强制登出（即使 API 失败也清除本地状态）
  forceLogout: () => void;
}

export const useLogoutStore = create<LogoutState>()((set) => ({
  isLoading: false,

  logout: async () => {
    set({ isLoading: true });

    try {
      // auth store 的 logout 内部已尽力调用后端登出并清理本地会话
      await useAuthStore.getState().logout();
    } finally {
      set({ isLoading: false });
    }
  },

  forceLogout: () => {
    useAuthStore.getState().logout().catch(console.error);
    set({ isLoading: false });
  },
}));
