"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import { moderatorApi } from "@/lib/api/modules/moderator";
import {
  ShieldCheckIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  ArrowTopRightOnSquareIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ExclamationTriangleIcon,
} from "@heroicons/react/24/outline";

// 申请状态类型
type ApplicationStatus = "pending" | "approved" | "rejected" | "canceled";

// 扩展的申请卡片数据
interface ExtendedApplication {
  id: number;
  board_id: number;
  board_name: string;
  board_slug: string;
  reason: string;
  status: ApplicationStatus;
  review_note?: string;
  reviewer_id?: number;
  created_at: string;
  reviewed_at?: string | null;
  requested_perms?: {
    delete_post: boolean;
    pin_post: boolean;
    edit_any_post: boolean;
    manage_moderator: boolean;
    ban_user: boolean;
  };
}

// API 返回的申请数据类型
interface ApiApplication {
  id: number;
  board_id: number;
  board?: {
    name: string;
    slug: string;
  };
  reason: string;
  status: ApplicationStatus;
  review_note?: string;
  reviewed_by?: number;
  created_at: string;
  reviewed_at?: string | null;
  req_delete_post?: boolean;
  req_pin_post?: boolean;
  req_edit_any_post?: boolean;
  req_manage_moderator?: boolean;
  req_ban_user?: boolean;
}

const STATUS_MAP: Record<
  ApplicationStatus,
  {
    icon: React.ElementType;
    label: string;
    cls: string;
    borderCls: string;
  }
> = {
  pending: {
    icon: ClockIcon,
    label: "审核中",
    cls: "text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400",
    borderCls: "border-yellow-200 dark:border-yellow-800",
  },
  approved: {
    icon: CheckCircleIcon,
    label: "已通过",
    cls: "text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400",
    borderCls: "border-green-200 dark:border-green-800",
  },
  rejected: {
    icon: XCircleIcon,
    label: "已拒绝",
    cls: "text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400",
    borderCls: "border-red-200 dark:border-red-800",
  },
  canceled: {
    icon: XCircleIcon,
    label: "已撤销",
    cls: "text-gray-600 bg-gray-100 dark:bg-gray-900/30 dark:text-gray-400",
    borderCls: "border-gray-200 dark:border-gray-700",
  },
};

const formatDate = (dateString: string, locale = "zh-CN"): string => {
  try {
    return new Date(dateString).toLocaleString(locale);
  } catch {
    return "日期无效";
  }
};

