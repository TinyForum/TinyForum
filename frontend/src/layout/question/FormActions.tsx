// components/question/FormActions.tsx
"use client";

import { useRouter } from "next/navigation";
import {
  ArrowLeftIcon,
  PaperAirplaneIcon,
  ExclamationTriangleIcon,
} from "@heroicons/react/24/outline";
import { useTranslations } from "next-intl";

interface FormActionsProps {
  loading: boolean;
  onCancel?: () => void;
  isDirty?: boolean;
  isValid?: boolean;
  submitText?: string;
  cancelText?: string;
}

export function FormActions({
  loading,
  onCancel,
  isDirty = true,
  isValid = true,
  submitText,
  cancelText,
}: FormActionsProps) {
  const router = useRouter();
  const t = useTranslations("Questions");

  const finalSubmitText = submitText || t("publish_question");
  const finalCancelText = cancelText || t("cancel");

  const handleCancel = () => {
    if (onCancel) {
      onCancel();
    } else {
      // 如果有未保存的内容，提示确认
      if (isDirty) {
        const confirmed = confirm(t("unsaved_confirm"));
        if (!confirmed) return;
      }
      router.back();
    }
  };

  const isDisabled = loading || !isDirty || !isValid;

  return (
    <div className="flex flex-col sm:flex-row justify-end gap-3 pt-6 border-t border-base-200">
      {/* 提示信息 */}
      {isDirty && !isValid && (
        <div className="flex-1 flex items-center gap-2 text-sm text-warning">
          <ExclamationTriangleIcon className="w-4 h-4" />
          <span>{t("fill_required")}</span>
        </div>
      )}

      {/* {isDirty && isValid && !loading && (
        <div className="flex-1 flex items-center gap-2 text-sm text-success">
          <CheckCircleIcon className="w-4 h-4" />
          <span>表单验证通过，可以发布了</span>
        </div>
      )} */}

      <div className="flex gap-3">
        {/* 取消按钮 */}
        <button
          type="button"
          onClick={handleCancel}
          className="btn btn-ghost gap-2"
          disabled={loading}
        >
          <ArrowLeftIcon className="w-4 h-4" />
          {finalCancelText}
        </button>

        {/* 提交按钮 */}
        <button
          type="submit"
          disabled={isDisabled}
          className={`btn btn-primary gap-2 min-w-[120px] ${
            isDisabled ? "opacity-50 cursor-not-allowed" : ""
          }`}
        >
          {loading ? (
            <>
              <span className="loading loading-spinner loading-sm"></span>
              {t("publishing")}
            </>
          ) : (
            <>
              <PaperAirplaneIcon className="w-4 h-4" />
              {finalSubmitText}
            </>
          )}
        </button>
      </div>
    </div>
  );
}
