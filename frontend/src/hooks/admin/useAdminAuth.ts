// hooks/admin/useAdminAuth.ts
import { useEffect, useState } from "react";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";

export function useAdminAuth() {
  const { user, isAuthenticated, isHydrated } = useAuthStore();
  const router = useRouter();
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);

  useEffect(() => {
    // 等待 store hydration 完成
    if (!isHydrated) return;

    // hydration 完成，停止检查状态
    setIsCheckingAuth(false);

    // 注意：middleware 已经处理了认证重定向
    // 这里不需要再重定向，避免冲突
    // 只需要返回权限状态即可
  }, [isHydrated, router]);

  const isAdmin =
    isAuthenticated && (user?.role === "admin" || user?.role === "super_admin");

  return {
    isCheckingAuth,
    isAdmin,
    user,
    isAuthenticated,
  };
}
