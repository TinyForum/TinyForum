"use client";

import { useDeleteAccountStore } from "@/store/delete";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";

export default function DangerZone() {
  const t = useTranslations("Settings");
  const router = useRouter();
  const {
    isModalOpen,
    isLoading: isDeleting,
    error,
    confirmText,
    setConfirmText,
    setModalOpen,
    deleteAccount,
    resetForm,
  } = useDeleteAccountStore();

  const handleDeleteAccount = async (): Promise<void> => {
    const result = await deleteAccount();
    if (result.success) {
      router.push("/");
      router.refresh();
    }
  };

  const handleCloseModal = (): void => {
    setModalOpen(false);
    resetForm();
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">{t("danger_zone")}</h1>
      <div className="card bg-red-50 dark:bg-red-900/10 border border-red-200 shadow-sm p-6">
        <div className="flex items-center justify-between flex-wrap gap-4">
          <div>
            <h3 className="text-lg font-semibold text-red-700 dark:text-red-400">
              {t("delete_account")}
            </h3>
            <p className="text-sm text-red-600 dark:text-red-300 mt-1">
              {t("delete_account_warning")}
            </p>
          </div>
          <button
            onClick={() => setModalOpen(true)}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
          >
            {t("delete_account")}
          </button>
        </div>
      </div>

      {/* 确认模态框 */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-red-600 dark:text-red-400 mb-4">
              {t("delete_account_warning_title")}
            </h2>
            <p className="text-gray-700 dark:text-gray-300 mb-4">
              {t("delete_account_warning_desc")}
            </p>
            <ul className="list-disc list-inside text-sm text-gray-600 dark:text-gray-400 mb-4 space-y-1">
              <li>{t("delete_item_profile")}</li>
              <li>{t("delete_item_history")}</li>
              <li>{t("delete_item_subscriptions")}</li>
            </ul>
            <p className="text-sm text-red-600 dark:text-red-400 font-semibold mb-4">
              {t("delete_irreversible")}
            </p>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {t.rich("enter_delete_to_confirm", {
                  code: (chunks) => (
                    <code className="px-1 py-0.5 bg-gray-100 dark:bg-gray-700 rounded">
                      {chunks}
                    </code>
                  ),
                })}
              </label>
              <input
                type="text"
                value={confirmText}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setConfirmText(e.target.value)
                }
                placeholder={t("delete_placeholder")}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500 dark:bg-gray-700 dark:text-white"
                autoFocus
              />
            </div>

            {error && <p className="text-red-600 text-sm mb-4">{t(error)}</p>}

            <div className="flex gap-3 justify-end">
              <button
                onClick={handleCloseModal}
                className="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 rounded-lg transition-colors"
              >
                {t("cancel")}
              </button>
              <button
                onClick={handleDeleteAccount}
                disabled={isDeleting}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-400 text-white rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
              >
                {isDeleting ? t("deleting") : t("confirm_delete")}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
