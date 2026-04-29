// components/question/AnswerForm.tsx
"use client";

import { questionApi } from "@/shared/api";
import { useState } from "react";
import { toast } from "react-hot-toast";
import {
  PaperAirplaneIcon,
  XMarkIcon,
  InformationCircleIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  SparklesIcon,
} from "@heroicons/react/24/outline";
import { ApiResponse } from "@/shared/api/types";

interface AnswerFormProps {
  questionId: number;
  questionAuthorId?: number;
  rewardScore?: number;
  hasAccepted?: boolean;
  onSuccess: () => void;
  onCancel?: () => void;
}

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

interface CreateAnswerResponse {
  id: number;
  content: string;
  created_at: string;
}

export function AnswerForm({
  questionId,
  questionAuthorId,
  rewardScore = 0,
  hasAccepted = false,
  onSuccess,
  onCancel,
}: AnswerFormProps) {
  const [content, setContent] = useState<string>("");
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [showPreview, setShowPreview] = useState<boolean>(false);

  const contentLength: number = content.trim().length;
  const isContentValid: boolean = contentLength >= 10;
  const isContentTooLong: boolean = contentLength > 50000;

  const handleSubmit = async (): Promise<void> => {
    if (!content.trim()) {
      toast.error("请输入回答内容");
      return;
    }

    if (contentLength < 10) {
      toast.error("回答内容至少需要 10 个字符");
      return;
    }

    if (contentLength > 50000) {
      toast.error("回答内容不能超过 50000 个字符");
      return;
    }

    setSubmitting(true);
    try {
      const response: { data: ApiResponse<CreateAnswerResponse> } =
        await questionApi.createAnswer(questionId, {
          content: content.trim(),
        });

      // 统一使用 code === 0 判断成功
      if (response.data.code === 0) {
        setContent("");
        setShowPreview(false);
        toast.success("回答发布成功");

        // 如果有悬赏，显示额外提示
        if (rewardScore > 0 && questionAuthorId) {
          toast.success(`💡 回答被采纳可获得 ${rewardScore} 积分悬赏`, {
            duration: 5000,
          });
        }

        onSuccess();
      } else {
        toast.error(response.data.message || "发布失败");
      }
    } catch (err: unknown) {
      console.error("发布回答失败:", err);
      const error = err as ErrorResponse;
      const errorMsg =
        error.response?.data?.message ||
        error.message ||
        "发布失败，请稍后重试";
      toast.error(errorMsg);
    } finally {
      setSubmitting(false);
    }
  };

  // 如果问题已有采纳答案，显示提示
  if (hasAccepted) {
    return (
      <div className="card bg-base-100 shadow-md border border-base-200">
        <div className="card-body p-6 text-center">
          <div className="w-12 h-12 mx-auto mb-3 rounded-full bg-success/10 flex items-center justify-center">
            <CheckCircleIcon className="w-6 h-6 text-success" />
          </div>
          <h3 className="text-lg font-semibold text-base-content mb-2">
            问题已解决
          </h3>
          <p className="text-sm text-base-content/60">
            这个问题已经有采纳的回答了，你可以查看其他问题或发布新问题
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="card bg-base-100 shadow-md border border-base-200 overflow-hidden">
      {/* 悬赏提示 */}
      {rewardScore > 0 && (
        <div className="bg-gradient-to-r from-amber-500 to-orange-500 px-4 py-2">
          <div className="flex items-center gap-2">
            <SparklesIcon className="w-4 h-4 text-white" />
            <span className="text-sm font-medium text-white">
              本问题悬赏 {rewardScore} 积分，优质回答有机会获得悬赏！
            </span>
          </div>
        </div>
      )}

      <div className="card-body p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
              <PaperAirplaneIcon className="w-4 h-4 text-primary" />
            </div>
            <h3 className="text-lg font-semibold text-base-content">
              你的回答
            </h3>
          </div>

          {/* 预览切换 */}
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => setShowPreview(!showPreview)}
              className="btn btn-ghost btn-xs"
              disabled={submitting}
            >
              {showPreview ? "编辑" : "预览"}
            </button>
          </div>
        </div>

        {/* 编辑器 / 预览 */}
        {showPreview ? (
          <div className="min-h-[200px] p-4 bg-base-200 rounded-lg prose prose-sm max-w-none">
            {content ? (
              <div
                dangerouslySetInnerHTML={{
                  __html: content.replace(/\n/g, "<br/>"),
                }}
              />
            ) : (
              <p className="text-base-content/40 text-center">暂无内容</p>
            )}
          </div>
        ) : (
          <>
            <textarea
              value={content}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                setContent(e.target.value)
              }
              rows={8}
              placeholder={`写下你的回答...

💡 优质回答小贴士：
• 直接回答问题核心
• 提供可运行的代码示例
• 给出具体的操作步骤
• 附上相关文档或参考资料
• 使用 Markdown 格式化内容`}
              className={`textarea w-full resize-y font-mono text-sm ${
                !isContentValid && contentLength > 0
                  ? "textarea-error"
                  : isContentTooLong
                    ? "textarea-error"
                    : "textarea-bordered"
              } focus:textarea-primary`}
              disabled={submitting}
            />

            {/* 字数统计 */}
            <div className="flex justify-between items-center mt-2">
              <div className="flex items-center gap-2">
                {contentLength > 0 && !isContentValid && (
                  <span className="text-xs text-warning flex items-center gap-1">
                    <ExclamationTriangleIcon className="w-3 h-3" />
                    至少需要 10 个字符
                  </span>
                )}
                {isContentValid && !isContentTooLong && (
                  <span className="text-xs text-success flex items-center gap-1">
                    <CheckCircleIcon className="w-3 h-3" />
                    内容完整
                  </span>
                )}
                {isContentTooLong && (
                  <span className="text-xs text-error flex items-center gap-1">
                    <ExclamationTriangleIcon className="w-3 h-3" />
                    内容过长，请精简到 50000 字符以内
                  </span>
                )}
              </div>
              <span
                className={`text-xs ${
                  contentLength > 45000
                    ? "text-error"
                    : contentLength > 0
                      ? "text-base-content/60"
                      : "text-base-content/40"
                }`}
              >
                {contentLength} / 50000
              </span>
            </div>
          </>
        )}

        {/* 按钮区域 */}
        <div className="flex justify-end gap-3 mt-4">
          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              disabled={submitting}
              className="btn btn-ghost gap-2"
            >
              <XMarkIcon className="w-4 h-4" />
              取消
            </button>
          )}
          <button
            onClick={handleSubmit}
            disabled={submitting || !isContentValid || isContentTooLong}
            className="btn btn-primary gap-2 min-w-[120px] shadow-md hover:shadow-lg transition-all"
          >
            {submitting ? (
              <>
                <span className="loading loading-spinner loading-sm"></span>
                发布中...
              </>
            ) : (
              <>
                <PaperAirplaneIcon className="w-4 h-4" />
                发布回答
              </>
            )}
          </button>
        </div>

        {/* 提示信息 */}
        <div className="mt-4 p-3 bg-base-200/50 rounded-lg">
          <div className="flex items-start gap-2">
            <InformationCircleIcon className="w-4 h-4 text-base-content/40 shrink-0 mt-0.5" />
            <div className="text-xs text-base-content/60 space-y-1">
              <p className="font-medium">💡 回答提示：</p>
              <ul className="list-disc list-inside space-y-0.5 ml-2">
                <li>支持 Markdown 格式，代码块请使用 ``` 包裹</li>
                <li>优质回答有机会获得悬赏积分和采纳</li>
                <li>尊重他人，友善发言</li>
                <li>回答被采纳后，悬赏积分将自动发放</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
