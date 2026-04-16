// hooks/useModeratorAuth.ts
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useLocale } from "next-intl";
import { useAuthStore } from "@/store/auth";

export function useModeratorAuth() {
  const router = useRouter();
  const locale = useLocale();
  const { user, isAuthenticated, isHydrated } = useAuthStore();

  const isModerator = 
    user?.role === "moderator" || 
    user?.role === "admin" || 
    user?.role === "super_admin";

  const isCheckingAuth = !isHydrated;

  useEffect(() => {
    if (isHydrated) {
      if (!isAuthenticated) {
        router.replace(`/${locale}/auth/login?redirect=/dashboard/moderator`);
      } else if (!isModerator) {
        router.replace(`/${locale}`);
      }
    }
  }, [isHydrated, isAuthenticated, isModerator, router, locale]);

  return {
    isCheckingAuth,
    isModerator,
    user,
  };
}