// components/question/AnswerForm.tsx
'use client';

import { postApi } from '@/lib/api';
import { useState } from 'react';
import { toast } from 'react-hot-toast';

interface AnswerFormProps {
  questionId: number;
  onSuccess: () => void;
  onCancel?: () => void;
}

export function AnswerForm({ questionId, onSuccess, onCancel }: AnswerFormProps) {
  const [content, setContent] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!content.trim()) {
      toast.error('请输入回答内容');
      return;
    }

    if (content.trim().length < 10) {
      toast.error('回答内容至少需要10个字符');
      return;
    }

    setSubmitting(true);
    try {
      const response = await postApi.createAnswer(questionId, { content: content.trim() });
      
      if (response.data.code === 200 || response.data.code === 201 || response.data.code === 0) {
        setContent('');
        toast.success('回答发布成功');
        onSuccess();
      } else {
        toast.error(response.data.message || '发布失败');
      }
    } catch (error: any) {
      console.error('发布回答失败:', error);
      const errorMsg = error.response?.data?.message || error.message || '发布失败，请稍后重试';
      toast.error(errorMsg);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">
        你的回答
      </h3>
      
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={6}
        placeholder={`写下你的回答...

建议：
1. 直接回答问题
2. 提供代码示例
3. 给出具体步骤
4. 附上相关文档链接`}
        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none resize-y"
        disabled={submitting}
      />
      
      <div className="flex justify-end gap-3 mt-4">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={submitting}
            className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
          >
            取消
          </button>
        )}
        <button
          onClick={handleSubmit}
          disabled={submitting || !content.trim()}
          className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {submitting ? (
            <span className="flex items-center gap-2">
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              发布中...
            </span>
          ) : (
            '发布回答'
          )}
        </button>
      </div>
      
      {/* 提示信息 */}
      <div className="mt-3 text-xs text-gray-400">
        <p>💡 提示：</p>
        <ul className="list-disc list-inside ml-2">
          <li>支持 Markdown 格式</li>
          <li>代码块请使用 ``` 包裹</li>
          <li>优质回答有机会获得悬赏积分</li>
        </ul>
      </div>
    </div>
  );
}