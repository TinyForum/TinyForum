"use client";

import { AlertTriangle, RefreshCw, LogOut, Trash2 } from "lucide-react";
import { useTranslations } from "next-intl";

interface DeletionStatus {
  is_deleted: boolean;
  deleted_at?: string;
  can_restore: boolean;
  remaining_days?: number;
}

interface RestoreDialogProps {
  isOpen: boolean;
  deletionStatus: DeletionStatus | null;
  onRestore: () => void;
  onPermanentDelete: () => void;
  onLogout: () => void;
  isLoading: boolean;
}

export default function RestoreDialog({
  isOpen,
  deletionStatus,
  onRestore,
  onPermanentDelete,
  onLogout,
  isLoading,
}: RestoreDialogProps) {
  const t = useTranslations("Auth");

  if (!isOpen || !deletionStatus) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6 animate-in fade-in zoom-in duration-200">
        <div className="text-center mb-4">
          <div className="w-16 h-16 rounded-full bg-warning/20 flex items-center justify-center mx-auto mb-4">
            <AlertTriangle className="w-8 h-8 text-warning" />
          </div>
          <h2 className="text-xl font-bold text-warning">{t("account_marked_for_deletion")}</h2>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            {t.rich("account_deletion_notice", {
              date: deletionStatus.deleted_at
                ? new Date(deletionStatus.deleted_at).toLocaleDateString()
                : "",
              days: deletionStatus.remaining_days ?? 0,
              strong: (chunks) => (
                <span className="text-warning font-bold">{chunks}</span>
              ),
            })}
          </p>
          <p className="text-sm text-red-600 mt-2 font-semibold">
            {t("expired_permanent_delete_note")}
          </p>
        </div>

        <div className="flex flex-col gap-3 mt-6">
          <button
            onClick={onRestore}
            disabled={isLoading}
            className="w-full px-4 py-2 bg-success hover:bg-success/90 disabled:bg-success/50 text-white rounded-lg transition-colors flex items-center justify-center gap-2"
          >
            {isLoading ? (
              <span className="loading loading-spinner loading-sm" />
            ) : (
              <RefreshCw className="w-4 h-4" />
            )}
            {t("restore_account")}
          </button>

          <button
            onClick={onPermanentDelete}
            disabled={isLoading}
            className="w-full px-4 py-2 bg-error hover:bg-error/90 disabled:bg-error/50 text-white rounded-lg transition-colors flex items-center justify-center gap-2"
          >
            <Trash2 className="w-4 h-4" />
            {t("permanent_delete_now")}
          </button>

          <button
            onClick={onLogout}
            disabled={isLoading}
            className="w-full px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 rounded-lg transition-colors flex items-center justify-center gap-2"
          >
            <LogOut className="w-4 h-4" />
            {t("logout_for_now")}
          </button>
        </div>
      </div>
    </div>
  );
}
