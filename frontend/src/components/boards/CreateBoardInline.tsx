// ─── CreateBoardInline ─────────────────────────────────────────────────────────

import { Board, boardApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { ExclamationCircleIcon } from "@heroicons/react/24/outline";
import Link from "next/link";
import { ArrowLeftIcon, FolderPlusIcon, MailIcon } from "lucide-react";
import { useState } from "react";

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
  const [form, setForm] = useState<FormData>({
    name: "",
    slug: slug,
    description: "",
  });
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

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
      setError("请填写板块名称");
      return;
    }
    if (!form.slug.trim()) {
      setError("请填写板块标识");
      return;
    }
    if (!/^[a-z0-9-]+$/.test(form.slug)) {
      setError("标识只能包含小写字母、数字和连字符");
      return;
    }

    setSubmitting(true);
    try {
      const res = await boardApi.create({
        name: form.name.trim(),
        slug: form.slug.trim(),
        description: form.description.trim(),
      });

      // 修复：检查 res.data.data 是否存在
      if (res.data.data) {
        onCreated(res.data.data);
      } else {
        setError("创建成功但未返回板块数据");
      }
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      setError(error.response?.data?.message || "创建失败，请稍后重试");
    } finally {
      setSubmitting(false);
    }
  };

  // 处理申请提交
  const handleApplySubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!applyReason.trim()) {
      setApplyError("请填写申请理由");
      return;
    }

    setApplySubmitting(true);
    setApplyError("");

    try {
      // TODO: 替换为实际的申请API接口
      await new Promise((resolve) => setTimeout(resolve, 1000));

      setApplySuccess(true);
      // 3秒后关闭申请表单
      setTimeout(() => {
        setShowApplyForm(false);
        setApplySuccess(false);
        setApplyReason("");
      }, 3000);
    } catch (err: unknown) {
      const error = err as ErrorResponse;
      setApplyError(error.response?.data?.message || "提交失败，请稍后重试");
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
            板块不存在
          </h1>
          <p className="mt-2 text-gray-500 dark:text-gray-400">
            <span className="font-mono bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded text-sm text-blue-600 dark:text-blue-400">
              /{slug}
            </span>{" "}
            尚未创建。
          </p>
        </div>

        {!user ? (
          // 未登录状态
          <div className="text-center space-y-3">
            <p className="text-gray-500">需要登录才能创建板块</p>
            <Link
              href={`/auth/login?redirect=/boards/${slug}`}
              className="inline-block px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-medium transition-colors"
            >
              去登录
            </Link>
            <div className="pt-2">
              <Link
                href="/boards"
                className="text-sm text-gray-400 hover:text-blue-500 transition-colors"
              >
                ← 返回板块列表
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
                板块名称 <span className="text-red-500">*</span>
              </label>
              <input
                name="name"
                value={form.name}
                onChange={handleChange}
                placeholder="例如：前端技术交流"
                className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                板块标识 <span className="text-red-500">*</span>
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
                仅小写字母、数字、连字符
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                板块描述
              </label>
              <textarea
                name="description"
                value={form.description}
                onChange={handleChange}
                placeholder="简单介绍这个板块的用途..."
                rows={3}
                className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition resize-none"
              />
            </div>

            {error && (
              <p className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 px-4 py-2 rounded-lg">
                {error}
              </p>
            )}

            <div className="flex items-center justify-between pt-1">
              <Link
                href="/boards"
                className="text-sm text-gray-400 hover:text-blue-500 transition-colors flex items-center gap-1"
              >
                <ArrowLeftIcon className="w-3.5 h-3.5" />
                返回列表
              </Link>
              <button
                type="submit"
                disabled={submitting}
                className="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed text-white rounded-xl font-medium transition-colors"
              >
                <FolderPlusIcon className="w-4 h-4" />
                {submitting ? "创建中..." : "创建板块"}
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
                    申请创建板块
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    只有管理员可以创建新板块。你可以提交申请，管理员审核后将为你创建。
                  </p>
                </div>
                <div className="pt-2 space-y-3">
                  <button
                    onClick={() => setShowApplyForm(true)}
                    className="inline-block px-6 py-2.5 bg-amber-600 hover:bg-amber-700 text-white rounded-xl font-medium transition-colors"
                  >
                    去申请
                  </button>
                  <div>
                    <Link
                      href="/boards"
                      className="text-sm text-gray-400 hover:text-blue-500 transition-colors flex items-center justify-center gap-1"
                    >
                      <ArrowLeftIcon className="w-3.5 h-3.5" />
                      返回板块列表
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
                    申请已提交
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    管理员将尽快审核你的申请，请耐心等待。
                  </p>
                </div>
                <div>
                  <Link
                    href="/boards"
                    className="text-sm text-blue-500 hover:text-blue-600 transition-colors flex items-center justify-center gap-1"
                  >
                    <ArrowLeftIcon className="w-3.5 h-3.5" />
                    返回板块列表
                  </Link>
                </div>
              </div>
            ) : (
              // 申请表单
              <form onSubmit={handleApplySubmit} className="space-y-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                    申请创建板块
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
                    请填写申请理由，管理员审核后将为你创建板块{" "}
                    <span className="font-mono text-blue-600">/{slug}</span>
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                    板块名称
                  </label>
                  <input
                    value={form.name}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setForm({ ...form, name: e.target.value })
                    }
                    type="text"
                    placeholder="请输入板块名称"
                    className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed"
                    disabled
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                    申请理由 <span className="text-red-500">*</span>
                  </label>
                  <textarea
                    value={applyReason}
                    onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => {
                      setApplyReason(e.target.value);
                      setApplyError("");
                    }}
                    placeholder="请说明为什么需要创建这个板块，以及它的用途..."
                    rows={4}
                    className="w-full px-4 py-2.5 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-transparent transition resize-none"
                  />
                  {applyError && (
                    <p className="mt-1 text-sm text-red-500">{applyError}</p>
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
                    取消
                  </button>
                  <button
                    type="submit"
                    disabled={applySubmitting}
                    className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-amber-600 hover:bg-amber-700 disabled:opacity-60 disabled:cursor-not-allowed text-white rounded-xl font-medium transition-colors"
                  >
                    <MailIcon className="w-4 h-4" />
                    {applySubmitting ? "提交中..." : "提交申请"}
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
