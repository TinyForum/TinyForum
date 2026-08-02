// hooks/useLogin.ts
import { useState, useCallback } from "react";
import { toast } from "react-hot-toast";
import { useLoginStore } from "@/store/login";
import { useDeletionStatus } from "./useDeletionStatus";
import { useAccountActions } from "./useAccountActions";

interface LoginForm {
  email: string;
  password: string;
}

// 登录流程 Hook：封装登录、删除状态检查与强制退出
export function useLogin() {
  const { setEmail, setPassword, login: storeLogin } = useLoginStore();
  const [showRestoreDialog, setShowRestoreDialog] = useState(false);

  // 删除状态查询（登录成功后手动触发）
  const deletionQuery = useDeletionStatus(false);
  const { handleForceLogout } = useAccountActions();

  // 检查账户删除状态
  const checkDeletionStatus = useCallback(async (): Promise<boolean> => {
    try {
      const result = await deletionQuery.refetch();
      const status = result.data;

      // 添加安全检查，确保 status 存在
      if (!status) {
        return false;
      }

      if (status.is_deleted && status.can_restore) {
        setShowRestoreDialog(true);
        return true;
      } else if (status.is_deleted && !status.can_restore) {
        toast.error("您的账户已被永久删除，请联系管理员");
        await handleForceLogout();
        return false;
      }
      return false;
    } catch (error) {
      console.error("获取删除状态失败:", error);
      return false;
    }
  }, [deletionQuery, handleForceLogout]);

  // 提交登录：先写入表单再调用 store 登录，成功后检查删除状态
  const handleSubmitLogin = useCallback(
    async (data: LoginForm): Promise<{ success: boolean }> => {
      setEmail(data.email);
      setPassword(data.password);

      const result = await storeLogin();

      if (result.success) {
        // 检查账户删除状态
        await checkDeletionStatus();
      }
      return result;
    },
    [setEmail, setPassword, storeLogin, checkDeletionStatus],
  );

  return {
    showRestoreDialog,
    setShowRestoreDialog,
    handleSubmitLogin,
  };
}
