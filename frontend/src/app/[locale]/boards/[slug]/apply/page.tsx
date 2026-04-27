"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Board, boardApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import {
  ShieldCheckIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  ArrowLeftIcon,
  ExclamationTriangleIcon,
  PencilSquareIcon,
} from "@heroicons/react/24/outline";
import {
  ApplyModeratorForm,
  moderatorApi,
  ModeratorApplication,
} from "@/lib/api/modules/moderator";

type AppStatus = ModeratorApplication["status"];

function StatusBadge({ status }: { status: AppStatus }) {
  const map: Record<
    AppStatus,
    { icon: React.ReactNode; label: string; cls: string }
  > = {
    pending: {
      icon: <ClockIcon className="w-4 h-4" />,
      label: "审核中",
      cls: "text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400",
    },
    approved: {
      icon: <CheckCircleIcon className="w-4 h-4" />,
      label: "已通过",
      cls: "text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400",
    },
    rejected: {
      icon: <XCircleIcon className="w-4 h-4" />,
      label: "已拒绝",
      cls: "text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400",
    },
    canceled: {
      icon: <XCircleIcon className="w-4 h-4" />,
      label: "已取消",
      cls: "text-gray-600 bg-gray-100 dark:bg-gray-700/30 dark:text-gray-400",
    },
  };

  const { icon, label, cls } = map[status];
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium ${cls}`}
    >
      {icon}
      {label}
    </span>
  );
}

export default function ApplyModeratorPage() {
  const params = useParams();
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();
  const slug = params.slug as string;

  const [board, setBoard] = useState<Board | null>(null);
  const [loadingBoard, setLoadingBoard] = useState(true);
  const [existing, setExisting] = useState<ModeratorApplication | null>(null);
  const [checkingStatus, setCheckingStatus] = useState(true);
  const [reason, setReason] = useState("");
  const [reqDeletePost, setReqDeletePost] = useState(false);
  const [reqPinPost, setReqPinPost] = useState(false);
  const [reqEditAnyPost, setReqEditAnyPost] = useState(false);
  const [reqManageModerator, setReqManageModerator] = useState(false);
  const [reqBanUser, setReqBanUser] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [reapply, setReapply] = useState(false);

  // Redirect unauthenticated users
  useEffect(() => {
    if (!isAuthenticated) {
      router.push(`/auth/login?redirect=/boards/${slug}/apply`);
    }
  }, [isAuthenticated, slug, router]);

  const loadBoard = useCallback(async () => {
    setLoadingBoard(true);
    try {
      const res = await boardApi.getBySlug(slug);
      setBoard(res.data.data);
    } catch {
      router.push("/boards");
    } finally {
      setLoadingBoard(false);
    }
  }, [slug, router]);

  // Check application status — depends on board being loaded
  const checkStatus = useCallback(
    async (boardId: number) => {
      setCheckingStatus(true);
      try {
        const res = await moderatorApi.listApplications({
          board_id: boardId,
          page: 1,
          page_size: 100,
        });
        const applications = res.data.data?.list || [];
        // 查找当前用户的申请
        const userApplication = applications.find(
          (app) => app.user_id === user?.id,
        );
        if (userApplication) {
          setExisting(userApplication);
        }
      } catch {
        // ignore
      } finally {
        setCheckingStatus(false);
      }
    },
    [user?.id],
  );

  useEffect(() => {
    if (isAuthenticated) loadBoard();
  }, [loadBoard, isAuthenticated]);
  useEffect(() => {
    if (board) checkStatus(board.id);
  }, [board, checkStatus]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!board) return;
    const trimmed = reason.trim();
    if (!trimmed) {
      setError("请填写申请理由");
      return;
    }
    if (trimmed.length < 20) {
      setError("申请理由至少需要 20 个字符");
      return;
    }

    setError("");
    setSubmitting(true);

    const formData: ApplyModeratorForm = {
      reason: trimmed,
      req_delete_post: reqDeletePost,
      req_pin_post: reqPinPost,
      req_edit_any_post: reqEditAnyPost,
      req_manage_moderator: reqManageModerator,
      req_ban_user: reqBanUser,
    };

    try {
      await moderatorApi.applyModerator(board.id, formData);
      router.push(`/boards/${slug}?applied=1`);
    } catch (err: any) {
      setError(err.response?.data?.message || "提交失败，请稍后重试");
    } finally {
      setSubmitting(false);
    }
  };

  const isLoading = loadingBoard || checkingStatus;
  const showForm = !existing || (existing.status === "rejected" && reapply);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500" />
      </div>
    );
  }

  if (!board) return null;

  return (
    <div className="max-w-2xl mx-auto px-4 py-10">
      {/* Back */}
      <Link
        href={`/boards/${slug}`}
        className="inline-flex items-center gap-1.5 text-sm text-gray-400 hover:text-blue-500 transition-colors mb-8"
      >
        <ArrowLeftIcon className="w-4 h-4" />
        返回 {board.name}
      </Link>

      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-blue-600 to-indigo-600 px-6 py-7 flex items-center gap-4">
          <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center shrink-0">
            <ShieldCheckIcon className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">申请版主</h1>
            <p className="text-blue-100 text-sm mt-0.5">{board.name}</p>
          </div>
        </div>

        <div className="p-6 space-y-6">
          {/* Existing application notice */}
          {existing && !reapply && (
            <div className="rounded-xl border border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/40 p-5">
              <div className="flex items-start gap-3">
                <ExclamationTriangleIcon className="w-5 h-5 text-yellow-500 shrink-0 mt-0.5" />
                <div className="space-y-2 flex-1">
                  <p className="font-medium text-gray-900 dark:text-white">
                    你已提交过申请
                  </p>
                  <div className="text-sm text-gray-500 dark:text-gray-400 space-y-1">
                    <p>
                      提交时间：
                      {new Date(existing.created_at).toLocaleString("zh-CN")}
                    </p>
                    <div className="flex items-center gap-2">
                      状态：
                      <StatusBadge status={existing.status} />
                    </div>
                  </div>
                  {existing.status === "rejected" && existing.review_note && (
                    <div className="mt-2 text-sm text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 rounded-lg px-3 py-2">
                      拒绝理由：{existing.review_note}
                    </div>
                  )}
                  {existing.status === "approved" && (
                    <p className="text-sm text-green-600 dark:text-green-400">
                      🎉 恭喜！你已成为该板块的版主。
                    </p>
                  )}
                  {existing.status === "rejected" && (
                    <button
                      onClick={() => {
                        setReapply(true);
                        setReason("");
                        setError("");
                        setReqDeletePost(false);
                        setReqPinPost(false);
                        setReqEditAnyPost(false);
                        setReqManageModerator(false);
                        setReqBanUser(false);
                      }}
                      className="mt-2 inline-flex items-center gap-1.5 text-sm text-blue-600 hover:text-blue-700 font-medium transition-colors"
                    >
                      <PencilSquareIcon className="w-4 h-4" />
                      重新申请
                    </button>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Form */}
          {showForm && (
            <form onSubmit={handleSubmit} className="space-y-5">
              {/* Board stats */}
              <div className="grid grid-cols-3 gap-3">
                {[
                  { label: "帖子总数", value: board.post_count },
                  { label: "今日发帖", value: board.today_count },
                  { label: "主题数", value: board.thread_count },
                ].map(({ label, value }) => (
                  <div
                    key={label}
                    className="text-center bg-gray-50 dark:bg-gray-700/50 rounded-xl py-3"
                  >
                    <p className="text-lg font-bold text-gray-900 dark:text-white">
                      {value}
                    </p>
                    <p className="text-xs text-gray-400 mt-0.5">{label}</p>
                  </div>
                ))}
              </div>

              {/* Responsibilities */}
              <div className="bg-blue-50 dark:bg-blue-900/20 rounded-xl p-4">
                <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  版主职责
                </p>
                <ul className="text-sm text-gray-500 dark:text-gray-400 space-y-1">
                  {[
                    "维护板块秩序，删除违规内容",
                    "管理优质内容，设置精华和置顶",
                    "及时回复用户问题，提供帮助",
                    "处理举报，维护社区和谐",
                  ].map((r) => (
                    <li key={r} className="flex items-center gap-2">
                      <span className="w-1 h-1 rounded-full bg-blue-400 shrink-0" />
                      {r}
                    </li>
                  ))}
                </ul>
              </div>

              {/* 期望权限（可选） */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  期望权限（可选）
                </label>
                <div className="space-y-2">
                  <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <input
                      type="checkbox"
                      checked={reqDeletePost}
                      onChange={(e) => setReqDeletePost(e.target.checked)}
                      className="rounded"
                    />
                    删除帖子权限
                  </label>
                  <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <input
                      type="checkbox"
                      checked={reqPinPost}
                      onChange={(e) => setReqPinPost(e.target.checked)}
                      className="rounded"
                    />
                    置顶帖子权限
                  </label>
                  <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <input
                      type="checkbox"
                      checked={reqEditAnyPost}
                      onChange={(e) => setReqEditAnyPost(e.target.checked)}
                      className="rounded"
                    />
                    编辑任意帖子权限
                  </label>
                  <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <input
                      type="checkbox"
                      checked={reqManageModerator}
                      onChange={(e) => setReqManageModerator(e.target.checked)}
                      className="rounded"
                    />
                    管理版主权限
                  </label>
                  <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <input
                      type="checkbox"
                      checked={reqBanUser}
                      onChange={(e) => setReqBanUser(e.target.checked)}
                      className="rounded"
                    />
                    禁言用户权限
                  </label>
                </div>
                <p className="text-xs text-gray-400 mt-2">
                  注：最终权限由管理员审批决定
                </p>
              </div>

              {/* Reason textarea */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  申请理由 <span className="text-red-500">*</span>
                </label>
                <textarea
                  value={reason}
                  onChange={(e) => {
                    setReason(e.target.value);
                    setError("");
                  }}
                  placeholder="请详细说明你为什么想成为版主，你有哪些优势，以及你打算如何管理这个板块..."
                  rows={5}
                  className="w-full px-4 py-3 rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition resize-none text-sm"
                />
                <div className="flex justify-between mt-1 text-xs">
                  {error ? (
                    <span className="text-red-500">{error}</span>
                  ) : (
                    <span className="text-gray-400">至少 20 个字符</span>
                  )}
                  <span
                    className={
                      reason.trim().length < 20
                        ? "text-gray-300"
                        : "text-green-500"
                    }
                  >
                    {reason.trim().length} 字
                  </span>
                </div>
              </div>

              {/* Actions */}
              <div className="flex items-center justify-between pt-1">
                {reapply ? (
                  <button
                    type="button"
                    onClick={() => {
                      setReapply(false);
                      setError("");
                    }}
                    className="text-sm text-gray-400 hover:text-gray-600 transition-colors"
                  >
                    取消重新申请
                  </button>
                ) : (
                  <Link
                    href={`/boards/${slug}`}
                    className="text-sm text-gray-400 hover:text-blue-500 transition-colors"
                  >
                    取消
                  </Link>
                )}
                <button
                  type="submit"
                  disabled={submitting || reason.trim().length < 20}
                  className="flex items-center gap-2 px-6 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-xl text-sm font-medium transition-colors"
                >
                  <ShieldCheckIcon className="w-4 h-4" />
                  {submitting ? "提交中..." : "提交申请"}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>

      {/* Tips */}
      <div className="mt-5 bg-white dark:bg-gray-800 rounded-xl border border-gray-100 dark:border-gray-700 p-4">
        <p className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">
          温馨提示
        </p>
        <ul className="text-xs text-gray-400 space-y-1">
          {[
            "申请后请耐心等待管理员审核，通常需要 1–3 个工作日",
            "审核通过后你将获得版主权限，可在板块内进行管理操作",
            "如申请被拒绝，可根据反馈修改后重新申请",
            "成为版主后请认真履行职责，长期不活跃可能被撤销权限",
          ].map((tip) => (
            <li key={tip} className="flex items-start gap-1.5">
              <span className="mt-1.5 w-1 h-1 rounded-full bg-gray-300 dark:bg-gray-600 shrink-0" />
              {tip}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
