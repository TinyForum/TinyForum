// components/question/QuestionForm.tsx
"use client";

import { useForm } from "react-hook-form";
import {
  InformationCircleIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
} from "@heroicons/react/24/outline";
import { useTranslations } from "next-intl";
import { AskFormData } from "@/features/qustion/hooks/useQuestionForm";

interface QuestionFormProps {
  register: ReturnType<typeof useForm<AskFormData>>["register"];
  errors: ReturnType<typeof useForm<AskFormData>>["formState"]["errors"];
  content: string;
  setContent: (content: string) => void;
}

export function QuestionForm({
  register,
  errors,
  content,
  setContent,
}: QuestionFormProps) {
  const t = useTranslations("Questions");
  const contentLength = content.length;
  const isContentValid = contentLength > 0;
  const minContentLength = 20;

  return (
    <div className="space-y-6">
      {/* 标题 */}
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">
            {t("title_label")} <span className="text-red-500">*</span>
          </span>
          <span className="label-text-alt text-base-content/60">
            {t("title_char_range")}
          </span>
        </label>
        <input
          {...register("title", {
            required: t("title_required"),
            minLength: { value: 5, message: t("title_min_5") },
            maxLength: { value: 100, message: t("title_max_100") },
          })}
          type="text"
          placeholder={t("title_placeholder")}
          className={`input input-bordered w-full focus:input-primary ${
            errors.title ? "input-error" : ""
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
              {t("title_hint")}
            </span>
          </label>
        )}
      </div>

      {/* 内容 */}
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">
            {t("content_label")} <span className="text-red-500">*</span>
          </span>
          <span className="label-text-alt text-base-content/60">
            {t("content_length_hint", { count: contentLength, min: minContentLength })}
          </span>
        </label>

        <div className="relative">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={12}
            placeholder={t("content_placeholder")}
            className={`textarea textarea-bordered w-full font-mono text-sm resize-y ${
              !isContentValid && contentLength > 0 ? "textarea-error" : ""
            } focus:textarea-primary`}
          />

          {/* 内容质量提示 */}
          {contentLength > 0 && contentLength < minContentLength && (
            <div className="absolute bottom-2 right-2">
              <div
                className="tooltip tooltip-left"
                data-tip={t("content_too_short_tip")}
              >
                <ExclamationTriangleIcon className="w-4 h-4 text-warning" />
              </div>
            </div>
          )}
        </div>

        {!isContentValid && contentLength === 0 && (
          <label className="label">
            <span className="label-text-alt text-error flex items-center gap-1">
              <ExclamationTriangleIcon className="w-3 h-3" />
              {t("content_required")}
            </span>
          </label>
        )}

        {isContentValid && contentLength >= minContentLength && (
          <label className="label">
            <span className="label-text-alt text-success flex items-center gap-1">
              <CheckCircleIcon className="w-3 h-3" />
              {t("content_complete")}
            </span>
          </label>
        )}

        <label className="label">
          <span className="label-text-alt text-base-content/40 flex items-center gap-1">
            <InformationCircleIcon className="w-3 h-3" />
            {t("markdown_hint")}
          </span>
        </label>
      </div>

      {/* 摘要 */}
      <div className="form-control w-full">
        <label className="label">
          <span className="label-text font-medium">{t("summary_label")}</span>
          <span className="label-text-alt text-base-content/60">
            {t("summary_optional_hint")}
          </span>
        </label>
        <textarea
          {...register("summary")}
          rows={2}
          placeholder={t("summary_placeholder")}
          className="textarea textarea-bordered w-full resize-none focus:textarea-primary"
        />
        <label className="label">
          <span className="label-text-alt text-base-content/40">
            {t("summary_auto_hint")}
          </span>
        </label>
      </div>

      {/* 提示卡片 */}
      <div className="alert alert-info bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800/30">
        <InformationCircleIcon className="w-4 h-4 text-blue-500" />
        <div className="text-xs text-blue-600 dark:text-blue-400">
          <p className="font-medium mb-1">{t("tips_title")}</p>
          <ul className="space-y-0.5">
            <li>{t("tip_title_concise")}</li>
            <li>{t("tip_background_details")}</li>
            <li>{t("tip_code_and_env")}</li>
            <li>{t("tip_reward")}</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
