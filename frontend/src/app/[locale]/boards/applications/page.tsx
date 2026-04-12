'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { boardApi } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import type { ModeratorApplication } from '@/types';
import {
  ShieldCheckIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  ArrowTopRightOnSquareIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/react/24/outline';

type AppStatus = ModeratorApplication['status'];

const STATUS_MAP: Record<AppStatus, { icon: React.ElementType; label: string; cls: string }> = {
  pending:  { icon: ClockIcon,        label: '审核中', cls: 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400' },
  approved: { icon: CheckCircleIcon,  label: '已通过', cls: 'text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400' },
  rejected: { icon: XCircleIcon,      label: '已拒绝', cls: 'text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400' },
};

function StatusBadge({ status }: { status: AppStatus }) {
  const { icon: Icon, label, cls } = STATUS_MAP[status];
  return (
    <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium ${cls}`}>
      <Icon className="w-4 h-4" />{label}
    </span>
  );
}

function ApplicationCard({ app }: { app: ModeratorApplication }) {
  const boardSlug = app.board?.slug || String(app.board_id);

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 hover:shadow-md transition-shadow p-5">
      {/* Top row */}
      <div className="flex items-start justify-between gap-4 mb-4">
        <div>
          <Link
            href={`/boards/${boardSlug}`}
            className="font-semibold text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
          >
            {app.board?.name || `板块 #${app.board_id}`}
          </Link>
          <p className="text-xs text-gray-400 mt-1">
            提交于 {new Date(app.created_at).toLocaleString('zh-CN')}
            {app.handled_by && app.updated_at !== app.created_at && (
              <> · 处理于 {new Date(app.updated_at).toLocaleString('zh-CN')}</>
            )}
          </p>
        </div>
        <StatusBadge status={app.status} />
      </div>

      {/* Reason */}
      <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg px-4 py-3 mb-3">
        <p className="text-xs font-medium text-gray-400 mb-1">申请理由</p>
        <p className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap line-clamp-3">
          {app.reason}
        </p>
      </div>

      {/* Handle note */}
      {app.handle_note && (
        <div className={`rounded-lg px-4 py-3 text-sm ${
          app.status === 'approved'
            ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300'
            : 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300'
        }`}>
          <p className="font-medium mb-0.5">{app.status === 'approved' ? '通过说明' : '拒绝理由'}</p>
          <p>{app.handle_note}</p>
          {app.handler && (
            <p className="text-xs opacity-60 mt-1">处理人：{app.handler.username}</p>
          )}
        </div>
      )}

      {/* Enter board link for approved */}
      {app.status === 'approved' && app.board && (
        <div className="flex justify-end mt-3">
          <Link
            href={`/boards/${boardSlug}`}
            className="inline-flex items-center gap-1 text-sm text-blue-600 hover:text-blue-700 font-medium transition-colors"
          >
            进入板块
            <ArrowTopRightOnSquareIcon className="w-3.5 h-3.5" />
          </Link>
        </div>
      )}
    </div>
  );
}

export default function MyApplicationsPage() {
  const { isAuthenticated } = useAuthStore();
  const router = useRouter();
  const [applications, setApplications] = useState<ModeratorApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const PAGE_SIZE = 20;

  useEffect(() => {
    if (!isAuthenticated) router.push('/auth/login?redirect=/boards/applications');
  }, [isAuthenticated, router]);

  const load = useCallback(async () => {
    if (!isAuthenticated) return;
    setLoading(true);
    try {
      const res = await boardApi.getMyApplications({ page, page_size: PAGE_SIZE });
      setApplications(res.data.data.list);
      setTotal(res.data.data.total);
    } catch (err) {
      console.error('Failed to load applications:', err);
    } finally {
      setLoading(false);
    }
  }, [page, isAuthenticated]);

  useEffect(() => { load(); }, [load]);

  const totalPages = Math.ceil(total / PAGE_SIZE);

  if (!isAuthenticated) return null;

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Header */}
      <div className="flex items-center gap-3 mb-8">
        <div className="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
          <ShieldCheckIcon className="w-5 h-5 text-blue-600" />
        </div>
        <div>
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">我的版主申请</h1>
          <p className="text-sm text-gray-400">查看全部申请记录</p>
        </div>
      </div>

      {loading ? (
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 h-40 animate-pulse" />
          ))}
        </div>
      ) : applications.length === 0 ? (
        <div className="text-center py-20 bg-white dark:bg-gray-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
          <ShieldCheckIcon className="w-12 h-12 text-gray-200 dark:text-gray-600 mx-auto mb-3" />
          <p className="text-gray-500 dark:text-gray-400 mb-4">还没有提交过版主申请</p>
          <Link
            href="/boards"
            className="inline-block px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl font-medium transition-colors"
          >
            浏览板块
          </Link>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {applications.map(app => <ApplicationCard key={app.id} app={app} />)}
          </div>

          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-8">
              <button
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page === 1}
                className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
              >
                <ChevronLeftIcon className="w-4 h-4" />
              </button>
              <span className="text-sm text-gray-500">
                第 <strong className="text-gray-900 dark:text-white">{page}</strong> / {totalPages} 页
              </span>
              <button
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages}
                className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
              >
                <ChevronRightIcon className="w-4 h-4" />
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}