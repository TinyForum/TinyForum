// components/question/AnswerForm.tsx
"use client";

import { useState } from "react";
import { toast } from "react-hot-toast";
import { useTranslations } from "next-intl";
import { useAnswerSubmit } from "@/features/question/hooks/useAnswerSubmit";
import { sanitizeHtml } from "@/shared/lib/utils";
import {
  PaperAirplaneIcon,
  XMarkIcon,
  InformationCircleIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  SparklesIcon,
} from "@heroicons/react/24/outline";

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
  const submitAnswer = useAnswerSubmit(questionId);
  const t = useTranslations("Questions");

  const contentLength: number = content.trim().length;
  const isContentValid: boolean = contentLength >= 10;
  const isContentTooLong: boolean = contentLength > 50000;

  const handleSubmit = async (): Promise<void> => {
    if (!content.trim()) {
      toast.error(t("answer_content_required"));
      return;
    }

    if (contentLength < 10) {
      toast.error(t("answer_min_length"));
      return;
    }

    if (contentLength > 50000) {
      toast.error(t("answer_max_length"));
      return;
    }

    setSubmitting(true);
    try {
      await submitAnswer.mutateAsync(content.trim());

      // 统一使用 code === 0 判断成功
      setContent("");
      setShowPreview(false);
      toast.success(t("answer_published"));

      // 如果有悬赏，显示额外提示
      if (rewardScore > 0 && questionAuthorId) {
        toast.success(t("answer_reward_tip", { score: rewardScore }), {
          duration: 5000,
        });
      }

      onSuccess();
    } catch (err: unknown) {
      console.error("发布回答失败:", err);
      const error = err as ErrorResponse;
      const errorMsg =
        error.response?.data?.message ||
        error.message ||
        t("publish_failed_retry");
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
            {t("question_resolved")}
          </h3>
          <p className="text-sm text-base-content/60">
            {t("question_resolved_desc")}
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
              {t("reward_banner", { score: rewardScore })}
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
              {t("your_answer")}
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
              {showPreview ? t("edit") : t("preview")}
            </button>
          </div>
        </div>

        {/* 编辑器 / 预览 */}
        {showPreview ? (
          <div className="min-h-[200px] p-4 bg-base-200 rounded-lg prose prose-sm max-w-none">
            {content ? (
              <div
                dangerouslySetInnerHTML={{
                  __html: sanitizeHtml(content.replace(/\n/g, "<br/>")),
                }}
              />
            ) : (
              <p className="text-base-content/40 text-center">{t("no_content")}</p>
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
              placeholder={t("answer_placeholder")}
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
                    {t("answer_min_hint")}
                  </span>
                )}
                {isContentValid && !isContentTooLong && (
                  <span className="text-xs text-success flex items-center gap-1">
                    <CheckCircleIcon className="w-3 h-3" />
                    {t("answer_complete")}
                  </span>
                )}
                {isContentTooLong && (
                  <span className="text-xs text-error flex items-center gap-1">
                    <ExclamationTriangleIcon className="w-3 h-3" />
                    {t("answer_too_long")}
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
              {t("cancel")}
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
                {t("publishing")}
              </>
            ) : (
              <>
                <PaperAirplaneIcon className="w-4 h-4" />
                {t("post_answer")}
              </>
            )}
          </button>
        </div>

        {/* 提示信息 */}
        <div className="mt-4 p-3 bg-base-200/50 rounded-lg">
          <div className="flex items-start gap-2">
            <InformationCircleIcon className="w-4 h-4 text-base-content/40 shrink-0 mt-0.5" />
            <div className="text-xs text-base-content/60 space-y-1">
              <p className="font-medium">{t("answer_tips_title")}</p>
              <ul className="list-disc list-inside space-y-0.5 ml-2">
                <li>{t("answer_tip_markdown")}</li>
                <li>{t("answer_tip_reward")}</li>
                <li>{t("answer_tip_respect")}</li>
                <li>{t("answer_tip_accepted")}</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
