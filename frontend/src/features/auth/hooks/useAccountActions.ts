// hooks/useAccountActions.ts
import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { authApi } from "@/shared/api/modules/auth";
import { useAuthStore } from "@/store/auth";

// 错误响应类型
interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

interface UseAccountActionsOptions {
  redirectOnLogout?: string;
  onRestoreSuccess?: () => void;
  onDeleteSuccess?: () => void;
}

// 账户操作 Hook：强制退出 / 恢复账户 / 永久删除 / 退出登录
export function useAccountActions(options?: UseAccountActionsOptions) {
  const { redirectOnLogout = "/", onRestoreSuccess, onDeleteSuccess } =
    options ?? {};
  const router = useRouter();
  const t = useTranslations("Auth");
  const { logout } = useAuthStore();

  const [isDialogOpen, setIsDialogOpen] = useState<boolean>(false);
  const [hasShown, setHasShown] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  // 强制退出（用于永久删除的账户）
  const handleForceLogout = useCallback(async (): Promise<void> => {
    await logout();
    router.push(redirectOnLogout);
    router.refresh();
  }, [logout, redirectOnLogout, router]);

  // 恢复账户
  const handleRestore = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    const loadingToast = toast.loading(t("restoring_account"));

    try {
      await authApi.cancelDeletion();
      toast.success(t("account_restored"), { id: loadingToast });

      // 刷新用户状态
      await useAuthStore.getState().refreshUser();

      // 关闭对话框
      setIsDialogOpen(false);
      setHasShown(false);

      // 执行回调
      if (onRestoreSuccess) {
        onRestoreSuccess();
      } else {
        // 默认刷新页面
        setTimeout(() => {
          window.location.reload();
        }, 1500);
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(t(error?.response?.data?.message || "restore_failed"), {
        id: loadingToast,
      });
    } finally {
      setIsLoading(false);
    }
  }, [onRestoreSuccess, t]);

  // 立即永久删除
  const handlePermanentDelete = useCallback(async (): Promise<void> => {
    const confirmDelete = window.confirm(t("permanent_delete_warning"));

    if (!confirmDelete) return;

    setIsLoading(true);
    const loadingToast = toast.loading(t("permanently_deleting_account"));

    try {
      await authApi.confirmDeletion({ confirm: "PERMANENT_DELETE" });
      toast.success(t("account_permanently_deleted"), { id: loadingToast });

      // 退出登录
      await logout();

      // 关闭对话框
      setIsDialogOpen(false);
      setHasShown(false);

      // 执行回调
      if (onDeleteSuccess) {
        onDeleteSuccess();
      }

      // 跳转到首页
      router.push(redirectOnLogout);
      router.refresh();
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      toast.error(t(error?.response?.data?.message || "delete_failed"), {
        id: loadingToast,
      });
    } finally {
      setIsLoading(false);
    }
  }, [logout, onDeleteSuccess, redirectOnLogout, router, t]);

  // 退出登录
  const handleLogout = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    const loadingToast = toast.loading(t("logging_out"));

    try {
      await logout();
      setIsDialogOpen(false);
      setHasShown(false);
      toast.success(t("logout_success"), { id: loadingToast });
      router.push(redirectOnLogout);
      router.refresh();
    } catch {
      // 移除未使用的 error 参数
      toast.error(t("logout_failed"), { id: loadingToast });
    } finally {
      setIsLoading(false);
    }
  }, [logout, redirectOnLogout, router, t]);

  return {
    isDialogOpen,
    setIsDialogOpen,
    hasShown,
    setHasShown,
    isLoading,
    setIsLoading,
    handleForceLogout,
    handleRestore,
    handlePermanentDelete,
    handleLogout,
  };
}
