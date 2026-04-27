// components/question/FormActions.tsx
"use client";

import { useRouter } from "next/navigation";
import {
  ArrowLeftIcon,
  PaperAirplaneIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
} from "@heroicons/react/24/outline";

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
  submitText = "发布问题",
  cancelText = "取消",
}: FormActionsProps) {
  const router = useRouter();

  const handleCancel = () => {
    if (onCancel) {
      onCancel();
    } else {
      // 如果有未保存的内容，提示确认
      if (isDirty) {
        const confirmed = confirm("有未保存的内容，确定要离开吗？");
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
          <span>请填写所有必填项后再发布</span>
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
          {cancelText}
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
              发布中...
            </>
          ) : (
            <>
              <PaperAirplaneIcon className="w-4 h-4" />
              {submitText}
            </>
          )}
        </button>
      </div>
    </div>
  );
}
