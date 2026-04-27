// hooks/useModeratorAuth.ts
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useLocale } from "next-intl";
import { useAuthStore } from "@/store/auth";

// 修正后的版本
export function useModeratorAuth() {
  const router = useRouter();
  const locale = useLocale();
  const { user, isAuthenticated, isHydrated } = useAuthStore();

  useEffect(() => {
    // 等待 hydration 完成
    if (!isHydrated) return;

    // 未认证 → 跳转登录
    if (!isAuthenticated) {
      router.replace(`/${locale}/auth/login?redirect=/dashboard/moderator`);
      return;
    }

    // 已认证但权限不足 → 跳转首页
    const isModerator = ["moderator", "admin", "super_admin"].includes(
      user?.role || "",
    );
    if (!isModerator) {
      router.replace(`/${locale}`);
    }
  }, [isHydrated, isAuthenticated, user?.role, router, locale]);

  return {
    isCheckingAuth: !isHydrated,
    isModerator: ["moderator", "admin", "super_admin"].includes(
      user?.role || "",
    ),
    user,
  };
}
