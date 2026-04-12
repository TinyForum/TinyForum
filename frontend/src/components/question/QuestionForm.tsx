// components/question/QuestionForm.tsx
'use client';

import { useForm } from 'react-hook-form';
import { AskFormData } from '@/hooks/useQuestionForm';

interface QuestionFormProps {
  register: ReturnType<typeof useForm<AskFormData>>['register'];
  errors: ReturnType<typeof useForm<AskFormData>>['formState']['errors'];
  content: string;
  setContent: (content: string) => void;
}

export function QuestionForm({ register, errors, content, setContent }: QuestionFormProps) {
  return (
    <div className="space-y-6">
      {/* 标题 */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          标题 <span className="text-red-500">*</span>
        </label>
        <input
          {...register('title', { 
            required: '请输入标题', 
            minLength: { value: 5, message: '标题至少5个字符' },
            maxLength: { value: 100, message: '标题最多100个字符' }
          })}
          type="text"
          placeholder="例如：如何在 Next.js 中实现动态路由？"
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
        />
        {errors.title && (
          <p className="mt-1 text-sm text-red-500">{errors.title.message}</p>
        )}
      </div>

      {/* 内容 */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          问题描述 <span className="text-red-500">*</span>
        </label>
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={12}
          placeholder={`详细描述你的问题...

建议包含以下内容：
1. 你想要实现什么功能？
2. 你尝试过哪些方法？
3. 遇到了什么具体的错误？`}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none font-mono text-sm"
        />
        {!content && (
          <p className="mt-1 text-sm text-red-500">请输入问题内容</p>
        )}
        <p className="mt-1 text-sm text-gray-400">
          支持 Markdown 格式，{content.length} 字符
        </p>
      </div>

      {/* 摘要 */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          问题摘要
        </label>
        <textarea
          {...register('summary')}
          rows={2}
          placeholder="简要描述问题（可选，将显示在列表中）"
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
        />
        <p className="mt-1 text-sm text-gray-400">
          最多500个字符
        </p>
      </div>
    </div>
  );
}