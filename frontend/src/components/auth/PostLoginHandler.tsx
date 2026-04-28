"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import { authApi } from "@/lib/api";
import toast from "react-hot-toast";
import RestoreDialog from "./RestoreDialog";
import { ApiResponse } from "@/lib/api/types";

interface DeletionStatus {
  is_deleted: boolean;
  deleted_at?: string;
  can_restore: boolean;
  remaining_days?: number;
}

interface PostLoginHandlerProps {
  children?: React.ReactNode;
  onRestoreSuccess?: () => void;
  onDeleteSuccess?: () => void;
  redirectOnLogout?: string;
}

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

export default function PostLoginHandler({
  children,
  onRestoreSuccess,
  onDeleteSuccess,
  redirectOnLogout = "/",
}: PostLoginHandlerProps) {
  const router = useRouter();
  const { user, isAuthenticated, logout } = useAuthStore();
  const [hasShown, setHasShown] = useState<boolean>(false);
  const [deletionStatus, setDeletionStatus] = useState<DeletionStatus | null>(
    null,
  );
  const [isDialogOpen, setIsDialogOpen] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  // 强制退出（用于永久删除的账户）
  const handleForceLogout = useCallback(async (): Promise<void> => {
    await authApi.logout();
    logout();
    router.push(redirectOnLogout);
    router.refresh();
  }, [logout, redirectOnLogout, router]);

  // 检查删除状态
  const checkDeletionStatus = useCallback(async (): Promise<void> => {
    try {
      const response: { data: ApiResponse<DeletionStatus> } = 
        await authApi.getDeletionStatus();
      const status = response.data.data;
      
      if (status) {
        setDeletionStatus(status);

        if (status.is_deleted && status.can_restore) {
          // 显示恢复对话框
          setIsDialogOpen(true);
          setHasShown(true);
        } else if (status.is_deleted && !status.can_restore) {
          // 账户已永久删除，强制登出
          toast.error("您的账户已被永久删除，请联系管理员");
          await handleForceLogout();
        }
      }
    } catch (error: unknown) {
      console.error("获取删除状态失败:", error);
    }
  }, [handleForceLogout]);

  useEffect(() => {
    // 只在已认证且未显示过提示时执行
    if (isAuthenticated && user && !hasShown) {
      checkDeletionStatus();
    }
  }, [isAuthenticated, user, hasShown, checkDeletionStatus]);

  // 恢复账户
  const handleRestore = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    const loadingToast = toast.loading("正在恢复账户...");

    try {
      await authApi.cancelDeletion();
      toast.success("账户已成功恢复！", { id: loadingToast });

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
      toast.error(error?.response?.data?.message || "恢复失败，请重试", {
        id: loadingToast,
      });
    } finally {
      setIsLoading(false);
    }
  }, [onRestoreSuccess]);

  // 立即永久删除
  const handlePermanentDelete = useCallback(async (): Promise<void> => {
    const confirmDelete = window.confirm(
      "⚠️ 警告：此操作将永久删除您的账户，所有数据无法恢复！\n\n确认永久删除吗？",
    );

    if (!confirmDelete) return;

    setIsLoading(true);
    const loadingToast = toast.loading("正在永久删除账户...");

    try {
      await authApi.confirmDeletion({ confirm: "PERMANENT_DELETE" });
      toast.success("账户已永久删除", { id: loadingToast });

      // 退出登录
      await authApi.logout();
      logout();

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
      toast.error(error?.response?.data?.message || "删除失败，请重试", {
        id: loadingToast,
      });
    } finally {
      setIsLoading(false);
    }
  }, [logout, onDeleteSuccess, redirectOnLogout, router]);

  // 退出登录
  const handleLogout = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    const loadingToast = toast.loading("正在退出登录...");

    try {
      await authApi.logout();
      logout();
      setIsDialogOpen(false);
      setHasShown(false);
      toast.success("您已退出登录", { id: loadingToast });
      router.push(redirectOnLogout);
      router.refresh();
    } catch {
      // 移除未使用的 error 参数
      toast.error("退出登录失败", { id: loadingToast });
    } finally {
      setIsLoading(false);
    }
  }, [logout, redirectOnLogout, router]);

  return (
    <>
      {children}
      <RestoreDialog
        isOpen={isDialogOpen}
        deletionStatus={deletionStatus}
        onRestore={handleRestore}
        onPermanentDelete={handlePermanentDelete}
        onLogout={handleLogout}
        isLoading={isLoading}
      />
    </>
  );
}