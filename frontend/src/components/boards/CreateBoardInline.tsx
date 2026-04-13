// ─── CreateBoardInline ─────────────────────────────────────────────────────────

import { Board, boardApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import { ExclamationCircleIcon } from "@heroicons/react/24/outline";
import { Link, ArrowLeftIcon, FolderPlusIcon } from "lucide-react";
import { useState } from "react";

export function CreateBoardInline({ slug, onCreated }: { slug: string; onCreated: (board: Board) => void }) {
  const { user } = useAuthStore();
  const [form, setForm] = useState({
    name: '',
    slug: slug,
    description: '',
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
    setError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name.trim()) { setError('请填写板块名称'); return; }
    if (!form.slug.trim()) { setError('请填写板块标识'); return; }
    if (!/^[a-z0-9-]+$/.test(form.slug)) { setError('标识只能包含小写字母、数字和连字符'); return; }

    setSubmitting(true);
    try {
      const res = await boardApi.create({
        name: form.name.trim(),
        slug: form.slug.trim(),
        description: form.description.trim(),
      });
      onCreated(res.data.data);
    } catch (err: any) {
      setError(err.response?.data?.message || '创建失败，请稍后重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="w-full max-w-lg">
        {/* Icon + heading */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-blue-50 dark:bg-blue-900/30 mb-4">
            <ExclamationCircleIcon className="w-8 h-8 text-blue-500" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">板块不存在</h1>
          <p className="mt-2 text-gray-500 dark:text-gray-400">
            <span className="font-mono bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded text-sm text-blue-600 dark:text-blue-400">/{slug}</span>
            {' '}尚未创建，你可以现在创建它。
          </p>
        </div>

        {!user ? (
          <div className="text-center space-y-3">
            <p className="text-gray-500">需要登录才能创建板块</p>
            <Link
              href={`/auth/login?redirect=/boards/${slug}`}
              className="inline-block px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-medium transition-colors"
            >
              去登录
            </Link>
            <div className="pt-2">
              <Link href="/boards" className="text-sm text-gray-400 hover:text-blue-500 transition-colors">
                ← 返回板块列表
              </Link>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 space-y-5">
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
                <span className="pl-4 pr-1 text-gray-400 text-sm select-none">/boards/</span>
                <input
                  name="slug"
                  value={form.slug}
                  onChange={handleChange}
                  placeholder="frontend"
                  className="flex-1 py-2.5 pr-4 bg-transparent text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none text-sm"
                />
              </div>
              <p className="mt-1 text-xs text-gray-400">仅小写字母、数字、连字符</p>
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
              <Link href="/boards" className="text-sm text-gray-400 hover:text-blue-500 transition-colors flex items-center gap-1">
                <ArrowLeftIcon className="w-3.5 h-3.5" />
                返回列表
              </Link>
              <button
                type="submit"
                disabled={submitting}
                className="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed text-white rounded-xl font-medium transition-colors"
              >
                <FolderPlusIcon className="w-4 h-4" />
                {submitting ? '创建中...' : '创建板块'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}