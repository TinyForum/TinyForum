'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { boardApi, ModeratorApplication } from '@/lib/api';
import { useAuthStore } from '@/store/auth';
import {
  ShieldCheckIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  ArrowTopRightOnSquareIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';

type AppStatus = ModeratorApplication['status'];

interface ExtendedModeratorApplication extends ModeratorApplication {
  board?: {
    id: number;
    name: string;
    slug: string;
  };
  handler?: {
    id: number;
    username: string;
  };
  handle_note?: string;
  updated_at?: string;
  handled_by?: number;
}

const STATUS_MAP: Record<AppStatus, { 
  icon: React.ElementType; 
  label: string; 
  cls: string;
  borderCls: string;
}> = {
  pending: { 
    icon: ClockIcon, 
    label: '审核中', 
    cls: 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400',
    borderCls: 'border-yellow-200 dark:border-yellow-800'
  },
  approved: { 
    icon: CheckCircleIcon, 
    label: '已通过', 
    cls: 'text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400',
    borderCls: 'border-green-200 dark:border-green-800'
  },
  rejected: { 
    icon: XCircleIcon, 
    label: '已拒绝', 
    cls: 'text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400',
    borderCls: 'border-red-200 dark:border-red-800'
  },
};

const formatDate = (dateString: string, locale = 'zh-CN'): string => {
  try {
    return new Date(dateString).toLocaleString(locale);
  } catch {
    return '日期无效';
  }
};

function StatusBadge({ status }: { status: AppStatus }) {
  const { icon: Icon, label, cls } = STATUS_MAP[status];
  return (
    <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium transition-colors ${cls}`}>
      <Icon className="w-4 h-4" aria-hidden="true" />
      {label}
    </span>
  );
}

function ApplicationCard({ app }: { app: ExtendedModeratorApplication }) {
  const boardSlug = app.board?.slug || `board-${app.board_id}`;
  const isHandled = app.handled_by && app.updated_at !== app.created_at;
  const hasHandleNote = !!app.handle_note;

  return (
    <div className={`bg-white dark:bg-gray-800 rounded-xl border ${STATUS_MAP[app.status].borderCls} hover:shadow-lg transition-all duration-200 p-5`}>
      {/* Top row */}
      <div className="flex items-start justify-between gap-4 mb-4">
        <div className="flex-1 min-w-0">
          <Link
            href={`/boards/${boardSlug}`}
            className="font-semibold text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors inline-flex items-center gap-1"
          >
            {app.board?.name || `板块 #${app.board_id}`}
          </Link>
          <p className="text-xs text-gray-400 mt-1">
            <time dateTime={app.created_at}>
              提交于 {formatDate(app.created_at)}
            </time>
            {isHandled && (
              <>
                {' · '}
                <time dateTime={app.updated_at!}>
                  处理于 {formatDate(app.updated_at!)}
                </time>
              </>
            )}
          </p>
        </div>
        <StatusBadge status={app.status} />
      </div>

      {/* Reason */}
      <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg px-4 py-3 mb-3">
        <p className="text-xs font-medium text-gray-400 mb-1">申请理由</p>
        <p className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap break-words">
          {app.reason || '未提供理由'}
        </p>
      </div>

      {/* Handle note */}
      {hasHandleNote && (
        <div className={`rounded-lg px-4 py-3 text-sm ${
          app.status === 'approved'
            ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 border-l-4 border-green-500'
            : 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 border-l-4 border-red-500'
        }`}>
          <p className="font-medium mb-1">
            {app.status === 'approved' ? '通过说明' : '拒绝理由'}
          </p>
          <p className="whitespace-pre-wrap break-words">{app.handle_note}</p>
          {app.handler && (
            <p className="text-xs opacity-70 mt-2">
              处理人：{app.handler.username}
            </p>
          )}
        </div>
      )}

      {/* Action buttons */}
      <div className="flex justify-end gap-3 mt-3">
        {app.status === 'approved' && app.board && (
          <Link
            href={`/boards/${boardSlug}`}
            className="inline-flex items-center gap-1 text-sm text-blue-600 hover:text-blue-700 font-medium transition-colors group"
          >
            进入板块
            <ArrowTopRightOnSquareIcon className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
          </Link>
        )}
      </div>
    </div>
  );
}

function LoadingSkeleton() {
  return (
    <div className="space-y-4" aria-label="加载中">
      {[...Array(3)].map((_, i) => (
        <div 
          key={i} 
          className="bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 h-40 animate-pulse"
          aria-hidden="true"
        />
      ))}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="text-center py-20 bg-white dark:bg-gray-800 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
      <ShieldCheckIcon className="w-12 h-12 text-gray-200 dark:text-gray-600 mx-auto mb-3" aria-hidden="true" />
      <p className="text-gray-500 dark:text-gray-400 mb-4">还没有提交过版主申请</p>
      <Link
        href="/boards"
        className="inline-block px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-xl font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
      >
        浏览板块
      </Link>
    </div>
  );
}

function ErrorState({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="text-center py-20 bg-white dark:bg-gray-800 rounded-2xl border border-red-200 dark:border-red-800">
      <ExclamationTriangleIcon className="w-12 h-12 text-red-400 mx-auto mb-3" aria-hidden="true" />
      <p className="text-gray-600 dark:text-gray-400 mb-4">加载申请记录失败</p>
      <button
        onClick={onRetry}
        className="inline-block px-5 py-2 bg-gray-600 hover:bg-gray-700 text-white text-sm rounded-xl font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
      >
        重试
      </button>
    </div>
  );
}

export default function MyApplicationsPage() {
  const { isAuthenticated } = useAuthStore();
  const router = useRouter();
  const [applications, setApplications] = useState<ExtendedModeratorApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const PAGE_SIZE = 20;

  // 重定向未认证用户
  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/auth/login?redirect=/boards/applications');
    }
  }, [isAuthenticated, router]);

  const loadApplications = useCallback(async () => {
    if (!isAuthenticated) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const res = await boardApi.getMyApplications({ page, page_size: PAGE_SIZE });
      
      if (res?.data?.data) {
        setApplications(res.data.data.list || []);
        setTotal(res.data.data.total || 0);
      } else {
        throw new Error('响应数据格式错误');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '加载申请记录失败';
      setError(errorMessage);
      console.error('Failed to load applications:', err);
    } finally {
      setLoading(false);
    }
  }, [page, isAuthenticated]);

  useEffect(() => {
    loadApplications();
  }, [loadApplications]);

  const totalPages = useMemo(() => Math.ceil(total / PAGE_SIZE), [total]);

  const handlePageChange = useCallback((newPage: number) => {
    setPage(Math.min(Math.max(1, newPage), totalPages));
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, [totalPages]);

  const handleRetry = useCallback(() => {
    loadApplications();
  }, [loadApplications]);

  // 未认证时不渲染内容
  if (!isAuthenticated) return null;

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Header */}
      <div className="flex items-center gap-3 mb-8">
        <div className="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
          <ShieldCheckIcon className="w-5 h-5 text-blue-600" aria-hidden="true" />
        </div>
        <div>
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">
            我的版主申请
          </h1>
          <p className="text-sm text-gray-400">
            共 {total} 条申请记录
          </p>
        </div>
      </div>

      {/* Content */}
      {loading ? (
        <LoadingSkeleton />
      ) : error ? (
        <ErrorState onRetry={handleRetry} />
      ) : applications.length === 0 ? (
        <EmptyState />
      ) : (
        <>
          <div className="space-y-4">
            {applications.map((app) => (
              <ApplicationCard key={app.id} app={app} />
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <nav className="flex items-center justify-center gap-2 mt-8" aria-label="分页导航">
              <button
                onClick={() => handlePageChange(page - 1)}
                disabled={page === 1}
                className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
                aria-label="上一页"
              >
                <ChevronLeftIcon className="w-4 h-4" aria-hidden="true" />
              </button>
              
              <span className="text-sm text-gray-500">
                第 <strong className="text-gray-900 dark:text-white">{page}</strong> / {totalPages} 页
              </span>
              
              <button
                onClick={() => handlePageChange(page + 1)}
                disabled={page >= totalPages}
                className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
                aria-label="下一页"
              >
                <ChevronRightIcon className="w-4 h-4" aria-hidden="true" />
              </button>
            </nav>
          )}
        </>
      )}
    </div>
  );
}