// features/settings/hooks/useChangePassword.ts
import { useState } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { authApi } from "@/shared/api";
import { getErrorMessage } from "@/shared/lib/utils";

interface UseChangePasswordOptions {
  onSuccess?: () => void; // 成功后的额外回调
  onError?: (error: string) => void;
}

export function useChangePassword(options: UseChangePasswordOptions = {}) {
  const router = useRouter();
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

      toast.success("密码修改成功");
      options.onSuccess?.();

      // 3 秒后登出并跳转到登录页
      setTimeout(() => {
        toast.success("即将退出，请使用新密码重新登录");
        authApi.logout();
        router.push("/auth/login");
      }, 3000);

      return true;
    } catch (err: unknown) {
      const errorMsg = getErrorMessage(err);
      toast.error(errorMsg);
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
