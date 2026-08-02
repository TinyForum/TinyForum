"use client";

import { useEffect, useState } from "react";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";
import RestoreDialog from "./RestoreDialog";
import { useDeletionStatus } from "../hooks/useDeletionStatus";
import type { DeletionStatus } from "../hooks/useDeletionStatus";
import { useAccountActions } from "../hooks/useAccountActions";

interface PostLoginHandlerProps {
  children?: React.ReactNode;
  onRestoreSuccess?: () => void;
  onDeleteSuccess?: () => void;
  redirectOnLogout?: string;
}

export default function PostLoginHandler({
  children,
  onRestoreSuccess,
  onDeleteSuccess,
  redirectOnLogout = "/",
}: PostLoginHandlerProps) {
  const { user, isAuthenticated } = useAuthStore();
  const [deletionStatus, setDeletionStatus] = useState<DeletionStatus | null>(
    null,
  );

  const {
    isDialogOpen,
    hasShown,
    isLoading,
    setHasShown,
    setIsDialogOpen,
    handleForceLogout,
    handleRestore,
    handlePermanentDelete,
    handleLogout,
  } = useAccountActions({ onRestoreSuccess, onDeleteSuccess, redirectOnLogout });

  // 查询：拉取账户删除状态（已认证且未处理时自动获取）
  const { data: deletionStatusQuery } = useDeletionStatus(
    isAuthenticated && !!user && !hasShown,
  );

  // 根据删除状态驱动恢复对话框 / 强制退出（仅触发一次）
  useEffect(() => {
    if (!deletionStatusQuery || hasShown) return;
    setHasShown(true);

    setDeletionStatus(deletionStatusQuery);

    if (deletionStatusQuery.is_deleted && deletionStatusQuery.can_restore) {
      // 显示恢复对话框
      setIsDialogOpen(true);
    } else if (
      deletionStatusQuery.is_deleted &&
      !deletionStatusQuery.can_restore
    ) {
      // 账户已永久删除，强制登出
      toast.error("您的账户已被永久删除，请联系管理员");
      handleForceLogout();
    }
  }, [
    deletionStatusQuery,
    hasShown,
    setHasShown,
    setIsDialogOpen,
    handleForceLogout,
  ]);

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
