// components/admin/PostModeration.tsx
import { useState } from "react";
import {
  CheckCircle,
  XCircle,
  Eye,
  User as UserIcon,
  FileText,
  AlertTriangle,
} from "lucide-react";
import DOMPurify from "dompurify";
import Image from "next/image";
import { useTranslations } from "next-intl";
import { useReviewPendingPosts } from "../hooks/useReviewPendingPosts";
import type { PostWithRisk } from "../hooks/useReviewPendingPosts";

export function ReviewManagement() {
  const t = useTranslations("Admin");
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState("");
  const [selectedPost, setSelectedPost] = useState<PostWithRisk | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [reviewNote, setReviewNote] = useState(""); // 审核备注/拒绝原因

  // 查询待审核列表 + 审核通过/拒绝操作
  const { posts, total, isLoading, refetch, approveMutation, rejectMutation } =
    useReviewPendingPosts(page, keyword);

  const handleOpenModal = (post: PostWithRisk) => {
    setSelectedPost(post);
    setReviewNote("");
    setIsModalOpen(true);
  };

  const handleApprove = () => {
    if (!selectedPost) return;
    approveMutation.mutate(
      { id: selectedPost.id, note: reviewNote.trim() || undefined },
      {
        onSuccess: (_, variables) => {
          if (selectedPost && selectedPost.id === variables.id) {
            setIsModalOpen(false);
            setSelectedPost(null);
            setReviewNote("");
          }
        },
      },
    );
  };

  const handleReject = () => {
    if (!selectedPost) return;
    rejectMutation.mutate(
      { id: selectedPost.id, reason: reviewNote.trim() || undefined },
      {
        onSuccess: (_, variables) => {
          if (selectedPost && selectedPost.id === variables.id) {
            setIsModalOpen(false);
            setSelectedPost(null);
            setReviewNote("");
          }
        },
      },
    );
  };

  const totalPages = Math.ceil(total / 20);

  const sanitizeHtml = (html: string) => ({
    __html: DOMPurify.sanitize(html, {
      ALLOWED_TAGS: [
        "b",
        "i",
        "em",
        "strong",
        "a",
        "p",
        "br",
        "ul",
        "ol",
        "li",
        "img",
        "h1",
        "h2",
        "h3",
        "h4",
        "blockquote",
        "code",
        "pre",
      ],
      ALLOWED_ATTR: ["href", "target", "src", "alt", "class", "id"],
    }),
  });

  return (
    <div className="space-y-4">
      {/* 搜索栏 */}
      <div className="flex justify-between items-center gap-4">
        <div className="flex gap-2 flex-1">
          <input
            type="text"
            placeholder={t("search_title_or_author")}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="input input-bordered flex-1"
          />
          <button className="btn btn-outline" onClick={() => setPage(1)}>
            {t("search")}
          </button>
          <button className="btn btn-ghost" onClick={() => refetch()}>
            {t("refresh")}
          </button>
        </div>
      </div>

      {/* 列表 */}
      {isLoading ? (
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : (
        <div className="space-y-3">
          {posts.map((post) => (
            <div
              key={post.id}
              className="card bg-base-100 shadow-sm border border-base-200 cursor-pointer hover:shadow-md transition-shadow"
              onClick={() => handleOpenModal(post)}
            >
              <div className="card-body p-4">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-medium">
                        {post.creation.author?.username ||
                          `用户${post.creation.author?.id}`}
                      </span>
                      <span className="text-xs text-gray-400">
                        #{post.creation.author?.id}
                      </span>
                      {post.risk_score && post.risk_score > 0 && (
                        <span className="badge badge-warning badge-sm">
                          <AlertTriangle className="w-3 h-3 mr-1" />
                          {t("risk_score")} {post.risk_score}
                        </span>
                      )}
                    </div>
                    <h3 className="font-semibold">{post.creation.title}</h3>
                    {post.risk_reason && (
                      <p className="text-sm text-orange-600 mt-1">
                        ⚠️ {post.risk_reason}
                      </p>
                    )}
                  </div>
                  <Eye className="w-5 h-5 text-gray-400" />
                </div>
              </div>
            </div>
          ))}
          {posts.length === 0 && (
            <div className="text-center py-8 text-gray-500">{t("no_pending_posts")}</div>
          )}
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-2 pt-4">
          <button
            className="btn btn-sm"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            {t("back")}
          </button>
          <span className="btn btn-sm btn-ghost">
            {t("page_info", { page, totalPages })}
          </span>
          <button
            className="btn btn-sm"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            {t("next_page")}
          </button>
        </div>
      )}

      {/* 模态框 - 详情 */}
      {isModalOpen && selectedPost && (
        <dialog
          className="modal modal-open"
          onClick={() => setIsModalOpen(false)}
        >
          <div
            className="modal-box max-w-3xl"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="font-bold text-lg mb-4">{t("post_review")}</h3>

            {/* 用户信息 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <UserIcon className="w-4 h-4" />
                <span className="font-medium">{t("publisher_info")}</span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                  {t("username_label")}：
                  {selectedPost.creation.author?.username ||
                    `ID ${selectedPost.creation.author?.id}`}
                </div>
                <div>{t("user_id_label")}：{selectedPost.creation.author?.id}</div>
                {selectedPost.creation.author?.email && (
                  <div>{t("email_label")}：{selectedPost.creation.author.email}</div>
                )}
                {selectedPost.creation.author?.role && (
                  <div>{t("role_label")}：{selectedPost.creation.author.role}</div>
                )}
                {selectedPost.creation.author?.created_at && (
                  <div>
                    {t("registration_time")}：
                    {new Date(
                      selectedPost.creation.author.created_at,
                    ).toLocaleString()}
                  </div>
                )}
              </div>
            </div>

            {/* 风险日志 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="w-4 h-4 text-warning" />
                <span className="font-medium">{t("risk_detection_log")}</span>
              </div>
              {selectedPost.risk_logs && selectedPost.risk_logs.length > 0 ? (
                <div className="space-y-2 max-h-40 overflow-y-auto">
                  {selectedPost.risk_logs.map((log) => (
                    <div
                      key={log.id}
                      className="text-sm border-l-2 border-warning pl-2"
                    >
                      <div className="flex justify-between">
                        <span className="font-mono text-xs">{log.level}</span>
                        <span className="text-xs text-gray-500">
                          {new Date(log.created_at).toLocaleString()}
                        </span>
                      </div>
                      <div className="text-gray-700">{t("rule_label")}：{log.rule}</div>
                      <div className="text-gray-500 text-xs break-all">
                        {t("match_content")}：{log.matched_content}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-sm text-gray-500">{t("no_risk_logs")}</div>
              )}
            </div>

            {/* 帖子内容 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <FileText className="w-4 h-4" />
                <span className="font-medium">{t("post_content")}</span>
              </div>
              <h4 className="font-semibold mb-2">
                {selectedPost.creation.title}
              </h4>
              <div
                className="text-sm whitespace-pre-wrap max-h-60 overflow-y-auto bg-base-100 p-2 rounded"
                dangerouslySetInnerHTML={sanitizeHtml(
                  selectedPost.creation.content,
                )}
              />
              {selectedPost.creation.cover_url && (
                <Image
                  src={selectedPost.creation.cover_url}
                  alt="封面"
                  width={200}
                  height={128}
                  className="mt-2 max-h-32 rounded object-cover"
                />
              )}
            </div>

            {/* 审核备注输入框 */}
            <div className="mb-4">
              <label className="label-text">
                {t("review_note_description")}
              </label>
              <textarea
                className="textarea textarea-bordered w-full"
                rows={3}
                placeholder={t("review_note_optional")}
                value={reviewNote}
                onChange={(e) => setReviewNote(e.target.value)}
              />
            </div>

            {/* 操作按钮 */}
            <div className="modal-action">
              <button className="btn" onClick={() => setIsModalOpen(false)}>
                {t("close")}
              </button>
              <button
                className="btn btn-success"
                onClick={handleApprove}
                disabled={approveMutation.isPending || rejectMutation.isPending}
              >
                <CheckCircle className="w-4 h-4 mr-1" />
                {t("approve")}
              </button>
              <button
                className="btn btn-error"
                onClick={handleReject}
                disabled={approveMutation.isPending || rejectMutation.isPending}
              >
                <XCircle className="w-4 h-4 mr-1" />
                {t("reject")}
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}
