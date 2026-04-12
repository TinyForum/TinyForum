// components/question/FormActions.tsx
'use client';

import { useRouter } from 'next/navigation';

interface FormActionsProps {
  loading: boolean;
  onCancel?: () => void;
}

export function FormActions({ loading, onCancel }: FormActionsProps) {
  const router = useRouter();

  const handleCancel = () => {
    if (onCancel) {
      onCancel();
    } else {
      router.back();
    }
  };

  return (
    <div className="flex justify-end gap-3 pt-4 border-t">
      <button
        type="button"
        onClick={handleCancel}
        className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
      >
        取消
      </button>
      <button
        type="submit"
        disabled={loading}
        className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? (
          <span className="flex items-center gap-2">
            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            发布中...
          </span>
        ) : (
          '发布问题'
        )}
      </button>
    </div>
  );
}