function StatusBadge({ status }: { status: ApplicationStatus }) {
  const { icon: Icon, label, cls } = STATUS_MAP[status];
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium transition-colors ${cls}`}
    >
      <Icon className="w-4 h-4" aria-hidden="true" />
      {label}
    </span>
  );
}

function ApplicationCard({
  app,
  onCancel,
  cancelLoadingId,
}: {
  app: ExtendedApplication;
  onCancel?: (id: number) => void;
  cancelLoadingId?: number | null;
}) {
  const isHandled = app.status !== "pending";
  const hasReviewNote = !!app.review_note;
  const canCancel = app.status === "pending";
  const isLoading = cancelLoadingId === app.id;

  return (
    <div
      className={`bg-white dark:bg-gray-800 rounded-xl border ${STATUS_MAP[app.status].borderCls} hover:shadow-lg transition-all duration-200 p-5`}
    >
      {/* Top row */}
      <div className="flex items-start justify-between gap-4 mb-4">
        <div className="flex-1 min-w-0">
          <Link
            href={`/boards/${app.board_slug}`}
            className="font-semibold text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors inline-flex items-center gap-1"
          >
            {app.board_name}
          </Link>
          <p className="text-xs text-gray-400 mt-1">
            <time dateTime={app.created_at}>
              提交于 {formatDate(app.created_at)}
            </time>
            {isHandled && app.reviewed_at && (
              <>
                {" · "}
                <time dateTime={app.reviewed_at}>
                  处理于 {formatDate(app.reviewed_at)}
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
          {app.reason || "未提供理由"}
        </p>
      </div>

      {/* Requested Permissions */}
      {app.requested_perms && (
        <div className="mb-3">
          <p className="text-xs font-medium text-gray-400 mb-2">申请权限</p>
          <div className="flex flex-wrap gap-2">
            {app.requested_perms.delete_post && (
              <span className="text-xs px-2 py-1 rounded bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                删除帖子
              </span>
            )}
            {app.requested_perms.pin_post && (
              <span className="text-xs px-2 py-1 rounded bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                置顶帖子
              </span>
            )}
            {app.requested_perms.edit_any_post && (
              <span className="text-xs px-2 py-1 rounded bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                编辑任意帖子
              </span>
            )}
            {app.requested_perms.manage_moderator && (
              <span className="text-xs px-2 py-1 rounded bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                管理版主
              </span>
            )}
            {app.requested_perms.ban_user && (
              <span className="text-xs px-2 py-1 rounded bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                禁言用户
              </span>
            )}
          </div>
        </div>
      )}

      {/* Review note */}
      {hasReviewNote && (
        <div
          className={`rounded-lg px-4 py-3 text-sm ${
            app.status === "approved"
              ? "bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300 border-l-4 border-green-500"
              : "bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 border-l-4 border-red-500"
          }`}
        >
          <p className="font-medium mb-1">
            {app.status === "approved" ? "通过说明" : "拒绝理由"}
          </p>
          <p className="whitespace-pre-wrap break-words">{app.review_note}</p>
          {app.reviewer_id && (
            <p className="text-xs opacity-70 mt-2">
              处理人ID：{app.reviewer_id}
            </p>
          )}
        </div>
      )}

      {/* Action buttons */}
      <div className="flex justify-end gap-3 mt-3">
        {canCancel && onCancel && (
          <button
            onClick={() => onCancel(app.id)}
            disabled={isLoading}
            className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-red-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? "撤销中..." : "撤销申请"}
          </button>
        )}
        {app.status === "approved" && (
          <Link
            href={`/boards/${app.board_slug}`}
            className="inline-flex items-center gap-1 text-sm text-blue-600 hover:text-blue-700 font-medium transition-colors group"
          >
            进入板块
            <ArrowTopRightOnSquareIcon className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
          </Link>
        )}
        {app.status === "rejected" && (
          <Link
            href={`/boards/${app.board_slug}/apply-moderator`}
            className="inline-flex items-center gap-1 text-sm text-blue-600 hover:text-blue-700 font-medium transition-colors"
          >
            重新申请
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
      <ShieldCheckIcon
        className="w-12 h-12 text-gray-200 dark:text-gray-600 mx-auto mb-3"
        aria-hidden="true"
      />
      <p className="text-gray-500 dark:text-gray-400 mb-4">
        还没有提交过版主申请
      </p>
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
      <ExclamationTriangleIcon
        className="w-12 h-12 text-red-400 mx-auto mb-3"
        aria-hidden="true"
      />
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
  const [applications, setApplications] = useState<ExtendedApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [cancelLoadingId, setCancelLoadingId] = useState<number | null>(null);

  const PAGE_SIZE = 10;

  // 重定向未认证用户
  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/auth/login?redirect=/boards/my-applications");
    }
  }, [isAuthenticated, router]);

  // 转换 API 数据为组件需要的格式
  const transformApplication = (app: ApiApplication): ExtendedApplication => ({
    id: app.id,
    board_id: app.board_id,
    board_name: app.board?.name || `板块 #${app.board_id}`,
    board_slug: app.board?.slug || `board-${app.board_id}`,
    reason: app.reason,
    status: app.status,
    review_note: app.review_note,
    reviewer_id: app.reviewed_by,
    created_at: app.created_at,
    reviewed_at: app.reviewed_at,
    requested_perms: {
      delete_post: app.req_delete_post || false,
      pin_post: app.req_pin_post || false,
      edit_any_post: app.req_edit_any_post || false,
      manage_moderator: app.req_manage_moderator || false,
      ban_user: app.req_ban_user || false,
    },
  });

  // 加载用户的所有申请
  const loadApplications = useCallback(async () => {
    if (!isAuthenticated) return;

    setLoading(true);
    setError(null);

    try {
      const res = await moderatorApi.getMyApplications({
        page,
        page_size: PAGE_SIZE,
      });

      if (res?.data?.data) {
        // 使用类型安全的方式转换数据
        const apps: ExtendedApplication[] = (res.data.data.list || []).map(
          (app: ApiApplication) => transformApplication(app),
        );
        setApplications(apps);
        setTotal(res.data.data.total);
      }
    } catch {
      setError("加载申请记录失败");
    } finally {
      setLoading(false);
    }
  }, [page, isAuthenticated]);

  // 撤销申请
  const handleCancel = useCallback(
    async (applicationId: number) => {
      if (!confirm("确定要撤销这个申请吗？")) return;

      setCancelLoadingId(applicationId);
      try {
        await moderatorApi.cancelApplication(applicationId);
        await loadApplications();
      } catch (err) {
        console.error("撤销失败:", err);
        alert("撤销失败，请重试");
      } finally {
        setCancelLoadingId(null);
      }
    },
    [loadApplications],
  );

  useEffect(() => {
    loadApplications();
  }, [loadApplications]);

  const totalPages = useMemo(() => Math.ceil(total / PAGE_SIZE), [total]);

  const handlePageChange = useCallback(
    (newPage: number) => {
      setPage(Math.min(Math.max(1, newPage), totalPages));
      window.scrollTo({ top: 0, behavior: "smooth" });
    },
    [totalPages],
  );

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
          <ShieldCheckIcon
            className="w-5 h-5 text-blue-600"
            aria-hidden="true"
          />
        </div>
        <div>
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">
            我的版主申请
          </h1>
          <p className="text-sm text-gray-400">
            共 {applications.length} 条申请记录
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
              <ApplicationCard
                key={app.id}
                app={app}
                onCancel={handleCancel}
                cancelLoadingId={cancelLoadingId}
              />
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <nav
              className="flex items-center justify-center gap-2 mt-8"
              aria-label="分页导航"
            >
              <button
                onClick={() => handlePageChange(page - 1)}
                disabled={page === 1}
                className="p-2 rounded-lg border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
                aria-label="上一页"
              >
                <ChevronLeftIcon className="w-4 h-4" aria-hidden="true" />
              </button>

              <span className="text-sm text-gray-500">
                第{" "}
                <strong className="text-gray-900 dark:text-white">
                  {page}
                </strong>{" "}
                / {totalPages} 页
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