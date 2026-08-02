// hooks/useResetPassword.ts
import { useCallback } from "react";
import { authApi } from "@/shared/api/modules/auth";

// 忘记密码 Hook：发送重置密码邮件
export function useResetPassword() {
  // 发送重置密码链接（后端始终返回成功，不暴露邮箱是否存在）
  const forgotPassword = useCallback(async (email: string): Promise<boolean> => {
    const response = await authApi.forgotPassword({ email });

    // 后端始终返回 200，无论邮箱是否存在
    return response.data.code === 0 || response.status === 200;
  }, []);

  return { forgotPassword };
}
