// ─── CreateBoardInline ─────────────────────────────────────────────────────────

import { useAuthStore } from "@/store/auth";
import { useTranslations } from "next-intl";
import { ExclamationCircleIcon } from "@heroicons/react/24/outline";
import Link from "next/link";
import { ArrowLeftIcon, FolderPlusIcon, MailIcon } from "lucide-react";
import { useState } from "react";
import { Board } from "@/shared/api/types/board.model";
import { useBoardCreate } from "@/features/boards/hooks/useBoardCreate";
import { useBoardCreationApply } from "@/features/boards/hooks/useBoardCreationApply";

interface CreateBoardInlineProps {
  slug: string;
  onCreated: (board: Board) => void;
}

interface FormData {
  name: string;
  slug: string;
  description: string;
}

interface ErrorResponse {
  response?: {
    data?: {
      message?: string;
    };
  };
  message?: string;
}

export function CreateBoardInline({ slug, onCreated }: CreateBoardInlineProps) {
  const { user } = useAuthStore();
  const t = useTranslations("Board");
  const [form, setForm] = useState<FormData>({
    name: "",
    slug: slug,
    description: "",
  });
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const createBoard = useBoardCreate();
  const applyBoardCreation = useBoardCreationApply();

  // 申请表单相关状态
  const [showApplyForm, setShowApplyForm] = useState<boolean>(false);
  const [applyReason, setApplyReason] = useState<string>("");
  const [applySubmitting, setApplySubmitting] = useState<boolean>(false);
  const [applySuccess, setApplySuccess] = useState<boolean>(false);
  const [applyError, setApplyError] = useState<string>("");

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ): void => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
    setError("");
  };

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!form.name.trim()) {
      setError("validation_board_name_required");
      return;
    }
    if (!form.slug.trim()) {
      setError("validation_board_slug_required");
      return;
    }
    if (!/^[a-z0-9-]+$/.test(form.slug)) {
      setError("validation_slug_format");
      return;
    }

    setSubmitting(true);
    try {
      const board = await createBoard.mutateAsync({
        name: form.name.trim(),
        slug: form.slug.trim(),
        description: form.description.trim(),
      });

      // 修复：检查返回的板块数据是否存在
      if (board) {
        onCreated(board);
      } else {
        setError("create_success_no_data");
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      setError(
        error.response?.data?.message ||
          (error instanceof Error ? error.message : "") ||
          "create_failed_retry",
      );
    } finally {
      setSubmitting(false);
    }
  };

  // 处理申请提交
  const handleApplySubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!form.name.trim()) {
      setApplyError("validation_board_name_required");
      return;
    }
    if (!applyReason.trim()) {
      setApplyError("validation_reason_required");
      return;
    }

    setApplySubmitting(true);
    setApplyError("");

    try {
      await applyBoardCreation.mutateAsync({
        name: form.name.trim(),
        slug: form.slug.trim(),
        description: form.description.trim(),
        reason: applyReason.trim(),
      });

      setApplySuccess(true);
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      setApplyError(
        error.response?.data?.message ||
          (error instanceof Error ? error.message : "") ||
          "submit_failed_retry",
      );
    } finally {
      setApplySubmitting(false);
    }
  };

  // 判断是否为管理员
  const isAdmin = user?.role === "admin" || user?.role === "super_admin";

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="w-full max-w-lg">
        {/* Icon + heading */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-blue-50 dark:bg-blue-900/30 mb-4">
            <ExclamationCircleIcon className="w-8 h-8 text-blue-500" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t("board_not_exists")}
          </h1>
          <p className="mt-2 text-gray-500 dark:text-gray-400">
            <span className="font-mono bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded text-sm text-blue-600 dark:text-blue-400">
              /{slug}
            </span>{" "}
             {t("board_not_created")}
          </p>
        </div>

        {!user ? (
          // 未登录状态
          <div className="text-center space-y-3">
            <p className="text-gray-500">{t("login_to_create_board")}</p>
            <Link
              href={`/auth/login?redirect=/boards/${slug}`}
              className="inline-block px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-medium transition-colors"
            >
               {t("go_login")}
            </Link>
            <div className="pt-2">
              <Link
                href="/boards"
                className="text-sm text-gray-400 hover:text-blue-500 transition-colors"
              >
                 {t("back_to_board_list")}
              </Link>
            </div>
          </div>
        ) : isAdmin ? (
          // 管理员：显示创建表单
          <form
            onSubmit={handleSubmit}
            className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 space-y-5"
          >
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                 {t("board_name")} <span className="text-red-500">*</span>
              </label>
              <input
                name="name"
                value={form.name}
                onChange={handleChange}
                placeholder={t("board_name_placeholder")}
                className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                 {t("board_slug")} <span className="text-red-500">*</span>
              </label>
              <div className="flex items-center rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 overflow-hidden focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-transparent transition">
                <span className="pl-4 pr-1 text-gray-400 text-sm select-none">
                  /boards/
                </span>
                <input
                  name="slug"
                  value={form.slug}
                  onChange={handleChange}
                  placeholder="frontend"
                  className="flex-1 py-2.5 pr-4 bg-transparent text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none text-sm"
                />
              </div>
              <p className="mt-1 text-xs text-gray-400">
                 {t("board_slug_hint")}
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                 {t("board_description")}
              </label>
              <textarea
                name="description"
                value={form.description}
                onChange={handleChange}
                placeholder={t("board_description_placeholder")}
                rows={3}
                className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition resize-none"
              />
            </div>

            {error && (
              <p className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 px-4 py-2 rounded-lg">
                {t(error)}
              </p>
            )}

            <div className="flex items-center justify-between pt-1">
              <Link
                href="/boards"
                className="text-sm text-gray-400 hover:text-blue-500 transition-colors flex items-center gap-1"
              >
                <ArrowLeftIcon className="w-3.5 h-3.5" />
                 {t("back_to_list")}
              </Link>
              <button
                type="submit"
                disabled={submitting}
                className="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed text-white rounded-xl font-medium transition-colors"
              >
                <FolderPlusIcon className="w-4 h-4" />
                {submitting ? t("creating") : t("create_board")}
              </button>
            </div>
          </form>
        ) : (
          // 非管理员用户：显示申请创建板块界面
          <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6">
            {!showApplyForm ? (
              // 申请按钮界面
              <div className="text-center space-y-4">
                <div className="inline-flex items-center justify-center w-14 h-14 rounded-full bg-amber-50 dark:bg-amber-900/30">
                  <MailIcon className="w-7 h-7 text-amber-500" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                     {t("apply_create_board")}
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                     {t("admin_only_create")}
                  </p>
                </div>
                <div className="pt-2 space-y-3">
                  <button
                    onClick={() => setShowApplyForm(true)}
                    className="inline-block px-6 py-2.5 bg-amber-600 hover:bg-amber-700 text-white rounded-xl font-medium transition-colors"
                  >
                     {t("go_apply")}
                  </button>
                  <div>
                    <Link
                      href="/boards"
                      className="text-sm text-gray-400 hover:text-blue-500 transition-colors flex items-center justify-center gap-1"
                    >
                      <ArrowLeftIcon className="w-3.5 h-3.5" />
                       {t("back_to_board_list")}
                    </Link>
                  </div>
                </div>
              </div>
            ) : applySuccess ? (
              // 提交成功界面
              <div className="text-center space-y-4">
                <div className="inline-flex items-center justify-center w-14 h-14 rounded-full bg-green-50 dark:bg-green-900/30">
                  <svg
                    className="w-7 h-7 text-green-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                     {t("apply_submitted")}
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                     {t("apply_wait_review")}
                  </p>
                </div>
                <div>
                  <Link
                    href="/boards"
                    className="text-sm text-blue-500 hover:text-blue-600 transition-colors flex items-center justify-center gap-1"
                  >
                    <ArrowLeftIcon className="w-3.5 h-3.5" />
                     {t("back_to_board_list")}
                  </Link>
                </div>
              </div>
            ) : (
              // 申请表单
              <form onSubmit={handleApplySubmit} className="space-y-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                     {t("apply_create_board")}
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
                     {t("apply_create_reason_desc")}{" "}
                    <span className="font-mono text-blue-600">/{slug}</span>
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                     {t("board_name")} <span className="text-red-500">*</span>
                  </label>
                  <input
                    value={form.name}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      setForm({ ...form, name: e.target.value });
                      setApplyError("");
                    }}
                    type="text"
                    placeholder={t("board_name_placeholder")}
                    className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-transparent transition"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                     {t("apply_reason")} <span className="text-red-500">*</span>
                  </label>
                  <textarea
                    value={applyReason}
                    onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
                      setApplyReason(e.target.value);
                      setApplyError("");
                    }}
                    placeholder={t("apply_reason_placeholder_board")}
                    rows={4}
                    className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-transparent transition resize-none"
                  />
                  {applyError && (
                    <p className="mt-1 text-sm text-red-500">{t(applyError)}</p>
                  )}
                </div>

                <div className="flex items-center justify-between gap-3 pt-2">
                  <button
                    type="button"
                    onClick={() => {
                      setShowApplyForm(false);
                      setApplyError("");
                      setApplyReason("");
                    }}
                    className="flex-1 px-4 py-2.5 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                  >
                     {t("cancel")}
                  </button>
                  <button
                    type="submit"
                    disabled={applySubmitting}
                    className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-amber-600 hover:bg-amber-700 disabled:opacity-60 disabled:cursor-not-allowed text-white rounded-xl font-medium transition-colors"
                  >
                    <MailIcon className="w-4 h-4" />
                    {applySubmitting ? t("submitting") : t("submit_application")}
                  </button>
                </div>
              </form>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
