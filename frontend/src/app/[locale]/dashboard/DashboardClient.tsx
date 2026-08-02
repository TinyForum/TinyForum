// dashboard/page.tsx 客户端组件
"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
// import { useCurrentUserProfile } from "@/hooks/useUserProfile";
import { useLocale } from "next-intl";
import { useUserRole } from "@/features/user/hooks/useUserRole";

export default function DashboardClient() {
  const router = useRouter();
  const locale = useLocale();
  const {
    isLoading,
    isUser,
    isMember,
    isAdmin,
    isSuperAdmin,
    isModerator,
    isReviewer,
    isSystemMaintainer,
  } = useUserRole();

  useEffect(() => {
    if (isLoading) return;

    // 根据角色重定向到对应的后台
    if (isSuperAdmin || isAdmin) {
      router.replace(`/${locale}/dashboard/admin`);
    } else if (isModerator) {
      router.replace(`/${locale}/dashboard/moderator`);
    } else if (isReviewer) {
      router.replace(`/${locale}/dashboard/reviewer`);
    } else if (isMember) {
      router.replace(`/${locale}/dashboard/member`);
    } else if (isUser) {
      router.replace(`/${locale}/dashboard/user`);
    } else if (isSystemMaintainer) {
      router.replace(`/${locale}/dashboard/system`);
    } else {
      router.replace(`/${locale}/auth/login`);
    }
  }, [
    isUser,
    isMember,
    isModerator,
    isReviewer,
    isLoading,
    isAdmin,
    isSuperAdmin,
    isSystemMaintainer,
    router,
    locale,
  ]);

  // 加载中状态
  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <div className="loading loading-spinner loading-lg text-primary"></div>
        <p className="mt-4 text-base-content/60">验证权限中...</p>
      </div>
    );
  }

  // 重定向中（实际上不会渲染到这里，因为 useEffect 会执行重定向）
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <div className="loading loading-spinner loading-lg text-primary"></div>
      <p className="mt-4 text-base-content/60">正在跳转...</p>
    </div>
  );
}
