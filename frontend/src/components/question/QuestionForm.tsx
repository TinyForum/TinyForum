// components/question/QuestionForm.tsx
'use client';

import { useForm } from 'react-hook-form';
import { AskFormData } from '@/hooks/useQuestionForm';
import { 
  DocumentTextIcon, 
  TagIcon, 
  InformationCircleIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';

interface QuestionFormProps {
  register: ReturnType<typeof useForm<AskFormData>>['register'];
  errors: ReturnType<typeof useForm<AskFormData>>['formState']['errors'];
  content: string;
  setContent: (content: string) => void;
}

export function QuestionForm({ register, errors, content, setContent }: QuestionFormProps) {
  const contentLength = content.length;
  const isContentValid = contentLength > 0;
  const minContentLength = 20;

  return (
    <div className="space-y-6">
      {/* 标题 */}
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">
            标题 <span className="text-red-500">*</span>
          </span>
          <span className="label-text-alt text-base-content/60">
            5-100 字符
          </span>
        </label>
        <input
          {...register('title', { 
            required: '请输入标题', 
            minLength: { value: 5, message: '标题至少 5 个字符' },
            maxLength: { value: 100, message: '标题最多 100 个字符' }
          })}
          type="text"
          placeholder="例如：如何在 Next.js 中实现动态路由？"
          className={`input input-bordered w-full focus:input-primary ${
            errors.title ? 'input-error' : ''
          }`}
        />
        {errors.title ? (
          <label className="label">
            <span className="label-text-alt text-error flex items-center gap-1">
              <ExclamationTriangleIcon className="w-3 h-3" />
              {errors.title.message}
            </span>
          </label>
        ) : (
          <label className="label">
            <span className="label-text-alt text-base-content/40">
              一个好的标题能吸引更多回答者
            </span>
          </label>
        )}
      </div>

      {/* 内容 */}
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">
            问题描述 <span className="text-red-500">*</span>
          </span>
          <span className="label-text-alt text-base-content/60">
            {contentLength} / 建议 ≥{minContentLength}
          </span>
        </label>
        
        <div className="relative">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={12}
            placeholder={`详细描述你的问题...

💡 建议包含以下内容：
1. 你想要实现什么功能？
2. 你尝试过哪些方法？
3. 遇到了什么具体的错误？
4. 提供相关的代码片段`}
            className={`textarea textarea-bordered w-full font-mono text-sm resize-y ${
              !isContentValid && contentLength > 0 ? 'textarea-error' : ''
            } focus:textarea-primary`}
          />
          
          {/* 内容质量提示 */}
          {contentLength > 0 && contentLength < minContentLength && (
            <div className="absolute bottom-2 right-2">
              <div className="tooltip tooltip-left" data-tip="内容太简短，建议提供更多细节">
                <ExclamationTriangleIcon className="w-4 h-4 text-warning" />
              </div>
            </div>
          )}
        </div>
        
        {!isContentValid && contentLength === 0 && (
          <label className="label">
            <span className="label-text-alt text-error flex items-center gap-1">
              <ExclamationTriangleIcon className="w-3 h-3" />
              请输入问题内容
            </span>
          </label>
        )}
        
        {isContentValid && contentLength >= minContentLength && (
          <label className="label">
            <span className="label-text-alt text-success flex items-center gap-1">
              <CheckCircleIcon className="w-3 h-3" />
              内容完整，可以发布了
            </span>
          </label>
        )}
        
        <label className="label">
          <span className="label-text-alt text-base-content/40 flex items-center gap-1">
            <InformationCircleIcon className="w-3 h-3" />
            支持 Markdown 格式，代码块语法高亮
          </span>
        </label>
      </div>

      {/* 摘要 */}
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">
            问题摘要
          </span>
          <span className="label-text-alt text-base-content/60">
            可选，最多 500 字符
          </span>
        </label>
        <textarea
          {...register('summary')}
          rows={2}
          placeholder="简要描述问题（将显示在列表中，帮助用户快速了解）"
          className="textarea textarea-bordered w-full resize-none focus:textarea-primary"
        />
        <label className="label">
          <span className="label-text-alt text-base-content/40">
            如果留空，将自动从内容中提取摘要
          </span>
        </label>
      </div>

      {/* 提示卡片 */}
      <div className="alert alert-info bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800/30">
        <InformationCircleIcon className="w-4 h-4 text-blue-500" />
        <div className="text-xs text-blue-600 dark:text-blue-400">
          <p className="font-medium mb-1">📝 提问小贴士：</p>
          <ul className="space-y-0.5">
            <li>• 标题简洁明了，概括问题核心</li>
            <li>• 提供详细的背景和错误信息</li>
            <li>• 附上相关代码和运行环境</li>
            <li>• 设置悬赏积分可以吸引更多回答</li>
          </ul>
        </div>
      </div>
    </div>
  );
}