// components/admin/PostModeration.tsx
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import toast from "react-hot-toast";
import {
  CheckCircle,
  XCircle,
  Eye,
  User as UserIcon,
  FileText,
  AlertTriangle,
} from "lucide-react";
import { adminApi } from "@/lib/api";
import type { Post } from "@/lib/api/types";
import DOMPurify from "dompurify";
import Image from "next/image";

// 扩展 Post 类型以包含风险信息
interface PostWithRisk extends Post {
  risk_score?: number;
  risk_reason?: string;
  risk_logs?: Array<{
    id: number;
    level: string;
    rule: string;
    matched_content: string;
    created_at: string;
  }>;
}

// 后端返回的分页数据结构
interface PendingPostsResponse {
  list: PostWithRisk[];
  total: number;
}

// 获取待审核帖子列表 - 修复返回类型
const fetchPendingPosts = async (params: {
  page: number;
  page_size: number;
  keyword?: string;
}): Promise<PendingPostsResponse> => {
  const res = await adminApi.listPendingPosts(params);
  // 确保返回的数据结构正确，并处理可能的 undefined
  const data = res.data.data;
  return {
    list: data?.list || [],
    total: data?.total || 0,
  };
};

export function ReviewManagement() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState("");
  const [selectedPost, setSelectedPost] = useState<PostWithRisk | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [reviewNote, setReviewNote] = useState(""); // 审核备注/拒绝原因

  // 查询待审核列表 - 移除泛型，让 TypeScript 自动推断
  const { data, isLoading, refetch } = useQuery({
    queryKey: ["admin-pending-posts", page, keyword],
    queryFn: () => fetchPendingPosts({ page, page_size: 20, keyword }),
    placeholderData: (prev) => prev,
  });

  // 审核通过 mutation
  const approveMutation = useMutation({
    mutationFn: ({ id, note }: { id: number; note?: string }) =>
      adminApi.approvePost(id, note),
    onSuccess: (_, variables) => {
      toast.success("帖子已通过审核");
      queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
      if (selectedPost && selectedPost.id === variables.id) {
        setIsModalOpen(false);
        setSelectedPost(null);
        setReviewNote("");
      }
    },
    onError: () => toast.error("审核通过失败"),
  });

  // 审核拒绝 mutation
  const rejectMutation = useMutation({
    mutationFn: ({ id, reason }: { id: number; reason?: string }) =>
      adminApi.rejectPost(id, reason),
    onSuccess: (_, variables) => {
      toast.success("帖子已拒绝");
      queryClient.invalidateQueries({ queryKey: ["admin-pending-posts"] });
      if (selectedPost && selectedPost.id === variables.id) {
        setIsModalOpen(false);
        setSelectedPost(null);
        setReviewNote("");
      }
    },
    onError: () => toast.error("拒绝失败"),
  });

  const handleOpenModal = (post: PostWithRisk) => {
    setSelectedPost(post);
    setReviewNote("");
    setIsModalOpen(true);
  };

  const handleApprove = () => {
    if (!selectedPost) return;
    approveMutation.mutate({
      id: selectedPost.id,
      note: reviewNote.trim() || undefined,
    });
  };

  const handleReject = () => {
    if (!selectedPost) return;
    rejectMutation.mutate({
      id: selectedPost.id,
      reason: reviewNote.trim() || undefined,
    });
  };

  const posts = data?.list ?? [];
  const total = data?.total ?? 0;
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
            placeholder="搜索标题或作者"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="input input-bordered flex-1"
          />
          <button className="btn btn-outline" onClick={() => setPage(1)}>
            搜索
          </button>
          <button className="btn btn-ghost" onClick={() => refetch()}>
            刷新
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
                        {post.author?.username || `用户${post.author_id}`}
                      </span>
                      <span className="text-xs text-gray-400">
                        #{post.author_id}
                      </span>
                      {post.risk_score && post.risk_score > 0 && (
                        <span className="badge badge-warning badge-sm">
                          <AlertTriangle className="w-3 h-3 mr-1" />
                          风险分 {post.risk_score}
                        </span>
                      )}
                    </div>
                    <h3 className="font-semibold">{post.title}</h3>
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
            <div className="text-center py-8 text-gray-500">暂无待审核帖子</div>
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
            上一页
          </button>
          <span className="btn btn-sm btn-ghost">
            第 {page} / {totalPages} 页
          </span>
          <button
            className="btn btn-sm"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            下一页
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
            <h3 className="font-bold text-lg mb-4">帖子审核</h3>

            {/* 用户信息 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <UserIcon className="w-4 h-4" />
                <span className="font-medium">发布人信息</span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                  用户名：
                  {selectedPost.author?.username ||
                    `ID ${selectedPost.author_id}`}
                </div>
                <div>用户ID：{selectedPost.author_id}</div>
                {selectedPost.author?.email && (
                  <div>邮箱：{selectedPost.author.email}</div>
                )}
                {selectedPost.author?.role && (
                  <div>角色：{selectedPost.author.role}</div>
                )}
                {selectedPost.author?.created_at && (
                  <div>
                    注册时间：
                    {new Date(selectedPost.author.created_at).toLocaleString()}
                  </div>
                )}
              </div>
            </div>

            {/* 风险日志 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="w-4 h-4 text-warning" />
                <span className="font-medium">风险检测日志</span>
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
                      <div className="text-gray-700">规则：{log.rule}</div>
                      <div className="text-gray-500 text-xs break-all">
                        匹配内容：{log.matched_content}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-sm text-gray-500">暂无风险日志</div>
              )}
            </div>

            {/* 帖子内容 */}
            <div className="mb-4 p-3 bg-base-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <FileText className="w-4 h-4" />
                <span className="font-medium">帖子内容</span>
              </div>
              <h4 className="font-semibold mb-2">{selectedPost.title}</h4>
              <div
                className="text-sm whitespace-pre-wrap max-h-60 overflow-y-auto bg-base-100 p-2 rounded"
                dangerouslySetInnerHTML={sanitizeHtml(selectedPost.content)}
              />
              {selectedPost.cover && (
                <Image
                  src={selectedPost.cover}
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
                审核备注（拒绝原因或通过说明）
              </label>
              <textarea
                className="textarea textarea-bordered w-full"
                rows={3}
                placeholder="选填，将记录在审核日志中"
                value={reviewNote}
                onChange={(e) => setReviewNote(e.target.value)}
              />
            </div>

            {/* 操作按钮 */}
            <div className="modal-action">
              <button className="btn" onClick={() => setIsModalOpen(false)}>
                关闭
              </button>
              <button
                className="btn btn-success"
                onClick={handleApprove}
                disabled={approveMutation.isPending || rejectMutation.isPending}
              >
                <CheckCircle className="w-4 h-4 mr-1" />
                通过
              </button>
              <button
                className="btn btn-error"
                onClick={handleReject}
                disabled={approveMutation.isPending || rejectMutation.isPending}
              >
                <XCircle className="w-4 h-4 mr-1" />
                拒绝
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  );
}