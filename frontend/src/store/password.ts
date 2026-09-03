// features/settings/hooks/useChangePassword.ts
import { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { getErrorMessage } from "@/shared/lib/utils";
import { authApi } from "@/shared/api/modules/auth";

interface UseChangePasswordOptions {
  onSuccess?: () => void; // 成功后的额外回调
  onError?: (error: string) => void;
}

export function useChangePassword(options: UseChangePasswordOptions = {}) {
  const router = useRouter();
  const t = useTranslations("Settings");
  const [isLoading, setIsLoading] = useState(false);

  const changePassword = async (oldPassword: string, newPassword: string) => {
    if (isLoading) return false;

    setIsLoading(true);
    try {
      // 调用后端 API
      await authApi.changePassword({
        old_password: oldPassword,
        new_password: newPassword,
      });

      toast.success(t("password_changed_successfully"));
      options.onSuccess?.();

      // 3 秒后登出并跳转到登录页
      setTimeout(() => {
        toast.success(t("about_to_logout_relogin"));
        authApi.logout();
        router.push("/auth/login");
      }, 3000);

      return true;
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      toast.error(t(errorMsg));
      options.onError?.(errorMsg);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    changePassword,
    isLoading,
  };
}
